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

func TestCreateFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	folder, err := entity.NewFolderInfo(nil, "name", "/path/", false)
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `folders` (`parent_folder_id`,`name`,`path`,`is_hide`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).WithArgs(folder.ParentFolderID, folder.Name.Value, folder.Path.Value, folder.IsHide, database.AnyTime{}, database.AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fi := NewFolderInfoInfrastructure()

	result, err := fi.Create(db, folder)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(entity.FolderInfo{}, "ID", "CreatedAt", "UpdatedAt"),
	}

	if diff := cmp.Diff(folder, result, opts...); diff != "" {
		t.Error(diff)
	}

	if result.ID == 0 {
		t.Error("failed to insert id automatically")
	}

	if result.CreatedAt == database.NullTime {
		t.Error("failed to insert created_at automatically")
	}

	if result.UpdatedAt == database.NullTime {
		t.Error("failed to insert updated_at automatically")
	}
}

func TestUpdateFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	folder, err := entity.NewFolderInfo(nil, "name", "/path/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folder.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `folders` SET `parent_folder_id`=?,`name`=?,`path`=?,`is_hide`=?,`created_at`=?,`updated_at`=? WHERE `id` = ?")).WithArgs(folder.ParentFolderID, folder.Name.Value, folder.Path.Value, folder.IsHide, database.AnyTime{}, database.AnyTime{}, folder.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fi := NewFolderInfoInfrastructure()

	result, err := fi.Update(db, folder)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(entity.FolderInfo{}, "UpdatedAt"),
	}

	if diff := cmp.Diff(folder, result, opts...); diff != "" {
		t.Error(diff)
	}

	if result.UpdatedAt == database.NullTime {
		t.Error("failed to insert updated_at automatically")
	}
}

func TestRemoveFolder(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	folder, err := entity.NewFolderInfo(nil, "name", "/path/", false)
	if err != nil {
		t.Error(err.Error())
	}
	folder.ID = 1

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `folders` WHERE `folders`.`id` = ?")).WithArgs(folder.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	fi := NewFolderInfoInfrastructure()

	err = fi.Remove(db, folder)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}
}

func TestFindOneFolderByID(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE id = ? ORDER BY `folders`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", false, time.Now(), time.Now()))

	fi := NewFolderInfoInfrastructure()

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

func TestFindOneFolderByPath(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE path = ? ORDER BY `folders`.`id` LIMIT ?")).WithArgs("/path/", 1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", false, time.Now(), time.Now()))

	fi := NewFolderInfoInfrastructure()

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

func TestFindOneFolderByPathWithChildren(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE path = ? ORDER BY `folders`.`id` LIMIT ?")).WithArgs("/path/", 1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", false, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE `files`.`folder_id` = ?")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "folder_id", "name", "path", "mime_type", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", "mime/type", false, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE `folders`.`parent_folder_id` = ?")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", false, time.Now(), time.Now()))

	fi := NewFolderInfoInfrastructure()

	result, err := fi.FindOneByPathWithChildren(db, "/path/")
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find the file by path with children")
	}
}

func TestFindOneFolderByPathAndIsHideWithChildren(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE path = ? and is_hide = ? ORDER BY `folders`.`id` LIMIT ?")).WithArgs("/path/", true, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", true, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE `files`.`folder_id` = ? AND `is_hide` = ?")).WithArgs(1, true).WillReturnRows(sqlmock.NewRows([]string{"id", "folder_id", "name", "path", "mime_type", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", "mime/type", true, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE `folders`.`parent_folder_id` = ? AND `is_hide` = ?")).WithArgs(1, true).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", true, time.Now(), time.Now()))

	fi := NewFolderInfoInfrastructure()

	result, err := fi.FindOneByPathAndIsHideWithChildren(db, "/path/", true)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find the file by path and is_hide with children")
	}
}

func TestFindOneFolderByIDWithLower(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE id = ? ORDER BY `folders`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", false, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE `files`.`folder_id` = ?")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "folder_id", "name", "path", "mime_type", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", "mime/type", false, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE `folders`.`parent_folder_id` = ?")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", false, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE id = ? ORDER BY `folders`.`id` LIMIT ?")).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "parent_folder_id", "name", "path", "is_hide", "created_at", "updated_at"}).AddRow(1, 1, "name", "/path/", false, time.Now(), time.Now()))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `files` WHERE `files`.`folder_id` = ?")).WithArgs(1).WillReturnRows(sqlmock.NewRows(nil))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `folders` WHERE `folders`.`parent_folder_id` = ?")).WithArgs(1).WillReturnRows(sqlmock.NewRows(nil))

	fi := NewFolderInfoInfrastructure()

	result, err := fi.FindOneByIDWithLower(db, 1)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find the file by id with lower")
	}
}
