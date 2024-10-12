package usecase

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"file-server/internal/app/api/domain/service"
	"file-server/internal/app/api/usecase/dto"
	"file-server/internal/pkg/zip"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type FolderUsecase interface {
	Create(uint64, string, bool) (*dto.FolderInfoDTO, error)
	Update(uint64, string, bool, bool) (*dto.FolderInfoDTO, error)
	Remove(uint64, bool) error
	Move(uint64, uint64, bool) (*dto.FolderInfoDTO, error)
	Copy(uint64, uint64, bool) (*dto.FolderInfoDTO, error)
	FindOne(string, bool) (*dto.FolderInfoDTO, error)
	Read(uint64, bool) (*dto.FolderBodyDTO, error)
}

type folderUsecase struct {
	db                   *gorm.DB
	folderInfoRepository repository.FolderInfoRepository
	folderBodyRepository repository.FolderBodyRepository
	folderInfoService    service.FolderInfoService
}

func NewFolderUsecase(db *gorm.DB, folderInfoRepository repository.FolderInfoRepository, folderBodyRepository repository.FolderBodyRepository, folderInfoService service.FolderInfoService) FolderUsecase {
	return &folderUsecase{
		db:                   db,
		folderInfoRepository: folderInfoRepository,
		folderBodyRepository: folderBodyRepository,
		folderInfoService:    folderInfoService,
	}
}

func (fu *folderUsecase) Create(parentFolderID uint64, name string, isHide bool) (*dto.FolderInfoDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}

		path := parentFolder.Path.Value + name + "/"

		folderInfo, err = entity.NewFolderInfo(&parentFolderID, name, path, isHide)
		if err != nil {
			return err
		}

		if isExists, err := fu.folderInfoService.IsExists(tx, folderInfo); err != nil {
			return err
		} else if isExists {
			return fmt.Errorf("%s is already exists", folderInfo.Path.Value)
		}

		folderInfo, err = fu.folderInfoRepository.Create(tx, folderInfo)
		if err != nil {
			return err
		}

		folderBody := entity.NewFolderBody(path)

		return fu.folderBodyRepository.Create(folderBody)
	}); err != nil {
		return nil, err
	}

	return fu.convertToFolderInfoDTO(folderInfo), nil
}

func (fu *folderUsecase) Update(id uint64, name string, isHide bool, isDisplayHiddenObject bool) (*dto.FolderInfoDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		if folderInfo.IsRoot() {
			return fmt.Errorf("root directory is not updatable")
		}

		folderInfo.IsHide = isHide

		if name != folderInfo.Name.Value {
			oldPath := folderInfo.Path.Value
			path := oldPath[:strings.LastIndex(oldPath, folderInfo.Name.Value)] + name + "/"

			if err := folderInfo.SetName(name); err != nil {
				return err
			}

			if err := folderInfo.Move(oldPath, path); err != nil {
				return err
			}

			if isExists, err := fu.folderInfoService.IsExists(tx, folderInfo); err != nil {
				return err
			} else if isExists {
				return fmt.Errorf("%s is already exists", folderInfo.Path.Value)
			}

			if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
				return err
			}
		}

		folderInfo, err = fu.folderInfoRepository.Update(tx, folderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.convertToFolderInfoDTO(folderInfo), nil
}

func (fu *folderUsecase) Remove(id uint64, isDisplayHiddenObject bool) error {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		if folderInfo.IsRoot() {
			return fmt.Errorf("root directory is not removable")
		}

		if err := fu.folderBodyRepository.Remove(folderInfo.Path.Value); err != nil {
			return err
		}

		return fu.folderInfoRepository.Remove(tx, folderInfo)
	}); err != nil {
		return err
	}

	return nil
}

