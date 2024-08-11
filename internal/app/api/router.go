package api

import (
	"file-server/internal/app/api/infrastructure"
	"file-server/internal/app/api/interface/handler"
	"file-server/internal/app/api/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func route(r *gin.Engine, db *gorm.DB) {
	folderInfoRepository := infrastructure.NewFolderInfoInfrastructure()
	folderBodyRepository := infrastructure.NewFolderBodyInfrastructure()
	folderUsecase := usecase.NewFolderUsecase(db, folderInfoRepository, folderBodyRepository)
	folderHandler := handler.NewFolderHandler(folderUsecase)

	folders := r.Group("/folders")
	{
		folders.POST("/", folderHandler.Create)
		folders.GET("/*path", folderHandler.FindOne)
	}
}
