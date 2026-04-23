package Raft

import (
	"encoding/binary"
	"os"
	"path/filepath"
)

const (
	MagicNumber = 0x52415654 // "RAFT" in ASCII
	Version     = 1
	StateFile   = "raft_state.dat"
	LogFile     = "raft_log.dat"
)

type RaftWAL struct {
	file     *os.File
	logPath  string
	metaPath string
}

type WALState struct {
	Term     int
	VotedFor int
}

func NewRaftWAL(dir string) (*RaftWAL, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	wal := &RaftWAL{
		logPath:  filepath.Join(dir, LogFile),
		metaPath: filepath.Join(dir, StateFile),
	}

	f, err := os.OpenFile(wal.logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	wal.file = f

	return wal, nil
}

func (w *RaftWAL) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func (w *RaftWAL) SaveState(term int, votedFor int) error {
	state := WALState{Term: term, VotedFor: votedFor}

	f, err := os.Create(w.metaPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := binary.Write(f, binary.BigEndian, MagicNumber); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, Version); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, int64(state.Term)); err != nil {
		return err
	}
	if err := binary.Write(f, binary.BigEndian, int64(state.VotedFor)); err != nil {
		return err
	}

	return f.Sync()
}

func (w *RaftWAL) LoadState() (int, int, error) {
	f, err := os.Open(w.metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, -1, nil
		}
		return 0, -1, err
	}
	defer f.Close()

	var magic uint32
	if err := binary.Read(f, binary.BigEndian, &magic); err != nil {
		return 0, -1, err
	}
	if magic != MagicNumber {
		return 0, -1, nil
	}

	var version uint32
	if err := binary.Read(f, binary.BigEndian, &version); err != nil {
		return 0, -1, err
	}

	var term int64
	if err := binary.Read(f, binary.BigEndian, &term); err != nil {
		return 0, -1, err
	}

	var votedFor int64
	if err := binary.Read(f, binary.BigEndian, &votedFor); err != nil {
		return 0, -1, err
	}

	return int(term), int(votedFor), nil
}

func (w *RaftWAL) AppendLog(entry LogEntry) error {
	if err := binary.Write(w.file, binary.BigEndian, int64(entry.Index)); err != nil {
		return err
	}
	if err := binary.Write(w.file, binary.BigEndian, int64(entry.Term)); err != nil {
		return err
	}
	if err := binary.Write(w.file, binary.BigEndian, int64(len(entry.Command))); err != nil {
		return err
	}
	if _, err := w.file.Write(entry.Command); err != nil {
		return err
	}

	return w.file.Sync()
}

func (w *RaftWAL) LoadLogs() ([]LogEntry, error) {
	f, err := os.Open(w.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var logs []LogEntry
	for {
		var index int64
		if err := binary.Read(f, binary.BigEndian, &index); err != nil {
			break
		}
		var term int64
		if err := binary.Read(f, binary.BigEndian, &term); err != nil {
			break
		}
		var cmdLen int64
		if err := binary.Read(f, binary.BigEndian, &cmdLen); err != nil {
			break
		}
		cmd := make([]byte, cmdLen)
		if _, err := f.Read(cmd); err != nil {
			break
		}

		logs = append(logs, LogEntry{
			Index:   int(index),
			Term:    int(term),
			Command: cmd,
		})
	}

	return logs, nil
}

func (w *RaftWAL) Clear() error {
	w.Close()
	if err := os.Remove(w.logPath); err != nil {
		return err
	}
	if err := os.Remove(w.metaPath); err != nil {
		return err
	}

	f, err := os.Create(w.logPath)
	if err != nil {
		return err
	}
	w.file = f

	return nil
}
