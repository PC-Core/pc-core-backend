package serrors

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/minio/minio-go/v7"
)

func MinIOCast(err error) errors.PCCError {
	return NewMinIOError(minio.ToErrorResponse(err))
}
