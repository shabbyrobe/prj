package main

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bbrks/wrap"
	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	"github.com/shabbyrobe/golib/errtools"
	"github.com/shabbyrobe/prj"
)

const (
	displayShort = "short"
	displayFull  = "full"
)

type logCommand struct {
	display string
}

func (cmd *logCommand) Help() cmdy.Help {
	return cmdy.Synopsis("Show the commit log for this project")
}

func (cmd *logCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	flags.StringVar(&cmd.display, "display", "short", "Display mode (short, full)")
}

func (cmd *logCommand) Run(ctx cmdy.Context) (rerr error) {
	project, _, err := loadSimpleProject("")
	if err != nil {
		return err
	}

	iter := project.Log()
	defer errtools.DeferClose(&rerr, iter)

	out := ctx.Stdout()

	var entry prj.LogEntry

	switch cmd.display {
	case displayShort:
		w := tabwriter.NewWriter(out, 8, 4, 2, ' ', 0)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "TIME", "AUTHOR", "BYTES", "FILES", "MSG")

		for iter.Next(&entry) {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n",
				entry.Time.Format(time.RFC3339),
				fmt.Sprintf("%s@%s", entry.Author, entry.Machine),
				entry.Size,
				entry.FilesCount,
				truncate(entry.Message, 50))
		}
		w.Flush()

	case displayFull:
		for iter.Next(&entry) {
			fmt.Fprintf(out, ""+
				"date:     %s\n"+
				"hash:     %s\n"+
				"contents: %s, %d byte(s), %d file(s)\n"+
				"author:   %s\n"+
				"\n%s\n",

				entry.Time,
				entry.Hash,
				bytesHuman(entry.Size, 3), entry.Size, entry.FilesCount,
				fmt.Sprintf("%s@%s", entry.Author, entry.Machine),
				indent(entry.Message))
		}

	default:
		return cmdy.UsageErrorf("unknown -display %q", cmd.display)
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
