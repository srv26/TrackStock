package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sourav/TrackStock/pkg/models"
)

type Storage struct {
	Db *sql.DB
}

// SetUp the sql Conection
func SetupSqlStorage() (models.Repository, error) {
	db, err := sql.Open("sqlserver", "server=LAPTOP-NQDHTU17\\SQLEXPRESS; database=Battleground;")
	err1 := db.Ping()
	if err1 != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	log.Print("Database conected")
	return &Storage{Db: db}, nil
}

// Close the connection
func (r *Storage) Close() {
	r.Db.Close()
}

// Adds a given item in the database
func (storage *Storage) AddItem(prod *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO PRODUCT (ProductId, PrdtName, Stock) VALUES (@id, @name, @stock)"
	_, err1 := storage.Db.ExecContext(ctx,
		query,
		sql.Named("id", prod.ID),
		sql.Named("name", prod.Name),
		sql.Named("stock", prod.Stock))
	if err1 != nil {
		fmt.Println("loll", err1)
		return err1
	}
	return nil
}

// Gets an item from the database of the given id, return error if item doesnt exist
func (storage *Storage) GetItemById(Id int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	prod := new(models.Product)
	rows, err := storage.Db.QueryContext(ctx,
		`SELECT ProductId, PrdtName, Stock FROM PRODUCT WHERE ProductId = @id `,
		sql.Named("id", Id))
	if err != nil {
		return &models.Product{}, err
	}
	if !rows.Next() {
		return &models.Product{}, models.ErrProductNotFound
	}
	err1 := rows.Scan(&prod.ID, &prod.Name, &prod.Stock)
	if err1 != nil {
		return &models.Product{}, err1
	}
	return prod, nil
}

// Update the existing item in the database else return error if item doesnt exist
func (storage *Storage) UpdateItem(prod *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row, err := storage.Db.ExecContext(ctx,
		`UPDATE PRODUCT SET ProductId = @id, PrdtName = @name ,Stock = @stock  WHERE ProductId = @id`,
		sql.Named("id", prod.ID),
		sql.Named("name", prod.Name),
		sql.Named("stock", prod.Stock),
		sql.Named("inputId", prod.ID))
	if err != nil {
		return err
	}
	count, _ := row.RowsAffected()
	if count == 0 {
		return models.ErrProductNotFound
	}
	return nil
}

// Delete a given item with the context, return error if item not present
func (storage *Storage) DeleteItemById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row, err := storage.Db.ExecContext(ctx,
		`DELETE FROM PRODUCT WHERE ProductId = @id`,
		sql.Named("id", id))
	if err != nil {
		return err
	}
	count, _ := row.RowsAffected()
	if count == 0 {
		return models.ErrProductNotFound
	}
	return nil
}
