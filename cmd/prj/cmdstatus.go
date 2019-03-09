package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
)

type statusCommand struct{}

func (cmd *statusCommand) Synopsis() string { return "Show the list of changed files" }

func (cmd *statusCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	return set
}

func (cmd *statusCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	return set
}

func (cmd *statusCommand) Run(ctx cmdy.Context) error {
	project, _, err := loadProject()
	if err != nil {
		return err
	}

	start := time.Now()
	status, err := project.Status(ctx, "")
	if err != nil {
		return err
	}

	taken := time.Since(start)
	spew.Dump(len(status.Files))
	spew.Dump(status.Size)
	spew.Dump(status.Hash)
	spew.Dump(taken)

	return nil
}
