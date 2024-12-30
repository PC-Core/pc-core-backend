package database

import (
	"database/sql"
)

const (
	LaptopCharsTable = "LaptopChars"
)

// An interface for Product Characterics
type ProductChars any

// Database controller
type DbController struct {
	db *sql.DB
}

// Creates a new Database Controller
//
// Params:
//
//	`driver` - the name of the database driver
//	`conn` - the connection string to the database
//
// Returns:
//
//	`*DbController` or `error`
func NewDbController(driver string, conn string) (*DbController, error) {
	db, err := sql.Open(driver, conn)

	if err != nil {
		return nil, err
	}

	return &DbController{
		db,
	}, nil
}

// func (c *DbController) AddLaptop(name string, cpu string, ram int16, gpu string, price string, discount int16) (*models.LaptopChars, error) {
// 	var id int

// 	err := c.db.QueryRow("INSERT INTO laptops (Name, Cpu, Ram, Gpu, Price, Discount) VALUES ($1, $2, $3, $4, $5, $6) returning id", name, cpu, ram, gpu, price, discount).Scan(&id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	bfPrice, err := helpers.StringToBigFloat(price)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return models.NewLaptop(id, name, cpu, ram, gpu, *bfPrice, discount), nil
// }

// func (c *DbController) RemoveLaptop(id int) (error) {
// 	_, err := c.db.Exec("DELETE FROM laptops WHERE id = $1", id)

// 	if err != nil {
// 		return err;
// 	}

// 	return nil
// }
