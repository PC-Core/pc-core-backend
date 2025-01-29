package database

import (
	"github.com/Core-Mouse/cm-backend/internal/database/dberrors"
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

func (c *DPostgresDbController) RegisterUser(name string, email string, password string) (*models.User, errors.PCCError) {
	var id int
	var role string

	passwordHash := helpers.Sha256(password)

	err := c.db.QueryRow("INSERT INTO users (Name, Email, PasswordHash) VALUES ($1, $2, $3) returning id, Role", name, email, passwordHash).Scan(&id, &role)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(err)
	}

	return models.NewUser(id, name, email, models.UserRole(role), passwordHash), nil
}

func (c *DPostgresDbController) LoginUser(email string, password string) (*models.User, errors.PCCError) {
	var (
		id            int
		name          string
		remail        string
		role          string
		rpasswordHash string
	)

	passwordHash := helpers.Sha256(password)

	row, err := c.db.Query("SELECT * FROM users WHERE Email = $1 AND PasswordHash = $2", email, passwordHash)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(err)
	}

	defer row.Close()

	if !row.Next() {
		return nil, dberrors.NewPQDbErrorWOInner(errors.EC_DB_NOT_FOUND_ERROR, errors.EK_DATABASE)
	}

	if err := row.Scan(&id, &name, &remail, &role, &rpasswordHash); err != nil {
		return nil, dberrors.PQDbErrorCaster(err)
	}

	return models.NewUser(id, name, remail, models.UserRole(role), rpasswordHash), nil
}

func (c *DPostgresDbController) AuthentificateWithRole(email string, password string, required_role models.UserRole) errors.PCCError {
	user, err := c.LoginUser(email, password)

	if err != nil {
		return err
	}

	if user.Role != required_role {
		return dberrors.NewPQDbErrorWOInner(errors.EC_DB_ROLE_ERROR, errors.EK_DATABASE)
	}

	return nil
}

func (c *DPostgresDbController) GetUserByID(id int) (*models.User, errors.PCCError) {
	res := c.db.QueryRow("SELECT id, name, email, role FROM Users WHERE id = $1", id)

	if err := res.Err(); err != nil {
		return nil, dberrors.PQDbErrorCaster(err)
	}

	var user models.User

	err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Role)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(err)
	}

	return &user, nil
}
