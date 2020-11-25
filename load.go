package prj

import (
	"fmt"
	"os"
	"path/filepath"
)

var defaultPriority = []ProjectKind{ProjectSimple}

func Load(searchPath string, priority []ProjectKind) (Project, error) {
	path, kind, err := FindRoot(searchPath, priority)
	if err != nil {
		return nil, err
	}

	project, err := kind.Load(path)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func FindRoot(in string, priority []ProjectKind) (dir string, kind ProjectKind, err error) {
	if priority == nil {
		priority = defaultPriority
	}
	if !filepath.IsAbs(in) {
		return "", 0, fmt.Errorf("prj: input %q is not absolute", in)
	}

	fstat, err := os.Stat(in)
	if err != nil {
		return "", 0, err
	} else if !fstat.IsDir() {
		return "", 0, fmt.Errorf("prj: input %q is not a directory", in)
	}

	for _, kind := range priority {
		cur := in
		for {
			has, err := kind.Contains(cur)
			if err != nil {
				return "", 0, err
			} else if has {
				return cur, kind, nil
			}

			next := filepath.Dir(cur)
			if next == cur {
				break
			}
			cur = next
		}
	}

	return "", 0, &errProjectNotFound{in, priority}
}
