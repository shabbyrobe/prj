package main

import (
	"fmt"
	"time"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
	prj "github.com/shabbyrobe/prj"
)

type hashCommand struct {
	path string
}

func (cmd *hashCommand) Synopsis() string { return "Display current hash" }

func (cmd *hashCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	set.StringOptional(&cmd.path, "path", "", "Limit status check to child path, if passed")
	return set
}

func (cmd *hashCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	return set
}

func (cmd *hashCommand) Run(ctx cmdy.Context) error {
	project, _, err := loadProject()
	if err != nil {
		return err
	}

	config, err := project.Config()
	if err != nil {
		return err
	}

	path := prj.NewResourcePath(cmd.path)

	start := time.Now()
	status, err := project.Status(ctx, path, start)
	if err != nil {
		return err
	}
	taken := time.Since(start)

	out := ctx.Stdout()

	fmt.Fprintf(out, ""+
		"project:  %s\n"+
		"projid:   %s\n"+
		"\n"+
		"date:     %s\n"+
		"hash:     %s\n"+
		"path:     %q\n"+
		"contents: %d byte(s), %d file(s)\n"+
		"taken:    %s\n",

		config.Name,
		config.ID,
		status.ModTime,
		status.Hash,
		path,
		status.Size, len(status.Files),
		taken)

	return nil
}
