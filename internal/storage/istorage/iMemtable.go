package istorage

type IMemTable interface {
	// 基本操作
	Size() int
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error

	// Flush 相关操作
	StartFlush()    // 触发 flush 信号
	Flush()         // 执行 flush 操作
	FlushWorker()   // 启动后台 flush 协程
}
