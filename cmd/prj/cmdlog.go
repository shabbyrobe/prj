package main

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bbrks/wrap"
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
)

const (
	displayShort = "short"
	displayFull  = "full"
)

type logCommand struct {
	display string
}

func (cmd *logCommand) Synopsis() string { return "Show the commit log for this project" }

func (cmd *logCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	return set
}

func (cmd *logCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	set.StringVar(&cmd.display, "display", "short", "Display mode (short, full)")
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
		return entries[i].ModTime.Before(entries[j].ModTime)
	})

	out := ctx.Stdout()

	switch cmd.display {
	case displayShort:
		w := tabwriter.NewWriter(out, 8, 4, 2, ' ', 0)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "TIME", "AUTHOR", "BYTES", "FILES", "MSG")

		for _, msg := range entries {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n",
				msg.ModTime.Format(time.RFC3339),
				fmt.Sprintf("%s@%s", msg.Author, msg.Machine),
				msg.Size,
				msg.FileCount,
				truncate(msg.Message, 50))
		}
		w.Flush()

	case displayFull:
		for _, msg := range entries {
			fmt.Fprintf(out, ""+
				"date:     %s\n"+
				"hash:     %s\n"+
				"contents: %d byte(s), %d file(s)\n"+
				"author:   %s\n"+
				"\n%s\n",

				msg.ModTime,
				msg.Hash,
				msg.Size, msg.FileCount,
				fmt.Sprintf("%s@%s", msg.Author, msg.Machine),
				indent(msg.Message))
		}

	default:
		return cmdy.NewUsageErrorf("unknown -display %q", cmd.display)
	}

	return nil
}

func truncate(str string, sz int) string {
	first := strings.IndexAny(strings.TrimSpace(str), "\n\r")
	if first >= 0 {
		str = str[:first]
	}
	if len(str) > sz {
		str = str[:sz] + "..."
	}
	return str
}

func indent(str string) string {
	w := wrap.NewWrapper()
	w.OutputLinePrefix = "    "
	out := w.Wrap(strings.TrimSpace(str), 100)
	return out
}
