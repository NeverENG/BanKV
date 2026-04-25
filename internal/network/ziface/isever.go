package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msg uint32, router IRouter)
	GetConnMgr() IConnManager

	SetConnStartFunc(func(connect IConnect))
	SetConnStopFunc(func(connect IConnect))
	CallConnStartFunc(connect IConnect)
	CallConnStopFunc(connect IConnect)
}
