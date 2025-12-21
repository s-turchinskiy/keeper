package repository

import (
	"errors"
	"fmt"
)

var ErrSecretAlreadyExist = errors.New("secret already exist")

type NotFoundError struct {
	Entity string
	UUID   string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Entity, e.UUID)
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func NewSecretNotFoundError(uuid string) error {
	return &NotFoundError{Entity: "secret", UUID: uuid}
}
