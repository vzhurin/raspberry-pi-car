package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

func Run(c *Container) error {
	defer binsdl.Load().Unload()
	defer sdl.Quit()
	err := sdl.Init(sdl.INIT_GAMEPAD)
	if err != nil {
		return err
	}

	logger := c.Logger().With(slog.String("where", "app"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	c.Watcher().Watch(ctx)
	c.JoystickDispatcher().Dispatch(ctx)
	c.JoystickListener().Run(ctx)

	select {
	case <-quit:
		logger.Info("shutdown signal received")
	case err := <-c.Watcher().ErrNotify():
		logger.Error("watcher error", slog.Any("error", err))
	case err := <-c.JoystickDispatcher().ErrNotify():
		logger.Error("dispatcher error", slog.Any("error", err))
	case err := <-c.JoystickListener().ErrNotify():
		logger.Error("listener error", slog.Any("error", err))
	}

	logger.Info("shutting down")
	cancel()
	c.Watcher().Shutdown()
	c.JoystickDispatcher().Shutdown()
	c.JoystickListener().Shutdown()

	return nil
}
