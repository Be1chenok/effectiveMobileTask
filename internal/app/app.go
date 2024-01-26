package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Be1chenok/effectiveMobileTask/internal/config"
	appHandler "github.com/Be1chenok/effectiveMobileTask/internal/delivery/http/handler"
	appServer "github.com/Be1chenok/effectiveMobileTask/internal/delivery/http/server"
	appRepository "github.com/Be1chenok/effectiveMobileTask/internal/repository"
	"github.com/Be1chenok/effectiveMobileTask/internal/repository/postgres"
	appService "github.com/Be1chenok/effectiveMobileTask/internal/service"
	appLogger "github.com/Be1chenok/effectiveMobileTask/pkg/logger"
	"go.uber.org/zap"
)

func Run() {
	logger, err := appLogger.New()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatalf("failed to sync logger: %v", err)
		}
	}()
	appLog := logger.With(zap.String("component", "app"))

	conf, err := config.Init()
	if err != nil {
		appLog.Fatalf("failed to init config: %v", err)
	}
	appLog.Info("config initialized")

	pgCtx, pgCancel := context.WithTimeout(context.Background(), 20*time.Second)
	postgres, err := postgres.New(pgCtx, conf)
	if err != nil {
		appLog.Fatalf("failed to connect database: %v", err)
	}
	pgCancel()
	appLog.Infof("database connected")

	repository := appRepository.New(postgres)
	service := appService.New(repository, conf)
	handler := appHandler.New(conf, logger, service)
	server := appServer.New(conf, handler.InitRoutes())

	go func() {
		if err := server.Start(); err != nil {
			appLog.Fatalf("failed to start server: %v", err)
		}
	}()

	appLog.Infof("server running on port: %d", conf.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	appLog.Info("shutting down")

	srvCtx, srvCancel := context.WithTimeout(context.Background(), 20*time.Second)
	if err := server.Shutdown(srvCtx); err != nil {
		appLog.Fatalf("failed to shutdown server: %v", err)
	}
	srvCancel()

	if err := postgres.Close(); err != nil {
		appLog.Fatalf("failed to close postgres: %v", err)
	}
}
