package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	prj "github.com/shabbyrobe/prj"
)

type initCommand struct {
	name string
	dest string
}

func (cmd *initCommand) Help() cmdy.Help {
	return cmdy.Synopsis("Initialise a project in the given directory")
}

func (cmd *initCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	flags.StringVar(&cmd.name, "name", "", "Name for this project (defaults to the last part of the directory")
	args.StringOptional(&cmd.dest, "dest", "", "Initialise in this destination. Uses current directory if empty.")
}

func (cmd *initCommand) Run(ctx cmdy.Context) error {
	dest := cmd.dest

	if dest == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		dest = wd
	}

	dest, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	name := cmd.name
	if name == "" {
		name = strings.TrimRight(dest, string(filepath.Separator))
		_, name = filepath.Split(name)
	}

	session, err := prj.NewOSSession()
	if err != nil {
		return err
	}

	_, config, err := prj.InitSimpleProject(ctx, session, dest, name, time.Now())
	if err != nil {
		return err
	}

	fmt.Printf("Project %q initialised in %q\n", config.Name, dest)

	return nil
}