func (fu *folderUsecase) Move(id uint64, parentFolderID uint64, isDisplayHiddenObject bool) (*dto.FolderInfoDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if isDisplayHiddenObject {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			folderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		if folderInfo.IsRoot() {
			return fmt.Errorf("root directory is not updatable")
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}

		oldPath := folderInfo.Path.Value
		if strings.Contains(parentFolder.Path.Value, oldPath) {
			return fmt.Errorf("cannot move to lower directory")
		}
		path := parentFolder.Path.Value + folderInfo.Name.Value + "/"

		if err := folderInfo.Move(oldPath, path); err != nil {
			return err
		}
		folderInfo.ParentFolderID = &parentFolderID

		if isExists, err := fu.folderInfoService.IsExists(tx, folderInfo); err != nil {
			return err
		} else if isExists {
			return fmt.Errorf("%s is already exists", folderInfo.Path.Value)
		}

		if err := fu.folderBodyRepository.Update(oldPath, path); err != nil {
			return err
		}

		folderInfo, err = fu.folderInfoRepository.Update(tx, folderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.convertToFolderInfoDTO(folderInfo), nil
}

func (fu *folderUsecase) Copy(id uint64, parentFolderID uint64, isDisplayHiddenObject bool) (*dto.FolderInfoDTO, error) {
	var folderInfo *entity.FolderInfo
	if err := fu.db.Transaction(func(tx *gorm.DB) error {
		var sourceFolderInfo *entity.FolderInfo
		var err error
		if isDisplayHiddenObject {
			sourceFolderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(tx, id)
		} else {
			sourceFolderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(tx, id, false)
		}
		if err != nil {
			return err
		}

		sourceFolderBody, err := fu.folderBodyRepository.Read(sourceFolderInfo.Path.Value)
		if err != nil {
			return err
		}

		parentFolder, err := fu.folderInfoRepository.FindOneByID(tx, parentFolderID)
		if err != nil {
			return err
		}

		path := parentFolder.Path.Value + sourceFolderInfo.Name.Value + "/"
		targetFolderInfo, err := sourceFolderInfo.Copy(path)
		if err != nil {
			return err
		}
		targetFolderInfo.ParentFolderID = &parentFolderID

		if isExists, err := fu.folderInfoService.IsExists(tx, targetFolderInfo); err != nil {
			return err
		} else if isExists {
			return fmt.Errorf("%s is already exists", targetFolderInfo.Path.Value)
		}

		targetFolderBody := sourceFolderBody.Copy(path)
		if err := fu.folderBodyRepository.Create(targetFolderBody); err != nil {
			return err
		}

		folderInfo, err = fu.folderInfoRepository.Create(tx, targetFolderInfo)
		return err
	}); err != nil {
		return nil, err
	}

	return fu.convertToFolderInfoDTO(folderInfo), nil
}

func (fu *folderUsecase) FindOne(path string, isDisplayHiddenObject bool) (*dto.FolderInfoDTO, error) {
	var folderInfo *entity.FolderInfo
	var err error
	if isDisplayHiddenObject {
		folderInfo, err = fu.folderInfoRepository.FindOneByPathWithChildren(fu.db, path)
	} else {
		folderInfo, err = fu.folderInfoRepository.FindOneByPathAndIsHideWithChildren(fu.db, path, false)
	}
	if err != nil {
		return nil, err
	}

	return fu.convertToFolderInfoDTO(folderInfo), nil
}

func (fu *folderUsecase) Read(id uint64, isDisplayHiddenObject bool) (*dto.FolderBodyDTO, error) {
	var folderInfo *entity.FolderInfo
	var err error
	if isDisplayHiddenObject {
		folderInfo, err = fu.folderInfoRepository.FindOneByIDWithLower(fu.db, id)
	} else {
		folderInfo, err = fu.folderInfoRepository.FindOneByIDAndIsHideWithLower(fu.db, id, false)
	}
	if err != nil {
		return nil, err
	}

	zipFile, err := zip.Compress(folderInfo.Path.Value)
	if err != nil {
		return nil, err
	}

	return dto.NewFolderBodyDTO("application/zip", zipFile), nil
}

func (fu *folderUsecase) convertToFolderInfoDTO(folder *entity.FolderInfo) *dto.FolderInfoDTO {
	folders := make([]dto.FolderInfoDTO, len(folder.Folders))
	for i, v := range folder.Folders {
		folders[i] = *fu.convertToFolderInfoDTO(&v)
	}

	files := make([]dto.FileInfoDTO, len(folder.Files))
	for i, v := range folder.Files {
		files[i] = *dto.NewFileInfoDTO(v.ID, v.FolderID, v.Name.Value, v.Path.Value, v.MimeType.Value, v.IsHide, v.CreatedAt, v.UpdatedAt)
	}

	return dto.NewFolderInfoDTO(folder.ID, folder.ParentFolderID, folder.Name.Value, folder.Path.Value, folder.IsHide, folders, files, folder.CreatedAt, folder.UpdatedAt)
}
