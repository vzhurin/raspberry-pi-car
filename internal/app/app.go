package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(c *Container) {
	c.HTTPServer().Run()

	// TODO check
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		fmt.Println("signal received, shutting down") // TODO logger
	case err := <-c.HTTPServer().ErrNotify():
		fmt.Println("http server error:", err) // TODO logger
	}

	fmt.Println("shutting down")

	if err := c.HTTPServer().Shutdown(); err != nil {
		fmt.Println("http server shutdown error :", err) // TODO logger
	}
}
