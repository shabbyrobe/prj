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
	ErrSimpleProjectNotFound ErrorCode = iota + 1
)

type errSimpleProjectNotFound struct {
	Path string
}

func (err *errSimpleProjectNotFound) Is(target error) bool {
	return target == ErrSimpleProjectNotFound
}

func (err *errSimpleProjectNotFound) Error() string {
	return fmt.Sprintf("prj: simple project not found in %q or any of its parents", err.Path)
}
