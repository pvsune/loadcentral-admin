package views

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pvsune/loadcentral-admin/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type LoadCentralResponse struct {
	PhoneNumber string
	Pcode       string
	Resp        string `xml:"RESP"`
	TID         string `xml:"TID"`
	Bal         string `xml:"BAL"`
	Err         string `xml:"ERR"`
}
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

	ch := make(chan *LoadCentralResponse)
	defer close(ch)
	for i, phone_number := range form.PhoneNumber {
		go doSendLoad(phone_number, form.Pcode[i], ch)
	}

	var res []string
	for range form.PhoneNumber {
		c := <-ch
		res = append(res, fmt.Sprintf("%+v", c))
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"error": res})
}

func md5Hex(data []string) string {
	h := md5.New()
	for _, d := range data {
		io.WriteString(h, d)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func doSendLoad(phone_number string, pcode string, ch chan *LoadCentralResponse) {
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

	res := LoadCentralResponse{PhoneNumber: phone_number, Pcode: pcode}
	err = xml.Unmarshal([]byte(fmt.Sprintf("<root>%s</root>", body)), &res)
	if err != nil {
		panic(fmt.Sprintf("LoadCentral Error: %s", err))
	}
	ch <- &res
}
