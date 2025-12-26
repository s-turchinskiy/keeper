package mongo_generic_repository

import (
	"errors"
	"fmt"
)

var (
	ErrIdInvalid = errors.New("invalid ID")
)

type NotFoundError struct {
	Entity string
	ID     string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Entity, e.ID)
}

func IsNotFound(err error) bool {
	var notFoundError *NotFoundError
	ok := errors.As(err, &notFoundError)
	return ok
}

func NewEntityNotFoundError(entity, id string) error {
	return &NotFoundError{Entity: entity, ID: id}
}
