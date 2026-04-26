package utils

import (
	"github.com/NeverENG/BanKV/internal/network/znet"
)

type DataPack struct {
	*znet.DataPack
}

func NewDataPack() *DataPack {
	return &DataPack{
		znet.NewDataPack(),
	}
}
