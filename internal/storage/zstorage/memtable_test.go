package zstorage

import (
	"fmt"
	"testing"
)

func BenchmarkMemTable_Set(b *testing.B) {
	memTable := NewMemTable()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}
}

func BenchmarkMemTable_Get(b *testing.B) {
	memTable := NewMemTable()

	// 预填充固定数量的数据
	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}

	b.ResetTimer()
	var value []byte
	var ok bool
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i%prefillSize))
		value, ok = memTable.Get(key)
	}
	_ = value
	_ = ok
}

func BenchmarkMemTable_Delete(b *testing.B) {
	memTable := NewMemTable()

	// 预填充固定数量的数据
	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}

	b.ResetTimer()
	var deleted bool
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i%prefillSize))
		deleted = memTable.Delete(key)
	}
	_ = deleted
}

func BenchmarkMemTable_SetStringKey(b *testing.B) {
	memTable := NewMemTable()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}
}

func BenchmarkMemTable_GetStringKey(b *testing.B) {
	memTable := NewMemTable()

	// 预填充固定数量的数据
	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}

	b.ResetTimer()
	var value []byte
	var ok bool
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i%prefillSize))
		value, ok = memTable.Get(key)
	}
	_ = value
	_ = ok
}
