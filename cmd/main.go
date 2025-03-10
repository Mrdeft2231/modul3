package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"rest/internal/http/route"
)

func main() {
	r := gin.Default()

	r.Static("/public", "./public")
	r.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	route.UserTransport()
	userCtrl := route.UserTransport()
	r.GET("/")
	r.POST("/createAuth", userCtrl.CreateUser)
	r.POST("/Auth", userCtrl.UserAuth)
	r.POST("/ChangePassword")

	if err := r.Run(); err != nil {
		log.Fatal("Не удалось запустить сервер", err)
	}
}
