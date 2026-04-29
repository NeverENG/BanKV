package Raft

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MinElectionTimeout = 150 * time.Millisecond
	MaxElectionTimeout = 300 * time.Millisecond
	HeartbeatInterval  = 50 * time.Millisecond
)

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type LogEntry struct {
	Index   int
	Term    int
	Command []byte
}

type Raft struct {
	peers    []string
	me       int
	state    State
	votedFor int
	Term     int
	mu       sync.Mutex

	electionTimeout time.Duration
	timer           *time.Timer
	heartbeatTicker *time.Ticker

	commitIndex int
	lastApplied int
	nextIndex   []int
	matchIndex  []int
	log         []LogEntry

	electionCh  chan bool
	heartbeatCh chan bool
	ApplyCh     chan LogEntry

	LastIncludedIndex int64
	LastIncludedTerm  int64

	wal     *RaftWAL
	addrMap map[int]string

	commitCond *sync.Cond
}

func NewRaft(peers []string, me int) *Raft {
	addrMap := make(map[int]string)
	for i, addr := range peers {
		addrMap[i] = addr
	}

	r := &Raft{
		peers:           peers,
		me:              me,
		state:           Follower,
		votedFor:        -1,
		Term:            0,
		electionTimeout: MinElectionTimeout + time.Duration(rand.Int63n(int64(MaxElectionTimeout-MinElectionTimeout))),
		commitIndex:     -1,
		lastApplied:     -1,
		nextIndex:       make([]int, len(peers)),
		matchIndex:      make([]int, len(peers)),
		log:             make([]LogEntry, 0),
		electionCh:      make(chan bool),
		heartbeatCh:     make(chan bool),
		ApplyCh:         make(chan LogEntry, 100),
		addrMap:         addrMap,
	}

	wal, _ := NewRaftWAL("raft_data")

	r.wal = wal
	term, votedFor, _ := r.wal.LoadState()
	r.Term = term
	r.votedFor = votedFor
	logs, _ := wal.LoadLogs()
	r.log = logs

	r.commitCond = sync.NewCond(&r.mu)

	go r.electionLoop()

	return r
}

func (r *Raft) Start() {
	if r.state == Leader {
		r.startHeartbeatLoop()
	}
}

func (r *Raft) electionLoop() {
	for {
		timeout := MinElectionTimeout + time.Duration(rand.Int63n(int64(MaxElectionTimeout-MinElectionTimeout)))
		r.timer = time.NewTimer(timeout)

		select {
		case <-r.timer.C:
			r.startElection()
		case <-r.heartbeatCh:
			r.timer.Reset(timeout)
		case <-r.electionCh:
			r.timer.Reset(timeout)
		}
	}
}

func (r *Raft) startElection() {
	r.mu.Lock()

	if r.state == Leader {
		r.mu.Unlock()
		return
	}

	fmt.Printf("[RAFT] Starting election, current state=%v, Term=%d\n", r.state, r.Term)

	r.state = Candidate
	r.Term++
	r.votedFor = r.me
	r.wal.SaveState(r.Term, r.votedFor)

	lastLogIndex := -1
	lastLogTerm := 0
	if len(r.log) > 0 {
		lastLogIndex = len(r.log) - 1
		lastLogTerm = r.log[lastLogIndex].Term
	}

	args := &RequestVoteArgs{
		Term:         r.Term,
		CandidateID:  r.me,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
	}

	peerCount := len(r.peers) - 1
	votes := 1
	voteCh := make(chan bool, peerCount+1)
	voteCh <- true

	for i := range r.peers {
		if i == r.me {
			continue
		}

		go func(peerID int) {
			reply, err := r.SendRequestVote(r.addrMap[peerID], args)
			if err != nil {
				voteCh <- false
				return
			}

			r.mu.Lock()
			defer r.mu.Unlock()

			if reply.Term > r.Term {
				r.Term = reply.Term
				r.state = Follower
				r.votedFor = -1
				voteCh <- false
				return
			}

			if reply.Term == r.Term && reply.VoteGranted {
				voteCh <- true
			} else {
				voteCh <- false
			}
		}(i)
	}

	r.mu.Unlock()

	// 自己投自己一票

	// 单节点模式：自己一票就超过半数，直接成为 Leader
	if votes > len(r.peers)/2 {
		r.mu.Lock()
		if r.state == Candidate {
			r.becomeLeader()
		}
		r.mu.Unlock()
		return
	}

	// 等待投票结果或超时
	timeout := time.After(500 * time.Millisecond)
	for j := 0; j < peerCount; j++ {
		select {
		case voteGranted := <-voteCh:
			if voteGranted {
				votes++
				// 获得多数票，成为 Leader
				if votes > len(r.peers)/2 {
					r.mu.Lock()
					if r.state == Candidate {
						r.becomeLeader()
					}
					r.mu.Unlock()
					return
				}
			}
		case <-timeout:
			// 选举超时，重置为 Follower
			r.mu.Lock()
			if r.state == Candidate {
				r.state = Follower
				r.votedFor = -1
			}
			r.mu.Unlock()
			return
		}
	}
}

func (r *Raft) becomeLeader() {
	fmt.Printf("[RAFT] Becoming Leader, Term=%d\n", r.Term)
	r.state = Leader

	for i := range r.peers {
		r.nextIndex[i] = len(r.log)
		r.matchIndex[i] = -1
	}

	fmt.Printf("[RAFT] Started heartbeat loop\n")
	r.startHeartbeatLoop()
}

