package prj

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shabbyrobe/golib/errtools"
	"github.com/shabbyrobe/golib/pathtools"
)

type ResourcePath string

type Project interface {
	Status(ctx context.Context, path ResourcePath) (*ProjectStatus, error)
	Mark(ctx context.Context, session *Session, message string, status *ProjectStatus) (*ProjectStatus, error)
	Config() (*ProjectConfig, error)
	Log() ([]*LogEntry, error)
}

func loadConfigFromDir(dir string) (*ProjectConfig, error) {
	return loadConfigFile(filepath.Join(dir, ProjectPath, ProjectConfigFile))
}

func loadConfigFile(file string) (*ProjectConfig, error) {
	var p ProjectConfig
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bts, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

type SimpleProject struct {
	Root string
}

var _ Project = (*SimpleProject)(nil)

func (s *SimpleProject) logFile() string {
	return filepath.Join(s.Root, ProjectPath, ProjectLogFile)
}

func (s *SimpleProject) configFile() string {
	return filepath.Join(s.Root, ProjectPath, ProjectConfigFile)
}

func (s *SimpleProject) saveConfig(pc *ProjectConfig) error {
	tmpFile := s.configFile() + ".tmp"

	bts, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmpFile, bts, 0600); err != nil {
		return err
	}

	return os.Rename(tmpFile, s.configFile())
}

func (s *SimpleProject) Config() (*ProjectConfig, error) {
	return loadConfigFromDir(s.Root)
}

func (s *SimpleProject) Log() ([]*LogEntry, error) {
	logFile := s.logFile()
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	f, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []*LogEntry

	scn := bufio.NewScanner(f)
	for scn.Scan() {
		bts := scn.Bytes()
		if len(bts) == 0 {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal(bts, &entry); err != nil {
			return nil, err
		}

		entries = append(entries, &entry)
	}
	if err := scn.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *SimpleProject) Mark(ctx context.Context, session *Session, message string, status *ProjectStatus) (rstatus *ProjectStatus, rerr error) {
	if status == nil {
		var err error
		status, err = s.Status(ctx, "")
		if err != nil {
			return nil, err
		}
	}

	config, err := s.Config()
	if err != nil {
		return nil, err
	}

	logEntry := status.LogEntry(session, message)
	if config.LastEntry != nil {
		if ok, err := config.LastEntry.Hash.Equal(logEntry.Hash); err != nil {
			return status, err
		} else if ok {
			return status, fmt.Errorf("prj: project is unchanged since %q", config.LastEntry.ModTime)
		}
	}

	bts, err := json.Marshal(logEntry)
	if err != nil {
		return nil, err
	}
	bts = append(bts, '\n')

	// FIXME: flock
	f, err := os.OpenFile(s.logFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	defer errtools.DeferClose(&rerr, f)

	if _, err := f.Write(bts); err != nil {
		return nil, err
	}

	config.LastEntry = logEntry
	if err := s.saveConfig(config); err != nil {
		return nil, err
	}

	return status, nil
}

func (s *SimpleProject) Status(ctx context.Context, path ResourcePath) (*ProjectStatus, error) {
	var files []ProjectFile

	if err := filepath.Walk(filepath.Join(s.Root, string(path)), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
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
			return err
		}

		ok, _, left, err := pathtools.FilepathPrefix(path, s.Root)
		if err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("prj: path %q escaped root %q", path, s.Root)
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

	status := NewProjectStatus(files)

	return status, nil
}
