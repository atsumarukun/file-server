package api

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestServe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Error("")
}
