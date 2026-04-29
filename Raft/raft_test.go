package Raft

import (
	"testing"
	"time"
)

func TestNewRaft(t *testing.T) {
	peers := []string{"localhost:8000", "localhost:8001", "localhost:8002"}
	r := NewRaft(peers, 0)

	if r == nil {
		t.Fatal("NewRaft returned nil")
	}

	if r.me != 0 {
		t.Errorf("Expected me to be 0, got %d", r.me)
	}

	if r.state != Follower {
		t.Errorf("Expected initial state to be Follower, got %v", r.state)
	}

	if r.Term != 0 {
		t.Errorf("Expected initial term to be 0, got %d", r.Term)
	}

	if len(r.peers) != 3 {
		t.Errorf("Expected 3 peers, got %d", len(r.peers))
	}
}

func TestGetState(t *testing.T) {
	peers := []string{"localhost:8000"}
	r := NewRaft(peers, 0)

	state, term := r.GetState()
	if state != Follower {
		t.Errorf("Expected state to be Follower, got %v", state)
	}
	if term != 0 {
		t.Errorf("Expected term to be 0, got %d", term)
	}
}

func TestGetLog(t *testing.T) {
	peers := []string{"localhost:8000"}
	r := NewRaft(peers, 0)

	log := r.GetLog()
	if len(log) != 0 {
		t.Errorf("Expected empty log, got %d entries", len(log))
	}
}

func TestAppendEntry(t *testing.T) {
	peers := []string{"localhost:8000"}
	r := NewRaft(peers, 0)

	index := r.AppendEntry([]byte("test command"))
	if index != -1 {
		t.Errorf("Expected -1 for non-leader, got %d", index)
	}
}

func TestElectionTimeout(t *testing.T) {
	peers := []string{"localhost:8000"}
	r := NewRaft(peers, 0)

	time.Sleep(400 * time.Millisecond)

	state, term := r.GetState()
	if state == Follower && term == 0 {
		t.Error("Expected election to start after timeout")
	}
}

func TestLeaderAppendsLog(t *testing.T) {
	peers := []string{"localhost:8000"}
	r := NewRaft(peers, 0)

	time.Sleep(400 * time.Millisecond)

	state, _ := r.GetState()
	if state != Leader {
		t.Skip("Not leader, skipping log append test")
	}

	index := r.AppendEntry([]byte("test command"))
	if index != 0 {
		t.Errorf("Expected index 0, got %d", index)
	}

	log := r.GetLog()
	if len(log) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(log))
	}

	if string(log[0].Command) != "test command" {
		t.Errorf("Expected command 'test command', got '%s'", string(log[0].Command))
	}
}

func TestLeaderSendsHeartbeats(t *testing.T) {
	peers := []string{"localhost:8000"}
	r := NewRaft(peers, 0)

	time.Sleep(400 * time.Millisecond)

	state, _ := r.GetState()
	if state != Leader {
		t.Skip("Not leader, skipping heartbeat test")
	}

	time.Sleep(100 * time.Millisecond)

	state, _ = r.GetState()
	if state != Leader {
		t.Error("Expected state to remain Leader after heartbeats")
	}
}
