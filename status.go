package prj

import (
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

type ProjectFile struct {
	Name    ResourcePath
	Hash    Hash
	Size    int64
	ModTime time.Time
}

type ProjectStatus struct {
	Files   []ProjectFile
	Hash    Hash
	ModTime time.Time
	Size    int64
}

func (status *ProjectStatus) LogEntry(session *Session, message string) *LogEntry {
	le := &LogEntry{
		Author:     session.User,
		Machine:    session.Machine,
		Hash:       status.Hash,
		ModTime:    status.ModTime,
		Size:       status.Size,
		FileCount:  len(status.Files),
		Message:    message,
		StatusFile: statusFileName(status.ModTime, status.Hash),
	}
	return le
}

func statusFileName(modTime time.Time, hash Hash) string {
	return fmt.Sprintf("%s-%s.json",
		modTime.Format("20060102150405"),
		hex.EncodeToString([]byte(hash.String()))[:16])
}

func NewProjectStatus(files []ProjectFile, at time.Time) *ProjectStatus {
	ps := &ProjectStatus{
		Files: files,
	}

	sort.Slice(ps.Files, func(i, j int) bool {
		return ps.Files[i].Name < ps.Files[j].Name
	})

	for _, file := range ps.Files {
		if file.ModTime.After(ps.ModTime) {
			ps.ModTime = file.ModTime
		}
		ps.Size += file.Size
	}

	if ps.ModTime.IsZero() {
		ps.ModTime = at
	}

	const projectHashDelimiter = "/"

	hasher, _ := DefaultHashAlgorithm.CreateHasher()
	for _, file := range ps.Files {
		hasher.Write([]byte(file.Name + projectHashDelimiter))
		hasher.Write([]byte(file.Hash.Algorithm + projectHashDelimiter))
		hasher.Write(file.Hash.Value)
		hasher.Write([]byte(projectHashDelimiter))
	}

	ps.Hash = DefaultHashAlgorithm.Sum(hasher, nil)

	return ps
}
