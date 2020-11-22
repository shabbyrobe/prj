package main

import (
	"fmt"
	"time"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
	prj "github.com/shabbyrobe/prj"
)

const hashUsage = cmdy.DefaultUsage + `
NOTE: this does not yet work with Git or Mercurial projects.
`

type hashCommand struct {
	child   string
	rawPath string
}

func (cmd *hashCommand) Help() cmdy.Help {
	return cmdy.Help{
		Synopsis: "Show hash of the current state of a project",
		Usage:    hashUsage,
	}
}

func (cmd *hashCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	flags.StringVar(&cmd.rawPath, "raw", "", "Hash path at -raw, even if it is not a 'prj' project")
	args.StringOptional(&cmd.child, "child", "", "Limit status check to child path, if passed")
}

func (cmd *hashCommand) Run(ctx cmdy.Context) error {
	project, _, done, err := loadSimpleProjectWithTemporaryFallback(ctx, "", cmd.rawPath)
	if err != nil {
		return err
	}
	defer done()

	path := prj.NewResourcePath(cmd.child)

	start := time.Now()
	status, err := project.Status(ctx, path, start)
	if err != nil {
		return err
	}
	taken := time.Since(start)

	out := ctx.Stdout()

	fmt.Fprintf(out, ""+
		"project:  %s\n"+
		"projid:   %s\n"+
		"\n"+
		"lastmod:  %s\n"+
		"hash:     %s\n"+
		"path:     %q\n"+
		"contents: %s, %d byte(s), %d file(s)\n"+
		"taken:    %s\n",

		project.Name(),
		project.ID(),
		status.ModTime,
		status.Hash,
		path,
		bytesHuman(status.Size, 3), status.Size, len(status.Files),
		taken)

	return nil
}
