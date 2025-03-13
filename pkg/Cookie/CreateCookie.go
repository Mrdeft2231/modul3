package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rest/pkg/jwt"
)

func CreateToken(id int, c *gin.Context) {
	token, err := jwt.CreateJWT(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать jwt токен"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь успешно зарегестрирован",
		"token":   token,
	})
}
