package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("в middelware")
		// Получаем токен из cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			fmt.Println("токен отсувствет")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен отсутствует"})
			c.Abort()
			return
		}

		// Парсим токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем алгоритм токена
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("неверный алгоритм подписи")
				return nil, fmt.Errorf("неверный алгоритм подписи")
			}
			return []byte("secret_key"), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("недействительный токен")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		// Декодируем `user_id`
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("неверный формат токена")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		fmt.Println("токен в мидле", claims)

		userID, ok := claims["user_id"].(float64) // JWT может хранить ID как float64
		if !ok {
			fmt.Println("не удалось получить ID пользователя")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "не удалось получить ID пользователя"})
			c.Abort()
			return
		}
		fmt.Println("айди в мидле", userID)
		// Сохраняем `user_id` в контексте Gin
		c.Set("user_id", int(userID))

		c.Next() // Передаём управление следующему обработчику
	}
}
