package main

import (
	"fmt"
	"os"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
	prj "github.com/shabbyrobe/go-prj"
)

type findCommand struct {
	paths []string
}

func (cmd *findCommand) Synopsis() string { return "Find projects on the filesystem" }

func (cmd *findCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	set.Remaining(&cmd.paths, "paths", args.AnyLen, "List of paths to search for projects. Uses CWD if empty")
	return set
}

func (cmd *findCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	return set
}

func (cmd *findCommand) Run(ctx cmdy.Context) error {
	if len(cmd.paths) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		cmd.paths = []string{wd}
	}

	for _, path := range cmd.paths {
		projects, err := prj.Find(path)
		if err != nil {
			return err
		}

		fmt.Println(projects)
	}

	return nil
}
