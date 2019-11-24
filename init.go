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

type initOptions struct {
	metaPath string
}

type InitOption func(opts *initOptions)

func InitWithSeparateMetaPath(metaPath string) InitOption {
	return func(opts *initOptions) { opts.metaPath = metaPath }
}

func InitSimpleProject(ctx context.Context, session *Session, projectPath string, name string, at time.Time, options ...InitOption) (Project, *SimpleProjectConfig, error) {
	var opts = initOptions{
		metaPath: projectPath,
	}
	for _, o := range options {
		o(&opts)
	}

	config, err := initSimpleProjectConfig(opts.metaPath, name, at)
	if err != nil {
		return nil, nil, err
	}

	project, err := loadSimpleProjectWithSeparateMeta(projectPath, opts.metaPath)
	if err != nil {
		return nil, nil, err
	}

	if _, err := project.Mark(ctx, session, "Initial", at, nil); err != nil {
		return nil, nil, err
	}

	return project, config, nil
}

func initSimpleProjectConfig(metaPath string, name string, at time.Time) (*SimpleProjectConfig, error) {
	if !filepath.IsAbs(metaPath) {
		return nil, fmt.Errorf("prj: input %q is not absolute", metaPath)
	}

	if exists, err := ContainsSimpleProject(metaPath); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("prj: project already exists at dest %q", metaPath)
	}

	config := &SimpleProjectConfig{
		ID:       createProjectID(),
		Name:     name,
		InitDate: at,
	}

	projectPath := filepath.Join(metaPath, ProjectPath)
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
