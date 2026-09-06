package main

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"serverigh/internal/apperror"
	"serverigh/internal/config"
	"serverigh/internal/httpserver"
)

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("serverigh stopped", "error", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := config.Default()
	if err != nil {
		return err
	}

	app := newCLI(&cfg, serve)

	return app.Run(args)
}

func serve(cfg config.Config) error {
	app, err := httpserver.New(cfg)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)

	go func() {
		slog.Info(
			"starting serverigh",
			"addr",
			cfg.Address(),
			"root",
			cfg.Root,
			"write",
			cfg.Write,
		)

		errCh <- app.Listen(cfg.Address())
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownCh)

	select {
	case sig := <-shutdownCh:
		slog.Info("shutting down serverigh", "signal", sig.String())

		if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
			return apperror.WrapOperation(
				apperror.CodeServer,
				"server error",
				"shutdown server",
				err,
			)
		}

		return nil
	case err := <-errCh:
		if err == nil {
			return nil
		}

		if errors.Is(err, os.ErrClosed) {
			return nil
		}

		return apperror.WrapOperation(
			apperror.CodeServer,
			"server error",
			"listen",
			err,
		)
	}
}
