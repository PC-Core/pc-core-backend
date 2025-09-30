package models

type LaptopChars struct {
	ID  uint64    `json:"id"`
	Cpu *CpuChars `json:"cpu"`
	Ram int16     `json:"ram"`
	Gpu *GpuChars `json:"gpu"`
}

func NewLaptopChars(id uint64, cpu *CpuChars, ram int16, gpu *GpuChars) *LaptopChars {
	return &LaptopChars{
		id, cpu, ram, gpu,
	}
}
