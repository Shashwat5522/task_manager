package utils

func GenerateToken(userID, email, secret string, expiryHours int) (string, error) {
	// TODO: Implement JWT token generation
	return "", nil
}

func ValidateToken(tokenString, secret string) (string, error) {
	// TODO: Implement JWT token validation
	// Returns userID
	return "", nil
}
