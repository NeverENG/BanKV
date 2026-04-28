package istorage

type ISSTable interface {
	LoadFromDisk() []*SSTableMata
	AddMata()
	RemoveMata()
	GetLevelFiles()
	GetAllMata() []*SSTableMata

	WriteToSSTable(entry []LogEntry) (SSTableMata, error)
	ReadFromSSTable(key []byte) []byte
	ReadAllEntries()
	MergeSSTables()
	DeleteSSTable()
}
