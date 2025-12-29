package errorsutils

import (
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"runtime"
)

func WrapError(err error) error {

	if err == nil {
		return nil
	}

	_, filename, line, _ := runtime.Caller(1)
	return fmt.Errorf("[error] %s %d: %w", filename, line, err)
}

func IsConnectionError(err error) bool {

	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	if pgErr == nil {
		return false
	}

	return pgerrcode.IsConnectionException(pgErr.Code)

}
