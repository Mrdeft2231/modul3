package route

import (
	"context"
	"rest/internal/config"
	"rest/internal/http/controllers"
	"rest/internal/http/repository"
	"rest/internal/http/service"
)

func UserTransport(ctx context.Context) *controllers.UserController {
	db := config.DB(ctx)
	//defer db.Close()
	repo := repository.NewRepository(db)
	userService := service.NewService(repo)
	ctrl := controllers.NewUserController(userService)
	return ctrl
}
