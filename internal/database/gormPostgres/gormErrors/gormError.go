package gormerrors

import (
	"errors"

	ierrors "github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	CART_QUANTITY_ERROR = "You're trying to add too many products"
	UNKNOWN             = "Unknown database error"
	RECORD_NOT_FOUND    = "Not found"
)

const KIND = ierrors.EK_DATABASE

type GormError struct {
	// code contains the error code in terms of this project
	code ierrors.ErrorCode
	// kind contains the error kind in terms of this project
	kind    ierrors.ErrorKind
	details any
	message string
}

var UNKNOWN_ERROR = GormError{
	code:    ierrors.EC_DB_OTHER,
	kind:    KIND,
	details: nil,
	message: UNKNOWN,
}

func parseInnerError(pgErr *pgconn.PgError) *GormError {
	if pgErr.Message == "Quantity exceeds available stock" {
		return &GormError{
			code:    ierrors.EC_DB_CART_QUANTITY_ERROR,
			kind:    KIND,
			details: nil,
			message: CART_QUANTITY_ERROR,
		}
	} else {
		return &UNKNOWN_ERROR
	}
}

func NewGormError(err error) *GormError {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return parseInnerError(pgErr)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return &GormError{
			code:    ierrors.EC_DB_NOT_FOUND_ERROR,
			kind:    KIND,
			details: nil,
			message: RECORD_NOT_FOUND,
		}
	} else {
		return &UNKNOWN_ERROR
	}
}

func (g *GormError) Error() string {
	return g.message
}

func (g *GormError) GetErrorKind() ierrors.ErrorKind {
	return g.kind
}

func (g *GormError) GetErrorCode() ierrors.ErrorCode {
	return g.code
}

func (g *GormError) IntoPublic() *ierrors.PublicPCCError {
	return ierrors.NewPublicPCCError(g.code, g.kind, g.details, g.message)
}