func (r *Raft) startHeartbeatLoop() {
	if r.heartbeatTicker != nil {
		r.heartbeatTicker.Stop()
	}

	r.heartbeatTicker = time.NewTicker(HeartbeatInterval)
	go func() {
		for r.state == Leader {
			<-r.heartbeatTicker.C
			r.SendHeartBeat()
		}
	}()
}

func (r *Raft) SendHeartBeat() {
	r.mu.Lock()
	if r.state != Leader {
		r.mu.Unlock()
		return
	}

	for i := range r.peers {
		if i == r.me {
			continue
		}

		prevLogIndex := r.nextIndex[i] - 1
		prevLogTerm := 0
		if prevLogIndex >= 0 && prevLogIndex < len(r.log) {
			prevLogTerm = r.log[prevLogIndex].Term
		}

		args := &AppendEntriesArgs{
			Term:         r.Term,
			LeaderID:     r.me,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			Entries:      []LogEntry{},
			LeaderCommit: r.commitIndex,
		}

		r.mu.Unlock()

		go func(peerID int, args *AppendEntriesArgs) {
			reply, err := r.SendAppendEntries(r.addrMap[peerID], args)
			if err != nil {
				r.mu.Lock()
				if r.state == Leader {
					r.nextIndex[peerID]--
				}
				r.mu.Unlock()
				return
			}

			r.mu.Lock()
			defer r.mu.Unlock()

			if reply.Term > r.Term {
				r.Term = reply.Term
				r.state = Follower
				r.votedFor = -1
				r.heartbeatTicker.Stop()
				return
			}

			if reply.Success {
				r.nextIndex[peerID] = len(r.log)
				r.matchIndex[peerID] = len(r.log) - 1
				r.updateCommitIndex()
			} else {
				r.nextIndex[peerID]--
			}
		}(i, args)

		r.mu.Lock()
	}
	r.mu.Unlock()
}

func (r *Raft) updateCommitIndex() {
	if r.state != Leader {
		return
	}

	for n := len(r.log) - 1; n > r.commitIndex; n-- {
		count := 1
		for i := range r.peers {
			if i != r.me && r.matchIndex[i] >= n {
				count++
			}
		}
		if count > len(r.peers)/2 && r.log[n].Term == r.Term {
			r.commitIndex = n
			r.applyCommittedLogs()
			r.commitCond.Broadcast()
			break
		}
	}
}

func (r *Raft) applyCommittedLogs() {
	for r.lastApplied < r.commitIndex {
		r.lastApplied++
		if r.ApplyCh != nil {
			r.ApplyCh <- r.log[r.lastApplied]
		}
	}
}

func (r *Raft) AppendEntry(command []byte) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != Leader {
		fmt.Printf("[RAFT] AppendEntry failed: not leader, state=%v\n", r.state)
		return -1
	}

	entry := LogEntry{
		Index:   len(r.log),
		Term:    r.Term,
		Command: command,
	}
	r.log = append(r.log, entry)
	r.wal.AppendLog(entry)

	fmt.Printf("[RAFT] Appended entry: Index=%d, Term=%d\n", entry.Index, entry.Term)

	// 单节点模式：立即提交
	if len(r.peers) == 1 {
		r.commitIndex = entry.Index
		fmt.Printf("[RAFT] Single node mode, committed index: %d\n", r.commitIndex)
		r.applyCommittedLogs()
		r.commitCond.Broadcast()
	} else {
		r.replicateLog()
	}

	return entry.Index
}

func (r *Raft) replicateLog() {
	if r.state != Leader {
		return
	}

	for i := range r.peers {
		if i == r.me {
			continue
		}

		prevLogIndex := r.nextIndex[i] - 1
		prevLogTerm := 0
		if prevLogIndex >= 0 && prevLogIndex < len(r.log) {
			prevLogTerm = r.log[prevLogIndex].Term
		}

		var entries []LogEntry
		if r.nextIndex[i] < len(r.log) {
			entries = r.log[r.nextIndex[i]:]
		}

		args := &AppendEntriesArgs{
			Term:         r.Term,
			LeaderID:     r.me,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			Entries:      entries,
			LeaderCommit: r.commitIndex,
		}

		go func(peerID int, args *AppendEntriesArgs) {
			reply, err := r.SendAppendEntries(r.addrMap[peerID], args)
			if err != nil {
				r.mu.Lock()
				if r.state == Leader {
					r.nextIndex[peerID]--
				}
				r.mu.Unlock()
				return
			}

			r.mu.Lock()
			defer r.mu.Unlock()

			if reply.Term > r.Term {
				r.Term = reply.Term
				r.state = Follower
				r.votedFor = -1
				r.heartbeatTicker.Stop()
				return
			}

			if reply.Success {
				r.nextIndex[peerID] = len(r.log)
				r.matchIndex[peerID] = len(r.log) - 1
				r.updateCommitIndex()
			} else {
				r.nextIndex[peerID]--
			}
		}(i, args)
	}
}

func (r *Raft) SendInstallSnapshot(peer int) {
	r.mu.Lock()
	// 构建 rpc
	r.mu.Unlock()

}

func (r *Raft) WaitCommitIndex(index int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for r.commitIndex < index {
		r.commitCond.Wait()
	}
}

func (r *Raft) GetState() (State, int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state, r.Term
}

func (r *Raft) GetLog() []LogEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	logCopy := make([]LogEntry, len(r.log))
	copy(logCopy, r.log)
	return logCopy
}

func (r *Raft) GetApplyCh() chan LogEntry {
	return r.ApplyCh
}

// GetCommitIndex 获取当前提交索引
func (r *Raft) GetCommitIndex() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.commitIndex
}
