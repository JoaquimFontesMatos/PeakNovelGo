package interfaces

type AuthRepositoryInterface interface {
	CheckIfTokenRevoked(refreshToken string) bool
	RevokeToken(refreshToken string) error
}