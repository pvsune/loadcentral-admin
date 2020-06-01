package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct{}

func (ac AuthController) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}
