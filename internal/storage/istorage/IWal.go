package istorage

type IWal interface {
	Write(entry LogEntry) error
	Read() ([]LogEntry, error)
	Close() error
	Sync() error
	Clear() error
}
