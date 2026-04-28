package istorage

type IWal interface {
	Write(entry LogEntry) error
	Read(apply func(LogEntry) error)
	Close() error
	Sync() error
	Clear() error
}
