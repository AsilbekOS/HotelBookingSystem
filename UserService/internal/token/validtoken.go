package tokens

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	_ "github.com/joho/godotenv/autoload"
)

func ValidateToken(tokenString string) bool {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return false
	}
	return true
}
