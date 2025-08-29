package database

import (
	"database/sql"

	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

type DbController struct {
	db *sql.DB
}

func NewDbController(driver string, conn string) (*DbController, error) {
	db, err := sql.Open(driver, conn)

	if err != nil {
		return nil, err
	}

	return &DbController{
		db,
	}, nil
}

func (c *DbController) RegisterUser(name string, email string, role models.UserRole, password string) (*models.User, error) {
	var id int

	passwordHash := helpers.Sha256(password)

	err := c.db.QueryRow("INSERT INTO users (Name, Email, Role, PasswordHash) VALUES ($1, $2, $3, $4) returning id", name, email, role, passwordHash).Scan(&id)

	if err != nil {
		return nil, err
	}

	return models.NewUser(id, name, email, role, passwordHash), nil
}
