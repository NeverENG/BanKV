package storage

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/NeverENG/BanKV/config"
	"github.com/NeverENG/BanKV/internal/storage/istorage"
	"github.com/NeverENG/BanKV/internal/storage/zstorage"
)

type StorageCommand struct {
	Type  string
	Key   []byte
	Value []byte
}

type Engine struct {
	memTable istorage.IMemTable
	mu       sync.RWMutex
	applyCh  chan StorageCommand
}

func NewEngine(memTable istorage.IMemTable) *Engine {
	e := &Engine{
		memTable: memTable,
		applyCh:  make(chan StorageCommand, 100),
	}
	go e.applyWorker()
	return e
}

func (e *Engine) Put(key []byte, value []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	err := e.memTable.Put(key, value)
	if err != nil {
		return err
	}
	if e.memTable.Size() > config.G.MaxMemTableSize {
		e.memTable.StartFlush()
	}
	return nil
}

func (e *Engine) Get(key []byte) ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.memTable.Get(key)
}

func (e *Engine) Delete(key []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.memTable.Delete(key)
}

func (e *Engine) Apply(cmd StorageCommand) error {
	e.applyCh <- cmd
	return nil
}

func (e *Engine) GetApplyCh() chan StorageCommand {
	return e.applyCh
}

func (e *Engine) applyWorker() {
	for cmd := range e.applyCh {
		switch cmd.Type {
		case "Put":
			e.Put(cmd.Key, cmd.Value)
		case "Delete":
			e.Delete(cmd.Key)
		}
	}
}

func (e *Engine) ApplySnapshot(data []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var state map[string]string
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal snapshot: %v", err)
	}

	e.memTable = NewMemTable()

	count := 0
	for key, value := range state {
		if err := e.memTable.Put([]byte(key), []byte(value)); err != nil {
			return fmt.Errorf("failed to restore key %s: %v", key, err)
		}
		count++
	}

	fmt.Printf("[ENGINE] Snapshot applied, restored %d keys\n", count)
	return nil
}

func (e *Engine) GetSnapshotData() ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	allData := make(map[string]string)

	entries := e.CollectAllEntries()
	for _, entry := range entries {
		allData[string(entry.Key)] = string(entry.Value)
	}

	return json.Marshal(allData)
}

func (e *Engine) CollectAllEntries() []istorage.LogEntry {
	return e.memTable.CollectAllEntries()
}
