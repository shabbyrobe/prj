package prj

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pmezard/go-difflib/difflib"
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

func (status *ProjectStatus) Filter(childPath ResourcePath, at time.Time) *ProjectStatus {
	childPathStr := string(childPath)

	var files = make([]ProjectFile, 0)
	for _, file := range status.Files {
		if strings.HasPrefix(string(file.Name), childPathStr) {
			files = append(files, file)
		}
	}
	return NewProjectStatus(files, at)
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

func (status *ProjectStatus) CompareTo(previous *ProjectStatus) (*ProjectDiff, error) {
	if previous == nil {
		previous = &ProjectStatus{}
	}

	currentFiles := make([]string, len(status.Files))
	currentIndex := make(map[ResourcePath]*ProjectFile, len(status.Files))
	for i, f := range status.Files {
		currentFiles[i] = string(f.Name)
		currentIndex[f.Name] = &f
	}

	prevFiles := make([]string, len(previous.Files))
	prevIndex := make(map[ResourcePath]*ProjectFile, len(previous.Files))
	for i, f := range previous.Files {
		prevFiles[i] = string(f.Name)
		prevIndex[f.Name] = &f
	}

	sort.Strings(currentFiles)
	sort.Strings(prevFiles)

	df := difflib.NewMatcher(currentFiles, prevFiles)
	matches := df.GetMatchingBlocks()

	var diff ProjectDiff
	diff.Current = status
	diff.Previous = previous

	lastCurrent, lastPrev := 0, 0
	for _, m := range matches {
		currentMatchIndex := m.A
		prevMatchIndex := m.B

		// differences prior to match
		for i := lastCurrent; i < currentMatchIndex; i++ {
			diff.Added = append(diff.Added, ResourcePath(currentFiles[i]))
		}
		for i := lastPrev; i < prevMatchIndex; i++ {
			diff.Removed = append(diff.Removed, ResourcePath(prevFiles[i]))
		}

		// differences inside match
		for i := 0; i < m.Size; i++ {
			res := ResourcePath(currentFiles[currentMatchIndex+i])
			currentFile, prevFile := currentIndex[res], prevIndex[res]
			if eq, err := currentFile.Hash.Equal(prevFile.Hash); err != nil {
				return nil, err
			} else if eq {
				diff.Same = append(diff.Same, res)
			} else {
				diff.Modified = append(diff.Modified, res)
			}
		}

		lastCurrent, lastPrev = currentMatchIndex+m.Size, prevMatchIndex+m.Size
	}

	// differences following match
	for i, j := lastCurrent, len(currentFiles); i < j; i++ {
		diff.Added = append(diff.Added, ResourcePath(currentFiles[i]))
	}
	for i, j := lastPrev, len(prevFiles); i < j; i++ {
		diff.Removed = append(diff.Removed, ResourcePath(prevFiles[i]))
	}

	return &diff, nil
}

type ProjectDiff struct {
	Current  *ProjectStatus
	Previous *ProjectStatus

	Added    []ResourcePath
	Removed  []ResourcePath
	Modified []ResourcePath
	Same     []ResourcePath
}

type ProjectDiffItem struct {
	Path   ResourcePath
	Status DiffStatus
}

func (diff *ProjectDiff) Items() []ProjectDiffItem {
	items := make([]ProjectDiffItem, len(diff.Added)+len(diff.Removed)+len(diff.Modified)+len(diff.Same))

	n := 0
	for _, p := range diff.Added {
		items[n].Path = p
		items[n].Status = DiffAdded
		n++
	}

	for _, p := range diff.Removed {
		items[n].Path = p
		items[n].Status = DiffRemoved
		n++
	}

	for _, p := range diff.Modified {
		items[n].Path = p
		items[n].Status = DiffModified
		n++
	}

	for _, p := range diff.Same {
		items[n].Path = p
		items[n].Status = DiffSame
		n++
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Path < items[j].Path
	})

	return items[:n]
}

type DiffStatus byte

const (
	DiffAdded    DiffStatus = 'A'
	DiffRemoved  DiffStatus = 'D'
	DiffModified DiffStatus = 'M'
	DiffSame     DiffStatus = '='
)

func statusFileName(modTime time.Time, hash Hash) string {
	return fmt.Sprintf("%s-%s.json",
		modTime.Format("20060102150405"),
		hex.EncodeToString([]byte(hash.String()))[:16])
}
