package istorage

import (
	"encoding/binary"
	"os"
	"sync"
)

type LogEntry struct {
	Key   []byte
	Value []byte
}

type SSTableMata struct {
	Level    int
	Filepath string
	MinKey   []byte
	MaxKey   []byte
	Size     int64

	mu           sync.Once
	MaxKeyLoaded bool
}

func (meta *SSTableMata) EnsureMeta() {
	meta.mu.Do(func() { // sync.Once 保证只执行一次
		if meta.MaxKeyLoaded {
			return
		}

		file, _ := os.Open(meta.Filepath)
		defer file.Close()

		var maxKey []byte

		// 遍历整个文件，找到最后一个 key
		for {
			var keyLen uint32
			if err := binary.Read(file, binary.BigEndian, &keyLen); err != nil {
				break // EOF
			}
			keyBytes := make([]byte, keyLen)
			file.Read(keyBytes)

			var valueLen uint32
			binary.Read(file, binary.BigEndian, &valueLen)
			file.Seek(int64(valueLen), 1) // 跳过 value

			maxKey = keyBytes // 不断更新，最后一条就是 MaxKey
		}

		meta.MaxKey = maxKey
		meta.MaxKeyLoaded = true
	})
}
