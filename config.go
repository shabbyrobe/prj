package prj

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	ProjectPath       = ".prj"
	ProjectConfigFile = "config.json"
	ProjectLogFile    = "log.jsonl"
)

type ProjectConfig struct {
	ID       string
	Name     string
	InitDate time.Time

	LastEntry *LogEntry
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

func FindSimpleProjectRoot(in string) (path string, err error) {
	if !filepath.IsAbs(in) {
		return "", fmt.Errorf("prj: input %q is not absolute", in)
	}

	fstat, err := os.Stat(in)
	if err != nil {
		return "", err
	} else if !fstat.IsDir() {
		return "", fmt.Errorf("prj: input %q is not a directory", in)
	}

	cur := in
	for {
		has, err := ContainsSimpleProject(cur)
		if err != nil {
			return "", err
		} else if has {
			return cur, nil
		}

		next := filepath.Dir(cur)
		if next == cur {
			return "", fmt.Errorf("prj: simple project not found in %q or any of its parents", in)
		}
		cur = next
	}
}
