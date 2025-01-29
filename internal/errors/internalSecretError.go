package errors

// ISE_MESSAGE stands for Internal Secret Error Message.
// ISE_MESSAGE is a safe message of an error
const ISE_MESSAGE = "An internal error occured. Details are secret"

type InternalSecretError struct {
	kind    ErrorKind
	code    ErrorCode
	message string
}

func NewInternalSecretError() *InternalSecretError {
	return &InternalSecretError{
		kind:    EK_INTERNAL,
		code:    EC_INTERNAL,
		message: ISE_MESSAGE,
	}
}

func (e *InternalSecretError) Error() string {
	return "An internal error occured. Details are secret"
}

func (e *InternalSecretError) GetErrorKind() ErrorKind {
	return e.kind
}

func (e *InternalSecretError) GetErrorCode() ErrorCode {
	return e.code
}

func (e *InternalSecretError) IntoPublic() *PublicPCCError {
	return NewPublicPCCError(e.code, e.kind, nil, e.message)
}
