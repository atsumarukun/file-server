package service

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/internal/app/api/domain/repository"
	"fmt"

	"gorm.io/gorm"
)

type FolderInfoService interface {
	Exists(*gorm.DB, *entity.FolderInfo) error
	Move(*gorm.DB, *entity.FolderInfo, string) error
	Copy(*gorm.DB, *entity.FolderInfo, string) (*entity.FolderInfo, error)
}

type folderInfoService struct {
	folderInfoRepository repository.FolderInfoRepository
	folderBodyRepository repository.FolderBodyRepository
	fileInfoRepository   repository.FileInfoRepository
	fileBodyRepository   repository.FileBodyRepository
}

func NewFolderInfoService(folderInfoRepository repository.FolderInfoRepository, folderBodyRepository repository.FolderBodyRepository, fileInfoRepository repository.FileInfoRepository, fileBodyRepository repository.FileBodyRepository) FolderInfoService {
	return &folderInfoService{
		folderInfoRepository: folderInfoRepository,
		folderBodyRepository: folderBodyRepository,
		fileInfoRepository:   fileInfoRepository,
		fileBodyRepository:   fileBodyRepository,
	}
}

func (fs *folderInfoService) Exists(db *gorm.DB, folder *entity.FolderInfo) error {
	path := folder.GetPath()
	if checkFolder, err := fs.folderInfoRepository.FindOneByPath(db, path); err != nil {
		return err
	} else if checkFolder != nil {
		return fmt.Errorf("%s is already exists", path)
	}
	return nil
}

func (fs *folderInfoService) Move(db *gorm.DB, folder *entity.FolderInfo, path string) error {
	oldPath := folder.GetPath()
	folder.SetPath(path)

	if err := fs.Exists(db, folder); err != nil {
		return err
	}

	lowerFolders, err := fs.folderInfoRepository.FindByIDNotAndPathLike(db, folder.GetID(), oldPath)
	if 0 < len(lowerFolders) {
		if err != nil {
			return err
		}
		for i := 0; i < len(lowerFolders); i++ {
			if err := lowerFolders[i].Move(oldPath, path); err != nil {
				return err
			}
		}
		if _, err := fs.folderInfoRepository.Saves(db, lowerFolders); err != nil {
			return err
		}
	}

	lowerFiles, err := fs.fileInfoRepository.FindByPathLike(db, oldPath)
	if 0 < len(lowerFiles) {
		if err != nil {
			return err
		}
		for i := 0; i < len(lowerFiles); i++ {
			if err := lowerFiles[i].Move(oldPath, path); err != nil {
				return err
			}
		}
		if _, err := fs.fileInfoRepository.Saves(db, lowerFiles); err != nil {
			return err
		}
	}

	return nil
}

func (fs *folderInfoService) Copy(db *gorm.DB, folder *entity.FolderInfo, path string) (*entity.FolderInfo, error) {
	newFolder, err := entity.NewFolderInfo(nil, folder.GetName(), path, folder.GetIsHide())
	if err != nil {
		return nil, err
	}

	if err := fs.Exists(db, newFolder); err != nil {
		return nil, err
	}

	if err := fs.folderBodyRepository.Create(path); err != nil {
		return nil, err
	}

	children := folder.GetFolders()
	if 0 < len(children) {
		newChildren := make([]entity.FolderInfo, len(children))
		for i := 0; i < len(children); i++ {
			child, err := fs.folderInfoRepository.FindOneByIDWithChildren(db, children[i].GetID())
			if err != nil {
				return nil, err
			}
			newChild, err := fs.Copy(db, child, path+child.GetName()+"/")
			if err != nil {
				return nil, err
			}
			newChildren[i] = *newChild
		}
		newFolder.SetFolders(newChildren)
	}

	files := folder.GetFiles()
	if 0 < len(files) {
		newFiles := make([]entity.FileInfo, len(files))
		for i, file := range files {
			fileBody, err := fs.fileBodyRepository.Read(file.GetPath())
			if err != nil {
				return nil, err
			}
			fileBody.SetPath(path + file.GetName())
			fs.fileBodyRepository.Create(fileBody)

			newFile, err := entity.NewFileInfo(0, file.GetName(), path+file.GetName(), file.GetMimeType(), file.GetIsHide())
			if err != nil {
				return nil, err
			}
			newFiles[i] = *newFile
		}
		newFolder.SetFiles(newFiles)
	}

	return newFolder, nil
}
