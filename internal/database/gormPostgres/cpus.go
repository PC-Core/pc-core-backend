package gormpostgres

// func (c *GormPostgresController) GetCpuChars(charId uint64) (*models.CpuChars, errors.PCCError) {
// 	var chars DbCpuChars

// 	err := c.db.Where("id = ?", charId).First(&chars).Error

// 	if err != nil {
// 		// TODO: error type
// 		return nil, errors.NewInternalSecretError()
// 	}

// 	return chars.IntoCpuChars(), nil
// }

// func (c *GormPostgresController) AddCpu(cpu *inputs.AddCpuInput) (*models.Product, *models.CpuChars, errors.PCCError) {
// 	chars := DbCpuChars{
// 		Name:         cpu.CpuName,
// 		PCores:       cpu.PCores,
// 		ECores:       cpu.ECores,
// 		Threads:      cpu.Threads,
// 		BasePFreqMHz: cpu.BasePFreqMHz,
// 		MaxPFreqMHz:  cpu.MaxPFreqMHz,
// 		BaseEFreqMHz: cpu.BaseEFreqMHz,
// 		MaxEFreqMHz:  cpu.MaxEFreqMHz,
// 		Socket:       cpu.Socket,
// 		L1KB:         cpu.L1KB,
// 		L2KB:         cpu.L2KB,
// 		L3KB:         cpu.L3KB,
// 		TecProcNM:    cpu.TecProcNM,
// 		TDPWatt:      cpu.TDPWatt,
// 		ReleaseYear:  cpu.ReleaseYear,
// 	}

// 	err := c.db.Create(chars).Error

// 	if err != nil {
// 		// TODO: error type
// 		return nil, nil, errors.NewInternalSecretError()
// 	}

// 	// TODO:
// }
