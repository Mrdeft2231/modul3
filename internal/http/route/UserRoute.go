package route

import (
	"rest/internal/config"
	"rest/internal/http/controllers"
	"rest/internal/http/repository"
	"rest/internal/http/service"
)

func UserTransport() *controllers.UserController {
	db := config.DB()
	repo := repository.NewRepository(db)
	userService := service.NewService(repo)
	ctrl := controllers.NewUserController(userService)
	return ctrl
}
