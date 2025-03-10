package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secret_Key = []byte("secret_key")

func CreateJWT(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	})

	return token.SignedString(secret_Key)
}
