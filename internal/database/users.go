package database

import (
	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/internal/models"
	"github.com/PC-Core/pc-core-backend/internal/models/inputs"
)

func (c *DPostgresDbController) RegisterUser(register *inputs.RegisterUserInput) (*models.User, errors.PCCError) {
	var id int
	var role string

	passwordHash := helpers.Sha256(register.Password)

	err := c.db.QueryRow("INSERT INTO users (Name, Email, PasswordHash) VALUES ($1, $2, $3) returning id, Role", register.Name, register.Email, passwordHash).Scan(&id, &role)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewUser(id, register.Name, register.Email, models.UserRole(role), passwordHash), nil
}

func (c *DPostgresDbController) LoginUser(login *inputs.LoginUserInput) (*models.User, errors.PCCError) {
	var (
		id            int
		name          string
		remail        string
		role          string
		rpasswordHash string
	)

	passwordHash := helpers.Sha256(login.Password)

	row, err := c.db.Query("SELECT * FROM users WHERE Email = $1 AND PasswordHash = $2", login.Email, passwordHash)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer row.Close()

	if !row.Next() {
		return nil, dberrors.NewInvalidLoginDataError()
	}

	if err := row.Scan(&id, &name, &remail, &role, &rpasswordHash); err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewUser(id, name, remail, models.UserRole(role), rpasswordHash), nil
}

func (c *DPostgresDbController) GetUserByID(id int) (*models.User, errors.PCCError) {
	res := c.db.QueryRow("SELECT id, name, email, role FROM Users WHERE id = $1", id)

	if err := res.Err(); err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	var user models.User

	err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return &user, nil
}
