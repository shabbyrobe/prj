package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/shabbyrobe/cmdy"
	"github.com/shabbyrobe/cmdy/args"
	prj "github.com/shabbyrobe/prj"
)

type markCommand struct {
	message string
	force   bool
}

func (cmd *markCommand) Synopsis() string { return "Mark the current hash of the project in the log" }

func (cmd *markCommand) Args() *args.ArgSet {
	set := args.NewArgSet()
	return set
}

func (cmd *markCommand) Flags() *cmdy.FlagSet {
	set := cmdy.NewFlagSet()
	set.StringVar(&cmd.message, "m", "", "Mark message")
	set.BoolVar(&cmd.force, "f", false, "Force mark")
	return set
}

func (cmd *markCommand) Run(ctx cmdy.Context) error {
	if cmd.message == "" {
		if !cmdy.ReaderIsPipe(ctx.Stdin()) {
			msg, err := editMarkMessage()
			if err != nil {
				return err
			}
			if msg == "" {
				return cmdy.NewUsageErrorf("prj: mark message was empty; aborting")
			}
			cmd.message = msg

		} else {
			return cmdy.NewUsageErrorf("prj: missing mark message; pass with -m")
		}
	}

	project, session, err := loadProject()
	if err != nil {
		return err
	}

	options := &prj.MarkOptions{
		Force: cmd.force,
	}

	status, err := project.Mark(ctx, session, cmd.message, time.Now(), options)
	if err != nil {
		return err
	}
	_ = status

	return nil
}

func editMarkMessage() (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "", nil
	}

	temp, err := ioutil.TempFile("", "prj-mark-msg")
	if err != nil {
		return "", err
	}
	defer os.Remove(temp.Name())
	if err := temp.Close(); err != nil {
		return "", err
	}

	editCmd := exec.Command(editor, temp.Name())
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr
	if err := editCmd.Run(); err != nil {
		return "", err
	}

	msgData, err := ioutil.ReadFile(temp.Name())
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(msgData)), nil
}
