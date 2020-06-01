package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pvsune/loadcentral-admin/config"
	authMiddleware "github.com/pvsune/loadcentral-admin/middlewares/auth"
	"github.com/pvsune/loadcentral-admin/views"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	config.Init()
	authMiddleware.Init()

	authView := new(views.Auth)
	adminView := new(views.Admin)

	r.GET("/login", authView.Login)
	r.POST("/login", authMiddleware.Login)
	r.GET("/logout", authMiddleware.Logout)

	authorized := r.Group("/", authMiddleware.MiddlewareFunc())
	authorized.GET("/", adminView.Index)

	r.Run(":8000")
}
