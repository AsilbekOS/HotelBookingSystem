package tokens

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/joho/godotenv/autoload"
)

func CreateToken(userid int) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userid,
			"exp":    time.Now().Add(time.Second * 60).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
