package storage_test

import (
	"fmt"
	"testing"

	"github.com/NeverENG/BanKV/internal/storage"
)

func BenchmarkMemTable_Set(b *testing.B) {
	memTable := storage.NewMemTable[int, int](func(a, b int) int {
		return a - b
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memTable.Set(i, i)
	}
}

func BenchmarkMemTable_Get(b *testing.B) {
	memTable := storage.NewMemTable[int, int](func(a, b int) int {
		return a - b
	})

	// 预填充固定数量的数据
	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		memTable.Set(i, i)
	}

	b.ResetTimer()
	var value int
	var ok bool
	for i := 0; i < b.N; i++ {
		value, ok = memTable.Get(i % prefillSize)
	}
	_ = value
	_ = ok
}

func BenchmarkMemTable_Delete(b *testing.B) {
	memTable := storage.NewMemTable[int, int](func(a, b int) int {
		return a - b
	})

	// 预填充固定数量的数据
	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		memTable.Set(i, i)
	}

	b.ResetTimer()
	var deleted bool
	for i := 0; i < b.N; i++ {
		deleted = memTable.Delete(i % prefillSize)
	}
	_ = deleted
}

func BenchmarkMemTable_SetStringKey(b *testing.B) {
	memTable := storage.NewMemTable[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		memTable.Set(key, value)
	}
}

func BenchmarkMemTable_GetStringKey(b *testing.B) {
	memTable := storage.NewMemTable[string, string](func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// 预填充固定数量的数据
	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := fmt.Sprintf("value_%d", i)
		memTable.Set(key, value)
	}

	b.ResetTimer()
	var value string
	var ok bool
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%prefillSize)
		value, ok = memTable.Get(key)
	}
	_ = value
	_ = ok
}
