package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Auth struct{}

func (auth Auth) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}
