package models

type LaptopChars struct {
	ID  uint64 `json:"id"`
	Cpu string `json:"cpu"`
	Ram int16  `json:"ram"`
	Gpu string `json:"gpu"`
}

func NewLaptopChars(id uint64, cpu string, ram int16, gpu string) *LaptopChars {
	return &LaptopChars{
		id, cpu, ram, gpu,
	}
}
