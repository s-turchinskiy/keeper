package models

import (
	"time"
)

type RemoteSecret struct {
	Name         string
	LastModified time.Time
	Hash         string
	Data         []byte
}
