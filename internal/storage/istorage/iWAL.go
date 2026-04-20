package istorage

type LogEntry struct {
	Key   []byte
	Value []byte
}

type IWal interface {
	Write(entry LogEntry) error
	Read(func(entry LogEntry) error)

	Close() error
	Sync() error
	Clear() error
}
