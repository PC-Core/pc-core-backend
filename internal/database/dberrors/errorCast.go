package dberrors

import (
	"database/sql"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/lib/pq"
)

// pqDbErrorCaster casts the database error type to the errors.PCCError type.
// This implementation is made for the Postgres database and github.com/lib/pq driver
// In normal situation it always takes *pq.Error and returns errors.InternalDatabaseError
// If the passed parameter was not the *pq.Error, it returns errors.InternalSecretError
func PQDbErrorCaster(db *sql.DB, err error) errors.PCCError {
	if err == nil {
		return nil
	}

	inner, ok := err.(*pq.Error)

	if !ok {
		return errors.NewInternalSecretError()
	}

	return newPQDbError(db, inner)
}
