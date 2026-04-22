package storage

import (
	"sync"

	"github.com/NeverENG/BanKV/config"
	"github.com/NeverENG/BanKV/internal/storage/istorage"
)

type Engine struct {
	memTable istorage.IMemTable

	mu sync.RWMutex
}

func NewEngine(memTable istorage.IMemTable) *Engine {
	return &Engine{
		memTable: memTable,
	}
}

func (e *Engine) Put(key []byte, value []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	err := e.memTable.Put(key, value)
	if err != nil {
		return err
	}
	if e.memTable.Size() > config.Global.MaxMemTableSize {
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
