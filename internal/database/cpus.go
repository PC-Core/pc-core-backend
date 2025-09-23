package database

import (
	"fmt"

	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
)

func (c *DPostgresDbController) GetCpuChars(charId uint64) (*models.CpuChars, errors.PCCError) {
	var chars models.CpuChars

	row := c.db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", CpuCharsTable), charId)

	err := row.Scan(&chars.ID, &chars.Name, &chars.PCores, &chars.ECores, &chars.Threads, &chars.BasePFreqMHz, &chars.MaxPFreqMHz, &chars.BaseEFreqMHz, &chars.MaxEFreqMHz, &chars.Socket, &chars.L1KB, &chars.L2KB, &chars.L3KB, &chars.TecProcNM, &chars.TDPWatt, &chars.ReleaseYear)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return &chars, nil
}

func (c *DPostgresDbController) AddCpu(cpu *inputs.AddCpuInput) (*models.Product, *models.CpuChars, errors.PCCError) {
	var (
		charId    uint64
		productId uint64
	)

	tx, err := c.db.Begin()

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer tx.Rollback()

	err = tx.QueryRow(fmt.Sprintf("INSERT INTO %s (name, pcores, ecores, threads, base_p_freq_mhz, max_p_freq_mhz, base_e_freq_mhz, max_e_freq_mhz, socket, l1_kb, l2_kb, l3_kb, tecproc_nm, tdp_watt, release_year) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning id", CpuCharsTable), cpu.CpuName, cpu.PCores, cpu.ECores, cpu.Threads, cpu.BasePFreqMHz, cpu.MaxPFreqMHz, cpu.BaseEFreqMHz, cpu.MaxEFreqMHz, cpu.Socket, cpu.L1KB, cpu.L2KB, cpu.L3KB, cpu.TecProcNM, cpu.TDPWatt, cpu.ReleaseYear).Scan(&charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	productId, medias, err := c.AddProduct(tx, cpu.Name, cpu.Price, 0, cpu.Stock, cpu.Medias, CpuCharsTable, charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewProduct(productId, cpu.Name, cpu.Price, 0, cpu.Stock, medias, CpuCharsTable, charId),
		models.NewCpuChars(charId, cpu.CpuName, cpu.PCores, cpu.ECores, cpu.Threads, cpu.BasePFreqMHz, cpu.MaxPFreqMHz, cpu.BaseEFreqMHz, cpu.MaxEFreqMHz, cpu.Socket, cpu.L1KB, cpu.L2KB, cpu.L3KB, cpu.TecProcNM, cpu.TDPWatt, cpu.ReleaseYear),
		nil
}

// ----------------------GPU----------------------------------------
func (c *DPostgresDbController) GetGpuChars() ([]models.Gpu, errors.PCCError) {
	gpus := make([]models.Gpu, 0)

	res, err := c.db.Query(`
		SELECT id, name, memory_gb, memory_type, bus_width_bit, 
		base_freq_mhz, boost_freq_mhz, tecproc_nm, tdp_watt, release_year 
		FROM GpuChars
	`)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer res.Close()

	for res.Next() {
		var gpu models.Gpu

		if err := res.Scan(&gpu.ID, &gpu.Name, &gpu.MemoryGB, &gpu.MemoryType,
			&gpu.BusWidthBit, &gpu.BaseFreqMHz, &gpu.BoostFreqMHz,
			&gpu.TecprocNm, &gpu.TDPWatt, &gpu.RealeseYear); err != nil {
			return nil, dberrors.PQDbErrorCaster(c.db, err)
		}

		gpus = append(gpus, gpu)
	}

	return gpus, nil
}

func (c *DPostgresDbController) GetGpuByID(id uint64) (*models.Gpu, errors.PCCError){
	var gpu models.Gpu

	err := c.db.QueryRow(`SELECT id, name, memory_gb, memory_type, bus_width_bit, 
		base_freq_mhz, boost_freq_mhz, tecproc_nm, tdp_watt, release_year 
		FROM GpuChars WHERE id = $1`, id).Scan(&gpu.ID, &gpu.Name, &gpu.MemoryGB, &gpu.MemoryType, 
		&gpu.BusWidthBit, &gpu.BaseFreqMHz, &gpu.BoostFreqMHz, 
		&gpu.TecprocNm, &gpu.TDPWatt, &gpu.RealeseYear)

		if err != nil{ 
			return nil, dberrors.PQDbErrorCaster(c.db, err)
		}

		return &gpu, nil
}