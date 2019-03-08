package main

import (
	"fmt"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
)

type commitCommand struct {
	message string
}

func (cmd *commitCommand) Synopsis() string { return "Commit the current hash to the log" }

func (cmd *commitCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	return set
}

func (cmd *commitCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	set.StringVar(&cmd.message, "m", "", "Commit message")
	return set
}

func (cmd *commitCommand) Run(ctx cmdy.Context) error {
	if cmd.message == "" {
		// FIXME: run editor
		return fmt.Errorf("prj: missing commit message")
	}

	project, session, err := loadProject()
	if err != nil {
		return err
	}

	status, err := project.Commit(ctx, session, cmd.message, nil)
	if err != nil {
		return err
	}
	_ = status

	return nil
}
