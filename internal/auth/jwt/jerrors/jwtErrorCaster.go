package jerrors

import "github.com/PC-Core/pc-core-backend/internal/errors"

func JwtErrorCaster(err error) errors.PCCError {
	return NewJwtError(err)
}
