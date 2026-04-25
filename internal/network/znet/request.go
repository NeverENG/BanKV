package znet

import "github.com/NeverENG/BanKV/internal/network/ziface"

type Request struct {
	msg  ziface.IMessage
	conn ziface.IConnect
}

var _ ziface.IRequest = &Request{}

func (req *Request) GetMsgData() []byte {
	return req.msg.GetData()
}

func (req *Request) GetMsgID() uint32 {
	return req.msg.GetMsgID()
}

func (req *Request) GetConnection() ziface.IConnect {
	return req.conn
}
