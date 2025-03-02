package database

import (
	"database/sql"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models"
)

const (
	LaptopCharsTable = "LaptopChars"
)

// An interface for Product Characterics
type ProductChars any

type RowLike interface {
	Scan(...any) error
}

type DbController interface {
	GetCartByUserID(userID uint64) (*models.Cart, errors.PCCError)
	AddToCart(product_id, user_id, quantity uint64) (uint64, errors.PCCError)
	RemoveFromCart(product_id, user_id uint64) (uint64, errors.PCCError)
	ChangeQuantity(product_id, user_id uint64, val int64) (uint64, errors.PCCError)
	GetCategories() ([]models.Category, errors.PCCError)
	GetLaptopChars(charId uint64) (*models.LaptopChars, errors.PCCError)
	AddLaptop(name string, price float64, selled uint64, stock uint64, cpu string, ram int16, gpu string, imedias []models.InputMedia) (*models.Product, *models.LaptopChars, errors.PCCError)
	GetProducts(start uint64, count uint64) ([]models.Product, errors.PCCError)
	ScanProduct(rows RowLike) (*models.Product, errors.PCCError)
	GetProductCharsByProductID(productId uint64) (ProductChars, errors.PCCError)
	GetProductById(id uint64) (*models.Product, errors.PCCError)
	LoadProductsRangeAsCartItem(tempCart []models.TempCartItem) ([]models.CartItem, errors.PCCError)
	RegisterUser(name string, email string, password string) (*models.User, errors.PCCError)
	LoginUser(email string, password string) (*models.User, errors.PCCError)
	GetUserByID(id int) (*models.User, errors.PCCError)
	AddProduct(tx *sql.Tx, name string, price float64, selled uint64, stock uint64, medias []models.InputMedia, charTable string, charID uint64) (uint64, []models.Media, errors.PCCError)
	AddMedias(imedias []models.InputMedia) ([]models.Media, errors.PCCError)
}

// Database controller
type DPostgresDbController struct {
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
func NewDPostgresDbController(driver string, conn string) (*DPostgresDbController, error) {
	db, err := sql.Open(driver, conn)

	if err != nil {
		return nil, err
	}

	return &DPostgresDbController{
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
