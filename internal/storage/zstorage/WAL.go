package zstorage

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"log/slog"
	"os"

	"github.com/NeverENG/BanKV/config"
	"github.com/NeverENG/BanKV/internal/storage/istorage"
)

const HEADER_LENGTH = 12

var _ istorage.IWal = &WAL{}

type WAL struct {
	file      *os.File
	headerBuf [HEADER_LENGTH]byte
}

func NewWAL() *WAL {
	file, err := os.OpenFile(config.G.WALPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("[ERROR]:OPEN WAL LOG ERROR !")
		return nil
	}
	return &WAL{file: file}
}

func (w *WAL) Write(entry istorage.LogEntry) error {

	hasher := crc32.NewIEEE()
	hasher.Write(entry.Key)
	hasher.Write(entry.Value)
	crc := hasher.Sum32()

	binary.BigEndian.PutUint32(w.headerBuf[:], crc)
	binary.BigEndian.PutUint32(w.headerBuf[4:], uint32(len(entry.Key)))
	binary.BigEndian.PutUint32(w.headerBuf[8:], uint32(len(entry.Value)))

	_, err := w.file.Write(w.headerBuf[:])
	if err != nil {
		slog.Error("[ERROR]:WRITE WAL LOG ERROR !")
		return err
	}

	_, err = w.file.Write(entry.Key)
	if err != nil {
		slog.Error("[ERROR]:WRITE WAL LOG ERROR !")
		return err
	}
	_, err = w.file.Write(entry.Value)
	if err != nil {
		slog.Error("[ERROR]:WRITE WAL LOG ERROR !")
		return err
	}

	return w.Sync()
}

func (w *WAL) Read(apply func(istorage.LogEntry) error) {
	if _, err := w.file.Seek(0, io.SeekStart); err != nil {
		return
	}
	for {
		header := make([]byte, HEADER_LENGTH)
		_, err := w.file.Read(header)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) {
				return
			}
			return
		}

		crc := binary.BigEndian.Uint32(header[:])
		keyLen := binary.BigEndian.Uint32(header[4:])
		valueLen := binary.BigEndian.Uint32(header[8:])

		key := make([]byte, keyLen)
		_, err = w.file.Read(key)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			return
		}
		value := make([]byte, valueLen)
		_, err = w.file.Read(value)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			return
		}

		haser := crc32.NewIEEE()
		haser.Write(key)
		haser.Write(value)
		if crc != haser.Sum32() {
			slog.Error("[ERROR]:THE DATA ERROR !")
			return
		}
		err = apply(istorage.LogEntry{key, value})
		if err != nil {
			return
		}
	}
}

func (w *WAL) Close() error {
	return w.file.Close()
}
func (w *WAL) Sync() error {
	return w.file.Sync()
}

// 采用日志滚动的模式来启动 Clear
func (w *WAL) Clear() error {

	if err := w.Close(); err != nil {
		slog.Error("[ERROR]:CLOSE WAL LOG ERROR !")
		return err
	}

	if err := os.Remove(w.file.Name()); err != nil {
		return err
	}

	f, err := os.OpenFile(w.file.Name(), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("[ERROR]:OPEN WAL LOG ERROR !")
		return err
	}
	w.file = f
	return w.Sync()
}
