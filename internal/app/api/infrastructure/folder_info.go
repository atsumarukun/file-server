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

func (fi *folderInfoInfrastructure) Removes(db *gorm.DB, folders []entity.FolderInfo) error {
	folderModels := fi.entitiesToModels(folders)
	return db.Delete(folderModels).Error
}

func (fi *folderInfoInfrastructure) FindOneByID(db *gorm.DB, id uint64) (*entity.FolderInfo, error) {
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

func (fi *folderInfoInfrastructure) FindOneByIDWithRelationship(db *gorm.DB, id uint64) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.Preload("Folders").Preload("Files").First(&folderModel, "id = ?", id).Error; err != nil {
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
	if err := db.Preload("Folders").Preload("Files").First(&folderModel, "path = ?", path).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindByIDNotAndPathLike(db *gorm.DB, id uint64, path string) ([]entity.FolderInfo, error) {
	var folderModels []model.FolderModel
	if err := db.Find(&folderModels, "id <> ? AND path LIKE ?", id, path+"%").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return make([]entity.FolderInfo, 0), nil
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
	var files []model.FileModel
	if folder.GetFiles() != nil {
		files = make([]model.FileModel, len(folder.GetFiles()))
		for i, v := range folder.GetFiles() {
			files[i] = model.FileModel{
				ID:        v.GetID(),
				FolderID:  v.GetFolderID(),
				Name:      v.GetName(),
				Path:      v.GetPath(),
				MimeType:  v.GetMimeType(),
				IsHide:    v.GetIsHide(),
				CreatedAt: v.GetCreatedAt(),
				UpdatedAt: v.GetUpdatedAt(),
			}
		}
	}
	return &model.FolderModel{
		ID:             folder.GetID(),
		ParentFolderID: folder.GetParentFolderID(),
		Name:           folder.GetName(),
		Path:           folder.GetPath(),
		IsHide:         folder.GetIsHide(),
		Folders:        folders,
		Files:          files,
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
	var files []entity.FileInfo
	if folder.Files != nil {
		files = make([]entity.FileInfo, len(folder.Files))
		for i, v := range folder.Files {
			f := &entity.FileInfo{}
			f.SetID(v.ID)
			f.SetFolderID(v.FolderID)
			if err := f.SetName(v.Name); err != nil {
				return nil, err
			}
			if err := f.SetPath(v.Path); err != nil {
				return nil, err
			}
			if err := f.SetMimeType(v.MimeType); err != nil {
				return nil, err
			}
			f.SetIsHide(v.IsHide)
			f.SetCreatedAt(v.CreatedAt)
			f.SetUpdatedAt(v.UpdatedAt)
			files[i] = *f
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
	folderEntity.SetFiles(files)
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
