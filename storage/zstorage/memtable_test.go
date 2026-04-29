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

	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}

	b.ResetTimer()
	var value []byte
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i%prefillSize))
		value, _ = memTable.Get(key)
	}
	_ = value
}

func BenchmarkMemTable_Delete(b *testing.B) {
	memTable := NewMemTable()

	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i%prefillSize))
		memTable.Delete(key)
	}
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

	const prefillSize = 10000
	for i := 0; i < prefillSize; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := []byte(fmt.Sprintf("value_%d", i))
		memTable.Put(key, value)
	}

	b.ResetTimer()
	var value []byte
	for i := 0; i < b.N; i++ {
		key := []byte(fmt.Sprintf("key_%d", i%prefillSize))
		value, _ = memTable.Get(key)
	}
	_ = value
}
