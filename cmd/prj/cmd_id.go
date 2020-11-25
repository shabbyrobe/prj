package main

import (
	"fmt"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	prj "github.com/shabbyrobe/prj"
)

type idCommand struct{}

func (cmd *idCommand) Help() cmdy.Help { return cmdy.Synopsis("Show project ID") }

func (cmd *idCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {}

func (cmd *idCommand) Run(ctx cmdy.Context) error {
	allKinds := []prj.ProjectKind{prj.ProjectSimple, prj.ProjectGit, prj.ProjectHg}
	project, _, err := loadProject("", allKinds)
	if err != nil {
		return err
	}

	fmt.Println(project.ID())

	return nil
}
