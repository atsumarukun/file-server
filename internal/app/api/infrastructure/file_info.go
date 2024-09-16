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

func (fi *fileInfoInfrastructure) Create(db *gorm.DB, file *entity.FileInfo) (*entity.FileInfo, error) {
	fileModel := fi.entityToModel(file)
	if err := db.Create(fileModel).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(fileModel)
}

func (fi *fileInfoInfrastructure) Creates(db *gorm.DB, files []entity.FileInfo) ([]entity.FileInfo, error) {
	fileModels := fi.entitiesToModels(files)
	if err := db.Create(fileModels).Error; err != nil {
		return nil, err
	}
	return fi.modelsToEntities(fileModels)
}

func (fi *fileInfoInfrastructure) Update(db *gorm.DB, file *entity.FileInfo) (*entity.FileInfo, error) {
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

func (fi *fileInfoInfrastructure) FindOneByIDAndIsHide(db *gorm.DB, id uint64, isHide bool) (*entity.FileInfo, error) {
	var fileModel model.FileModel
	if err := db.First(&fileModel, "id = ? and is_hide = ?", id, isHide).Error; err != nil {
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
		ID:        file.ID,
		FolderID:  file.FolderID,
		Name:      file.Name.Value,
		Path:      file.Path.Value,
		MimeType:  file.MimeType.Value,
		IsHide:    file.IsHide,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}

func (fi *fileInfoInfrastructure) entitiesToModels(files []entity.FileInfo) []model.FileModel {
	models := make([]model.FileModel, len(files))
	for i, file := range files {
		model := fi.entityToModel(&file)
		models[i] = *model
	}
	return models
}

func (fi *fileInfoInfrastructure) modelToEntity(file *model.FileModel) (*entity.FileInfo, error) {
	fileEntity := &entity.FileInfo{}
	fileEntity.ID = file.ID
	fileEntity.FolderID = file.FolderID
	if err := fileEntity.SetName(file.Name); err != nil {
		return nil, err
	}
	if err := fileEntity.SetPath(file.Path); err != nil {
		return nil, err
	}
	if err := fileEntity.SetMimeType(file.MimeType); err != nil {
		return nil, err
	}
	fileEntity.IsHide = file.IsHide
	fileEntity.CreatedAt = file.CreatedAt
	fileEntity.UpdatedAt = file.UpdatedAt
	return fileEntity, nil
}

func (fi *fileInfoInfrastructure) modelsToEntities(files []model.FileModel) ([]entity.FileInfo, error) {
	entities := make([]entity.FileInfo, len(files))
	for i, file := range files {
		entity, err := fi.modelToEntity(&file)
		if err != nil {
			return nil, err
		}
		entities[i] = *entity
	}
	return entities, nil
}
