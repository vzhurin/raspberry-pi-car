package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	server    *http.Server
	errNotify chan error
}

type Config struct {
	Host string
	Port string

	GinMode string
}

func NewServer(cfg *Config, handler *Handler) (*Server, error) {
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	registerHandler(router, handler)

	server := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: router,
	}

	return &Server{
		server:    server,
		errNotify: make(chan error, 1), // TODO check
	}, nil
}

func (s *Server) Run() {
	go func() {
		s.errNotify <- s.run()
		close(s.errNotify)
	}()
}

func (s *Server) ErrNotify() <-chan error {
	return s.errNotify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	// TODO logger
	fmt.Println("Server shutdown complete")

	return nil
}

func (s *Server) run() error {
	return s.server.ListenAndServe()
}
