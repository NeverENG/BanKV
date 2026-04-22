package istorage

type LogEntry struct {
	Key   []byte
	Value []byte
}

type IWal interface {
	Write(entry LogEntry) error
	Read(apply func(LogEntry) error)
	Close() error
	Sync() error
	Clear() error
}

type IMemTable interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Size() int
	StartFlush()
}
