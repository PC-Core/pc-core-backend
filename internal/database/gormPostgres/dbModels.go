package gormpostgres

import (
	"time"

	"github.com/PC-Core/pc-core-backend/pkg/models"
)

type DbCart struct {
	ID        uint64         `gorm:"primaryKey"`
	UserID    uint64         `gorm:"column:user_id"`
	ProductID uint64         `gorm:"column:product_id"`
	Product   models.Product `gorm:"foreignKey:ProductID"`
	Quantity  uint           `gorm:"not null"`
	AddedAt   time.Time      `gorm:"autoCreateTime"`
}

func (DbCart) TableName() string {
	return "Cart"
}

func DbCartIntoCart(cart []DbCart) *models.Cart {
	items := make([]models.CartItem, 0, len(cart))
	var user_id uint64

	for _, c := range cart {
		user_id = c.UserID
		items = append(items, *models.NewCartItem(c.Product, c.Quantity, c.AddedAt))
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
	return "Products"
}

type DbCategories struct {
	ID          uint64 `gorm:"primaryKey"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	Icon        string `gorm:"column:icon"`
	Slug        string `gorm:"column:slug"`
}

func (DbCategories) TableName() string {
	return "Categories"
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
	return "CpuChars"
}

// type MediasType struct {
// 	ID  uint64 `json:"id"`
// 	Url string `json:"url"`
// }

type DbProductWithMedias struct {
	ID             uint64         `gorm:"column:id"`
	Name           string         `gorm:"column:name"`
	Price          float64        `gorm:"column:price"`
	Selled         uint64         `gorm:"column:selled"`
	Stock          uint64         `gorm:"column:stock"`
	CharsTableName string         `gorm:"column:chars_table_name"`
	CharsID        uint64         `gorm:"column:chars_id"`
	Medias         []models.Media `gorm:"column:medias;type:json"`
}

func (p *DbProductWithMedias) IntoProduct() *models.Product {
	return models.NewProduct(
		p.ID,
		p.Name,
		p.Price,
		p.Selled,
		p.Stock,
		p.Medias,
		p.CharsTableName,
		p.CharsID)

}

type DbProduct struct {
	ID             uint64         `gorm:"primaryKey"`
	Name           string         `gorm:"column:name"`
	Price          float64        `gorm:"column:price"`
	Selled         uint64         `gorm:"column:selled"`
	Stock          uint64         `gorm:"column:stock"`
	CharsTableName string         `gorm:"column:chars_table_name"`
	CharsID        uint64         `gorm:"column:chars_id"`
	Medias         [] `gorm:"column:medias;type:json"`
}
