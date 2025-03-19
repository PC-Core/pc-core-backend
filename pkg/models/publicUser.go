package models

type PublicUser struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

func NewPublicUser(id int, name string, email string, role UserRole) *PublicUser {
	return &PublicUser{
		id, name, email, role,
	}
}

func NewPublicUserFromUser(user *User) *PublicUser {
	return NewPublicUser(user.ID, user.Name, user.Email, user.Role)
}
