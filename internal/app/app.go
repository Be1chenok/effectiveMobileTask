package app

import (
	"context"
	"log"
	"time"

	"github.com/Be1chenok/effectiveMobileTask/internal/config"
	"github.com/Be1chenok/effectiveMobileTask/internal/repository/postgres"
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	postgres, err := postgres.New(ctx, conf)
	if err != nil {
		appLog.Fatalf("failed to connect database: %v", err)
	}
	cancel()
}
