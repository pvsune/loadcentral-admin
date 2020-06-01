package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Admin struct{}

func (admin Admin) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}
