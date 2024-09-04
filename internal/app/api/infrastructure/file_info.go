package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/infrastructure/model"

	"gorm.io/gorm"
)

type fileInfoInfrastructure struct{}

func NewFileInfoInfrastructure() repository.FileInfoRepository {
	return &fileInfoInfrastructure{}
}

func (fi *fileInfoInfrastructure) Save(db *gorm.DB, file *entity.FileInfo) (*entity.FileInfo, error) {
	fileModel := fi.entityToModel(file)
	if err := db.Save(fileModel).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(fileModel)
}

func (fi *fileInfoInfrastructure) Remove(db *gorm.DB, file *entity.FileInfo) error {
	fileModel := fi.entityToModel(file)
	return db.Delete(fileModel).Error
}

func (fi *fileInfoInfrastructure) FindOneByID(db *gorm.DB, id uint64) (*entity.FileInfo, error) {
	var fileModel model.FileModel
	if err := db.First(&fileModel, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(&fileModel)
}

func (fi *fileInfoInfrastructure) FindOneByPath(db *gorm.DB, path string) (*entity.FileInfo, error) {
	var fileModel model.FileModel
	if err := db.First(&fileModel, "path = ?", path).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(&fileModel)
}

func (fi *fileInfoInfrastructure) entityToModel(file *entity.FileInfo) *model.FileModel {
	return &model.FileModel{
		ID:        file.GetID(),
		FolderID:  file.GetFolderID(),
		Name:      file.GetName(),
		Path:      file.GetPath(),
		MimeType:  file.GetMimeType(),
		IsHide:    file.GetIsHide(),
		CreatedAt: file.GetCreatedAt(),
		UpdatedAt: file.GetUpdatedAt(),
	}
}

func (fi *fileInfoInfrastructure) modelToEntity(file *model.FileModel) (*entity.FileInfo, error) {
	fileEntity := &entity.FileInfo{}
	fileEntity.SetID(file.ID)
	fileEntity.SetFolderID(file.FolderID)
	if err := fileEntity.SetName(file.Name); err != nil {
		return nil, err
	}
	if err := fileEntity.SetPath(file.Path); err != nil {
		return nil, err
	}
	if err := fileEntity.SetMimeType(file.MimeType); err != nil {
		return nil, err
	}
	fileEntity.SetIsHide(file.IsHide)
	fileEntity.SetCreatedAt(file.CreatedAt)
	fileEntity.SetUpdatedAt(file.UpdatedAt)
	return fileEntity, nil
}
