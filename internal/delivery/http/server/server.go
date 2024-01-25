package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Be1chenok/effectiveMobileTask/internal/config"
)

type Server struct {
	httpServer http.Server
}

func New(conf *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:           fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port),
			MaxHeaderBytes: 1024 * 1024, // 1 MB
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			Handler:        handler,
		},
	}
}

func (srv *Server) Start() error {
	if err := srv.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (srv *Server) Shutdown(ctx context.Context) error {
	if err := srv.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
