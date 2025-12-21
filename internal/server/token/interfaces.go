package token

type TokenManager interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (string, error)
}
