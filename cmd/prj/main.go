package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shabbyrobe/cmdy"
)

func main() {
	if err := run(); err != nil {
		cmdy.Fatal(err)
	}
}

func run() error {
	var wd string

	mainGroup := func() cmdy.Command {
		return cmdy.NewGroup(
			"prj: your friendly arbitrary project folder helper",
			cmdy.Builders{
				"diff": func() cmdy.Command { return &diffCommand{} },
				"find": func() cmdy.Command { return &findCommand{} },
				"hash": func() cmdy.Command { return &hashCommand{} },
				"init": func() cmdy.Command { return &initCommand{} },
				"log":  func() cmdy.Command { return &logCommand{} },
				"mark": func() cmdy.Command { return &markCommand{} },
			},
			cmdy.GroupFlags(func() *cmdy.FlagSet {
				flags := cmdy.NewFlagSet()
				flags.StringVar(&wd, "C", "", "Run subcommand inside this working directory (instead of cwd)")
				return flags
			}),
			cmdy.GroupBefore(func(ctx cmdy.Context) error {
				if wd != "" {
					if err := os.Chdir(wd); err != nil {
						return fmt.Errorf("-C option invalid, chdir failed: %w", wd)
					}
				}
				return nil
			}),
			// cmdy.GroupPrefixMatcher(2),
		)
	}

	return cmdy.Run(context.Background(), os.Args[1:], mainGroup)
}
