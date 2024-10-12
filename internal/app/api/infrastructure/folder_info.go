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

func (fi *folderInfoInfrastructure) Update(db *gorm.DB, folder *entity.FolderInfo) (*entity.FolderInfo, error) {
	folderModel := fi.entityToModel(folder)
	if err := db.Save(folderModel).Error; err != nil {
		return nil, err
	}
	if 0 < len(folderModel.Folders) {
		for _, v := range folder.Folders {
			if _, err := fi.Update(db, &v); err != nil {
				return nil, err
			}
		}
	}
	if 0 < len(folderModel.Files) {
		if err := db.Save(folderModel.Files).Error; err != nil {
			return nil, err
		}
	}
	return fi.modelToEntity(folderModel)
}

func (fi *folderInfoInfrastructure) Remove(db *gorm.DB, folder *entity.FolderInfo) error {
	folderModel := fi.entityToModel(folder)
	return db.Delete(folderModel).Error
}

func (fi *folderInfoInfrastructure) FindOneByID(db *gorm.DB, id uint64) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.First(&folderModel, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindOneByPath(db *gorm.DB, path string) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.First(&folderModel, "path = ?", path).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindOneByPathWithChildren(db *gorm.DB, path string) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.Preload("Folders").Preload("Files").First(&folderModel, "path = ?", path).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindOneByPathAndIsHideWithChildren(db *gorm.DB, path string, isHide bool) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.Preload("Folders", "is_hide", isHide).Preload("Files", "is_hide", isHide).First(&folderModel, "path = ? and is_hide = ?", path, isHide).Error; err != nil {
		return nil, err
	}
	return fi.modelToEntity(&folderModel)
}

func (fi *folderInfoInfrastructure) FindOneByIDWithLower(db *gorm.DB, id uint64) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.Preload("Folders").Preload("Files").First(&folderModel, "id = ?", id).Error; err != nil {
		return nil, err
	}
	folder, err := fi.modelToEntity(&folderModel)
	if err != nil {
		return nil, err
	}
	folders := folder.Folders
	if 0 < len(folders) {
		for i, v := range folders {
			f, err := fi.FindOneByIDWithLower(db, v.ID)
			if err != nil {
				return nil, err
			}
			folders[i] = *f
		}
		folder.Folders = folders
	}
	return folder, nil
}

func (fi *folderInfoInfrastructure) FindOneByIDAndIsHideWithLower(db *gorm.DB, id uint64, isHide bool) (*entity.FolderInfo, error) {
	var folderModel model.FolderModel
	if err := db.Preload("Folders", "is_hide", isHide).Preload("Files", "is_hide", isHide).First(&folderModel, "id = ? and is_hide = ?", id, isHide).Error; err != nil {
		return nil, err
	}
	folder, err := fi.modelToEntity(&folderModel)
	if err != nil {
		return nil, err
	}
	folders := folder.Folders
	if 0 < len(folders) {
		for i, v := range folders {
			f, err := fi.FindOneByIDAndIsHideWithLower(db, v.ID, isHide)
			if err != nil {
				return nil, err
			}
			folders[i] = *f
		}
		folder.Folders = folders
	}
	return folder, nil
}

func (fi *folderInfoInfrastructure) entityToModel(folder *entity.FolderInfo) *model.FolderModel {
	var folders []model.FolderModel
	if folder.Folders != nil {
		folders = make([]model.FolderModel, len(folder.Folders))
		for i, v := range folder.Folders {
			folders[i] = *fi.entityToModel(&v)
		}
	}
	var files []model.FileModel
	if folder.Files != nil {
		files = make([]model.FileModel, len(folder.Files))
		for i, v := range folder.Files {
			files[i] = model.FileModel{
				ID:        v.ID,
				FolderID:  v.FolderID,
				Name:      v.Name.Value,
				Path:      v.Path.Value,
				MimeType:  v.MimeType.Value,
				IsHide:    v.IsHide,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			}
		}
	}
	return &model.FolderModel{
		ID:             folder.ID,
		ParentFolderID: folder.ParentFolderID,
		Name:           folder.Name.Value,
		Path:           folder.Path.Value,
		IsHide:         folder.IsHide,
		Folders:        folders,
		Files:          files,
		CreatedAt:      folder.CreatedAt,
		UpdatedAt:      folder.UpdatedAt,
	}
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
			f.ID = v.ID
			f.FolderID = v.FolderID
			if err := f.SetName(v.Name); err != nil {
				return nil, err
			}
			if err := f.SetPath(v.Path); err != nil {
				return nil, err
			}
			if err := f.SetMimeType(v.MimeType); err != nil {
				return nil, err
			}
			f.IsHide = v.IsHide
			f.CreatedAt = v.CreatedAt
			f.UpdatedAt = v.UpdatedAt
			files[i] = *f
		}
	}
	folderEntity := &entity.FolderInfo{}
	folderEntity.ID = folder.ID
	folderEntity.ParentFolderID = folder.ParentFolderID
	if err := folderEntity.SetName(folder.Name); err != nil {
		return nil, err
	}
	if err := folderEntity.SetPath(folder.Path); err != nil {
		return nil, err
	}
	folderEntity.IsHide = folder.IsHide
	folderEntity.Folders = folders
	folderEntity.Files = files
	folderEntity.CreatedAt = folder.CreatedAt
	folderEntity.UpdatedAt = folder.UpdatedAt
	return folderEntity, nil
}
