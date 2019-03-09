package prj

import (
	"os"
	"path/filepath"

	"github.com/shabbyrobe/go-prj/internal/fastwalk"
)

type FoundProject struct {
	Path   string
	Config *ProjectConfig
	Err    error
}

func Find(path string) (projects []FoundProject, err error) {
	if err := fastwalk.Walk(path, func(path string, mode os.FileMode) error {
		if !mode.IsDir() {
			return nil
		}

		if _, dir := filepath.Split(path); false ||
			dir == ".git" ||
			dir == ".hg" ||
			dir == ".svn" {
			return nil
		}

		if ok, err := containsSimpleProjectUnchecked(path); err != nil {
			return err
		} else if ok {
			config, err := loadConfigFromDir(path)
			projects = append(projects, FoundProject{
				Path:   path,
				Config: config,
				Err:    err,
			})
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return projects, nil
}
