package ziface

type IRequest interface {
	GetConnection() IConnect
	GetMsgData() []byte
	GetMsgID() uint32
}
