package main

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"time"

	prj "github.com/shabbyrobe/prj"
)

func loadSimpleProject(searchPath string) (prj.Project, *prj.Session, error) {
	if searchPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, nil, err
		}
		searchPath = wd
	}

	path, err := prj.FindSimpleProjectRoot(searchPath)
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

func loadSimpleProjectWithTemporaryFallback(ctx context.Context, searchPath string, fallbackPath string) (p prj.Project, sess *prj.Session, done func(), err error) {
	done = func() {}
	defer func() {
		if err != nil {
			done()
			done = func() {}
		}
	}()

	p, sess, err = loadSimpleProject(searchPath)
	if errors.Is(err, prj.ErrSimpleProjectNotFound) {
		if fallbackPath == "" {
			return p, sess, done, err
		}

		p, sess, done, err = loadTemporaryProject(ctx, fallbackPath)
	}

	return p, sess, done, err
}

func loadTemporaryProject(ctx context.Context, path string) (p prj.Project, sess *prj.Session, done func(), err error) {
	done = func() {}

	sess, err = prj.NewOSSession()
	if err != nil {
		return p, sess, done, err
	}

	metaPath, err := ioutil.TempDir("", "")
	if err != nil {
		return p, sess, done, err
	}
	done = func() {
		os.RemoveAll(metaPath)
	}

	defer func() {
		if err != nil {
			done()
			done = func() {}
		}
	}()

	p, _, err = prj.InitSimpleProject(
		ctx, sess, path, path, time.Now(),
		prj.InitWithSeparateMetaPath(metaPath))

	return p, sess, done, err
}
