package main

import (
	"fmt"
	"time"

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
	status, err := project.Status(ctx, "", time.Now())
	if err != nil {
		return err
	}

	taken := time.Since(start)
	fmt.Println(len(status.Files))
	fmt.Println(status.Size)
	fmt.Println(status.Hash)
	fmt.Println(taken)

	return nil
}
