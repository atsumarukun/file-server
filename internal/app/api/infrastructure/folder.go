package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/infrastructure/model"

	"gorm.io/gorm"
)

type folderInfrastructure struct{}

func NewFolderInfrastructure() repository.FolderRepository {
	return &folderInfrastructure{}
}

func (fi *folderInfrastructure) Create(db *gorm.DB, folder *entity.Folder) (*entity.Folder, error) {
	folderModel := fi.entityToModel(folder)
	if err := db.Create(folderModel).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(folderModel)
}

func (fi *folderInfrastructure) entityToModel(folder *entity.Folder) *model.FolderModel {
	return &model.FolderModel{
		ID:             folder.GetID(),
		ParentFolderID: folder.GetParentFolderID(),
		Name:           folder.GetName(),
		Path:           folder.GetPath(),
		IsHide:         folder.GetIsHide(),
		CreatedAt:      folder.GetCreatedAt(),
		UpdatedAt:      folder.GetUpdatedAt(),
	}
}

func (fi *folderInfrastructure) modelToEntity(folder *model.FolderModel) (*entity.Folder, error) {
	folderEntity := &entity.Folder{}
	folderEntity.SetID(folder.ID)
	folderEntity.SetParentFolderID(folder.ParentFolderID)
	if err := folderEntity.SetName(folder.Name); err != nil {
		return nil, err
	}
	if err := folderEntity.SetPath(folder.Path); err != nil {
		return nil, err
	}
	folderEntity.SetIsHide(folder.IsHide)
	folderEntity.SetCreatedAt(folder.CreatedAt)
	folderEntity.SetUpdatedAt(folder.UpdatedAt)
	return folderEntity, nil
}
