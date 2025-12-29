package models

import "time"

type Secret struct {
	ID           string
	UserID       string
	LastModified time.Time
	Hash         string
	Data         []byte
	Deleted      bool
}

type User struct {
	ID           string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}
