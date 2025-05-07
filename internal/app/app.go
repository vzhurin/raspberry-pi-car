package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(c *Container) {
	logger := c.Logger().With(slog.String("where", "app"))
	
	c.HTTPServer().Run()

	// TODO check
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Info("shutdown signal received")
	case err := <-c.HTTPServer().ErrNotify():
		logger.Error("http server error", slog.Any("error", err))
	}

	logger.Info("shutting down")

	if err := c.HTTPServer().Shutdown(); err != nil {
		logger.Error("http server shutdown error", slog.Any("error", err))
	}
}
