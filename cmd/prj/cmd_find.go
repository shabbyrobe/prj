package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	"github.com/shabbyrobe/cmdy/flags"
	"github.com/shabbyrobe/prj"
)

type findCommand struct {
	paths            []string
	showID           bool
	showHash         bool
	nested           bool
	kinds            prj.ProjectKindSet
	exclude          flags.StringList
	noDefaultExclude bool
}

func (cmd *findCommand) Help() cmdy.Help { return cmdy.Synopsis("Find projects on the filesystem") }

func (cmd *findCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	kinds := ""
	for i := prj.ProjectKind(0); int(i) < prj.ProjectKindSize; i++ {
		if i.IsValid() {
			if kinds != "" {
				kinds += ", "
			}
			kinds += i.String()
		}
	}

	flags.BoolVar(&cmd.showID, "id", false, "Show ID")
	flags.BoolVar(&cmd.showHash, "hash", false, "Show Hash")
	flags.BoolVar(&cmd.nested, "nested", false, "Find nested projects (i.e. .git within .git)")
	flags.Var(&cmd.kinds, "kind", "Show these kinds, all by default. Can pass multiple times. ("+kinds+")")
	flags.Var(&cmd.exclude, "exclude", "List of regexps to exclude. Must include anchors if desired. `/` matches `\\` as well.")
	flags.BoolVar(&cmd.noDefaultExclude, "no-default-exclude", false, "Exclude the default list of exclude paths")
	args.Remaining(&cmd.paths, "paths", arg.AnyLen, "List of paths to search for projects. Uses CWD if empty")
}

func (cmd *findCommand) Run(ctx cmdy.Context) error {
	if len(cmd.paths) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		cmd.paths = []string{wd}
	}

	if cmd.kinds.Count() == 0 {
		cmd.kinds.SetAll()
	}

	out := ctx.Stdout()

	var exclude = []string{}

	// FIXME: add this to a global configuration
	if !cmd.noDefaultExclude {
		exclude = append(exclude,
			`/\.cargo\/`,
			`/\.cache\/`,
			`/\.npm\/`,
		)
		if gopath := os.Getenv("GOPATH"); gopath != "" {
			exclude = append(exclude, regexp.QuoteMeta(gopath))
		}
	}

	// FIXME: this is yet another experiment to make tabular CLI layout code
	// flexible without jank. it failed, spectacularly.
	//
	// tabwriter doesn't work here as we want the results to appear as they are
	// found by the Scanner, not in one big hit.
	//
	// maybe if this was wrapped up in a simple struct it wouldn't be so bad.

	type col struct {
		name   string
		width  int
		tpl    string
		hdrtpl string
		hide   bool
	}

	const (
		colID = iota
		colKind
		colName
		colLastMod
		colPath
		colHash
	)

	var cols = []col{
		colID:      {"ID", 36, "%-*s", "%-*s", !cmd.showID},
		colKind:    {"KIND", 6, "%-*s", "%-*s", false},
		colName:    {"PROJECT NAME", 30, "%-*s", "%-*s", false},
		colLastMod: {"LASTMOD", 26, "%-*s", "%-*s", false},
		colPath:    {"PATH", 40, "%-*s", "%-*s", false},
		colHash:    {"HASH", 0, "%-*s", "%-*s", !cmd.showHash},
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
	var last int = -1
	for idx, col := range cols {
		if !col.hide {
			last = idx
		}
	}
	if last >= 0 {
		cols[last].width = 0
	}

	for _, col := range cols {
		if !col.hide {
			row = append(row, col.width, col.name)
		}
	}
	fmt.Fprintf(out, hdrTpl, row...)

	var failed []*prj.FoundProject
	var opts []prj.ScanOption
	if cmd.nested {
		opts = append(opts, prj.ScanNested())
	}
	if len(exclude) > 0 {
		opts = append(opts, prj.ScanExcludePattern(exclude...))
	}

	for _, path := range cmd.paths {
		scn := prj.Scan(path, opts...)
		for scn.Next() {
			found := scn.Current()

			if found.Project == nil {
				failed = append(failed, found)
				continue
			}
			if !cmd.kinds[found.Project.Kind()] {
				continue
			}

			lastEntry, err := found.Project.LastEntry()
			if err != nil {
				return err
			}

			row = row[:0]
			if !cols[colID].hide {
				row = append(row, cols[colID].width, found.Project.ID())
			}
			if !cols[colKind].hide {
				row = append(row, cols[colKind].width, found.Project.Kind())
			}
			if !cols[colName].hide {
				row = append(row, cols[colName].width, found.Project.Name())
			}
			if !cols[colLastMod].hide {
				if lastEntry != nil && !lastEntry.ModTime.IsZero() {
					row = append(row, cols[colLastMod].width, lastEntry.ModTime.Format(time.RFC3339))
				} else {
					row = append(row, cols[colLastMod].width, "<none>")
				}
			}
			if !cols[colPath].hide {
				row = append(row, cols[colPath].width, found.Path)
			}
			if !cols[colHash].hide {
				if lastEntry != nil {
					row = append(row, cols[colHash].width, lastEntry.Hash.String())
				} else {
					row = append(row, cols[colHash].width, "<none>")
				}
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
