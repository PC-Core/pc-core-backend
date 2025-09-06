package gormerrors

import (
	ierrors "github.com/PC-Core/pc-core-backend/internal/errors"
)

func GormErrorCast(err error) ierrors.PCCError {
	return NewGormError(err)
}

func GormErrorCastUserOwn(err error) ierrors.PCCError {
	return NewGormErrorUserOwn(err)
}
