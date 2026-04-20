package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type GlobalConfig struct {
	Name string
	Port int
	Host string

	WALPath     string
	SSTablePath string

	MaxMemTableSize int
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
	}
	global.Init()
	return global
}

var Global = NewGlobalConfig()
