package main

import (
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
)

type markCommand struct {
	message string
}

func (cmd *markCommand) Synopsis() string { return "Mark the current hash of the project in the log" }

func (cmd *markCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	return set
}

func (cmd *markCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	set.StringVar(&cmd.message, "m", "", "Mark message")
	return set
}

func (cmd *markCommand) Run(ctx cmdy.Context) error {
	if cmd.message == "" {
		// FIXME: run editor
		return cmdy.NewUsageErrorf("prj: missing mark message; pass with -m")
	}

	project, session, err := loadProject()
	if err != nil {
		return err
	}

	status, err := project.Mark(ctx, session, cmd.message, nil)
	if err != nil {
		return err
	}
	_ = status

	return nil
}
