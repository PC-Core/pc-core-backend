package errors

import (
	"fmt"
	"strconv"
)

const (
	AE_MESSAGE_FORMAT = "Uint cast error of code %d"
	AE_SAFE_MESSAGE   = "Cast error"
)

type AtoiError struct {
	Code ErrorCode
	Kind ErrorKind
}

func NewAtoiError(err error) *AtoiError {
	ec := parseAtoiErrorToCode(err)

	return &AtoiError{
		ec,
		EK_ATOI,
	}
}

func (a *AtoiError) Error() string {
	return fmt.Sprintf(AE_MESSAGE_FORMAT, a.Code)
}

func (a *AtoiError) GetErrorCode() ErrorCode {
	return a.Code
}

func (a *AtoiError) GetErrorKind() ErrorKind {
	return a.Kind
}

func (a *AtoiError) IntoPublic() *PublicPCCError {
	return NewPublicPCCError(
		a.Code,
		a.Kind,
		nil,
		AE_SAFE_MESSAGE,
	)
}

func parseAtoiErrorToCode(err error) ErrorCode {
	switch err {
	case strconv.ErrSyntax:
		return EC_ATOI_NOT_AN_UINT
	case strconv.ErrRange:
		return EC_ATOI_RANGE_ERROR
	default:
		return EC_ATOI_UNKNOWN
	}
}
