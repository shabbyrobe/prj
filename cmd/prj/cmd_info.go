package main

import (
	"fmt"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
)

type infoCommand struct{}

func (cmd *infoCommand) Help() cmdy.Help { return cmdy.Synopsis("Show project info") }

func (cmd *infoCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {}

func (cmd *infoCommand) Run(ctx cmdy.Context) error {
	project, _, err := loadSimpleProject("")
	if err != nil {
		return err
	}

	fmt.Println("ID:", project.ID())
	fmt.Println("Kind:", project.Kind())
	fmt.Println("Name:", project.Name())
	fmt.Println("Path:", project.Path())

	return nil
}
