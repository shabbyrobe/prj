package prj

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func InitSimpleProject(dest string, name string, at time.Time) (*ProjectConfig, error) {
	if !filepath.IsAbs(dest) {
		return nil, fmt.Errorf("prj: input %q is not absolute", dest)
	}

	if exists, err := ContainsSimpleProject(dest); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("prj: project already exists at dest %q", dest)
	}

	config := &ProjectConfig{
		Name:     name,
		InitDate: at,
	}

	projectPath := filepath.Join(dest, ProjectPath)
	if err := os.Mkdir(projectPath, 0700); err != nil {
		return nil, err
	}

	bts, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(filepath.Join(projectPath, ProjectConfigFile), bts, 0600); err != nil {
		return nil, err
	}

	return config, nil
}
