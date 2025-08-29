package models

type UserRole string

const (
	Default UserRole = "Default"

	Admin = "Admin"
)

type User struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Role       UserRole `json:"user_role"`
	passwdHash string
}

func NewUser(id int, name string, email string, role UserRole, passwdHash string) *User {
	return &User{
		id, name, email, role, passwdHash,
	}
}
