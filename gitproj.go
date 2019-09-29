package prj

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GitProject struct {
	path string
	id   string
}

func ContainsGitProject(dir string) (ok bool, err error) {
	if !filepath.IsAbs(dir) {
		return false, fmt.Errorf("prj: input %q is not absolute", dir)
	}
	return containsGitProjectUnchecked(dir)
}

func containsGitProjectUnchecked(dir string) (ok bool, err error) {
	if _, err = os.Stat(filepath.Join(dir, ".git")); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func LoadGitProject(path string) (*GitProject, error) {
	return &GitProject{
		path: path,
		// FIXME: this can be pulled out of the git repo
		id: "",
	}, nil
}

func (g *GitProject) ID() string { return g.id }

func (g *GitProject) Name() string {
	_, name := filepath.Split(strings.TrimRight(g.path, string(filepath.Separator)))
	return name
}

func (g *GitProject) Path() string      { return g.path }
func (g *GitProject) Kind() ProjectKind { return ProjectGit }

func (g *GitProject) LastEntry() *LogEntry { return nil }

func (g *GitProject) Status(ctx context.Context, path ResourcePath, at time.Time) (*ProjectStatus, error) {
	return nil, fmt.Errorf("prj: not implemented")
}

func (g *GitProject) Diff(ctx context.Context, path ResourcePath, at time.Time) (*ProjectDiff, error) {
	return nil, fmt.Errorf("prj: not implemented")
}

func (g *GitProject) Mark(ctx context.Context, session *Session, message string, at time.Time, options *MarkOptions) (*ProjectStatus, error) {
	return nil, fmt.Errorf("prj: not implemented")
}

func (g *GitProject) Log() ([]*LogEntry, error) {
	return nil, nil
}
