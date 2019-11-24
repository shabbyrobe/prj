package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/tabwriter"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	prj "github.com/shabbyrobe/prj"
)

type listCommand struct {
	child  string
	format string
}

func (cmd *listCommand) Synopsis() string { return "List files tracked by the project" }

func (cmd *listCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	args.StringOptional(&cmd.child, "child", "", "Limit status check to child path, if passed")
	flags.StringVar(&cmd.format, "fmt", "list", "Output format (list, table, json)")
}

func (cmd *listCommand) Run(ctx cmdy.Context) error {
	project, _, err := loadSimpleProject("")
	if err != nil {
		return err
	}

	// FIXME: This kinda sucks a bit; the API in prj only lets you scan the FS
	// for the status or get the log without the status, but the .prj directory
	// contains a full cache of the last marked status (which is what we want)
	entry, err := project.LastEntry()
	if err != nil {
		return err
	}

	var status prj.ProjectStatus

	// FIXME: yeah, as above, this very sucks:
	bts, err := ioutil.ReadFile(filepath.Join(project.Path(), ".prj", "status", entry.StatusFile))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bts, &status); err != nil {
		return err
	}

	limit := prj.NewResourcePath(cmd.child)

	var filtered []prj.ProjectFile
	if limit == "" {
		filtered = status.Files
	} else {
		filtered = make([]prj.ProjectFile, 0, len(status.Files))

		for _, f := range status.Files {
			if limit == "" || f.Name.IsChildOf(limit) {
				filtered = append(filtered, f)
			}
		}
	}

	out := ctx.Stdout()
	if cmd.format == "list" {
		for _, f := range filtered {
			fmt.Fprintln(out, f.Name)
		}

	} else if cmd.format == "table" {
		w := tabwriter.NewWriter(out, 2, 2, 2, ' ', 0)
		fmt.Fprintf(w, "NAME\tSIZE\tMODTIME\tHASH\n")

		for _, f := range filtered {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", f.Name, f.Size, f.ModTime.Format("2006-01-02T15:04:05"), f.Hash)
		}
		if err := w.Flush(); err != nil {
			return err
		}

	} else if cmd.format == "json" {
		enc := json.NewEncoder(out)
		for _, f := range filtered {
			if err := enc.Encode(f); err != nil {
				return err
			}
		}

	} else {
		return fmt.Errorf("unknown -fmt")
	}

	return nil
}
