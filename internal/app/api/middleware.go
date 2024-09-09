package api

import (
	"bytes"
	"errors"
	"file-server/internal/app/api/interface/requests"
	"file-server/internal/pkg/config"
	"net/http"
	"strings"
	"sync"

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

func batchMiddleware(engine *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requests []requests.BatchRequest
		if err := c.ShouldBindJSON(&requests); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		responses := make([]gin.H, len(requests))

		var wg sync.WaitGroup
		for i, req := range requests {
			wg.Add(1)
			go func() {
				defer wg.Done()

				r, err := http.NewRequest(req.Method, req.Path, bytes.NewBuffer([]byte(req.Body)))
				if err != nil {
					c.JSON(http.StatusInternalServerError, err.Error())
					return
				}
				r.Header = c.Request.Header

				w := &responseWriter{header: make(http.Header)}

				engine.ServeHTTP(w, r)

				responses[i] = gin.H{
					"status":  w.status,
					"headers": w.header,
					"body":    string(w.body),
				}
			}()
		}
		wg.Wait()

		c.JSON(http.StatusOK, responses)
	}
}

type responseWriter struct {
	header http.Header
	status int
	body   []byte
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}
