package dbparse

import (
	"errors"
	"strings"
)

type Database struct {
	TypeDB   string
	Host     string
	Port     string
	DBName   string
	Login    string
	Password string
}

func ParsedConnectionString(s string) (*Database, error) {

	s = strings.Replace(s, "://", " ", 1)
	s = strings.Replace(s, ":", " ", 1)
	s = strings.Replace(s, "@", " ", 1)
	s = strings.Replace(s, "/", " ", 1)
	s = strings.Replace(s, "?", " ", 1)

	hp := strings.Split(s, " ")
	if len(hp) < 5 && len(hp) > 6 {
		return nil, errors.New("incorrect format database-dsn")
	}

	d := &Database{
		TypeDB:   hp[0],
		Login:    hp[1],
		Password: hp[2],
		DBName:   hp[4],
	}

	hostPort := strings.Split(hp[3], ":")
	if len(hp) == 0 {
		return nil, errors.New("incorrect format database-dsn")
	}

	d.Host = hostPort[0]
	if len(hostPort) == 2 {
		d.Port = hostPort[1]
	}

	return d, nil
}
