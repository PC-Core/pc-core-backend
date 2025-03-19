package inputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type AddCpuInput struct {
	Name         string              `json:"name"`
	Price        float64             `json:"price"`
	Stock        uint64              `json:"stock"`
	CpuName      string              `json:"cpu_name"`
	PCores       uint64              `json:"pcores"`
	ECores       uint64              `json:"ecores"`
	Threads      uint64              `json:"threads"`
	BasePFreqMHz uint64              `json:"base_p_freq_mhz"`
	MaxPFreqMHz  uint64              `json:"max_p_freq_mhz"`
	BaseEFreqMHz uint64              `json:"base_e_freq_mhz"`
	MaxEFreqMHz  uint64              `json:"max_e_freq_mhz"`
	Socket       models.CpuSocket    `json:"socket"`
	L1KB         uint64              `json:"l1_kb"`
	L2KB         uint64              `json:"l2_kb"`
	L3KB         uint64              `json:"l3_kb"`
	TecProcNM    uint64              `json:"tecproc_nm"`
	TDPWatt      uint64              `json:"tdp_watt"`
	ReleaseYear  uint64              `json:"release_year"`
	Medias       []models.InputMedia `json:"medias"`
}
