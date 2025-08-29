package conerrors

import "github.com/PC-Core/pc-core-backend/internal/errors"

const (
	// GCE_BIND_ERROR_MESSAGE contains the error message means that input data is invalid
	GCE_BIND_ERROR_MESSAGE   = "Error while binding: invalid data provided"
	GCE_NO_USER_DATA_MESSAGE = "Error while getting user data: no data provided"
	GCE_EMPTY_BODY           = "The body expected to be non-empty, was empty"
	GCE_UNKNOWN_BIND_ERROR   = "Unknown bind error"
)

// GinControllerError represents an error occured in controllers
type GinControllerError struct {
	Code    errors.ErrorCode
	Kind    errors.ErrorKind
	Message string
	Details any
}

func (g *GinControllerError) Error() string {
	return g.Message
}

func (g *GinControllerError) GetErrorCode() errors.ErrorCode {
	return g.Code
}

func (g *GinControllerError) GetErrorKind() errors.ErrorKind {
	return g.Kind
}

func (g *GinControllerError) IntoPublic() *errors.PublicPCCError {
	return errors.NewPublicPCCError(
		g.Code,
		g.Kind,
		g.Details,
		g.Message,
	)
}

// NewBindValidationError creates an instance of GinControllerError.
// Error represents bind input error
func NewBindValidationError(details []ValError) *GinControllerError {
	return NewGinControllersError(errors.EC_CTRLS_INPUT_ERROR, GCE_BIND_ERROR_MESSAGE, details)
}

func NewGinControllersError(code errors.ErrorCode, message string, details any) *GinControllerError {
	return &GinControllerError{
		code,
		errors.EK_CTRLS,
		message,
		details,
	}
}

// NewEmptyBodyError creates an instance of GinControllerError.
// Error represents unexpectedly empty request body
func NewEmptyBodyError() *GinControllerError {
	return NewGinControllersError(errors.EC_CTRLS_UNEXPECTED_EMPTY_BODY, GCE_EMPTY_BODY, nil)
}

// GetUserDataFromContextError creates an instance of GinControllerError.
// Error represents error while getting user data from gin.Context.
// User data should be provided in an Authentification Middleware
func GetUserDataFromContextError() *GinControllerError {
	return NewGinControllersError(errors.EC_CTRLS_NO_USER_DATA_ERROR, GCE_NO_USER_DATA_MESSAGE, nil)
}

// NewUnknownBindError creates an instance of GinControllerError.
// Error represents the unknown bind error
func NewUnknownInputError() *GinControllerError {
	return NewGinControllersError(errors.EC_CTRLS_INPUT_ERROR, GCE_UNKNOWN_BIND_ERROR, nil)
}
