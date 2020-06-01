package auth

import (
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pvsune/loadcentral-admin/config"
	"net/http"
	"time"
)

var auth *jwt.GinJWTMiddleware

type formLogin struct {
	Username string `form:"username" bind:"required"`
	Password string `form:"password" bind:"required"`
}

func Init() {
	conf := config.GetConfig()

	var err error
	auth, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:          conf.GetString("AUTH_REALM"),
		Key:            []byte(conf.GetString("AUTH_KEY")),
		Timeout:        time.Duration(conf.GetInt("AUTH_TIMEOUT")) * time.Second,
		SendCookie:     conf.GetBool("AUTH_SENDCOOKIE"),
		SecureCookie:   conf.GetBool("AUTH_SECURECOOKIE"),
		CookieHTTPOnly: conf.GetBool("AUTH_COOKIEHTTPONLY"),
		CookieDomain:   conf.GetString("AUTH_COOKIEDOMAIN"),
		TokenLookup:    conf.GetString("AUTH_TOKENLOOKUP"),
		Authenticator:  authenticator,
		LoginResponse:  loginResponse,
		Unauthorized:   unauthorized,
		LogoutResponse: logoutResponse,
	})
	if err != nil {
		panic(fmt.Sprintf("JWT Auth Error: %s\n", err.Error()))
	}
}

func MiddlewareFunc() gin.HandlerFunc {
	return auth.MiddlewareFunc()
}

func Login(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode) // Avoid CSRF!
	auth.LoginHandler(c)
}

func Logout(c *gin.Context) {
	auth.LogoutHandler(c)
}

func authenticator(c *gin.Context) (interface{}, error) {
	var data formLogin
	if err := c.ShouldBind(&data); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	if !(data.Username == "admin" && data.Password == "icanseeyou") {
		return nil, jwt.ErrFailedAuthentication
	}
	return nil, nil
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.Redirect(http.StatusMovedPermanently, "/")
}

func logoutResponse(c *gin.Context, code int) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func unauthorized(c *gin.Context, code int, message string) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}
