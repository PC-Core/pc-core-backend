package inputs

type AddLaptopInput struct {
	Name     string `json:"name"`
	Cpu      string `json:"cpu"`
	Ram      int16  `json:"ram"`
	Gpu      string `json:"gpu"`
	Price    string `json:"price"`
	Discount int16  `json:"discount"`
}
