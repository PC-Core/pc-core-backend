package database

import (
	"database/sql"
	"fmt"

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

func (c *DbController) RegisterUser(name string, email string, password string) (*models.User, error) {
	var id int
	var role string

	passwordHash := helpers.Sha256(password)

	err := c.db.QueryRow("INSERT INTO users (Name, Email, PasswordHash) VALUES ($1, $2, $3) returning id, Role", name, email, passwordHash).Scan(&id, &role)

	if err != nil {
		return nil, err
	}

	return models.NewUser(id, name, email, models.UserRole(role), passwordHash), nil
}

func (c *DbController) LoginUser(email string, password string) (*models.User, error) {
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
		return nil, err
	}

	defer row.Close()

	if !row.Next() {
		return nil, fmt.Errorf("user not found")
	}

	if err := row.Scan(&id, &name, &remail, &role, &rpasswordHash); err != nil {
		return nil, err
	}

	return models.NewUser(id, name, remail, models.UserRole(role), rpasswordHash), nil
}

func (c *DbController) AuthentificateWithRole(email string, password string, required_role models.UserRole) error {
	user, err := c.LoginUser(email, password)

	if err != nil {
		return err
	}

	if user.Role != required_role {
		return fmt.Errorf("wrong role")
	}

	return nil
}
