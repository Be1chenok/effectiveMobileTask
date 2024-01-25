package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Be1chenok/effectiveMobileTask/internal/domain"
)

type Person interface {
	Find(ctx context.Context, searchParams domain.PersonSearchParams) ([]domain.Person, error)
	FindById(ctx context.Context, personId int) (domain.Person, error)
	Add(ctx context.Context, person domain.Person) (int, error)
	DeleteById(ctx context.Context, personId int) error
	UpdateById(ctx context.Context, person domain.Person) error
}

type person struct {
	db *sql.DB
}

func NewPersonRepo(db *sql.DB) Person {
	return &person{
		db: db,
	}
}

func (p person) Find(ctx context.Context, searchParams domain.PersonSearchParams) ([]domain.Person, error) {
	rows, err := p.db.QueryContext(
		ctx,
		`SELECT id, name, surname, patronymic, age, gender, nationality
		FROM persons
		WHERE gender LIKE $1||'%'
		AND nationality LIKE $2||'%'
		ORDER BY id ASC
		OFFSET $3
		LIMIT $4`,
		searchParams.Gender,
		searchParams.Nationality,
		searchParams.Offset,
		searchParams.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SELECT query: %w", err)
	}

	var persons []domain.Person

	for rows.Next() {
		var person domain.Person
		if err := rows.Scan(
			&person.Id,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		); err != nil {
			return nil, domain.ErrNothingFound
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (p person) FindById(ctx context.Context, personId int) (domain.Person, error) {
	var person domain.Person

	if err := p.db.QueryRowContext(
		ctx,
		`SELECT id, name, surname, patronymic, age, gender, nationality
		FROM persons
		WHERE id = $1`,
		personId,
	).Scan(
		&person.Id,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
	); err != nil {
		return domain.Person{}, domain.ErrNothingFound
	}

	return person, nil
}

func (p person) Add(ctx context.Context, person domain.Person) (int, error) {
	var personId int
	if err := p.db.QueryRowContext(
		ctx,
		`INSERT INTO persons (name, surname, patronymic, age, gender, nationality) values ($1, $2, $3, $4, $5, $6) RETURNING id`,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	).Scan(&personId); err != nil {
		return 0, fmt.Errorf("failed to scan row: %w", err)
	}

	return personId, nil
}

func (p person) DeleteById(ctx context.Context, personId int) error {
	result, err := p.db.ExecContext(
		ctx,
		`DELETE FROM persons WHERE id=$1`,
		personId,
	)
	if err != nil {
		return fmt.Errorf("failed to execute DELETE query: %w", err)
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of deleted rows: %w", err)
	}
	if deletedRows == 0 {
		return domain.ErrNothingWasDeleted
	}

	return nil
}

func (p person) UpdateById(ctx context.Context, person domain.Person) error {
	result, err := p.db.ExecContext(
		ctx,
		`UPDATE persons
		SET name = $1,
		surname = $2,
		patronymic = $3,
		age = $4,
		gender = $5,
		nationality = $6
		WHERE id = $7`,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		person.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to execute UPDATE query: %w", err)
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of updated rows: %w", err)
	}
	if updatedRows == 0 {
		return domain.ErrNothingUpdated
	}

	return nil
}
