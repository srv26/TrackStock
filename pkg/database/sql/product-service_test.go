package sql

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/denisenkom/go-mssqldb"
	r "github.com/sourav/TrackStock/pkg/models"
	"github.com/stretchr/testify/assert"
)

var u = &r.Product{
	ID:    1,
	Name:  "Soap",
	Stock: 3,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}
func TestGetItemById(t *testing.T) {
	db, mock := NewMock()
	repo := &Storage{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT ProductId, PrdtName, Stock FROM PRODUCT WHERE ProductId = @id"

	rows := sqlmock.NewRows([]string{"ProductId", "PrdtName", "Stock"}).
		AddRow(u.ID, u.Name, u.Stock)

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	prdct, err := repo.GetItemById(u.ID)
	assert.NotNil(t, prdct)
	assert.NoError(t, err)
}

func TestGetItemByIDError(t *testing.T) {
	db, mock := NewMock()
	repo := &Storage{db}
	defer func() {
		repo.Close()
	}()

	query := "SELECT ProductId, PrdtName, Stock FROM PRODUCT WHERE ProductId = @id"

	rows := sqlmock.NewRows([]string{"ProductId", "PrdtName", "Stock"})

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

	user, err := repo.GetItemById(u.ID)
	assert.Empty(t, user)
	assert.Error(t, err)
}

func TestUpdateItemById(t *testing.T) {
	db, mock := NewMock()
	repo := &Storage{db}
	defer func() {
		repo.Close()
	}()

	query := "UPDATE PRODUCT SET ProductId = @id, PrdtName = @name ,Stock = @stock  WHERE ProductId = @id"

	mock.ExpectExec(query).WithArgs(u.ID, u.Name, 3, u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateItem(u)
	assert.NoError(t, err)
}

func TestDeleteItemById(t *testing.T) {
	db, mock := NewMock()
	repo := &Storage{db}
	defer func() {
		repo.Close()
	}()

	query := "DELETE FROM PRODUCT WHERE ProductId = @id"

	mock.ExpectExec(query).WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteItemById(u.ID)
	assert.NoError(t, err)
}
