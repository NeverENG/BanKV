package banNet

import (
	"context"
	"fmt"

	"github.com/NeverENG/BanKV/config"
	"github.com/NeverENG/BanKV/internal/network/banIface"
)

type MsgHandle struct {
	Arip           map[uint32]banIface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan banIface.IRequest
	ctx            context.Context
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Arip:           make(map[uint32]banIface.IRouter),
		WorkerPoolSize: config.G.WorkerPoolSize,
		TaskQueue:      make([]chan banIface.IRequest, config.G.WorkerPoolSize),
		ctx:            context.Background(),
	}
}

var _ banIface.IMsgHandle = &MsgHandle{}

func (m *MsgHandle) AddRouter(msgID uint32, r banIface.IRouter) {
	if _, ok := m.Arip[msgID]; ok {
		fmt.Println("duplicate Arip:", m.Arip)
		return
	}
	m.Arip[msgID] = r
}

func (m *MsgHandle) DoMsgHandle(request banIface.IRequest) {
	handler, ok := m.Arip[request.GetMsgID()]
	if !ok {
		fmt.Println("[ERROR] 该 Msgid 没有注册:", request.GetMsgID())
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan banIface.IRequest, config.G.MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request banIface.IRequest) {

	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	m.TaskQueue[workerID] <- request
}

func (m *MsgHandle) StartOneWorker(workerId int, taskQueue chan banIface.IRequest) {
	fmt.Println("Worker id:", workerId, "is started")
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandle(request)

		case <-m.ctx.Done():
			return
		}
	}
}
