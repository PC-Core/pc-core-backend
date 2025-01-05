package models

type PublicUser struct {
	ID    int
	Name  string
	Email string
	Role  UserRole
}

func NewPublicUser(id int, name string, email string, role UserRole) *PublicUser {
	return &PublicUser{
		id, name, email, role,
	}
}
