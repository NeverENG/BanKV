package service

import (
	"encoding/json"

	"github.com/NeverENG/BanKV/config"
	"github.com/NeverENG/BanKV/internal/Raft"
	"github.com/NeverENG/BanKV/internal/storage"
	"github.com/NeverENG/BanKV/internal/storage/zstorage"
)

type Command struct {
	Type  string
	Key   []byte
	Value []byte
}

type FSM struct {
	raft    *Raft.Raft
	storage *storage.Engine
	applyCh chan Raft.LogEntry
}

// NewFSM 创建 FSM，自动从全局配置初始化 Raft 和存储
// 从 config.Global.Peers 和 config.Global.Me 获取集群配置
func NewFSM() *FSM {
	// 从全局配置获取集群信息
	peers := config.G.Peers
	me := config.G.Me

	// 初始化 Raft
	raft := Raft.NewRaft(peers, me)

	// 初始化存储
	memTable := zstorage.NewMemTable()
	store := storage.NewEngine(memTable)

	fsm := &FSM{
		raft:    raft,
		storage: store,
		applyCh: make(chan Raft.LogEntry, 100),
	}

	raft.RegisterApplyCh(fsm.applyCh)

	return fsm
}

// Run 运行 FSM
func (f *FSM) Run() {
	for entry := range f.applyCh {
		f.apply(entry)
	}
}

// apply 应用日志到存储
func (f *FSM) apply(entry Raft.LogEntry) {
	var cmd Command
	if err := json.Unmarshal(entry.Command, &cmd); err != nil {
		return
	}

	switch cmd.Type {
	case "Put":
		f.storage.Put(cmd.Key, cmd.Value)
	case "Delete":
		f.storage.Delete(cmd.Key)
	}
}

// Get 从存储获取值
func (f *FSM) Get(key []byte) ([]byte, error) {
	return f.storage.Get(key)
}

// Put 直接写入存储（仅用于测试，生产环境应通过 Raft 写入）
func (f *FSM) Put(key []byte, value []byte) error {
	return f.storage.Put(key, value)
}

// Delete 直接删除存储中的值（仅用于测试，生产环境应通过 Raft 写入）
func (f *FSM) Delete(key []byte) error {
	return f.storage.Delete(key)
}

// GetRaft 获取 Raft 实例
func (f *FSM) GetRaft() *Raft.Raft {
	return f.raft
}

// AppendEntry 通过 Raft 追加日志
func (f *FSM) AppendEntry(cmd Command) (int, error) {
	cmdBytes, err := EncodeCommand(cmd)
	if err != nil {
		return -1, err
	}
	return f.raft.AppendEntry(cmdBytes), nil
}

// EncodeCommand 编码命令为 JSON
func EncodeCommand(cmd Command) ([]byte, error) {
	return json.Marshal(cmd)
}
