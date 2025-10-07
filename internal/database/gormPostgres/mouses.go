package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"gorm.io/gorm"
)

func (c *GormPostgresController) GetMouseChars() ([]models.MouseChars, errors.PCCError){
	var mouses []models.MouseChars

	err := c.db.Model(&models.MouseChars{}).Find(&mouses).Error
	if err != nil { 
		return nil, gormerrors.GormErrorCast(err)
	}

	return mouses, nil
}

func (c *GormPostgresController) GetMouseByID(id uint64) (*models.MouseChars, errors.PCCError){
	var mouse DbMouseChars

	err := c.db.Model(&DbMouseChars{}).Where("id = ?", id).First(&mouse).Error
	if err != nil{
		if err == gorm.ErrRecordNotFound{
			return nil, nil
		}
		return nil, gormerrors.GormErrorCast(err)
	}

	return mouse.IntoMouse(), nil
}

func (c *GormPostgresController) AddMouse(mouse *inputs.AddMouseInput) (*models.MouseChars, *models.Product, errors.PCCError){
	tx := c.db.Begin()

	if tx.Error != nil { 
		return nil, nil, errors.NewInternalSecretError()
	}

	defer tx.Rollback()

	chars := DbMouseChars{
		ID: mouse.ID,
		Name: mouse.Name,
		TypeMouses: mouse.TypeMouses,
		Dpi: mouse.Dpi,
		ReleaseYear: mouse.ReleaseYear,
	}

	err := tx.Create(&chars).Error

	if err != nil { 
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	medias, err := c.AddMedias(tx, mouse.Medias)

	if err != nil { 
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	product := DbProduct{
		Name: mouse.Name,
		Price: mouse.Price,
		Selled: 0,
		Stock: mouse.Stock,
		CharsTableName: database.MouseCharsTable,
		CharsID: chars.ID,
	}

	err = tx.Create(&product).Error

	if err != nil { 
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	tx.Commit()

	return chars.IntoMouse(), product.WithMediasIntoProduct(medias), nil
}