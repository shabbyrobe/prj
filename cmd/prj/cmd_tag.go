package main

import (
	"fmt"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/arg"
)

type tagCommand struct {
	tags []string
}

func (cmd *tagCommand) Help() cmdy.Help {
	return cmdy.Synopsis("Show or edit project tags")
}

func (cmd *tagCommand) Configure(flags *cmdy.FlagSet, args *arg.ArgSet) {
	args.Remaining(&cmd.tags, "tags", arg.AnyLen, ""+
		"Tags to add/remove from project. Prefix with '-' to remove "+
		"(you will need to use '--' to avoid interpreting as flags).")
}

func (cmd *tagCommand) Run(ctx cmdy.Context) error {
	project, _, err := loadSimpleProject("")
	if err != nil {
		return err
	}

	tagger := project.Tagger()
	if len(cmd.tags) > 0 {
		add := []string{}
		rem := []string{}
		for _, tag := range cmd.tags {
			if len(tag) == 0 {
				continue
			}
			if tag[0] == '-' {
				rem = append(rem, tag[1:])
			} else {
				add = append(add, tag)
			}
		}
		if err := tagger.Tag(add...); err != nil {
			return err
		}
		if err := tagger.Untag(rem...); err != nil {
			return err
		}
	}

	tags, err := tagger.Tags()
	if err != nil {
		return err
	}

	for _, t := range tags {
		fmt.Println(t)
	}

	return nil
}
