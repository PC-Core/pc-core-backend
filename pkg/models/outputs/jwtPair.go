package outputs

type JWTPair struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

func NewJWTPair(access, refresh string) *JWTPair {
	return &JWTPair{
		access, refresh,
	}
}
