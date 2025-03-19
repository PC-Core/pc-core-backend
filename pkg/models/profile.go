package models

type Profile struct {
	User *PublicUser `json:"user"`
}

func NewProfile(user *PublicUser) *Profile {
	return &Profile{
		user,
	}
}
