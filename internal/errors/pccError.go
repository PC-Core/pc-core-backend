package errors

type ErrorKind string
type ErrorCode uint
type ErrorCast func(error) PCCError

// Error kinds
const (
	// The error is internal and did not occur due to the fault of the user or the frontend
	EK_INTERNAL ErrorKind = "internal"
	// Error occured while executing the database query
	EK_DATABASE ErrorKind = "database"
	// Error occured while working with the JWT token
	EK_JWT ErrorKind = "jwt"
	// Error occured in controllers
	EK_CTRLS ErrorKind = "controllers"
	// Error occured in redis controller
	EK_REDIS ErrorKind = "redis"
	// Error occured while working with JSON
	EK_JSON ErrorKind = "json"
	// Error occured while working with HTTP headers
	EK_HEADER ErrorKind = "header"
	// Error occured while parsing string to int
	EK_ATOI ErrorKind = "atoi"
	// Error occured while working with roles
	EK_ROLES ErrorKind = "roles"
)

const (
	// Internal error code used in the InternalSecretError type
	EC_INTERNAL ErrorCode = iota
	// Database error code which means UNIQUE constraint has failed
	EC_DB_UNIQUE_FAIL
	// Database error code which means the required value was not found
	EC_DB_NOT_FOUND_ERROR
	// Error code which means a wrong chars table name is provided
	EC_DB_WRONG_CHARS_TABLE_NAME
	// Error code is which means there is an error with role
	EC_DB_ROLE_ERROR
	// Unknown / unimplemented database error
	EC_DB_OTHER
	// Input error
	EC_CTRLS_INPUT_ERROR
	// Error code means that there is no user data found in the context.
	// User data should be provided inside the Authentification Middleware
	EC_CTRLS_NO_USER_DATA_ERROR
	// Error code means that the return value is nil
	EC_REDIS_NIL
	// Unknown / unimplemented redis error
	EC_REDIS_OTHER
	// Error code means that redis returned the wrong value
	EC_REDIS_WRONG_VALUE
	// Error code means that json marshal failed
	EC_JSON_MARSHAL_ERROR
	// Error code means that json unmarshal failed
	EC_JSON_UNMARSHAL_ERROR
	// Error code means that the jwt token is malformed
	EC_JWT_TOKEN_MALFORMED
	// Error code means that the jwt token is expired
	EC_JWT_TOKEN_EXPIRED
	// Error code means that the jwt token is not valid yet
	EC_JWT_TOKEN_NOT_VALID_YET
	// Error code means that there is an unknown error with the jwt token
	EC_UNKNOWN_JWT_TOKEN_ERROR
	// Error code means that the provided token has the wrong type
	EC_JWT_ERROR_TOKEN_TYPE
	// Error code means that the header is missing
	EC_HEADER_MISSING
	// Error code means that the provided string was not a uint string
	EC_ATOI_NOT_AN_UINT
	// Error code means that the provided string overflows the type uint
	EC_ATOI_RANGE_ERROR
	// Unknown ATOI error
	EC_ATOI_UNKNOWN
	// Error code means that the presented role was lower than the required one
	EC_ROLE_LOWER
)

// PCCError - minimal error interface used in the PC Core project
//
// PCCError error interface is compatible with the Go's standard `error` interface
type PCCError interface {
	// Error returns stringified error message
	//
	// This function makes `PCCError` compatible with the Go's standard `error` interface
	Error() string
	// GetErrorKind returns the kind of an error
	GetErrorKind() ErrorKind
	// GetErrorCode returns the code of an error in terms of this project. Some types based
	// on `PCCError` may have the inner error code which is usually PRIVATE
	GetErrorCode() ErrorCode
	// IntoPublic translates the PCCError into the PublicPCCError
	IntoPublic() *PublicPCCError
}

// PublicPCCError represents the safe information about the error
// It will be sent to a frontend so it should NEVER contain the secret or sensetive data
type PublicPCCError struct {
	// Code represents the error code in terms of this project
	// It should NEVER contain the inner code from the database drive or something similar
	Code ErrorCode `json:"code"`
	// Kind represents the information about the module or the step on which the error occured
	Kind ErrorKind `json:"kind"`
	// Details may provide more details about the error but should NEVER contain the secret or
	// sensetive data
	Details any `json:"details"`
	// SafeMessage contains the summary message that should be safe
	SafeMessage string `json:"message"`
}

func NewPublicPCCError(code ErrorCode, kind ErrorKind, details any, safe_message string) *PublicPCCError {
	return &PublicPCCError{
		code, kind, details, safe_message,
	}
}
