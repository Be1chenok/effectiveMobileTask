package app

import (
	"context"
	"log"
	"time"

	"github.com/Be1chenok/effectiveMobileTask/internal/config"
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	postgres, err := postgres.New(ctx, conf)
	if err != nil {
		appLog.Fatalf("failed to connect database: %v", err)
	}
	cancel()
	appLog.Infof("database connected")

	repository := appRepository.New(postgres)
	service := appService.New(repository, conf)
}
