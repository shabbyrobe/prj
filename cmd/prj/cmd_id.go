package main

import (
	"fmt"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
)

type idCommand struct{}

func (cmd *idCommand) Help() cmdy.Help { return cmdy.Synopsis("Show project ID") }

func (cmd *idCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {}

func (cmd *idCommand) Run(ctx cmdy.Context) error {
	project, _, err := loadSimpleProject("")
	if err != nil {
		return err
	}

	fmt.Println(project.ID())

	return nil
}
