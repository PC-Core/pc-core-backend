package models

type Gpu struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	MemoryGB     int    `json:"memory_bg"`
	MemoryType   string `json:"memory_type"` //GDDR6 и тд
	BusWidthBit  int    `json:"bus_width_bit"`
	BaseFreqMHz  int    `json:"base_freq_mhz"`
	BoostFreqMHz int    `json:"boost_freq_mhz"`
	TecprocNm    int    `json:"tecproc_nm"` // нанометры (7, 6, 5 и т.д)
	TDPWatt      int    `json:"tdp_watt"`
	RealeseYear  int    `json:"realese_year"`
}
