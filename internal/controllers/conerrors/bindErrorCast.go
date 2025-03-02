package conerrors

import (
	"encoding/json"
	"errors"
	"io"

	inerrors "github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/go-playground/validator/v10"
)

// VFR - Validation Fault Reason
type VFR = uint64

const (
	// VFR_REQUIRED means that the required field was not present
	VFR_REQUIRED VFR = iota
	// VFR_EMAIL means that the email field is ill-formed
	VFR_EMAIL
	// VFR_UNKNOWN means that the reason is unknown.
	// Usually, this should NOT happen
	VFR_UNKNOWN
)

type ValError struct {
	Field  string `json:"field"`
	Tag    string `json:"tag"`
	Reason VFR    `json:"reason"`
}

func BindErrorCast(err error) inerrors.PCCError {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		return processValidationError(verr)
	}

	var syntax_err *json.SyntaxError
	if errors.As(err, &syntax_err) {
		return inerrors.NewJsonSyntaxError(syntax_err.Offset)
	}

	if errors.Is(err, io.EOF) {
		return NewEmptyBodyError()
	}

	return NewUnknownInputError()
}

func processValidationError(errs validator.ValidationErrors) inerrors.PCCError {
	return NewBindValidationError(processValidationErrorGetDetails(errs))
}

func processValidationErrorGetDetails(errs validator.ValidationErrors) []ValError {
	details := make([]ValError, 0)

	for _, err := range errs {
		vfr := processValidationVFR(err.Tag())

		details = append(details, ValError{err.Field(), err.Tag(), vfr})
	}

	return details
}

func processValidationVFR(tag string) VFR {
	switch tag {
	case "email":
		return VFR_EMAIL
	case "required":
		return VFR_REQUIRED
	default:
		return VFR_UNKNOWN
	}
}
