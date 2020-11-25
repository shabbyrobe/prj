package prj

import "fmt"

// ErrorCode is used to provide Is() checks to external consumers, for example:
//
//	if errors.Is(err, prj.ErrSimpleProjectNotFound) {
// 	    doStuff()
// 	}
//
// These errors will never be returned directly.
type ErrorCode int

// Don't use this directly, it's just here for errors.Is checks()
func (ec ErrorCode) Error() string {
	return fmt.Sprintf("error:%d", ec)
}

const (
	ErrProjectNotFound ErrorCode = iota + 1
)

type errProjectNotFound struct {
	Path  string
	Kinds []ProjectKind
}

func (err *errProjectNotFound) Is(target error) bool {
	return target == ErrProjectNotFound
}

func (err *errProjectNotFound) Error() string {
	return fmt.Sprintf("prj: project not found in %q or any of its parents", err.Path)
}
