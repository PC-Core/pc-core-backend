package serrors

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/minio/minio-go/v7"
)

type MinIOError struct {
	Code errors.ErrorCode
	Details any
	Message string
}

func getPccCode(mess minio.ErrorResponse) (any, errors.ErrorCode) {
	switch mess.Code {
	case "NoSuchKey":
		return mess.BucketName, errors.EC_MINIO_NOT_FOUNT
	default:
		return nil, errors.EC_UNKNOWN_MINIO_ERROR
	}
}

func NewMinIOError(mess minio.ErrorResponse) *MinIOError {
	details, code := getPccCode(mess)
	return &MinIOError{
		code, details, mess.Message,
	}
}

func (e *MinIOError) Error() string {
	return e.Message
}
func (e *MinIOError) GetErrorKind() errors.ErrorKind {
	return errors.EK_MINIO
}

func (e *MinIOError) GetErrorCode() errors.ErrorCode {
	return e.Code
}

func (e *MinIOError) IntoPublic() *errors.PublicPCCError {
	return errors.NewPublicPCCError(e.Code, errors.EK_MINIO, e.Details, e.Message)
}