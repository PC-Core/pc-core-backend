package models

import "fmt"

type AuthData struct {
	pub fmt.Stringer
	priv fmt.Stringer
}

func NewAuthData(pub fmt.Stringer, priv fmt.Stringer) *AuthData {
	return &AuthData {
		pub, priv,
	}
}

func (a *AuthData) GetPublic() fmt.Stringer {
	return a.pub
}

func (a *AuthData) GetPrivate() fmt.Stringer {
	return a.priv
}