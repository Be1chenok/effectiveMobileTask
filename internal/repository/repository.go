package repository

import (
	"database/sql"

	"github.com/Be1chenok/effectiveMobileTask/internal/repository/postgres"
	appLogger "github.com/Be1chenok/effectiveMobileTask/pkg/logger"
)

type Repository struct {
	PostgresPerson postgres.Person
}

func New(logger appLogger.Logger, db *sql.DB) *Repository {
	return &Repository{
		PostgresPerson: postgres.NewPersonRepo(logger, db),
	}
}
