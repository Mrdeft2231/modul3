package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rest/internal/http/service"
	"rest/internal/model"
	"rest/pkg/Cookie"
	"rest/pkg/auth"
	"strconv"
)

type UserController struct {
	UserService service.UserServiceInterface
}

func NewUserController(userService service.UserServiceInterface) *UserController {
	return &UserController{UserService: userService}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user *model.User

	// Пробуем связать JSON с моделью пользователя
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	fmt.Println("Контроллер", user.Login, user.Email, user.Password, user.Role, user.Status)

	// Пытаемся создать пользователя через сервис
	user, err := ctrl.UserService.ServiceCreateUsers(context.Background(), user.Login, user.Email, user.Password, user.Role, user.Status)
	if errors.Is(err, auth.ErrUniqName) {
		fmt.Println("Попал 2 раз")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Такое имя уже существует"})
		return
	}
	if err != nil {
		// Если ошибка при создании пользователя, отправляем ошибку в ответ
		fmt.Printf("не удалось отправить данные сервису: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}
}

func (ctrl *UserController) UserAuth(c *gin.Context) {
	var user *model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	users, err := ctrl.UserService.UserAuthService(context.Background(), user.Login, user.Email, user.Password)
	if errors.Is(err, auth.ErrBlockStatus) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "вас заблокировали, вход невозможен"})
		return
	}
	if errors.Is(err, auth.ErrBlockDate) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователю больше 30 дней, пользователь заблокирован"})
		return
	}
	if errors.Is(err, auth.ErrBlock) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "вы достигли маскимальное колчество попыток входа, пользователь заблокирован"})
		return
	}
	if errors.Is(err, auth.ErrPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный пароль, повторите попытку"})
		return
	}

	if err != nil {
		log.Printf("не удалось отправить данные, %v", err)
	}

	if users == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "неправильный логин или пароль"})
		return
	}
	fmt.Println("в контроллере", user.PasswordChanged)

	fmt.Println("id в контроллере", users.Id)
	cookie := Cookie.CreateJwtToken(users.Id)
	c.SetCookie("token", cookie, 3600, "/", "", false, false)

	if users.Role == "Администратор" {
		c.JSON(http.StatusOK, gin.H{
			"redirect": "auth/CreateUser",
		})
		return
	}

	if users.PasswordChanged == false {
		fmt.Println("Сработало перенаправление")
		c.JSON(http.StatusOK, gin.H{
			"redirect": "auth/ChangePassword",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"redirect": "/userForm",
	})
	return
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.UserService.GetUsers()
	if err != nil {
		log.Printf("не удалось достучаться до сервиса %v", err)
	}

	if users == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ничего не нашли"})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}
	err = ctrl.UserService.DeleteUser(context.Background(), userID)
	if err != nil {
		log.Printf("Не удалось удалить данные %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удалён"})
}

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	fmt.Println("в контроллере")
	var user *model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Получаем `user_id` из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не удалось получить ID пользователя"})
		log.Printf("не удалось получить данные")
		return
	}

	// Приводим userID к int
	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ID пользователя"})
		return
	}

	err := ctrl.UserService.ChangePassword(context.Background(), id, user.Password, user.OldPassword)
	if errors.Is(err, auth.ErrIncorrectPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный старый пароль"})
		return
	}
	if errors.Is(err, auth.ErrSamePassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Новый пароль не должен совпадать со старым"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirect": "/userForm"})
}

func (ctrl *UserController) UserStatus(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}
	err = ctrl.UserService.BlocketUser(context.Background(), userID)
	if err != nil {
		log.Printf("Не удалось изменить статус %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус изменён"})
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не удалось получить ID пользователя"})
		log.Printf("не удалось получить данные")
		return
	}

	// Приводим userID к int
	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ID пользователя"})
		return
	}

	user, err := ctrl.UserService.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	login := user.Login

	c.JSON(http.StatusOK, gin.H{"login": login})
}
