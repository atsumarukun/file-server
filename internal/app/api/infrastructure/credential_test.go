package infrastructure

import (
	"file-server/test/database"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFindOneCredential(t *testing.T) {
	db, mock, err := database.Open()
	if err != nil {
		t.Error(err.Error())
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `credentials` ORDER BY `credentials`.`id` LIMIT ?")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "password", "created_at", "updated_at"}).AddRow(1, "password", time.Now(), time.Now()))

	ci := NewCredentialInfrastructure()
	result, err := ci.FindOne(db)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}

	if result == nil {
		t.Error("failed to find the credential")
	}
}
