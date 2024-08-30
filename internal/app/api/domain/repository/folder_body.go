package repository

type FolderBodyRepository interface {
	Create(string) error
	Update(string, string) error
}
