package conerrors

import (
	"fmt"

	"github.com/Core-Mouse/cm-backend/internal/errors"
)

const MISSING_TOKEN_FORMAT = "Missing cookie: %s"

type MissingTokenCookieError struct {
	cookieName string
}

func NewMissingTokenCookieError(cookieName string) *MissingTokenCookieError {
	return &MissingTokenCookieError{
		cookieName,
	}
}

func (e *MissingTokenCookieError) Error() string {
	return fmt.Sprintf(MISSING_TOKEN_FORMAT, e.cookieName)
}

func (e *MissingTokenCookieError) GetErrorKind() errors.ErrorKind {
	return errors.EK_COOKIE
}

func (e *MissingTokenCookieError) GetErrorCode() errors.ErrorCode {
	return errors.EC_COOKIE_MISSING
}

func (e *MissingTokenCookieError) IntoPublic() *errors.PublicPCCError {
	return errors.NewPublicPCCError(e.GetErrorCode(), e.GetErrorKind(), map[string]string{
		"missing_cookie": e.cookieName,
	}, MISSING_TOKEN_FORMAT)
}
