package Middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	common "elementser/src/common"
	help "elementser/src/help"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/dev-api/vue-element-admin/user/logout" || path == "/dev-api/vue-element-admin/user/login" {
			c.Next()
			return
		}

		fmt.Println("CheckToken======", path)

		token, err := c.Cookie("Admin-Token")
		if err != nil {
			common.ResFail(c, "token err")
			return
		}

		_, err = help.ValidateJwtToken(token)
		if err != nil {
			common.ResFail(c, "token err")
			return
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		reqHead := c.Request.Header.Get("Access-Control-Request-Headers")
		println("reqHead:==========", reqHead)
		println(method, origin)
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			if reqHead == "content-type" {
				c.AbortWithStatus(http.StatusNoContent)
			}
			/*
				if reqHead == "x-token" {
					c.AbortWithStatus(http.StatusOK)
				}*/

		}
		c.Next()
	}
}
