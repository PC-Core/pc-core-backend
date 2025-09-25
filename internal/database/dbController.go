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
	GpuCharsTable    = "GpuChars"
)

// An interface for Product Characterics
type ProductChars any

type RowLike interface {
	Scan(...any) error
}

type DbController interface {
	GetCartByUserID(userID uint64) (*models.Cart, errors.PCCError)
	AddToCart(product_id, user_id, quantity uint64) (uint64, errors.PCCError)
	SetToCart(product_id, user_id, quantity uint64) (uint64, errors.PCCError)
	RemoveFromCart(product_id, user_id uint64) (uint64, errors.PCCError)
	ChangeQuantity(product_id, user_id uint64, val int64) (uint64, errors.PCCError)
	GetCategories() ([]models.Category, errors.PCCError)
	GetLaptopChars(charId uint64) (*models.LaptopChars, errors.PCCError)
	AddLaptop(laptop *inputs.AddLaptopInput) (*models.Product, *models.LaptopChars, errors.PCCError)
	GetProducts(start uint64, count uint64) ([]models.Product, uint64, errors.PCCError)
	GetProductCharsByProductID(productId uint64) (ProductChars, errors.PCCError)
	GetProductById(id uint64) (*models.Product, errors.PCCError)
	LoadProductsRangeAsCartItem(tempCart []models.TempCartItem) ([]models.CartItem, errors.PCCError)
	RegisterUser(register *inputs.RegisterUserInput) (*models.User, errors.PCCError)
	LoginUser(login *inputs.LoginUserInput) (*models.User, errors.PCCError)
	GetUserByID(id int) (*models.User, errors.PCCError)
	GetCpuChars(charId uint64) (*models.CpuChars, errors.PCCError)
	AddCpu(cpu *inputs.AddCpuInput) (*models.Product, *models.CpuChars, errors.PCCError)
	GetRootCommentsForProduct(product_id int64, userID *int64) ([]models.Comment, errors.PCCError)
	GetAnswersOnComment(product_id int64, userID *int64, comment_id int64) ([]models.Comment, errors.PCCError)
	AddComment(input *inputs.AddCommentInput, userID int64, product_id int64) (int64, errors.PCCError)
	EditComment(newText string, commentID int64, userID int64) (int64, errors.PCCError)
	DeleteComment(commentID int64, userID int64) (int64, errors.PCCError)
	AddGpu(gpu *inputs.AddGpuInput) (*models.GpuChars, *models.Product, errors.PCCError)
	SetReaction(commentID int64, userID int64, ty models.ReactionType) (int64, errors.PCCError)
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
