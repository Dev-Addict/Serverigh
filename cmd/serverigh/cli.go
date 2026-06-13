package main

import (
	"github.com/urfave/cli/v2"

	"serverigh/internal/config"
)

type serverAction func(config.Config) error

func newCLI(cfg *config.Config, action serverAction) *cli.App {
	return &cli.App{
		Name:            "serverigh",
		Usage:           "browse local files from a web UI",
		UsageText:       "serverigh [options]",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "root",
				Usage:       "filesystem `PATH` to expose",
				Value:       cfg.Root,
				EnvVars:     []string{"SERVERIGH_ROOT"},
				Destination: &cfg.Root,
			},
			&cli.StringFlag{
				Name:        "host",
				Usage:       "host `ADDRESS` to bind",
				Value:       cfg.Host,
				EnvVars:     []string{"SERVERIGH_HOST"},
				Destination: &cfg.Host,
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "port `NUMBER` to bind",
				Value:       cfg.Port,
				EnvVars:     []string{"SERVERIGH_PORT"},
				Destination: &cfg.Port,
			},
			&cli.BoolFlag{
				Name:        "write",
				Usage:       "enable mutating file actions",
				Value:       cfg.Write,
				EnvVars:     []string{"SERVERIGH_WRITE"},
				Destination: &cfg.Write,
			},
			&cli.BoolFlag{
				Name:        "show-hidden",
				Usage:       "show dotfiles by default",
				Value:       cfg.ShowHidden,
				EnvVars:     []string{"SERVERIGH_SHOW_HIDDEN"},
				Destination: &cfg.ShowHidden,
			},
			&cli.Int64Flag{
				Name:        "max-preview-bytes",
				Usage:       "maximum `BYTES` to read for text-like previews",
				Value:       cfg.MaxPreviewBytes,
				EnvVars:     []string{"SERVERIGH_MAX_PREVIEW_BYTES"},
				Destination: &cfg.MaxPreviewBytes,
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := cfg.Normalize(); err != nil {
				return err
			}

			return action(*cfg)
		},
	}
}
