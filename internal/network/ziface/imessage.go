package ziface

type IMessage interface {
	GetMsgID() uint32
	GetData() []byte
	GetMsgLen() uint32

	SetMsgLen(uint32)
	SetData([]byte)
	SetMsgID(uint32)
}
