package main

import (
	"os"

	prj "github.com/shabbyrobe/prj"
)

func loadProject() (prj.Project, *prj.Session, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	path, err := prj.FindSimpleProjectRoot(wd)
	if err != nil {
		return nil, nil, err
	}

	project := &prj.SimpleProject{Root: path}

	session, err := prj.NewOSSession()
	if err != nil {
		return nil, nil, err
	}

	return project, session, nil
}
