package views

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pvsune/loadcentral-admin/config"
	"io"
	"io/ioutil"
	"log"
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
	for i, phone_number := range form.PhoneNumber {
		doSendLoad(phone_number, form.Pcode[i])
	}
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func md5Hex(data []string) string {
	h := md5.New()
	for _, d := range data {
		io.WriteString(h, d)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func doSendLoad(phone_number string, pcode string) {
	log.Printf("Sending load \"%s\" to \"%s\"", pcode, phone_number)
	conf := config.GetConfig()

	rrn := uuid.New().String()
	auth := md5Hex([]string{
		md5Hex([]string{rrn}),
		md5Hex([]string{
			conf.GetString("LC_USERNAME"), conf.GetString("LC_PASSWORD"),
		}),
	})
	requestURL := fmt.Sprintf(
		"%s?uid=%s&auth=%s&rrn=%s&pcode=%s&to=%s",
		conf.GetString("LC_BASEURL"),
		conf.GetString("LC_USERNAME"),
		auth, rrn, pcode, phone_number,
	)
	log.Printf("Trying %s...", requestURL)

	resp, err := http.Get(requestURL)
	if err != nil {
		panic(fmt.Sprintf("LoadCentral Error: %s", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("LoadCentral Error: %s", err))
	}
	fmt.Printf("%s", body)
}
