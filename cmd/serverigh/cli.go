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
				Name:  "root",
				Usage: "filesystem `PATH` to expose",
				Value: cfg.Root,
			},
			&cli.StringFlag{
				Name:  "host",
				Usage: "host `ADDRESS` to bind",
				Value: cfg.Host,
			},
			&cli.IntFlag{
				Name:  "port",
				Usage: "port `NUMBER` to bind",
				Value: cfg.Port,
			},
			&cli.StringFlag{
				Name:  "theme",
				Usage: "initial UI `THEME`",
				Value: cfg.Theme,
			},
			&cli.BoolFlag{
				Name:  "write",
				Usage: "enable mutating file actions",
				Value: cfg.Write,
			},
			&cli.BoolFlag{
				Name:  "show-hidden",
				Usage: "show hidden files",
				Value: cfg.ShowHidden,
			},
			&cli.Int64Flag{
				Name:  "max-preview-bytes",
				Usage: "maximum `BYTES` to read for text-like previews",
				Value: cfg.MaxPreviewBytes,
			},
			&cli.BoolFlag{
				Name:  "column-size",
				Usage: "show the size column",
				Value: cfg.Columns.Size,
			},
			&cli.BoolFlag{
				Name:  "column-modified",
				Usage: "show the modified column",
				Value: cfg.Columns.Modified,
			},
			&cli.BoolFlag{
				Name:  "column-created",
				Usage: "show the created column",
				Value: cfg.Columns.Created,
			},
			&cli.BoolFlag{
				Name:  "column-mode",
				Usage: "show the mode column",
				Value: cfg.Columns.Mode,
			},
		},
		Action: func(ctx *cli.Context) error {
			loaded, err := config.Load(*cfg, flagOverrides(ctx))
			if err != nil {
				return err
			}
			*cfg = loaded

			if err := cfg.Normalize(); err != nil {
				return err
			}

			return action(*cfg)
		},
	}
}

func flagOverrides(ctx *cli.Context) config.Overrides {
	var overrides config.Overrides

	if ctx.IsSet("root") {
		value := ctx.String("root")
		overrides.Root = &value
	}
	if ctx.IsSet("host") {
		value := ctx.String("host")
		overrides.Host = &value
	}
	if ctx.IsSet("port") {
		value := ctx.Int("port")
		overrides.Port = &value
	}
	if ctx.IsSet("theme") {
		value := ctx.String("theme")
		overrides.Theme = &value
	}
	if ctx.IsSet("write") {
		value := ctx.Bool("write")
		overrides.Write = &value
	}
	if ctx.IsSet("show-hidden") {
		value := ctx.Bool("show-hidden")
		overrides.ShowHidden = &value
	}
	if ctx.IsSet("max-preview-bytes") {
		value := ctx.Int64("max-preview-bytes")
		overrides.MaxPreviewBytes = &value
	}

	if ctx.IsSet("column-size") {
		value := ctx.Bool("column-size")
		overrides.Columns.Size = &value
	}
	if ctx.IsSet("column-modified") {
		value := ctx.Bool("column-modified")
		overrides.Columns.Modified = &value
	}
	if ctx.IsSet("column-created") {
		value := ctx.Bool("column-created")
		overrides.Columns.Created = &value
	}
	if ctx.IsSet("column-mode") {
		value := ctx.Bool("column-mode")
		overrides.Columns.Mode = &value
	}

	return overrides
}
