package dtos

type AuthSession struct {
	User         UserDTO `json:"user"`
	RefreshToken string  `json:"refreshToken"`
	AccessToken  string  `json:"accessToken"`
}
