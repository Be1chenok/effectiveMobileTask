package service

import (
	"context"
	"fmt"

	"github.com/Be1chenok/effectiveMobileTask/internal/domain"
	"github.com/Be1chenok/effectiveMobileTask/internal/repository/postgres"
	appLogger "github.com/Be1chenok/effectiveMobileTask/pkg/logger"
	"go.uber.org/zap"
)

type Person interface {
	Add(ctx context.Context, person *domain.Person) (*domain.Person, error)
	Find(ctx context.Context, searchParams *domain.PersonSearchParams) (*[]domain.Person, error)
	FindById(ctx context.Context, personID int) (*domain.Person, error)
	UpdateById(ctx context.Context, person *domain.Person) error
	DeleteById(ctx context.Context, personID int) error
}

type person struct {
	enrichment     Enrichment
	postgresPerson postgres.Person
	logger         appLogger.Logger
}

func NewPerson(logger appLogger.Logger, postgresPerson postgres.Person, enrichment Enrichment) Person {
	return &person{
		enrichment:     enrichment,
		postgresPerson: postgresPerson,
		logger:         logger.With(zap.String("component", "service-person")),
	}
}

func (p person) Add(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	age, err := p.enrichment.GetAgeByName(ctx, person.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich age: %w", err)
	}

	person.Age = age

	gender, err := p.enrichment.GetGenderByName(ctx, person.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich gender: %w", err)
	}

	person.Gender = gender

	nationality, err := p.enrichment.GetNationalityByName(ctx, person.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich nationality: %w", err)
	}

	person.Nationality = nationality

	personID, err := p.postgresPerson.Add(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("failed to add person: %w", err)
	}

	person.ID = personID

	return person, nil
}

func (p person) Find(ctx context.Context, searchParams *domain.PersonSearchParams) (*[]domain.Person, error) {
	persons, err := p.postgresPerson.Find(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to find persons: %w", err)
	}

	return persons, nil
}

func (p person) FindById(ctx context.Context, personID int) (*domain.Person, error) {
	person, err := p.postgresPerson.FindById(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to find person by id: %w", err)
	}

	return person, nil
}

func (p person) UpdateById(ctx context.Context, person *domain.Person) error {
	if err := p.postgresPerson.UpdateById(ctx, person); err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}

	return nil
}

func (p person) DeleteById(ctx context.Context, personID int) error {
	if err := p.postgresPerson.DeleteById(ctx, personID); err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	return nil
}
