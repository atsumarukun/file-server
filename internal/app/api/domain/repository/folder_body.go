package repository

import "file-server/internal/app/api/domain/entity"

type FolderBodyRepository interface {
	Create(*entity.FolderBody) error
}
