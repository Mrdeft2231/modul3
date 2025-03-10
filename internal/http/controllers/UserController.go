package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rest/internal/http/service"
	"rest/internal/model"
	"rest/pkg/jwt"
)

type UserController struct {
	UserService service.UserServiceInterface
}

func NewUserController(userService service.UserServiceInterface) *UserController {
	return &UserController{UserService: userService}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user *model.User
	fmt.Println("сработал")

	// Пробуем связать JSON с моделью пользователя
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Пытаемся создать пользователя через сервис
	user, err := ctrl.UserService.ServiceCreateUsers(user.Login, user.Email, user.Password)
	if err != nil {
		// Если ошибка при создании пользователя, отправляем ошибку в ответ
		fmt.Printf("не удалось отправить данные сервису: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	token, err := jwt.CreateJWT(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать jwt токен"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь успешно зарегестрирован",
		"token":   token,
	})
}

func (ctrl *UserController) UserAuth(c *gin.Context) {
	var user *model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	users, err := ctrl.UserService.UserAuthService(user.Login, user.Email, user.Password)
	if err != nil {
		fmt.Errorf("не удалось отправить данные, %v", err)
	}

	if users == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "неправильный логин или пароль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": users})
}
