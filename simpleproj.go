package prj

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/shabbyrobe/golib/errtools"
	"github.com/shabbyrobe/golib/pathtools"
)

type SimpleProject struct {
	// Real project data files are stored here (user's files)
	dataRoot string

	// Project metadata is stored here (should be the same as dataRoot, except
	// in rare cases)
	metaRoot string

	config *SimpleProjectConfig
}

var _ Project = (*SimpleProject)(nil)

func LoadSimpleProject(projectPath string) (*SimpleProject, error) {
	sp := &SimpleProject{dataRoot: projectPath, metaRoot: projectPath}
	if err := sp.refreshConfig(); err != nil {
		return nil, err
	}
	return sp, nil
}

// loadSimpleProjectWithSeparateMeta is a bit of a hack to allow us to create
// virtual projects on an ad-hoc basis, i.e. when trying to hash or diff to
// an existing directory that does not contain a project.
func loadSimpleProjectWithSeparateMeta(dataPath string, metaPath string) (*SimpleProject, error) {
	sp := &SimpleProject{dataRoot: dataPath, metaRoot: metaPath}
	if err := sp.refreshConfig(); err != nil {
		return nil, err
	}
	return sp, nil
}

func (s *SimpleProject) ID() string        { return s.config.ID }
func (s *SimpleProject) Name() string      { return s.config.Name }
func (s *SimpleProject) Path() string      { return s.dataRoot }
func (s *SimpleProject) Kind() ProjectKind { return ProjectSimple }

func (s *SimpleProject) LastEntry() (*LogEntry, error) {
	return s.config.LastEntry, nil
}

func (s *SimpleProject) logFile() string {
	return filepath.Join(s.metaRoot, ProjectPath, ProjectLogFile)
}

func (s *SimpleProject) configFile() string {
	return filepath.Join(s.metaRoot, ProjectPath, ProjectConfigFile)
}

func (s *SimpleProject) statusPath() string {
	return filepath.Join(s.metaRoot, ProjectPath, projectStatusPath)
}

func (s *SimpleProject) ensureStatusPath() (string, error) {
	statusPath := s.statusPath()

	if _, err := os.Stat(statusPath); os.IsNotExist(err) {
		if err := os.Mkdir(statusPath, 0700); err != nil {
			return statusPath, err
		}
		return statusPath, nil

	} else if err != nil {
		return statusPath, err
	}

	return statusPath, nil
}

func (s *SimpleProject) saveConfig() error {
	tmpFile := s.configFile() + ".tmp"

	bts, err := json.MarshalIndent(s.config, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmpFile, bts, 0600); err != nil {
		return err
	}

	return os.Rename(tmpFile, s.configFile())
}

func (s *SimpleProject) refreshConfig() (err error) {
	s.config, err = loadConfigFromDir(s.metaRoot)
	return err
}

func (s *SimpleProject) Log() LogIterator {
	logFile := s.logFile()
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return &nilLogIterator{}
	} else if err != nil {
		return &errLogIterator{err}
	}

	f, err := os.Open(logFile)
	if err != nil {
		return &errLogIterator{err}
	}

	return &jsonlLogIterator{scn: bufio.NewScanner(f), cls: f}
}

func (s *SimpleProject) Mark(ctx context.Context, session *Session, message string, at time.Time, options *MarkOptions) (rstatus *ProjectStatus, rerr error) {
	if options == nil {
		options = markOptionsDefault
	}

	// FIXME: flock
	if err := s.refreshConfig(); err != nil {
		return nil, err
	}

	var status = options.Status

	if status == nil {
		var err error
		status, err = s.Status(ctx, "", at)
		if err != nil {
			return nil, err
		}
	}

	logEntry := status.LogEntry(session, message, at)

	if !options.Force && s.config.LastEntry != nil {
		if ok, err := s.config.LastEntry.Hash.Equal(logEntry.Hash); err != nil {
			return status, err
		} else if ok {
			return status, fmt.Errorf("prj: project is unchanged since %q", s.config.LastEntry.ModTime)
		}
	}

	// Prepare serialised data before writing anything:
	statusPath, err := s.ensureStatusPath()
	if err != nil {
		return nil, err
	}

	statusData, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return nil, err
	}

	logEntryData, err := json.Marshal(logEntry)
	if err != nil {
		return nil, err
	}
	logEntryData = append(logEntryData, '\n')

	{ // Append to log
		f, err := os.OpenFile(s.logFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		defer errtools.DeferClose(&rerr, f)

		if _, err := f.Write(logEntryData); err != nil {
			return nil, err
		}
	}

	{ // Update config with last entry
		s.config.LastEntry = logEntry
		if err := s.saveConfig(); err != nil {
			return nil, err
		}
	}

	{ // Write status file
		if err := ioutil.WriteFile(filepath.Join(statusPath, logEntry.StatusFile), statusData, 0600); err != nil {
			return nil, err
		}

	}

	return status, nil
}

