package errors

import "fmt"

const (
	JMR_MARSHAL_MESSAGE     = "JSON marshal error"
	JMR_UNMARSHAL_ERROR     = "JSON unmarshal error"
	JMR_SAFE_MESSAGE        = "JSON error"
	JMR_SYNTAX_ERROR_FORMAT = "JSON syntax error on: %d"
)

// JsonError represents json marshal / unmarshal errors
type JsonError struct {
	Code    ErrorCode
	Kind    ErrorKind
	Message string
	Details any
}

func NewJsonMarshalError() *JsonError {
	return &JsonError{
		EC_JSON_MARSHAL_ERROR,
		EK_JSON,
		JMR_MARSHAL_MESSAGE,
		nil,
	}
}

func NewJsonUnmarshalError() *JsonError {
	return &JsonError{
		EC_JSON_UNMARSHAL_ERROR,
		EK_JSON,
		JMR_UNMARSHAL_ERROR,
		nil,
	}
}

func NewJsonSyntaxError(offset int64) *JsonError {
	return &JsonError{
		EC_JSON_SYNTAX_ERROR,
		EK_JSON,
		fmt.Sprintf(JMR_SYNTAX_ERROR_FORMAT, offset),
		map[string]int64{
			"offset": offset,
		},
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
	return NewPublicPCCError(e.Code, e.Kind, e.Details, e.Message)
}
