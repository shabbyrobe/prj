package prj

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HgProject struct {
	path string
	id   string
}

func ContainsHgProject(dir string) (ok bool, err error) {
	if !filepath.IsAbs(dir) {
		return false, fmt.Errorf("prj: input %q is not absolute", dir)
	}
	return containsHgProjectUnchecked(dir)
}

func containsHgProjectUnchecked(dir string) (ok bool, err error) {
	if _, err = os.Stat(filepath.Join(dir, ".hg")); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func LoadHgProject(path string) (*HgProject, error) {
	return &HgProject{
		path: path,
		// FIXME: this can be pulled out of the hg repo
		id: "",
	}, nil
}

func (g *HgProject) ID() string {
	return g.id
}

func (g *HgProject) Name() string {
	_, name := filepath.Split(strings.TrimRight(g.path, string(filepath.Separator)))
	return name
}

func (g *HgProject) Path() string      { return g.path }
func (g *HgProject) Kind() ProjectKind { return ProjectHg }

func (g *HgProject) LastEntry() (*LogEntry, error) { return nil, nil }

func (g *HgProject) Status(ctx context.Context, path ResourcePath, at time.Time) (*ProjectStatus, error) {
	return nil, fmt.Errorf("prj: not implemented")
}

func (g *HgProject) Diff(ctx context.Context, path ResourcePath, at time.Time) (*ProjectDiff, error) {
	return nil, fmt.Errorf("prj: not implemented")
}

func (g *HgProject) Mark(ctx context.Context, session *Session, message string, at time.Time, options *MarkOptions) (*ProjectStatus, error) {
	return nil, fmt.Errorf("prj: not implemented")
}

func (g *HgProject) Log() ([]*LogEntry, error) {
	return nil, nil
}
