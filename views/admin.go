package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Admin struct{}

func (admin Admin) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func (admin Admin) SendLoad(c *gin.Context) {
	type sendLoadForm struct {
		PhoneNumber []string `form:"phone_number[]" binding:"required"`
		Pcode       []string `form:"pcode[]" binding:"required"`
	}

	var form sendLoadForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", nil)
}
