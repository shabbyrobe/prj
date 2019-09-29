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
	// This stuff is too slow:
	/*
		repo, err := git.PlainOpen(path)
		if err != nil {
			return nil, err
		}
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		id := ref.Hash().String()
	*/
	id := ""
	return &GitProject{
		path: path,
		id:   id,
	}, nil
}

func (g *GitProject) ID() string { return g.id }

func (g *GitProject) Name() string {
	_, name := filepath.Split(strings.TrimRight(g.path, string(filepath.Separator)))
	return name
}

func (g *GitProject) Path() string      { return g.path }
func (g *GitProject) Kind() ProjectKind { return ProjectGit }

func (g *GitProject) LastEntry() (*LogEntry, error) {
	return nil, nil
	// This stuff is too slow:
	/*
		repo, err := git.PlainOpen(g.path)
		if err != nil {
			return nil, err
		}
		ref, err := repo.Head()
		if err != nil {
			return nil, err
		}
		obj, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return nil, err
		}
		iter, err := obj.Files()
		if err != nil {
			return nil, err
		}

		n := 0
		for {
			_, err := iter.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			n++
		}
		iter.Close()

		hash := ref.Hash()
		var lastEntry = LogEntry{
			Author:    obj.Author.String(),
			Message:   obj.Message,
			Hash:      Hash{Value: hash[:]},
			FileCount: n,
			Time:      obj.Author.When,
		}
		return &lastEntry, nil
	*/
}

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
