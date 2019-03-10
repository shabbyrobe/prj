package prj

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/uuid"
)

func InitSimpleProject(ctx context.Context, session *Session, dest string, name string, at time.Time) (Project, *SimpleProjectConfig, error) {
	config, err := initSimpleProjectConfig(dest, name, at)
	if err != nil {
		return nil, nil, err
	}

	project, err := LoadSimpleProject(dest)
	if err != nil {
		return nil, nil, err
	}

	if _, err := project.Mark(ctx, session, "Initial", at, nil); err != nil {
		return nil, nil, err
	}

	return project, config, nil
}

func initSimpleProjectConfig(dest string, name string, at time.Time) (*SimpleProjectConfig, error) {
	if !filepath.IsAbs(dest) {
		return nil, fmt.Errorf("prj: input %q is not absolute", dest)
	}

	if exists, err := ContainsSimpleProject(dest); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("prj: project already exists at dest %q", dest)
	}

	config := &SimpleProjectConfig{
		ID:       createProjectID(),
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

func createProjectID() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Errorf("could not generate project ID"))
	}
	return u.String()
}
