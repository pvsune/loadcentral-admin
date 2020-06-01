package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pvsune/loadcentral-admin/views"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	auth := new(views.AuthController)
	r.GET("/login", auth.Login)

	r.Run(":8000")
}
