package main

import (
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
)

type logCommand struct{}

func (cmd *logCommand) Synopsis() string { return "Show the commit log for this project" }

func (cmd *logCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	return set
}

func (cmd *logCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	return set
}

func (cmd *logCommand) Run(ctx cmdy.Context) error {
	return nil
}
