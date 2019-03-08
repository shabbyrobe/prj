package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
	prj "github.com/shabbyrobe/go-prj"
)

type initCommand struct {
	name string
	dest string
}

func (cmd *initCommand) Synopsis() string { return "Initialise a project in the given directory" }

func (cmd *initCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	set.StringOptional(&cmd.dest, "dest", "", "Initialise in this destination. Uses current directory if empty.")
	return set
}

func (cmd *initCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	set.StringVar(&cmd.name, "name", "", "Name for this project (defaults to the last part of the directory")
	return set
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

	config, err := prj.InitSimpleProject(dest, name, time.Now())
	if err != nil {
		return err
	}

	fmt.Printf("Project %q initialised in %q\n", config.Name, dest)

	return nil
}
