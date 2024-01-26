package domain

import "errors"

var (
	ErrNothingUpdated      = errors.New("nothing updated")
	ErrNothingWasDeleted   = errors.New("nothing was deleted")
	ErrNothingFound        = errors.New("nothing found")
	ErrAgeNotFound         = errors.New("age not found")
	ErrGenderNotFound      = errors.New("gender not found")
	ErrNationalityNotFound = errors.New("nationality not found")
)
