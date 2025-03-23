package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
)

func (c *GormPostgresController) RegisterUser(register *inputs.RegisterUserInput) (*models.User, errors.PCCError) {
	passwordHash := helpers.Sha256(register.Password)

	user := DbUser{
		Name:         register.Name,
		Email:        register.Email,
		PasswordHash: passwordHash,
	}

	err := c.db.Create(&user).Error

	if err != nil {
		return nil, errors.NewInternalSecretError()
	}

	return user.IntoUser(), nil
}

func (c *GormPostgresController) LoginUser(login *inputs.LoginUserInput) (*models.User, errors.PCCError) {
	passwordHash := helpers.Sha256(login.Password)

	var user DbUser

	err := c.db.
		Where("email = ? AND passwordhash = ?", login.Email, passwordHash).
		First(&user).
		Error

	if err != nil {
		return nil, errors.NewInternalSecretError()
	}

	return user.IntoUser(), nil
}

func (c *GormPostgresController) GetUserByID(id int) (*models.User, errors.PCCError) {
	var user DbUser;

	err := c.db.
		Where("id = ?", id).
		First(&user).
		Error

	if err != nil {
		return nil, errors.NewInternalSecretError()
	}

	return user.IntoUser(), nil
}