package znet

import (
	"sync"

	"github.com/NeverENG/BanKV/internal/network/ziface"
)

type ConnManager struct {
	mu          sync.RWMutex
	connections map[uint32]ziface.IConnect
}

var _ ziface.IConnManager = &ConnManager{}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnect),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnect) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[conn.GetConnID()] = conn
}

func (cm *ConnManager) Remove(conn ziface.IConnect) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, conn.GetConnID())
}

func (cm *ConnManager) Get(connId uint32) ziface.IConnect {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if conn, ok := cm.connections[connId]; ok {
		return conn
	}
	return nil
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for connId, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connId)
	}
}
