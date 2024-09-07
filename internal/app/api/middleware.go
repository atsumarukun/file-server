package api

import (
	"errors"
	"file-server/internal/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != "" {
			token := strings.Split(c.Request.Header.Get("Authorization"), " ")
			if len(token) != 2 || token[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, "middleware: invalid token")
				c.Abort()
				return
			}

			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(token[1], claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.JWT_SECRET_KEY), nil
			})
			c.Set("isDisplayHiddenObject", true)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					c.Set("isDisplayHiddenObject", false)
				} else {
					c.JSON(http.StatusUnauthorized, "middleware: invalid token")
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
