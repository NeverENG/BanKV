package ziface

type IMsgHandle interface {
	AddRouter(msgId uint32, router IRouter)
	DoMsgHandle(request IRequest)
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}
