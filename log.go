package prj

import (
	"encoding/json"
	"time"
)

type LogEntry struct {
	Author     string
	Message    string
	Machine    string
	Hash       Hash
	Size       int64
	FileCount  int
	StatusFile string

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
