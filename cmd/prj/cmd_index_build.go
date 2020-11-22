package main

import (
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
)

type indexBuildCommand struct {
	app *App
}

func (cmd *indexBuildCommand) Help() cmdy.Help {
	return cmdy.Synopsis("Index all projects found under the configured directories")
}

func (cmd *indexBuildCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {}

func (cmd *indexBuildCommand) Run(ctx cmdy.Context) error {
	return nil
}
