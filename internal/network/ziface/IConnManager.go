package ziface

type IConnManager interface {
	Add(conn IConnect)
	Remove(conn IConnect)
	Get(connID uint32) IConnect
	Len() int
	ClearConn()
}
