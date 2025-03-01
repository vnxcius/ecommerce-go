package models

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrNoRecords          = errors.New("models: no records found")
)
