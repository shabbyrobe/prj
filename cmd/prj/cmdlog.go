package main

import (
	"fmt"
	"sort"
	"text/tabwriter"
	"time"

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
	project, _, err := loadProject()
	if err != nil {
		return err
	}

	entries, err := project.Log()
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ModTime.Before(entries[j].ModTime) || entries[i].Message < entries[j].Message
	})

	out := ctx.Stdout()
	w := tabwriter.NewWriter(out, 8, 4, 2, ' ', 0)

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "TIME", "AUTHOR", "BYTES", "FILES", "MSG")

	for _, msg := range entries {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n",
			msg.ModTime.Format(time.RFC3339),
			fmt.Sprintf("%s@%s", msg.Author, msg.Machine),
			msg.Size,
			msg.FileCount,
			msg.Message)
	}
	w.Flush()

	return nil
}
