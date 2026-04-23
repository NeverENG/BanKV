package config

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"strconv"
)

type GlobalConfig struct {
	Name string
	Port int
	Host string

	WALPath     string
	SSTablePath string

	MaxMemTableSize int

	// Raft 集群配置
	Peers []string // 集群中所有节点的地址
	Me    int      // 当前节点在 Peers 中的索引（0-based）
}

func (g *GlobalConfig) Init() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		slog.Error("[ERROR]:READ CONFIG ERROR !")
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		slog.Error("[ERROR]:CONFIG INIT ERROR")
	}
	slog.Info("[INFO]:CONFIG INIT SUCCESS")
}

func NewGlobalConfig() *GlobalConfig {
	global := &GlobalConfig{
		MaxMemTableSize: 1024,
		WALPath:         "../../../log/wal.log",
		SSTablePath:     "../../../log",
		Peers:           []string{"localhost:8080"}, // 默认单节点
		Me:              0,                          // 默认节点ID
	}
	global.Init()
	global.ParseFlags()
	return global
}

// ParseFlags 解析命令行参数
func (g *GlobalConfig) ParseFlags() {
	// 定义命令行参数
	meFlag := flag.Int("me", -1, "Current node index in peers list")

	flag.Parse()

	// 处理命令行参数
	if *meFlag >= 0 {
		g.Me = *meFlag
		slog.Info("[INFO]:ME SET BY FLAG", "me", g.Me)
	}

	// 处理环境变量（优先级低于命令行参数）
	if g.Me < 0 {
		if meEnv := os.Getenv("RAFT_ME"); meEnv != "" {
			if meInt, err := strconv.Atoi(meEnv); err == nil {
				g.Me = meInt
				slog.Info("[INFO]:ME SET BY ENV", "me", g.Me)
			}
		}
	}

	// 验证配置
	if g.Me < 0 || g.Me >= len(g.Peers) {
		slog.Error("[ERROR]:INVALID ME VALUE", "me", g.Me, "peers_len", len(g.Peers))
		os.Exit(1)
	}

	slog.Info("[INFO]:CONFIG FINALIZED", "peers", g.Peers, "me", g.Me)
}

var Global = NewGlobalConfig()
