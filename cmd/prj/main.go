package main

import (
	"context"
	"os"

	"github.com/shabbyrobe/cmdy"
)

func main() {
	if err := run(); err != nil {
		cmdy.Fatal(err)
	}
}

func run() error {
	mainGroup := func() (cmdy.Command, cmdy.Init) {
		return cmdy.NewGroup(
			"prj: your friendly arbitrary project folder helper",
			cmdy.Builders{
				"find":   func() (cmdy.Command, cmdy.Init) { return &findCommand{}, nil },
				"init":   func() (cmdy.Command, cmdy.Init) { return &initCommand{}, nil },
				"log":    func() (cmdy.Command, cmdy.Init) { return &logCommand{}, nil },
				"mark":   func() (cmdy.Command, cmdy.Init) { return &markCommand{}, nil },
				"status": func() (cmdy.Command, cmdy.Init) { return &statusCommand{}, nil },
			},
		), nil
	}

	return cmdy.Run(context.Background(), os.Args[1:], mainGroup)
}
