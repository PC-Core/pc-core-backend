package database

import (
	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
)

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