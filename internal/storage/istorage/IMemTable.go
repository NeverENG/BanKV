package istorage

type IMemTable interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Size() int
	StartFlush()
}
