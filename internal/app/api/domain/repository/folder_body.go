package repository

import "file-server/internal/app/api/domain/entity"

type FolderBodyRepository interface {
	Create(*entity.FolderBody) error
	Update(string, string) error
	Remove(string) error
	Read(string) (*entity.FolderBody, error)
}
