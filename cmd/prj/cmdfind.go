package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
	prj "github.com/shabbyrobe/prj"
)

type findCommand struct {
	paths []string
}

func (cmd *findCommand) Synopsis() string { return "Find projects on the filesystem" }

func (cmd *findCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	set.Remaining(&cmd.paths, "paths", args.AnyLen, "List of paths to search for projects. Uses CWD if empty")
	return set
}

func (cmd *findCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	return set
}

func (cmd *findCommand) Run(ctx cmdy.Context) error {
	if len(cmd.paths) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		cmd.paths = []string{wd}
	}

	out := ctx.Stdout()

	type col struct {
		name   string
		width  int
		tpl    string
		hdrtpl string
		hide   bool
	}

	var cols = []col{
		{"ID", 36, "%-*s", "%-*s", true},
		{"NAME", 24, "%-*s", "%-*s", false},
		{"LASTMOD", 32, "%-*s", "%-*s", false},
		{"HASH", 32, "%-*s", "%-*s", true},
		{"PATH", 0, "%-*s", "%-*s", false},
	}

	var hdrTpls = make([]string, 0, len(cols))
	var rowTpls = make([]string, 0, len(cols))
	for _, col := range cols {
		if !col.hide {
			hdrTpls = append(hdrTpls, col.hdrtpl)
			rowTpls = append(rowTpls, col.tpl)
		}
	}
	var hdrTpl = strings.Join(hdrTpls, " ") + "\n"
	var rowTpl = strings.Join(rowTpls, " ") + "\n"

	var row = []interface{}{}
	for _, col := range cols {
		if !col.hide {
			row = append(row, col.width, col.name)
		}
	}
	fmt.Fprintf(out, hdrTpl, row...)

	var failed []*prj.FoundProject

	for _, path := range cmd.paths {
		scn := prj.Scan(path)
		for scn.Next() {
			project := scn.Project()

			if project.Config == nil {
				failed = append(failed, project)
				continue
			}

			row = row[:0]
			if !cols[0].hide {
				row = append(row, cols[0].width, project.Config.ID)
			}
			if !cols[1].hide {
				row = append(row, cols[1].width, project.Config.Name)
			}
			if !cols[2].hide {
				if project.Config.LastEntry != nil && !project.Config.LastEntry.ModTime.IsZero() {
					row = append(row, cols[2].width, project.Config.LastEntry.ModTime)
				} else {
					row = append(row, cols[2].width, "<none>")
				}
			}
			if !cols[3].hide {
				if project.Config.LastEntry != nil {
					row = append(row, cols[3].width, project.Config.LastEntry.Hash.String())
				} else {
					row = append(row, cols[3].width, "<none>")
				}
			}
			if !cols[4].hide {
				row = append(row, cols[4].width, project.Path)
			}

			fmt.Fprintf(out, rowTpl, row...)
		}

		if err := scn.Close(); err != nil {
			return err
		}
	}

	if len(failed) > 0 {
		fmt.Fprintln(out)
		for _, fprj := range failed {
			fmt.Fprintf(out, "ERROR: could not load %q: %v\n", fprj.Path, fprj.Err)
		}
	}

	return nil
}
