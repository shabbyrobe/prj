package prj

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	ProjectPath       = ".prj"
	ProjectConfigFile = "config.json"
	ProjectLogFile    = "log.jsonl"
	projectStatusPath = "status" // Child of ProjectPath
)

type SimpleProjectConfig struct {
	ID       string
	Name     string
	InitDate time.Time

	LastEntry *LogEntry
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
			return "", &errProjectNotFound{Path: in, Kinds: []ProjectKind{ProjectSimple}}
		}
		cur = next
	}
}

func loadConfigFromDir(dir string) (*SimpleProjectConfig, error) {
	return loadConfigFile(filepath.Join(dir, ProjectPath, ProjectConfigFile))
}

func loadConfigFile(file string) (*SimpleProjectConfig, error) {
	var p SimpleProjectConfig
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bts, &p); err != nil {
		return nil, err
	}

	return &p, nil
}
