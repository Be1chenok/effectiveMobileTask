package handler

import "errors"

var (
	ErrSomethingWentWrong = errors.New("oops, Something went wrong")
	ErrInvalidInput       = errors.New("invalid input")
	ErrIdIsNotNumber      = errors.New("id is not number")
	ErrUncorrectedUrl     = errors.New("uncorrected url")
)
