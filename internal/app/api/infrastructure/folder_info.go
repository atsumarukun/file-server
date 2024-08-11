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

func (fi *folderInfoInfrastructure) Create(db *gorm.DB, folder *entity.FolderInfo) (*entity.FolderInfo, error) {
	folderModel := fi.entityToModel(folder)
	if err := db.Create(folderModel).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(folderModel)
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

func (fi *folderInfoInfrastructure) entityToModel(folder *entity.FolderInfo) *model.FolderModel {
	var folders []model.FolderModel
	if folder.GetFolders() != nil {
		folders = make([]model.FolderModel, 0)
		for _, v := range folder.GetFolders() {
			folders = append(folders, *fi.entityToModel(&v))
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

func (fi *folderInfoInfrastructure) modelToEntity(folder *model.FolderModel) (*entity.FolderInfo, error) {
	var folders []entity.FolderInfo
	if folder.Folders != nil {
		folders = make([]entity.FolderInfo, 0)
		for _, v := range folder.Folders {
			f, err := fi.modelToEntity(&v)
			if err != nil {
				return nil, err
			}
			folders = append(folders, *f)
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
