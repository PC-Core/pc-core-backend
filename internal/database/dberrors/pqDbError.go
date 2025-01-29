package dberrors

import (
	"fmt"

	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/lib/pq"
)

// parseDbErrorCode parses lib/pq ErrorCode to the errors.ErrorCode
func parseDbErrorCode(code pq.ErrorCode) errors.ErrorCode {
	switch code {
	case "23505":
		return errors.EC_DB_UNIQUE_FAIL
	default:
		return errors.EC_DB_OTHER
	}
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
}

// NewPQDbError creates a new instance of PQDbError struct.
// It evaluates the code and kind fields at once
func NewPQDbError(inner *pq.Error) *PQDbError {
	code := parseDbErrorCode(inner.Code)

	return &PQDbError{
		inner: inner,
		code:  code,
		kind:  errors.EK_DATABASE,
	}
}

func NewPQDbErrorWOInner(code errors.ErrorCode, kind errors.ErrorKind) *PQDbError {
	return &PQDbError{
		nil, code, kind,
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
	error_message := fmt.Sprintf("Database error with code: %d", e.code)
	return errors.NewPublicPCCError(e.code, e.kind, nil, error_message)
}
