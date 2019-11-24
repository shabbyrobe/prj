package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	prj "github.com/shabbyrobe/prj"
)

type listCommand struct {
	child string
}

func (cmd *listCommand) Synopsis() string { return "List files tracked by the project" }

func (cmd *listCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	args.StringOptional(&cmd.child, "child", "", "Limit status check to child path, if passed")
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

	out := ctx.Stdout()
	limit := prj.NewResourcePath(cmd.child)

	for _, f := range status.Files {
		if limit == "" || f.Name.IsChildOf(limit) {
			fmt.Fprintln(out, f.Name)
		}
	}

	return nil
}
