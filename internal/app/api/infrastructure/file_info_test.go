package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/test/database"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	file, err := entity.NewFileInfo(1, "name", "/path/", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `files` (`folder_id`,`name`,`path`,`mime_type`,`is_hide`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")).WithArgs(file.FolderID, file.Name.Value, file.Path.Value, file.MimeType.Value, file.IsHide, database.AnyTime{}, database.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fi := NewFileInfoInfrastructure()

	result, err := fi.Create(db, file)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(entity.FileInfo{}, "ID", "CreatedAt", "UpdatedAt"),
	}

	if diff := cmp.Diff(file, result, opts...); diff != "" {
		t.Error(diff)
	}

	if result.ID == 0 {
		t.Error("failed to insert id automatically")
	}

	if result.CreatedAt.IsZero() {
		t.Error("failed to insert created_at automatically")
	}

	if result.UpdatedAt.IsZero() {
		t.Error("failed to insert updated_at automatically")
	}
}

func TestCreateFiles(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	file, err := entity.NewFileInfo(1, "name", "/path/", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	files := []entity.FileInfo{*file}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `files` (`folder_id`,`name`,`path`,`mime_type`,`is_hide`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")).WithArgs(file.FolderID, file.Name.Value, file.Path.Value, file.MimeType.Value, file.IsHide, database.AnyTime{}, database.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fi := NewFileInfoInfrastructure()

	results, err := fi.Creates(db, files)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(entity.FileInfo{}, "ID", "CreatedAt", "UpdatedAt"),
	}

	if diff := cmp.Diff(files, results, opts...); diff != "" {
		t.Error(diff)
	}

	for _, result := range results {
		if result.ID == 0 {
			t.Error("failed to insert id automatically")
		}

		if result.CreatedAt.IsZero() {
			t.Error("failed to insert created_at automatically")
		}

		if result.UpdatedAt.IsZero() {
			t.Error("failed to insert updated_at automatically")
		}
	}
}

func TestUpdateFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	file, err := entity.NewFileInfo(1, "name", "/path/", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	file.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `files` SET `folder_id`=?,`name`=?,`path`=?,`mime_type`=?,`is_hide`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")).WithArgs(file.FolderID, file.Name.Value, file.Path.Value, file.MimeType.Value, file.IsHide, database.AnyTime{}, database.AnyTime{}, file.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fi := NewFileInfoInfrastructure()

	result, err := fi.Update(db, file)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(entity.FileInfo{}, "UpdatedAt"),
	}

	if diff := cmp.Diff(file, result, opts...); diff != "" {
		t.Error(diff)
	}

	if result.UpdatedAt.IsZero() {
		t.Error("failed to insert updated_at automatically")
	}
}

func TestRemoveFile(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	file, err := entity.NewFileInfo(1, "name", "/path/", "mime/type", false)
	if err != nil {
		t.Error(err.Error())
	}
	file.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `files` WHERE `files`.`id` = ?")).WithArgs(file.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fi := NewFileInfoInfrastructure()

	err = fi.Remove(db, file)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}
}

func TestFindOneFileByID(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE id = ? ORDER BY `files`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "folder_id", "name", "path", "mime_type", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", "mime/type", false, time.Now(), time.Now()))

	fi := NewFileInfoInfrastructure()

	result, err := fi.FindOneByID(db, 1)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find the file by id")
	}
}

func TestFindOneFileByIDAndIsHide(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE id = ? and is_hide = ? ORDER BY `files`.`id` LIMIT ?")).WithArgs(1, true, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "folder_id", "name", "path", "mime_type", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", "mime/type", false, time.Now(), time.Now()))

	fi := NewFileInfoInfrastructure()

	result, err := fi.FindOneByIDAndIsHide(db, 1, true)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find the file by id and is_hide")
	}
}

func TestFindOneFileByPath(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE path = ? ORDER BY `files`.`id` LIMIT ?")).WithArgs("/path/", 1).WillReturnRows(sqlmock.NewRows([]string{"id", "folder_id", "name", "path", "mime_type", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", "mime/type", false, time.Now(), time.Now()))

	fi := NewFileInfoInfrastructure()

	result, err := fi.FindOneByPath(db, "/path/")
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find the file by path")
	}
}
