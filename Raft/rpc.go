package Raft

import (
	"net/rpc"
)

type RequestVoteArgs struct {
	Term         int
	CandidateID  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term         int
	LeaderID     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

type InstallSnapshotArgs struct {
	Term              int
	LeaderID          int
	Data              []byte
	LastIncludedIndex int32
	LastIncludedTerm  int
}

type InstallSnapshotReply struct {
	Term int
}

type RaftRPC struct {
	raft *Raft
}

func NewRaftRPC(raft *Raft) *RaftRPC {
	return &RaftRPC{raft: raft}
}

func (r *RaftRPC) RegisterRPC(server *rpc.Server) {
	server.Register(r)
}

func (r *RaftRPC) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	r.raft.mu.Lock()
	defer r.raft.mu.Unlock()

	if args.Term < r.raft.Term {
		reply.Term = r.raft.Term
		reply.VoteGranted = false
		return nil
	}

	if args.Term > r.raft.Term {
		r.raft.Term = args.Term
		r.raft.state = Follower
		r.raft.votedFor = -1
	}

	votedForMe := r.raft.votedFor == -1 || r.raft.votedFor == args.CandidateID
	logUpToDate := r.isLogUpToDate(args.LastLogIndex, args.LastLogTerm)

	if votedForMe && logUpToDate {
		r.raft.votedFor = args.CandidateID
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}

	reply.Term = r.raft.Term
	return nil
}

func (r *RaftRPC) isLogUpToDate(candidateLastIndex, candidateLastTerm int) bool {
	if len(r.raft.log) == 0 {
		return true
	}

	lastLog := r.raft.log[len(r.raft.log)-1]
	if candidateLastTerm > lastLog.Term {
		return true
	}
	if candidateLastTerm == lastLog.Term && candidateLastIndex >= len(r.raft.log)-1 {
		return true
	}
	return false
}

func (r *RaftRPC) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error {
	r.raft.mu.Lock()
	defer r.raft.mu.Unlock()

	if args.Term < r.raft.Term {
		reply.Term = r.raft.Term
		reply.Success = false
		return nil
	}

	if args.Term > r.raft.Term {
		r.raft.Term = args.Term
		r.raft.state = Follower
		r.raft.votedFor = -1
	}

	if len(r.raft.log) > 0 && (args.PrevLogIndex >= len(r.raft.log) || r.raft.log[args.PrevLogIndex].Term != args.PrevLogTerm) {
		reply.Success = false
		reply.Term = r.raft.Term
		return nil
	}

	for i, entry := range args.Entries {
		logIndex := args.PrevLogIndex + i + 1
		if logIndex < len(r.raft.log) && r.raft.log[logIndex].Term != entry.Term {
			r.raft.log = r.raft.log[:logIndex]
		}
		if logIndex >= len(r.raft.log) {
			r.raft.log = append(r.raft.log, entry)
		}
	}

	if args.LeaderCommit > r.raft.commitIndex {
		r.raft.commitIndex = min(args.LeaderCommit, len(r.raft.log)-1)
		r.applyCommittedLogs()
	}

	reply.Success = true
	reply.Term = r.raft.Term
	return nil
}

func (r *RaftRPC) applyCommittedLogs() {
	for r.raft.lastApplied < r.raft.commitIndex {
		r.raft.lastApplied++
		if r.raft.ApplyCh != nil {
			r.raft.ApplyCh <- r.raft.log[r.raft.lastApplied]
		}
	}
}

// 被调用端
func (r *RaftRPC) InstallSnapshot(args *InstallSnapshotArgs) (*InstallSnapshotReply, error) {
	r.raft.mu.Lock()
	defer r.raft.mu.Unlock()

	return nil, nil
}

func (r *Raft) SendRequestVote(serverAddr string, args *RequestVoteArgs) (*RequestVoteReply, error) {
	client, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var reply RequestVoteReply
	err = client.Call("RaftRPC.RequestVote", args, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (r *Raft) SendAppendEntries(serverAddr string, args *AppendEntriesArgs) (*AppendEntriesReply, error) {
	client, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var reply AppendEntriesReply
	err = client.Call("RaftRPC.AppendEntries", args, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (r *Raft) SendInstallSnapshot(serverAddr string, args *InstallSnapshotArgs) (*InstallSnapshotReply, error) {
	client, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var reply InstallSnapshotReply
	err = client.Call("RaftRPC.InstallSnapshot", args, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
