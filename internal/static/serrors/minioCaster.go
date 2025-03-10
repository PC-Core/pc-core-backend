package serrors

import "github.com/PC-Core/pc-core-backend/internal/errors"

func MinIOCast(err error) errors.PCCError {
	return errors.NewInternalSecretError()
}
