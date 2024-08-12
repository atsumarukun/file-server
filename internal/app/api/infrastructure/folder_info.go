package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/infrastructure/model"

	"gorm.io/gorm"
)

type folderInfoInfrastructure struct{}

func NewFolderInfoInfrastructure() repository.FolderInfoRepository {
	return &folderInfoInfrastructure{}
}

func (fi *folderInfoInfrastructure) Save(db *gorm.DB, folder *entity.FolderInfo) (*entity.FolderInfo, error) {
	folderModel := fi.entityToModel(folder)
	if err := db.Save(folderModel).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(folderModel)
}

func (fi *folderInfoInfrastructure) Saves(db *gorm.DB, folders []entity.FolderInfo) ([]entity.FolderInfo, error) {
	folderModels := fi.entitiesToModels(folders)
	if err := db.Save(folderModels).Error; err != nil {
		return nil, err
	}
	return fi.modelsToEntities(folderModels)
}

func (fi *folderInfoInfrastructure) Remove(db *gorm.DB, folder *entity.FolderInfo) error {
	folderModel := fi.entityToModel(folder)
	return db.Delete(folderModel).Error
}

func (fi *folderInfoInfrastructure) FindOneByID(db *gorm.DB, id int64) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.First(&folderModel, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindOneByPath(db *gorm.DB, path string) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.First(&folderModel, "path = ?", path).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindOneByPathWithRelationship(db *gorm.DB, path string) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.Preload("Folders").First(&folderModel, "path = ?", path).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindByIDNotAndPathLike(db *gorm.DB, id int64, path string) ([]entity.FolderInfo, error) {
	var folderModels []model.FolderModel
	if err := db.Find(&folderModels, "id <> ? AND path LIKE ?", id, path+"%").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return fi.modelsToEntities(folderModels)
}

func (fi *folderInfoInfrastructure) entityToModel(folder *entity.FolderInfo) *model.FolderModel {
	var folders []model.FolderModel
	if folder.GetFolders() != nil {
		folders = make([]model.FolderModel, len(folder.GetFolders()))
		for i, v := range folder.GetFolders() {
			folders[i] = *fi.entityToModel(&v)
		}
	}
	return &model.FolderModel{
		ID:             folder.GetID(),
		ParentFolderID: folder.GetParentFolderID(),
		Name:           folder.GetName(),
		Path:           folder.GetPath(),
		IsHide:         folder.GetIsHide(),
		Folders:        folders,
		CreatedAt:      folder.GetCreatedAt(),
		UpdatedAt:      folder.GetUpdatedAt(),
	}
}

func (fi *folderInfoInfrastructure) entitiesToModels(folders []entity.FolderInfo) []model.FolderModel {
	folderModels := make([]model.FolderModel, len(folders))
	for i, folder := range folders {
		folderModels[i] = *fi.entityToModel(&folder)
	}
	return folderModels
}

func (fi *folderInfoInfrastructure) modelToEntity(folder *model.FolderModel) (*entity.FolderInfo, error) {
	var folders []entity.FolderInfo
	if folder.Folders != nil {
		folders = make([]entity.FolderInfo, len(folder.Folders))
		for i, v := range folder.Folders {
			f, err := fi.modelToEntity(&v)
			if err != nil {
				return nil, err
			}
			folders[i] = *f
		}
	}
	folderEntity := &entity.FolderInfo{}
	folderEntity.SetID(folder.ID)
	folderEntity.SetParentFolderID(folder.ParentFolderID)
	if err := folderEntity.SetName(folder.Name); err != nil {
		return nil, err
	}
	if err := folderEntity.SetPath(folder.Path); err != nil {
		return nil, err
	}
	folderEntity.SetIsHide(folder.IsHide)
	folderEntity.SetFolders(folders)
	folderEntity.SetCreatedAt(folder.CreatedAt)
	folderEntity.SetUpdatedAt(folder.UpdatedAt)
	return folderEntity, nil
}

func (fi *folderInfoInfrastructure) modelsToEntities(folders []model.FolderModel) ([]entity.FolderInfo, error) {
	folderEntities := make([]entity.FolderInfo, len(folders))
	for i, folder := range folders {
		folderEntity, err := fi.modelToEntity(&folder)
		if err != nil {
			return nil, err
		}
		folderEntities[i] = *folderEntity
	}
	return folderEntities, nil
}
