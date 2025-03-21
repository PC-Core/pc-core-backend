package database

import (
	"database/sql"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
)

const (
	LaptopCharsTable = "LaptopChars"
	CpuCharsTable    = "CpuChars"
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
	AddLaptop(laptop *inputs.AddLaptopInput) (*models.Product, *models.LaptopChars, errors.PCCError)
	GetProducts(start uint64, count uint64) ([]models.Product, errors.PCCError)
	ScanProduct(rows RowLike) (*models.Product, errors.PCCError)
	GetProductCharsByProductID(productId uint64) (ProductChars, errors.PCCError)
	GetProductById(id uint64) (*models.Product, errors.PCCError)
	LoadProductsRangeAsCartItem(tempCart []models.TempCartItem) ([]models.CartItem, errors.PCCError)
	RegisterUser(register *inputs.RegisterUserInput) (*models.User, errors.PCCError)
	LoginUser(login *inputs.LoginUserInput) (*models.User, errors.PCCError)
	GetUserByID(id int) (*models.User, errors.PCCError)
	AddMedias(imedias []models.InputMedia) ([]models.Media, errors.PCCError)
	GetCpuChars(charId uint64) (*models.CpuChars, errors.PCCError)
	AddCpu(cpu *inputs.AddCpuInput) (*models.Product, *models.CpuChars, errors.PCCError)
}

// Database controller
// Deprecated. Use GormPostgresController instead
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
