package app

import (
	"raspberry-pi-car/config"
	apiHTTP "raspberry-pi-car/internal/api/http"
	"sync"
)

type Container struct {
	cfg *config.Config
	mu  sync.Mutex

	httpServer  *apiHTTP.Server
	httpHandler *apiHTTP.Handler
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{cfg: cfg}
}

func (c *Container) HTTPServer() *apiHTTP.Server {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpServer == nil {
		httpServer, err := apiHTTP.NewServer(&apiHTTP.Config{
			Host:    c.cfg.RESTHost,
			Port:    c.cfg.RESTPort,
			GinMode: c.cfg.GinMode,
		}, c.HTTPHandler())

		if err != nil {
			panic(err)
		}

		c.httpServer = httpServer
	}

	return c.httpServer
}

func (c *Container) HTTPHandler() *apiHTTP.Handler {
	if c.httpHandler == nil {
		c.httpHandler = apiHTTP.NewHandler()
	}

	return c.httpHandler
}
