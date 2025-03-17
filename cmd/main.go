package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rest/internal/http/route"
	"rest/internal/middleware"
)

func main() {
	r := gin.Default()

	r.Static("../public", "../public")
	//r.StaticFS("../public", http.Dir("../public"))
	r.LoadHTMLFiles("../public/index.html", "../public/ChangePassword.html", "../public/CreateUser.html", "../public/userForm.html")

	route.UserTransport(context.Background())
	userCtrl := route.UserTransport(context.Background())
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/userForm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "userForm.html", nil)
	})

	auth := r.Group("/auth", middleware.LoginMiddleware())

	auth.GET("/CreateUser", func(c *gin.Context) {
		c.HTML(http.StatusOK, "CreateUser.html", nil)
	})
	auth.GET("/ChangePassword", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ChangePassword.html", nil)
	})
	r.POST("/CreateAuth", userCtrl.CreateUser)
	r.POST("/Auth", userCtrl.UserAuth)
	r.PUT("/ChangePassword", middleware.AuthMiddleware(), userCtrl.ChangePassword)
	r.GET("/GetUsers", userCtrl.GetUsers)
	r.DELETE("/DeleteUser/:id", userCtrl.DeleteUser)
	r.PUT("/StatusPut/:id", userCtrl.UserStatus)
	r.GET("/GetUser", middleware.AuthMiddleware(), userCtrl.GetUser)

	if err := r.Run(); err != nil {
		log.Fatal("Не удалось запустить сервер", err)
	}
}
