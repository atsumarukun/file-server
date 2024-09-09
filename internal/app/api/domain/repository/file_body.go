package repository

import "file-server/internal/app/api/domain/entity"

type FileBodyRepository interface {
	Create(*entity.FileBody) error
	Update(string, string) error
	Remove(string) error
	Read(string) (*entity.FileBody, error)
}
