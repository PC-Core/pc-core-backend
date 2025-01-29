package conerrors

import "github.com/Core-Mouse/cm-backend/internal/errors"

const (
	// GCE_BIND_ERROR_MESSAGE contains the error message means that input data is invalid
	GCE_BIND_ERROR_MESSAGE   = "Error while binding: invalid data provided"
	GCE_NO_USER_DATA_MESSAGE = "Error while getting user data: no data provided"
)

// GinControllerError represents an error occured in controllers
type GinControllerError struct {
	Code    errors.ErrorCode
	Kind    errors.ErrorKind
	Message string
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
		nil,
		g.Message,
	)
}

// BindError creates an instance of GinControllerError.
// Error represents bind input error
func BindError() *GinControllerError {
	return &GinControllerError{
		errors.EC_CTRLS_INPUT_ERROR,
		errors.EK_CTRLS,
		GCE_BIND_ERROR_MESSAGE,
	}
}

// GetUserDataFromContextError creates an instance of GinControllerError.
// Error represents error while getting user data from gin.Context.
// User data should be provided in an Authentification Middleware
func GetUserDataFromContextError() *GinControllerError {
	return &GinControllerError{
		errors.EC_CTRLS_NO_USER_DATA_ERROR,
		errors.EK_CTRLS,
		GCE_NO_USER_DATA_MESSAGE,
	}
}
