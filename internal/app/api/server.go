package api

import (
	"context"
	"file-server/internal/pkg/config"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Serve() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(config.MYSQL_DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println(db.Config)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello World!")
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.API_PORT),
		Handler: r,
	}

	go srv.ListenAndServe()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
