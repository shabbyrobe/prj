package prj

import (
	"encoding/json"
	"time"
)

type LogEntry struct {
	Author  string
	Message string
	Machine string
	Hash    Hash
	Size    int64

	StatusFile string

	// Number of files in the repo at this point in time, or -1 if the backend
	// can't report this information (Git/Hg may make this too expensive).
	FilesCount int

	// Number of files changed in the repo at this point in time, or -1 if the backend
	// can't report this information.
	FilesChanged int

	// Latest modification date of all files in the tree
	ModTime time.Time

	// Time of the log entry
	Time time.Time
}

func (le *LogEntry) UnmarshalJSON(bts []byte) error {
	// Strip away UnmarshalJSON method:
	type inner LogEntry

	var entry inner
	if err := json.Unmarshal(bts, &entry); err != nil {
		return err
	}

	// This stuff is necessary because in an earlier implementation I seemed to
	// think that the latest last mod time was an adequate standin for a log
	// entry date, totally oblivious to the fact that you might _remove_ that
	// file in the most recent version! Oops! Too many of these things in the
	// wild.
	if entry.Time.IsZero() {
		entry.Time = entry.ModTime
	}
	*le = LogEntry(entry)
	return nil
}

type LogIterator interface {
	Next(entry *LogEntry) bool
	Close() error
}

type errLogIterator struct {
	err error
}

func (e *errLogIterator) Next(*LogEntry) bool { return false }
func (e *errLogIterator) Close() error        { return e.err }

type nilLogIterator struct{}

func (*nilLogIterator) Next(*LogEntry) bool { return false }
func (*nilLogIterator) Close() error        { return nil }
