package storage

import (
	"fmt"
	"math/rand"
)

const (
	MAXL = 32
	P    = 0.5
)

// MemTable 基于跳表的内存表实现
type MemTable[K comparable, V any] struct {
	size   int
	level  int
	head   *SkipNode[K, V]
	compare func(K, K) int // 返回 -1, 0, 1 分别表示小于、等于、大于
}

// SkipNode 跳表节点
type SkipNode[K comparable, V any] struct {
	Next  []*SkipNode[K, V]
	Key   K
	Value V
}

// NewMemTable 创建新的 MemTable
func NewMemTable[K comparable, V any](compare func(K, K) int) *MemTable[K, V] {
	return &MemTable[K, V]{
		size:   0,
		level:  0,
		head:   newSkipNode[K, V](MAXL, *new(K), *new(V)),
		compare: compare,
	}
}

// newSkipNode 创建新的跳表节点
func newSkipNode[K comparable, V any](level int, key K, value V) *SkipNode[K, V] {
	return &SkipNode[K, V]{
		Next:  make([]*SkipNode[K, V], level),
		Key:   key,
		Value: value,
	}
}

// randomLevel 生成随机层级
func randomLevel() int {
	level := 1
	for rand.Float64() < P && level < MAXL {
		level++
	}
	return level
}

// Size 返回跳表中元素个数
func (m *MemTable[K, V]) Size() int {
	return m.size
}

// Get 获取指定 key 的值，如果不存在返回零值和 false
func (m *MemTable[K, V]) Get(key K) (V, bool) {
	var zeroValue V
	if m.head == nil {
		return zeroValue, false
	}

	p := m.head
	// 从最高层开始查找
	for i := m.level - 1; i >= 0; i-- {
		for p.Next[i] != nil && m.compare(p.Next[i].Key, key) < 0 {
			p = p.Next[i]
		}
	}

	// 检查下一层的节点是否匹配
	p = p.Next[0]
	if p != nil && m.compare(p.Key, key) == 0 {
		return p.Value, true
	}
	return zeroValue, false
}

// Set 插入或更新键值对
func (m *MemTable[K, V]) Set(key K, value V) {
	if m.head == nil {
		return
	}

	// update 数组记录每一层需要更新的节点
	update := make([]*SkipNode[K, V], MAXL)
	p := m.head

	// 从最高层开始查找插入位置
	for i := m.level - 1; i >= 0; i-- {
		for p.Next[i] != nil && m.compare(p.Next[i].Key, key) < 0 {
			p = p.Next[i]
		}
		update[i] = p
	}

	// 检查 key 是否已存在
	p = p.Next[0]
	if p != nil && m.compare(p.Key, key) == 0 {
		// key 已存在，更新值
		p.Value = value
		return
	}

	// 生成新节点的随机层级
	newLevel := randomLevel()
	if newLevel > m.level {
		// 如果新层级大于当前最大层级，更新 update 数组
		for i := m.level; i < newLevel; i++ {
			update[i] = m.head
		}
		m.level = newLevel
	}

	// 创建新节点
	newNode := newSkipNode(newLevel, key, value)

	// 在每一层插入新节点
	for i := 0; i < newLevel; i++ {
		newNode.Next[i] = update[i].Next[i]
		update[i].Next[i] = newNode
	}

	m.size++
}

// Delete 删除指定 key 的节点
func (m *MemTable[K, V]) Delete(key K) bool {
	if m.head == nil {
		return false
	}

	// update 数组记录每一层需要更新的节点
	update := make([]*SkipNode[K, V], MAXL)
	p := m.head

	// 从最高层开始查找要删除的节点
	for i := m.level - 1; i >= 0; i-- {
		for p.Next[i] != nil && m.compare(p.Next[i].Key, key) < 0 {
			p = p.Next[i]
		}
		update[i] = p
	}

	// 检查目标节点是否存在
	p = p.Next[0]
	if p == nil || m.compare(p.Key, key) != 0 {
		// key 不存在
		return false
	}

	// 在每一层删除节点
	for i := 0; i < m.level; i++ {
		if update[i].Next[i] != p {
			break
		}
		update[i].Next[i] = p.Next[i]
	}

	// 更新最大层级
	for m.level > 0 && m.head.Next[m.level-1] == nil {
		m.level--
	}

	m.size--
	return true
}

// String 返回跳表的字符串表示（用于调试）
func (m *MemTable[K, V]) String() string {
	return fmt.Sprintf("MemTable{size: %d, level: %d}", m.size, m.level)
}
