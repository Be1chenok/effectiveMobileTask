package service

import (
	"github.com/Be1chenok/effectiveMobileTask/internal/config"
	"github.com/Be1chenok/effectiveMobileTask/internal/repository"
)

type Service struct {
	Person
}

func New(repo *repository.Repository, conf *config.Config) *Service {
	enrichment := NewEnrichment(conf)

	return &Service{
		Person: NewPerson(repo.PostgresPerson, enrichment),
	}
}
