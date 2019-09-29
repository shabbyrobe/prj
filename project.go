package prj

import (
	"context"
	"path/filepath"
	"strings"
	"time"
)

type ResourcePath string

func NewResourcePath(p string) ResourcePath {
	return ResourcePath(strings.TrimLeft(p, string(filepath.Separator)))
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
	Log() ([]*LogEntry, error)
}
