package jerrors

import "github.com/Core-Mouse/cm-backend/internal/errors"

func JwtErrorCaster(err error) errors.PCCError {
	return NewJwtError(err)
}