func (s *SimpleProject) Status(ctx context.Context, childPath ResourcePath, at time.Time) (*ProjectStatus, error) {
	var files []ProjectFile

	if err := filepath.Walk(filepath.Join(s.dataRoot, string(childPath)), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if info.Mode()&os.ModeSocket != 0 {
			return nil
		}

		if info.IsDir() {
			if _, dir := filepath.Split(path); dir == ProjectPath {
				return filepath.SkipDir
			}

			// FIXME: what if the dir contains a sub-project?
			return nil
		}

		hash, err := DefaultHashAlgorithm.HashFile(path)
		if err != nil {
			return fmt.Errorf("prj: hash file %q failed: %w", path, err)
		}

		ok, _, left, err := pathtools.FilepathPrefix(path, s.dataRoot)
		if err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("prj: path %q escaped root %q", path, s.dataRoot)
		}

		files = append(files, ProjectFile{
			Name:    ResourcePath(left),
			Hash:    hash,
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})

		return nil

	}); err != nil {
		return nil, err
	}

	status := NewProjectStatus(files, at)

	return status, nil
}

func (s *SimpleProject) Diff(ctx context.Context, path ResourcePath, at time.Time) (*ProjectDiff, error) {
	if err := s.refreshConfig(); err != nil {
		return nil, err
	}

	currentStatus, err := s.Status(ctx, path, at)
	if err != nil {
		return nil, err
	}

	var lastStatus ProjectStatus
	if s.config.LastEntry != nil {
		if s.config.LastEntry.StatusFile == "" {
			return nil, fmt.Errorf("prj: no status file for last log entry, cannot diff")
		}

		bts, err := ioutil.ReadFile(filepath.Join(s.statusPath(), s.config.LastEntry.StatusFile))
		if err != nil {
			return nil, fmt.Errorf("prj: could not read status file for last log entry, cannot diff; previous error: %v", err)
		}

		if err := json.Unmarshal(bts, &lastStatus); err != nil {
			return nil, fmt.Errorf("prj: could not unmarshal status file for last log entry, cannot diff; previous error: %v", err)
		}

		if path != "" {
			lastStatus = *lastStatus.Filter(path, at)
		}
	}

	return currentStatus.CompareTo(&lastStatus)
}

func (s *SimpleProject) Tagger() Tagger {
	return fileTaggerFromDir(s.dataRoot)
}

func ContainsSimpleProject(dir string) (ok bool, err error) {
	if !filepath.IsAbs(dir) {
		return false, fmt.Errorf("prj: input %q is not absolute", dir)
	}
	return containsSimpleProjectUnchecked(dir)
}

func containsSimpleProjectUnchecked(dir string) (ok bool, err error) {
	_, err = os.Stat(filepath.Join(dir, ProjectPath, ProjectConfigFile))
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type jsonlLogIterator struct {
	scn *bufio.Scanner
	cls io.Closer
	err error
}

func (jl *jsonlLogIterator) Next(entry *LogEntry) bool {
	if jl.err != nil {
		return false
	}

	for jl.scn.Scan() {
		bts := jl.scn.Bytes()
		if len(bts) == 0 {
			continue
		}

		if err := json.Unmarshal(bts, entry); err != nil {
			jl.err = err
			return false
		}
		return true
	}

	return false
}

func (jl *jsonlLogIterator) Close() error {
	if cerr := jl.cls.Close(); jl.err == nil && cerr != nil {
		jl.err = cerr
	}
	if serr := jl.scn.Err(); jl.err == nil && serr != nil {
		jl.err = serr
	}
	return jl.err
}
