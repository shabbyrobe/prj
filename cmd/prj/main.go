package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/shabbyrobe/cmdy"
)

func main() {
	if err := run(); err != nil {
		cmdy.Fatal(err)
	}
}

type App struct {
	config             Config
	configFileOverride string
	wd                 string
	configPath         string
	cachePath          string
}

func (app App) ConfigFile() string {
	if app.configFileOverride != "" {
		return app.configFileOverride
	}
	return filepath.Join(app.configPath, "config.toml")
}

func run() error {
	var app App

	userConfigPath, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	userCachePath, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	app.configPath = filepath.Join(userConfigPath, "shabbyrobe", "prj")
	app.cachePath = filepath.Join(userCachePath, "shabbyrobe", "prj")

	mainGroup := func() cmdy.Command {
		indexGroup := func() cmdy.Command {
			return cmdy.NewGroup(
				"Tools to build and search an index of found projects",
				cmdy.Builders{
					"build": func() cmdy.Command { return &indexBuildCommand{app: &app} },
				},
			)
		}

		return cmdy.NewGroup(
			"prj: your friendly arbitrary project folder helper",

			cmdy.Builders{
				"diff":  func() cmdy.Command { return &diffCommand{} },
				"find":  func() cmdy.Command { return &findCommand{} },
				"hash":  func() cmdy.Command { return &hashCommand{} },
				"list":  func() cmdy.Command { return &listCommand{} },
				"init":  func() cmdy.Command { return &initCommand{} },
				"index": indexGroup,
				"log":   func() cmdy.Command { return &logCommand{} },
				"mark":  func() cmdy.Command { return &markCommand{} },
			},

			cmdy.GroupFlags(func() *cmdy.FlagSet {
				flags := cmdy.NewFlagSet()
				flags.StringVar(&app.wd, "C", "", "Run subcommand inside this working directory (instead of cwd)")
				flags.StringVar(&app.configFileOverride, "config", "", "Use this config file instead of the one in your user config dir")
				return flags
			}),

			cmdy.GroupBefore(func(ctx cmdy.Context) error {
				if app.wd != "" {
					if err := os.Chdir(app.wd); err != nil {
						return fmt.Errorf("-C option invalid, chdir failed: %w", app.wd)
					}
				}

				configFile := app.ConfigFile()
				if _, err := toml.DecodeFile(configFile, &app.config); err != nil {
					if !os.IsNotExist(err) || app.configFileOverride != "" {
						return err
					}
				}

				return nil
			}),

			// cmdy.GroupPrefixMatcher(2),
		)
	}

	return cmdy.Run(context.Background(), os.Args[1:], mainGroup)
}
