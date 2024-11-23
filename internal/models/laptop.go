package models

type Laptop struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Cpu   string  `json:"cpu"`
	Ram   int16   `json:"ram"`
	Gpu   string  `json:"gpu"`
	Price float64 `json:"price"`
}

func NewLaptop(id int, name string, cpu string, ram int16, gpu string, price float64) *Laptop {
	return &Laptop{
		id, name, cpu, ram, gpu, price,
	}
}
