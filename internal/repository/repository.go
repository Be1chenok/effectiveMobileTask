package repository

import (
	"database/sql"

	"github.com/Be1chenok/effectiveMobileTask/internal/repository/postgres"
)

type Repository struct {
	PostgresPerson postgres.Person
}

func New(db *sql.DB) *Repository {
	return &Repository{
		PostgresPerson: postgres.NewPersonRepo(db),
	}
}
