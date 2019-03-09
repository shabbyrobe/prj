package prj

import (
	"errors"
	"path/filepath"

	"github.com/karrick/godirwalk"
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

		err := godirwalk.Walk(path, &godirwalk.Options{
			Unsorted:            true,
			FollowSymbolicLinks: false,
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				// Skipping incoming errors; we actually don't care when scanning
				// if we can't traverse. The only thing seen so far here is
				// permissions errors while scanning from root; perhaps we should
				// log though.
				return godirwalk.SkipNode
			},
			Callback: func(path string, info *godirwalk.Dirent) error {
				if !info.IsDir() {
					return nil
				}
				if _, dir := filepath.Split(path); false ||
					dir == ".git" ||
					dir == ".hg" ||
					dir == ".svn" {
					return nil
				}
				if ok, _ := containsSimpleProjectUnchecked(path); !ok {
					return nil
				}

				config, err := loadConfigFromDir(path)
				found := &FoundProject{Path: path, Config: config, Err: err}

				select {
				case result <- found:
				case <-stop:
					return errStop
				}

				return nil
			},
		})

		if err != nil && err != errStop {
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
