package routes

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/auth"
	"github.com/zinclabs/zinc/pkg/config"
)

var JWTMiddleWare *jwt.GinJWTMiddleware

var handleJwtAuth gin.HandlerFunc

const PermissionKey = "permission"

func init() {
	const identityKey = "id"
	const roleKey = "role"
	var err error
	JWTMiddleWare, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:        config.Global.JWTRealm,
		Key:          []byte(config.Global.JWTSecret),
		Timeout:      time.Minute * 10,
		CookieMaxAge: time.Hour * 24 * 7,
		MaxRefresh:   time.Hour * 24 * 7,
		IdentityKey:  identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if u, ok := data.(LoginUser); ok {
				return jwt.MapClaims{
					identityKey: u.ID,
					roleKey:     u.Role,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: Login,
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if p, ok := c.Get(PermissionKey); ok {
				if c, ok := data.(jwt.MapClaims); ok {
					auth.VerifyRoleHasPermission(c[roleKey].(string), p.(string))
					return true
				}
			}
			return false
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			var user = c.Keys["user"]
			c.JSON(http.StatusOK, user)
		},
		RefreshResponse: func(c *gin.Context, code int, message string, time time.Time) {
			var user = c.Keys["user"]
			c.JSON(http.StatusOK, user)
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		CookieName:     "token",
		TokenLookup:    "cookie:token",
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost",
		CookieSameSite: http.SameSiteDefaultMode,
		TimeFunc:       time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	handleJwtAuth = JWTMiddleWare.MiddlewareFunc()
}
