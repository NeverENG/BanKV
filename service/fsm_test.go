package service

import (
	"os"
	"testing"

	"github.com/NeverENG/BanKV/config"
	"github.com/NeverENG/BanKV/internal/Raft"
)

func setupTest(t *testing.T) (*FSM, func()) {
	oldWALPath := config.G.WALPath
	oldMaxSize := config.G.MaxMemTableSize
	oldPeers := config.G.Peers
	oldMe := config.G.Me

	config.G.WALPath = "test_service_wal.log"
	config.G.MaxMemTableSize = 100
	config.G.Peers = []string{"localhost:9000"}
	config.G.Me = 0

	fsm := NewFSM()

	cleanup := func() {
		// 清理文件
		os.Remove("test_service_wal.log")
		config.G.WALPath = oldWALPath
		config.G.MaxMemTableSize = oldMaxSize
		config.G.Peers = oldPeers
		config.G.Me = oldMe
	}

	return fsm, cleanup
}

func TestFSM_BasicOperation(t *testing.T) {
	fsm, cleanup := setupTest(t)
	defer cleanup()

	cmd := Command{
		Type:  "Put",
		Key:   []byte("key1"),
		Value: []byte("value1"),
	}

	cmdBytes, err := EncodeCommand(cmd)
	if err != nil {
		t.Fatalf("EncodeCommand failed: %v", err)
	}

	entry := Raft.LogEntry{
		Index:   0,
		Term:    1,
		Command: cmdBytes,
	}

	fsm.apply(entry)

	val, err := fsm.Get([]byte("key1"))
	if err != nil {
		t.Fatalf("fsm.Get failed: %v", err)
	}
	if string(val) != "value1" {
		t.Errorf("Expected 'value1', got '%s'", string(val))
	}
}

func TestFSM_DeleteOperation(t *testing.T) {
	fsm, cleanup := setupTest(t)
	defer cleanup()

	putCmd := Command{
		Type:  "Put",
		Key:   []byte("key1"),
		Value: []byte("value1"),
	}
	putBytes, _ := EncodeCommand(putCmd)
	fsm.apply(Raft.LogEntry{Index: 0, Term: 1, Command: putBytes})

	val, _ := fsm.Get([]byte("key1"))
	if string(val) != "value1" {
		t.Errorf("Expected 'value1', got '%s'", string(val))
	}

	delCmd := Command{
		Type: "Delete",
		Key:  []byte("key1"),
	}
	delBytes, _ := EncodeCommand(delCmd)
	fsm.apply(Raft.LogEntry{Index: 1, Term: 1, Command: delBytes})

	val, err := fsm.Get([]byte("key1"))
	if err == nil && val != nil {
		t.Errorf("Expected nil after delete, got '%s'", string(val))
	}
}

func TestFSM_UpdateOperation(t *testing.T) {
	fsm, cleanup := setupTest(t)
	defer cleanup()

	cmd1 := Command{Type: "Put", Key: []byte("key1"), Value: []byte("value1")}
	cmdBytes1, _ := EncodeCommand(cmd1)
	fsm.apply(Raft.LogEntry{Index: 0, Term: 1, Command: cmdBytes1})

	cmd2 := Command{Type: "Put", Key: []byte("key1"), Value: []byte("value2")}
	cmdBytes2, _ := EncodeCommand(cmd2)
	fsm.apply(Raft.LogEntry{Index: 1, Term: 1, Command: cmdBytes2})

	val, _ := fsm.Get([]byte("key1"))
	if string(val) != "value2" {
		t.Errorf("Expected 'value2', got '%s'", string(val))
	}
}

func TestFSM_DirectGetPut(t *testing.T) {
	fsm, cleanup := setupTest(t)
	defer cleanup()

	err := fsm.Put([]byte("key1"), []byte("value1"))
	if err != nil {
		t.Fatalf("fsm.Put failed: %v", err)
	}

	val, err := fsm.Get([]byte("key1"))
	if err != nil {
		t.Fatalf("fsm.Get failed: %v", err)
	}
	if string(val) != "value1" {
		t.Errorf("Expected 'value1', got '%s'", string(val))
	}
}
