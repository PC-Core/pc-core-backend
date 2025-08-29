package merrors

import (
	"fmt"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
)

const (
	RE_ERROR_FORMAT = "Role error with code: %d"
	RE_SAFE_MESSAGE = "Role error"
)

type RoleError struct {
	Code          errors.ErrorCode
	Kind          errors.ErrorKind
	RequiredRole  models.UserRole
	PresentedRole models.UserRole
}

func NewLowerRoleError(required models.UserRole, presented models.UserRole) *RoleError {
	return &RoleError{
		errors.EC_ROLE_LOWER,
		errors.EK_ROLES,
		required,
		presented,
	}
}

func (r *RoleError) Error() string {
	return fmt.Sprintf(RE_ERROR_FORMAT, r.Code)
}

func (r *RoleError) GetErrorCode() errors.ErrorCode {
	return r.Code
}

func (r *RoleError) GetErrorKind() errors.ErrorKind {
	return r.Kind
}

func (r *RoleError) IntoPublic() *errors.PublicPCCError {
	return errors.NewPublicPCCError(
		r.Code,
		r.Kind,
		map[string]models.UserRole{
			"required":  r.RequiredRole,
			"presented": r.PresentedRole,
		},
		RE_SAFE_MESSAGE,
	)
}
