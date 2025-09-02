package gormpostgres

import (
	"time"

	"github.com/PC-Core/pc-core-backend/pkg/models"
)

type DbCart struct {
	ID        uint64              `gorm:"primaryKey"`
	UserID    uint64              `gorm:"column:user_id"`
	ProductID uint64              `gorm:"column:product_id"`
	Product   DbProductWithMedias `gorm:"foreignKey:ProductID"`
	Quantity  uint                `gorm:"not null"`
	AddedAt   time.Time           `gorm:"autoCreateTime"`
}

func (DbCart) TableName() string {
	return "cart"
}

func DbCartIntoCart(cart []DbCart) *models.Cart {
	items := make([]models.CartItem, 0, len(cart))
	var user_id uint64

	for _, c := range cart {
		user_id = c.UserID
		items = append(items, *models.NewCartItem(*c.Product.IntoProduct(), c.Quantity, c.AddedAt))
	}

	return models.NewCart(user_id, items)
}

type DbProducts struct {
	ID            uint64   `gorm:"primaryKey"`
	Name          string   `gorm:"column:name"`
	Price         float64  `gorm:"column:price"`
	Selled        int      `gorm:"column:selled"`
	Stock         int      `gorm:"column:stock"`
	CharTableName string   `gorm:"column:chars_table_name"`
	CharId        uint64   `gorm:"column:chars_id"`
	Medias        []uint64 `gorm:"column:medias"`
}

func (DbProducts) TableName() string {
	return "products"
}

type DbCategories struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	Icon        string `gorm:"column:icon"`
	Slug        string `gorm:"column:slug"`
}

func (DbCategories) TableName() string {
	return "categories"
}

type DbCpuChars struct {
	ID           uint64           `gorm:"primaryKey"`
	Name         string           `gorm:"column:name"`
	PCores       uint64           `gorm:"column:pcores"`
	ECores       uint64           `gorm:"column:ecores"`
	Threads      uint64           `gorm:"column:threads"`
	BasePFreqMHz uint64           `gorm:"column:base_p_freq_mhz"`
	MaxPFreqMHz  uint64           `gorm:"column:max_p_freq_mhz"`
	BaseEFreqMHz uint64           `gorm:"column:base_e_freq_mhz"`
	MaxEFreqMHz  uint64           `gorm:"column:max_e_freq_mhz"`
	Socket       models.CpuSocket `gorm:"column:socket"`
	L1KB         uint64           `gorm:"column:l1_kb"`
	L2KB         uint64           `gorm:"column:l2_kb"`
	L3KB         uint64           `gorm:"column:l3_kb"`
	TecProcNM    uint64           `gorm:"column:tecproc_nm"`
	TDPWatt      uint64           `gorm:"column:tdp_watt"`
	ReleaseYear  uint64           `gorm:"column:release_year"`
}

func (chars *DbCpuChars) IntoCpuChars() *models.CpuChars {
	return models.NewCpuChars(
		chars.ID,
		chars.Name,
		chars.PCores,
		chars.ECores,
		chars.Threads,
		chars.BasePFreqMHz,
		chars.MaxPFreqMHz,
		chars.BaseEFreqMHz,
		chars.MaxEFreqMHz,
		chars.Socket,
		chars.L1KB,
		chars.L2KB,
		chars.L3KB,
		chars.TecProcNM,
		chars.TDPWatt,
		chars.ReleaseYear,
	)
}

func (DbCpuChars) TableName() string {
	return "cpuchars"
}

type DbMedia struct {
	ID        uint64           `gorm:"primaryKey"`
	Url       string           `gorm:"column:url"`
	Type      models.MediaType `gorm:"column:type"`
	ProductID uint64           `gorm:"column:product_id"`
}

func (DbMedia) TableName() string {
	return "medias"
}

func (m *DbMedia) IntoMedia() *models.Media {
	return models.NewMedia(
		m.ID,
		m.Url,
		m.Type,
		m.ProductID,
	)
}

