package main

import (
	"os"

	prj "github.com/shabbyrobe/prj"
)

func loadProject(in string) (prj.Project, *prj.Session, error) {
	if in == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, nil, err
		}
		in = wd
	}

	path, err := prj.FindSimpleProjectRoot(in)
	if err != nil {
		return nil, nil, err
	}

	project, err := prj.LoadSimpleProject(path)
	if err != nil {
		return nil, nil, err
	}

	session, err := prj.NewOSSession()
	if err != nil {
		return nil, nil, err
	}

	return project, session, nil
}
