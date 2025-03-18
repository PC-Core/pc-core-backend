package inputs

import "github.com/PC-Core/pc-core-backend/internal/models"

type AddCpuInput struct {
	Name        string              `json:"name"`
	Price       float64             `json:"price"`
	Stock       uint64              `json:"stock"`
	CpuName     string              `json:"cpu_name"`
	PCores      uint64              `json:"pcores"`
	ECores      uint64              `json:"ecores"`
	Threads     uint64              `json:"threads"`
	BaseFreqMHz uint64              `json:"base_freq_mhz"`
	MaxFreqMHz  uint64              `json:"max_freq_mhz"`
	Socket      models.CpuSocket    `json:"socket"`
	L1KB        uint64              `json:"l1_kb"`
	L2KB        uint64              `json:"l2_kb"`
	L3KB        uint64              `json:"l3_kb"`
	TecProcNM   uint64              `json:"tecproc_nm"`
	TDPWatt     uint64              `json:"tdp_watt"`
	ReleaseYear uint64              `json:"release_year"`
	Medias      []models.InputMedia `json:"medias"`
}
