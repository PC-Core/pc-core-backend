package models

type GpuChars struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	MemoryGB     uint64 `json:"memory_gb"`
	MemoryType   string `json:"memory_type"`
	BusWidthBit  uint64 `json:"bus_width_bit"`
	BaseFreqMHz  uint64 `json:"base_freq_mhz"`
	BoostFreqMHz uint64 `json:"boost_freq_mhz"`
	TecprocNm    uint64 `json:"tecproc_nm"`
	TDPWatt      uint64 `json:"tdp_watt"`
	ReleaseYear  uint64 `json:"release_year"`
}

func NewGpuChars(id uint64, name, memorytype string, memorygb, buswidthbit, basefreqmhz, boostfreqmhz, tecprocnm, tdpwatt, releasedate uint64) *GpuChars {
	return &GpuChars{
		id, name, memorygb, memorytype, buswidthbit, basefreqmhz, boostfreqmhz, tecprocnm, tdpwatt, releasedate,
	}
}
