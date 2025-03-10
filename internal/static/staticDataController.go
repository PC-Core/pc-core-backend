package static

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
)

type StaticDataController interface {
	UploadFiles(files []StaticFile) ([]string, errors.PCCError)
}
