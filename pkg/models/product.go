package models

import (
	"fmt"
)

// ErrProductNotFound is an error raised when a product cannot be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Service which each database implements
type Repository interface {
	AddItem(prod *Product) error
	GetItemById(Id int) (*Product, error)
	UpdateItem(prod *Product) error
	DeleteItemById(id int) error
}

type Product struct {
	ID    int    `json:"id"`    // unique id of the item
	Name  string `json:"name"`  // Name of the item
	Stock int    `json:"stock"` // count of the item
}
