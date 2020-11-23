package prj

import (
	"context"
	"path/filepath"
	"strings"
	"time"
)

type ResourcePath string

func NewResourcePath(p string) ResourcePath {
	return ResourcePath(strings.TrimLeft(filepath.FromSlash(p), string(filepath.Separator)))
}

func (rp ResourcePath) IsChildOf(parent ResourcePath) bool {
	pp := strings.TrimRight(string(parent), string(filepath.Separator))
	cp := strings.TrimRight(string(rp), string(filepath.Separator))
	return pp != cp && strings.HasPrefix(cp, pp)
}

var markOptionsDefault = &MarkOptions{}

type MarkOptions struct {
	Force  bool // Ignore 'no change' error
	Status *ProjectStatus
}

type Project interface {
	ID() string
	Name() string
	Path() string
	Kind() ProjectKind
	LastEntry() (*LogEntry, error)
	Status(ctx context.Context, path ResourcePath, at time.Time) (*ProjectStatus, error)
	Diff(ctx context.Context, path ResourcePath, at time.Time) (*ProjectDiff, error)
	Mark(ctx context.Context, session *Session, message string, at time.Time, options *MarkOptions) (*ProjectStatus, error)
	Log() LogIterator
	Tagger() Tagger
}

type Tagger interface {
	Tags() ([]string, error)
	Tag(with ...string) error
	Untag(with ...string) error
}
