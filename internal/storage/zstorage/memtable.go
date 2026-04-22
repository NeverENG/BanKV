package zstorage

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/NeverENG/BanKV/internal/storage/istorage"
)

var _ istorage.IMemTable = &MemTable{}

const (
	MAXL = 32
	P    = 0.5
)

// MemTable 基于跳表的内存表实现
type MemTable struct {
	size  int
	level int
	head  *SkipNode

	FlushChan chan bool

	wal *WAL
}

// SkipNode 跳表节点
type SkipNode struct {
	Next  []*SkipNode
	Key   []byte
	Value []byte
}

// NewMemTable 创建新的 MemTable
func NewMemTable() *MemTable {
	return &MemTable{
		size:      0,
		level:     0,
		FlushChan: make(chan bool, 1),
		head:      newSkipNode(MAXL, nil, nil),
		wal:       NewWAL(),
	}
}

// newSkipNode 创建新的跳表节点
func newSkipNode(level int, key []byte, value []byte) *SkipNode {
	return &SkipNode{
		Next:  make([]*SkipNode, level),
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
func (m *MemTable) Size() int {
	return m.size
}

// Get 获取指定 key 的值，如果不存在返回零值和 false
func (m *MemTable) Get(key []byte) ([]byte, error) {

	if m.head == nil {
		return nil, errors.New("NO DATA IN MEM")
	}

	p := m.head
	// 从最高层开始查找
	for i := m.level - 1; i >= 0; i-- {
		for p.Next[i] != nil && bytes.Compare(p.Next[i].Key, key) < 0 {
			p = p.Next[i]
		}
	}

	// 检查下一层的节点是否匹配
	p = p.Next[0]
	if p != nil && bytes.Compare(p.Key, key) == 0 {
		return p.Value, nil
	}
	return nil, nil
}

// Set 插入或更新键值对
func (m *MemTable) Put(key []byte, value []byte) error {

	if m.head == nil {
		return errors.New("NO DATA IN MEMTABLE")
	}

	err := m.wal.Write(istorage.LogEntry{Key: key, Value: value})

	if err != nil {
		fmt.Println("Error writing to WAL:", err)
	}

	// update 数组记录每一层需要更新的节点
	update := make([]*SkipNode, MAXL)
	p := m.head

	// 从最高层开始查找插入位置
	for i := m.level - 1; i >= 0; i-- {
		for p.Next[i] != nil && bytes.Compare(p.Next[i].Key, key) < 0 {
			p = p.Next[i]
		}
		update[i] = p
	}

	// 检查 key 是否已存在
	p = p.Next[0]
	if p != nil && bytes.Compare(p.Key, key) == 0 {
		// key 已存在，更新值
		p.Value = value
		return nil
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
	fmt.Println("Put Sucessful")
	return nil
}

// Delete 删除指定 key 的节点
func (m *MemTable) Delete(key []byte) error {

	if m.head == nil {
		return errors.New("NO DATA IN MEMTABLE")
	}

	// update 数组记录每一层需要更新的节点
	update := make([]*SkipNode, MAXL)
	p := m.head

	// 从最高层开始查找要删除的节点
	for i := m.level - 1; i >= 0; i-- {
		for p.Next[i] != nil && bytes.Compare(p.Next[i].Key, key) < 0 {
			p = p.Next[i]
		}
		update[i] = p
	}

	// 检查目标节点是否存在
	p = p.Next[0]
	if p == nil || bytes.Compare(p.Key, key) != 0 {
		// key 不存在
		fmt.Println("the key is not exist")
		return errors.New("KEY NO Found")
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
	return nil
}

func (m *MemTable) Sync() error {
	return m.wal.Sync()
}

func (m *MemTable) Clear() error {
	return m.wal.Clear()
}

func (m *MemTable) Close() error {
	return m.wal.Close()
}

func (m *MemTable) StartFlush() {
	m.FlushChan <- true
}

func (m *MemTable) Flush() {
	fmt.Printf("Flushing MemTable with %d entries...\n", m.size)

	allEntries := m.collectAllEntry()

	err := m.writeToSSTable(allEntries)
	if err != nil {
		fmt.Printf("Flush error: %v\n", err)
		return
	}

	m.resetMemTable()

	fmt.Println("Flush completed successfully")
}

func (m *MemTable) FlushWorker() {
	for {
		select {
		case <-m.FlushChan:
			fmt.Println("Flush")
			m.Flush()
		}
	}
}

func (m *MemTable) collectAllEntry() []istorage.LogEntry {
	LogEntrys := make([]istorage.LogEntry, m.Size())

	p := m.head.Next[0]
	for i := 0; i < m.Size(); i++ {
		LogEntrys[i] = istorage.LogEntry{
			Key:   p.Key,
			Value: p.Value,
		}
		p = p.Next[0]
	}
	return LogEntrys
}

func (m *MemTable) writeToSSTable(entries []istorage.LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// 跳表本身是有序的，collectAllEntry 按顺序遍历，所以 entries 已经有序

	// 2. 生成文件名并创建目录
	filename := fmt.Sprintf("sstable_%d.sst", time.Now().UnixNano())
	dir := "data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create data directory failed: %v", err)
	}
	fullPath := filepath.Join(dir, filename)

	// 3. 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create SSTable file failed: %v", err)
	}
	defer file.Close()

	// 4. 写入数据
	// 格式: [KeyLen(4B)][Key][ValueLen(4B)][Value]
	for _, entry := range entries {
		keyLen := uint32(len(entry.Key))
		valueLen := uint32(len(entry.Value))

		// 写入 Key 长度和内容
		if err := binary.Write(file, binary.BigEndian, keyLen); err != nil {
			return err
		}
		if _, err := file.Write(entry.Key); err != nil {
			return err
		}

		// 写入 Value 长度和内容
		if err := binary.Write(file, binary.BigEndian, valueLen); err != nil {
			return err
		}
		if _, err := file.Write(entry.Value); err != nil {
			return err
		}
	}

	// 5. 确保数据刷入磁盘
	return file.Sync()
}

// resetMemTable 重置内存表
func (m *MemTable) resetMemTable() error {
	m.head = newSkipNode(MAXL, nil, nil)
	m.size = 0
	m.level = 0

	err := m.Clear()
	return err
}
