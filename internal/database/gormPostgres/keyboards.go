package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"gorm.io/gorm"
)

func (c *GormPostgresController) GetKeyBoardChars() ([]models.KeyboardChars, errors.PCCError){
	var keyboards []models.KeyboardChars

	err := c.db.Model(&models.KeyboardChars{}).Find(&keyboards).Error
	if err != nil { 
		return nil, gormerrors.GormErrorCast(err)
	}

	return keyboards, nil
}

func (c *GormPostgresController) GetKeyBoardByID(id uint64) (*models.KeyboardChars, errors.PCCError){
	var keyboard DbKeyboardChars

	err := c.db.Model(&DbKeyboardChars{}).Where("id = ?", id).First(&keyboard).Error
	if err != nil{
		if err == gorm.ErrRecordNotFound{
			return nil, nil
		}
		return nil, gormerrors.GormErrorCast(err)
	}

	return keyboard.IntoKeyBoard(), nil
}

func (c *GormPostgresController) AddKeyBoard(keyboard *inputs.AddKeyBoardInput) (*models.KeyboardChars, *models.Product, errors.PCCError){
	tx := c.db.Begin()

	if tx.Error != nil { 
		return nil, nil, errors.NewInternalSecretError()
	}

	defer tx.Rollback()

	chars := DbKeyboardChars{
		ID: keyboard.ID,
		Name: keyboard.Name,
		TypeKeyBoards: keyboard.TypeKeyBoards,
		Switches: keyboard.Switches,
		ReleaseYear: keyboard.ReleaseYear,
	}

	err := tx.Create(&chars).Error

	if err != nil { 
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	medias, err := c.AddMedias(tx, keyboard.Medias)

	if err != nil { 
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	product := DbProduct{
		Name: keyboard.Name,
		Price: keyboard.Price,
		Selled: 0,
		Stock: keyboard.Stock,
		CharsTableName: database.KeyboardCharsTable,
		CharsID: chars.ID,
	}

	err = tx.Create(&product).Error

	if err != nil { 
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	tx.Commit()

	return chars.IntoKeyBoard(), product.WithMediasIntoProduct(medias), nil
}