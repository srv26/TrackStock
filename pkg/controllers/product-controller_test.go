package controllers

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	z "github.com/sourav/TrackStock/pkg/database/sql"
	"github.com/sourav/TrackStock/pkg/models"
	"github.com/stretchr/testify/assert"
)

var u = &models.Product{
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

func TestGetItem(t *testing.T) {
	r, err := http.NewRequest("GET", "/product/1", nil)
	if err != nil {
		t.Fatalf("cound not sent %v", err)
	}
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	query := "SELECT ProductId, PrdtName, Stock FROM PRODUCT WHERE ProductId = @id"

	rows := sqlmock.NewRows([]string{"ProductId", "PrdtName", "Stock"}).
		AddRow(u.ID, u.Name, u.Stock)

	handler := &Handler{}
	dbStorage, mock := NewMock()
	sqldatabase := &z.Storage{Db: dbStorage}
	handler.DbStorage = sqldatabase

	mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)
	handler.GetItemById(w, r)

	b, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error in reading %v", err)
	}
	f := &Response{}
	json.Unmarshal([]byte(b), f)
	if f.Success != true {
		t.Errorf("expected status true  but got %v", f.Success)
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, f.Success, true)
}

func TestUpdateItem(t *testing.T) {

	handler := &Handler{}
	dbStorage, mock := NewMock()
	sqldatabase := &z.Storage{Db: dbStorage}
	handler.DbStorage = sqldatabase
	jstr := []byte(`{"id":1,"name":"Soap","stock":3}`)
	r, err := http.NewRequest("PUT", "/product/", bytes.NewBuffer(jstr))
	if err != nil {
		t.Fatalf("cound not sent %v", err)
	}
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	query := "UPDATE PRODUCT SET ProductId = @id, PrdtName = @name ,Stock = @stock  WHERE ProductId = @id"

	mock.ExpectExec(query).WithArgs(u.ID, u.Name, u.Stock, u.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	handler.UpdateItem(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteItemById(t *testing.T) {
	handler := &Handler{}
	dbStorage, mock := NewMock()
	sqldatabase := &z.Storage{Db: dbStorage}
	handler.DbStorage = sqldatabase
	r, err := http.NewRequest("DELETE", "/product/1", nil)
	if err != nil {
		t.Fatalf("cound not sent %v", err)
	}
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()
	query := "DELETE FROM PRODUCT WHERE ProductId = @id"

	mock.ExpectExec(query).WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	handler.DeleteItemById(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteItemByIdError(t *testing.T) {
	handler := &Handler{}
	dbStorage, mock := NewMock()
	sqldatabase := &z.Storage{Db: dbStorage}
	handler.DbStorage = sqldatabase
	r, err := http.NewRequest("DELETE", "/product/1", nil)
	if err != nil {
		t.Fatalf("cound not sent %v", err)
	}
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	w := httptest.NewRecorder()
	query := "DELETE FROM PRODUCT WHERE ProductId = @id"

	mock.ExpectExec(query).WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	handler.DeleteItemById(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
