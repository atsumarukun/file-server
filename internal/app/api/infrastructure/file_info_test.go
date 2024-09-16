package infrastructure

import (
	"file-server/internal/app/api/domain/entity"
	"file-server/test/database"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreate(t *testing.T) {
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

	if result.CreatedAt == database.NullTime {
		t.Error("failed to insert created_at automatically")
	}

	if result.UpdatedAt == database.NullTime {
		t.Error("failed to insert updated_at automatically")
	}
}
