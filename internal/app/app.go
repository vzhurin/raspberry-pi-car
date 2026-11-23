package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(c *Container) error {
	logger := c.Logger().With(slog.String("where", "app"))

	ctx := context.Background()
	c.Watcher().Watch(ctx)
	c.Dispatcher().Dispatch(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Info("shutdown signal received")
	case err := <-c.Watcher().ErrNotify():
		logger.Error("watcher error", slog.Any("error", err))
	case err := <-c.Dispatcher().ErrNotify():
		logger.Error("dispatcher error", slog.Any("error", err))
	}

	logger.Info("shutting down")
	c.Watcher().Shutdown()
	c.Dispatcher().Shutdown()

	return nil
}
