package storage

import (
	"sync"

	"github.com/NeverENG/BanKV/config"
	zstorge "github.com/NeverENG/BanKV/internal/storage/zstorage"
)

type Engine struct {
	memTable *zstorge.MemTable

	mu sync.RWMutex
}

func (e *Engine) Put(key []byte, value []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()

	err := e.memTable.Put(key, value)
	if err != nil {
		return
	}
	if e.memTable.Size() > config.Global.MaxMemTableSize {
		e.memTable = zstorge.NewMemTable()
	}
}
