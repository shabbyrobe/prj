package prj

import (
	"os"
	"os/user"
)

type Session struct {
	User    string
	Machine string
}

func NewOSSession() (*Session, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Session{
		User:    user.Username,
		Machine: host,
	}, nil
}
