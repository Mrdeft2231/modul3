package Cookie

import (
	"log"
	"rest/pkg/jwt"
)

func CreateJwtToken(id int) string {
	token, err := jwt.CreateJWT(id)
	if err != nil {
		log.Println("Не удалось создать куки", err)
	}
	return token
}
