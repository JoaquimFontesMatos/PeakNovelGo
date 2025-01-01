package interfaces

type AuthRepositoryInterface interface {
	CheckIfTokenRevoked(refreshToken string) (bool, error)
	RevokeToken(refreshToken string) error
}