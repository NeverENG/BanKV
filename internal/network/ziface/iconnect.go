package ziface

import "net"

type IConnect interface {
	Start()
	Stop()
	GetTcpConn() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error
	SendBuffMsg(msgId uint32, data []byte) error

	SetProperty(key string, value interface{})
	GetProperty(key string) interface{}
	RemoveProperty(key string)
}

// type HandleFunc func(*net.TCPConn, []byte, int) error
