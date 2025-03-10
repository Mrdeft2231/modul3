package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rest/internal/http/route"
)

func main() {
	r := gin.Default()

	r.Static("../public", "../public")
	r.LoadHTMLFiles("../public/index.html")

	route.UserTransport()
	userCtrl := route.UserTransport()
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/createAuth", userCtrl.CreateUser)
	r.POST("/Auth", userCtrl.UserAuth)
	r.POST("/ChangePassword")

	if err := r.Run(); err != nil {
		log.Fatal("Не удалось запустить сервер", err)
	}
}
