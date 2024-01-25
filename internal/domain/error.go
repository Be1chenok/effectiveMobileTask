package domain

import "errors"

var (
	ErrNothingUpdated    = errors.New("nothing updated")
	ErrNothingWasDeleted = errors.New("nothing was deleted")
	ErrNothingFound      = errors.New("nothing found")
)
