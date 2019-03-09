package prj

import (
	"time"
)

type LogEntry struct {
	Author     string
	Message    string
	Machine    string
	Hash       Hash
	Size       int64
	FileCount  int
	ModTime    time.Time
	StatusFile string
}
