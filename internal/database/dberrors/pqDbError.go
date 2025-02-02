package dberrors

import (
	"database/sql"

	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/lib/pq"
)

const (
	PQE_UNIQUE_VIOLATION = "Unique constraint violation"
	PQE_OTHER_ERROR = "Unknown error"
	PQE_LOGIN_INVALID_DATA_ERROR = "Email or password is invalid"
)

type DbUniqueDetails struct {
	Columns []string `json:"columns"`
}

// createPostgresError creates the error instance
func createPostgresError(db *sql.DB, err *pq.Error) errors.PCCError {
	switch err.Code {
	case "23505":
		return createUniqueError(db, err)
	default:
		return createOtherError(err)
	}
}

func createOtherError(err *pq.Error) errors.PCCError {
	return &PQDbError{
		err,
		errors.EC_DB_OTHER,
		errors.EK_DATABASE,
		nil,
		PQE_OTHER_ERROR,
	}
}

func createUniqueError(db *sql.DB, err *pq.Error) errors.PCCError {
	details, ierr := dbUniqueFailDetails(db, err)

	if ierr != nil {
		return PQDbErrorCaster(db, ierr)
	}

	return &PQDbError{
		err,
		errors.EC_DB_UNIQUE_FAIL,
		errors.EK_DATABASE,
		details,
		PQE_UNIQUE_VIOLATION,
	}
}

// dbUniqueFailDetails creates instance of DbUniqueDetails
func dbUniqueFailDetails(db *sql.DB, err *pq.Error) (*DbUniqueDetails, error) {
	columns, ierr := getConstraintColumns(db, err.Table, err.Constraint)

	if ierr != nil {
		return nil, ierr
	}

	return &DbUniqueDetails{
		columns,
	}, nil
}

// getConstraintColumns gets column name from Constraint
func getConstraintColumns(db *sql.DB, tableName, constraintName string) ([]string, error) {
    query := `
        SELECT a.attname
        FROM   pg_index i
        JOIN   pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
        WHERE  i.indrelid = $1::regclass
        AND    i.indexrelid = $2::regclass;
    `
    rows, err := db.Query(query, tableName, constraintName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var columns []string
    for rows.Next() {
        var col string
        if err := rows.Scan(&col); err != nil {
            return nil, err
        }
        columns = append(columns, col)
    }
    return columns, nil
}


// PQDbError represents the inner Postgres error from the github.com/lib/pq driver
// This error struct is private so it can contain secret data and should NEVER go beyond the server
type PQDbError struct {
	// inner contains the inner error in terms of lib/pq
	inner *pq.Error
	// code contains the error code in terms of this project
	code errors.ErrorCode
	// kind contains the error kind in terms of this project
	kind errors.ErrorKind
	details any
	message string
}

// newPQDbError creates a new instance of PQDbError struct.
// It evaluates the code and kind fields at once
func newPQDbError(db *sql.DB, inner *pq.Error) errors.PCCError {
	return createPostgresError(db, inner)
}

func NewInvalidLoginDataError() *PQDbError {
	return newPQDbErrorWOInner(errors.EC_DB_LOGIN_ERROR, errors.EK_DATABASE, PQE_LOGIN_INVALID_DATA_ERROR, nil)
}

func newPQDbErrorWOInner(code errors.ErrorCode, kind errors.ErrorKind, message string, details any) *PQDbError {
	return &PQDbError{
		nil, code, kind, details, message,
	}
}

func (e *PQDbError) Error() string {
	return e.inner.Error()
}

func (e *PQDbError) GetErrorKind() errors.ErrorKind {
	return e.kind
}

func (e *PQDbError) GetErrorCode() errors.ErrorCode {
	return e.code
}

func (e *PQDbError) IntoPublic() *errors.PublicPCCError {
	//error_message := fmt.Sprintf("Database error with code: %d", e.code)
	return errors.NewPublicPCCError(e.code, e.kind, e.details, e.message)
}
