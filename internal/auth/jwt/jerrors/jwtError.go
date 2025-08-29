package jerrors

import (
	"errors"

	ierrors "github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/golang-jwt/jwt/v5"
)

const (
	JE_MALFORMED_MESSAGE        = "JWT token is malformed"
	JE_EXPIRED_MESSAGE          = "JWT token is expired"
	JE_NOT_VALID_YET_MESSAGE    = "JWT token is not valid yet"
	JE_UNKNOWN_MESSAGE          = "Unknown error with the JWT token"
	JE_WRONG_TOKEN_TYPE_MESSAGE = "The token with the wrong type provided. See the provided type in details"
)

type JwtError struct {
	Inner   error
	Code    ierrors.ErrorCode
	Kind    ierrors.ErrorKind
	Message string
}

func NewJwtError(inner error) *JwtError {
	ec, message := parseJwtErrorCodeAndMessage(inner)

	return &JwtError{
		inner,
		ec,
		ierrors.EK_JWT,
		message,
	}
}

func (j *JwtError) Error() string {
	return j.Message
}

func (j *JwtError) GetErrorKind() ierrors.ErrorKind {
	return j.Kind
}

func (j *JwtError) GetErrorCode() ierrors.ErrorCode {
	return j.Code
}

func (j *JwtError) IntoPublic() *ierrors.PublicPCCError {
	return ierrors.NewPublicPCCError(
		j.Code,
		j.Kind,
		nil,
		j.Message,
	)
}

// parseJwtErrorCodeAndMessage returns error code in terms of this project and the message
func parseJwtErrorCodeAndMessage(inner error) (ierrors.ErrorCode, string) {
	if errors.Is(inner, jwt.ErrTokenMalformed) {
		return ierrors.EC_JWT_TOKEN_MALFORMED, JE_MALFORMED_MESSAGE
	} else if errors.Is(inner, jwt.ErrTokenExpired) {
		return ierrors.EC_JWT_TOKEN_EXPIRED, JE_EXPIRED_MESSAGE
	} else if errors.Is(inner, jwt.ErrTokenNotValidYet) {
		return ierrors.EC_JWT_TOKEN_NOT_VALID_YET, JE_NOT_VALID_YET_MESSAGE
	} else {
		return ierrors.EC_UNKNOWN_JWT_TOKEN_ERROR, JE_UNKNOWN_MESSAGE
	}
}

// JwtTokenTypeError represents the token type error.
// Token type error is occured when the token with the wrong
// type is provided
type JwtTokenTypeError struct {
	// tt - provided token type
	tt any
}

func NewJwtTokenTypeError(provided_type any) *JwtTokenTypeError {
	return &JwtTokenTypeError{
		provided_type,
	}
}

func (j *JwtTokenTypeError) Error() string {
	return JE_WRONG_TOKEN_TYPE_MESSAGE
}

func (j *JwtTokenTypeError) GetErrorCode() ierrors.ErrorCode {
	return ierrors.EC_JWT_ERROR_TOKEN_TYPE
}

func (j *JwtTokenTypeError) GetErrorKind() ierrors.ErrorKind {
	return ierrors.EK_JWT
}

func (j *JwtTokenTypeError) IntoPublic() *ierrors.PublicPCCError {
	return ierrors.NewPublicPCCError(ierrors.EC_JWT_ERROR_TOKEN_TYPE, ierrors.EK_JWT, map[string]any{"provided_type": j.tt}, JE_WRONG_TOKEN_TYPE_MESSAGE)
}