type DbMedias []DbMedia

func (ms DbMedias) IntoMedias() models.Medias {
	medias := make(models.Medias, 0, len(ms))

	for _, m := range ms {
		medias = append(medias, *m.IntoMedia())
	}

	return medias
}

type DbProductWithMedias struct {
	ID             uint64   `gorm:"column:id"`
	Name           string   `gorm:"column:name"`
	Price          float64  `gorm:"column:price"`
	Selled         uint64   `gorm:"column:selled"`
	Stock          uint64   `gorm:"column:stock"`
	CharsTableName string   `gorm:"column:chars_table_name"`
	CharsID        uint64   `gorm:"column:chars_id"`
	Medias         DbMedias `gorm:"foreignKey:ProductID"`
}

func (p *DbProductWithMedias) IntoProduct() *models.Product {
	return models.NewProduct(
		p.ID,
		p.Name,
		p.Price,
		p.Selled,
		p.Stock,
		p.Medias.IntoMedias(),
		p.CharsTableName,
		p.CharsID,
	)
}

func (DbProductWithMedias) TableName() string {
	return "products"
}

type DbProduct struct {
	ID             uint64  `gorm:"primaryKey"`
	Name           string  `gorm:"column:name"`
	Price          float64 `gorm:"column:price"`
	Selled         uint64  `gorm:"column:selled"`
	Stock          uint64  `gorm:"column:stock"`
	CharsTableName string  `gorm:"column:chars_table_name"`
	CharsID        uint64  `gorm:"column:chars_id"`
}

func (DbProduct) TableName() string {
	return "products"
}

func (p *DbProduct) WithMediasIntoProduct(medias models.Medias) *models.Product {
	return models.NewProduct(
		p.ID,
		p.Name,
		p.Price,
		p.Selled,
		p.Stock,
		medias,
		p.CharsTableName,
		p.CharsID,
	)
}

type DbLaptopChars struct {
	ID    uint64     `gorm:"primaryKey"`
	CpuID uint64     `gorm:"column:cpu_id"`
	Cpu   DbCpuChars `gorm:"foreignKey:CpuID"`
	Ram   int16      `gorm:"column:ram"`
	Gpu   string     `gorm:"column:gpu"`
}

func (DbLaptopChars) TableName() string {
	return "laptopchars"
}

func (c *DbLaptopChars) IntoLaptopChars() *models.LaptopChars {
	return models.NewLaptopChars(
		c.ID,
		c.Cpu.IntoCpuChars(),
		c.Ram,
		c.Gpu,
	)
}

type DbUser struct {
	ID           int             `gorm:"primaryKey"`
	Name         string          `gorm:"column:name"`
	Email        string          `gorm:"column:email"`
	Role         models.UserRole `gorm:"column:role;default:'Default'"`
	PasswordHash string          `gorm:"column:passwordhash"`
}

func (DbUser) TableName() string {
	return "users"
}

func (u *DbUser) IntoUser() *models.User {
	return models.NewUser(u.ID, u.Name, u.Email, u.Role, u.PasswordHash)
}

type DbComment struct {
	ID          int64      `gorm:"column:id;primaryKey"`
	UserID      int64      `gorm:"column:user_id"`
	ProductID   int64      `gorm:"column:product_id"`
	CommentText string     `gorm:"column:comment_text"`
	AnswerOn    *int64     `gorm:"column:answer_on"`
	Rating      int16      `gorm:"column:rating"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
	MediaIDs    []int64    `gorm:"column:media_ids"`

	User    DbUser    `gorm:"foreignKey:UserID;references:ID"`
	Product DbProduct `gorm:"foreignKey:ProductID;references:ID"`
}

type DbCommentReaction struct {
	UserID    int64               `gorm:"column:user_id"`
	CommentID int64               `gorm:"column:comment_id"`
	Type      models.ReactionType `gorm:"column:ty"`
	AddedAt   time.Time           `gorm:"column:added_at"`
}
