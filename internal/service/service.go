package service

import (
	"github.com/Be1chenok/effectiveMobileTask/internal/config"
	"github.com/Be1chenok/effectiveMobileTask/internal/repository"
	appLogger "github.com/Be1chenok/effectiveMobileTask/pkg/logger"
)

type Service struct {
	Person
}

func New(conf *config.Config, logger appLogger.Logger, repo *repository.Repository) *Service {
	enrichment := NewEnrichment(conf, logger)

	return &Service{
		Person: NewPerson(logger, repo.PostgresPerson, enrichment),
	}
}
