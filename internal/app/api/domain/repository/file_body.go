package repository

import "file-server/internal/app/api/domain/entity"

type FileBodyRepository interface {
	Create(*entity.FileBody) error
}
