package main

import (
	"fmt"
	"time"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	prj "github.com/shabbyrobe/prj"
)

type diffCommand struct {
	path  string
	stats bool
	all   bool
}

func (cmd *diffCommand) Synopsis() string { return "Show the list of changed files" }

func (cmd *diffCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	flags.BoolVar(&cmd.stats, "stats", false, "Print some stats at the end")
	flags.BoolVar(&cmd.all, "all", false, "Print identical files too")
	args.StringOptional(&cmd.path, "path", "", "Limit status check to child path, if passed")
}

func (cmd *diffCommand) Run(ctx cmdy.Context) error {
	project, _, err := loadProject("")
	if err != nil {
		return err
	}

	start := time.Now()
	diff, err := project.Diff(ctx, prj.NewResourcePath(cmd.path), time.Now())
	if err != nil {
		return err
	}
	taken := time.Since(start)

	out := ctx.Stdout()

	items := diff.Items()
	for _, item := range items {
		if cmd.all || item.Status != prj.DiffSame {
			fmt.Fprintf(out, " %c %s\n", item.Status, item.Path)
		}
	}

	if cmd.stats {
		fmt.Fprintln(ctx.Stderr(), "\ntime taken:", taken)
	}

	return nil
}
