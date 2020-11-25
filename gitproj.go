package prj

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/shabbyrobe/golib/errtools"
)

type GitProject struct {
	path string
	id   string
}

var _ Project = &GitProject{}

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
	h, err := gitReadID(path)
	if err != nil {
		return nil, err
	}

	return &GitProject{
		path: path,
		id:   hex.EncodeToString(h[:]),
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

func (g *GitProject) Log() LogIterator {
	panic("not implemented")
}

func (g *GitProject) Tagger() Tagger {
	return fileTaggerFromDir(g.path)
}

// Read the ID of a git repository (probably first revision's hash)
//
// Custom unrolling of git utilities from go-git is WAY faster than
// interfacing with go-git directly.
func gitReadID(path string) (id [20]byte, err error) {
	dot, err := openDotGit(path)
	if err != nil {
		return id, err
	}

	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	ref, err := storer.ResolveReference(s, plumbing.HEAD)
	if err == plumbing.ErrReferenceNotFound {
		// This means the repo has been initialised but does not contain a commit.
		// We might be better off catching this error higher up, and allowing the
		// git project to be returned with an empty, invalid ID, then making sure
		// we account for "invalid IDs" elsewhere where we might wish to dedupe, etc.
		return id, &errProjectNotFound{}
	} else if err != nil {
		return id, err
	}

	return ref.Hash(), nil
}

func openDotGit(path string) (billy.Filesystem, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("path must be absolute")
	}

	var fi os.FileInfo

	fs := osfs.New(path)
	fi, err := fs.Stat(".git")
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		dot, err := fs.Chroot(".git")
		return dot, err
	}

	return dotGitFileToOSFilesystem(path, fs)
}

func dotGitFileToOSFilesystem(path string, fs billy.Filesystem) (bfs billy.Filesystem, err error) {
	f, err := fs.Open(".git")
	if err != nil {
		return nil, err
	}
	defer errtools.DeferClose(&err, f)

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	line := string(b)
	const prefix = "gitdir: "
	if !strings.HasPrefix(line, prefix) {
		return nil, fmt.Errorf(".git file has no %s prefix", prefix)
	}

	gitdir := strings.Split(line[len(prefix):], "\n")[0]
	gitdir = strings.TrimSpace(gitdir)
	if filepath.IsAbs(gitdir) {
		return osfs.New(gitdir), nil
	}

	return osfs.New(fs.Join(path, gitdir)), nil
}
