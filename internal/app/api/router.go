package api

import "github.com/gin-gonic/gin"

func route(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/signin", authHandler.Signin)
	}

	folders := r.Group("/folders")
	{
		folders.POST("/", folderHandler.Create)

		folders.Use(authMiddleware())

		folders.GET("/*path", folderHandler.FindOne)
		folders.PUT("/:id", folderHandler.Update)
		folders.DELETE("/:id", folderHandler.Remove)
		folders.PUT("/:id/move", folderHandler.Move)
		folders.POST("/:id/copy", folderHandler.Copy)
	}

	files := r.Group("/files")
	{
		files.POST("/", fileHandler.Create)

		files.Use(authMiddleware())

		files.PUT("/:id", fileHandler.Update)
		files.DELETE("/:id", fileHandler.Remove)
		files.PUT("/:id/move", fileHandler.Move)
		files.POST("/:id/copy", fileHandler.Copy)
	}

	batch := r.Group("/batch")
	{
		batch.POST("/", batchMiddleware(r))
	}
}
