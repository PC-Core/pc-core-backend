package errors

const (
	JMR_MARSHAL_MESSAGE = "JSON marshal error"
	JMR_UNMARSHAL_ERROR = "JSON unmarshal error"
	JMR_SAFE_MESSAGE    = "JSON error"
)

// JsonError represents json marshal / unmarshal errors
type JsonError struct {
	Code    ErrorCode
	Kind    ErrorKind
	Message string
}

func NewJsonMarshalError() *JsonError {
	return &JsonError{
		EC_JSON_MARSHAL_ERROR,
		EK_JSON,
		JMR_MARSHAL_MESSAGE,
	}
}

func NewJsonUnmarshalError() *JsonError {
	return &JsonError{
		EC_JSON_UNMARSHAL_ERROR,
		EK_JSON,
		JMR_UNMARSHAL_ERROR,
	}
}

func (e *JsonError) Error() string {
	return e.Message
}

func (e *JsonError) GetErrorKind() ErrorKind {
	return e.Kind
}

func (e *JsonError) GetErrorCode() ErrorCode {
	return e.Code
}

func (e *JsonError) IntoPublic() *PublicPCCError {
	return NewPublicPCCError(e.Code, e.Kind, nil, e.Message)
}
