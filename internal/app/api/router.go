package api

import (
	"file-server/internal/app/api/infrastructure"
	"file-server/internal/app/api/interface/handler"
	"file-server/internal/app/api/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func route(r *gin.Engine, db *gorm.DB) {
	folderRepository := infrastructure.NewFolderInfrastructure()
	folderUsecase := usecase.NewFolderUsecase(db, folderRepository)
	folderHandler := handler.NewFolderHandler(folderUsecase)

	folders := r.Group("/folders")
	{
		folders.POST("/", folderHandler.Create)
	}
}
