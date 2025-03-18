package database

import (
	"fmt"

	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models"
)

func (c *DPostgresDbController) GetCpuChars(charId uint64) (*models.CpuChars, errors.PCCError) {
	var chars models.CpuChars

	row := c.db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", CpuCharsTable), charId)

	err := row.Scan(&chars.ID, &chars.Name, &chars.PCores, &chars.ECores, &chars.Threads, &chars.BaseFreqMHz, &chars.MaxFreqMHz, &chars.Socket, &chars.L1KB, &chars.L2KB, &chars.L3KB, &chars.TecProcNM, &chars.TDPWatt, &chars.ReleaseYear)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return &chars, nil
}

func (c *DPostgresDbController) AddCpu(name string, price float64, selled uint64, stock uint64, pcores, ecores, threads, bfmhz, mfmhz uint64, socket models.CpuSocket, l1, l2, l3, tpnm, tdp, ry uint64, imedias []models.InputMedia) (*models.Product, *models.CpuChars, errors.PCCError) {
	var (
		charId    uint64
		productId uint64
	)

	tx, err := c.db.Begin()

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer tx.Rollback()

	err = tx.QueryRow(fmt.Sprintf("INSERT INTO %s (name, pcores, ecores, threads, base_freq_mhz, max_freq_mhz, socket, l1_kb, l2_kb, l3_kb, tecproc_nm, tdp_watt, release_year) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning id", CpuCharsTable), name, pcores, ecores, threads, bfmhz, mfmhz, socket, l1, l2, l3, tpnm, tdp, ry).Scan(&charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	productId, medias, err := c.AddProduct(tx, name, price, selled, stock, imedias, CpuCharsTable, charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewProduct(productId, name, price, selled, stock, medias, CpuCharsTable, charId),
		models.NewCpuChars(charId, name, pcores, ecores, threads, bfmhz, mfmhz, socket, l1, l2, l3, tpnm, tdp, ry),
		nil
}
