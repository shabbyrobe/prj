package prj

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/shabbyrobe/prj/internal/fastwalk"
)

type FoundProject struct {
	Path   string
	Config *ProjectConfig
	Err    error
}

func Scan(path string) *Scanner {
	var result = make(chan *FoundProject, 2000)
	var errc = make(chan error, 1)
	var stop = make(chan struct{})
	var errStop = errors.New("stop")

	go func() {
		defer close(result)
		defer close(errc)

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
				found := &FoundProject{Path: path, Config: config, Err: err}

				select {
				case result <- found:
				case <-stop:
					return errStop
				}
			}
			return nil

		}); err != nil && err != errStop {
			errc <- err
		}
	}()

	return &Scanner{result: result, errc: errc, stop: stop}
}

type Scanner struct {
	result chan *FoundProject
	errc   chan error
	stop   chan struct{}

	done bool
	err  error
	cur  *FoundProject
}

func (scn *Scanner) Next() bool {
	if scn.err != nil || scn.done {
		return false
	}

	select {
	case err := <-scn.errc:
		scn.done, scn.err = true, err
		return false

	case found := <-scn.result:
		scn.done, scn.cur = found == nil, found
		return !scn.done
	}
}

func (scn *Scanner) Project() *FoundProject { return scn.cur }

func (scn *Scanner) Close() error {
	close(scn.stop)
	select {
	case err := <-scn.errc:
		scn.done, scn.err = true, err
	}
	return scn.err
}
