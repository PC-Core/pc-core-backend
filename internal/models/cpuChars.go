package models

type CpuSocket string

const (
	SOCKET_AM4       CpuSocket = "AM4"
	SOCKET_AM5       CpuSocket = "AM5"
	SOCKET_LGA775    CpuSocket = "LGA775"
	SOCKET_LGA1156   CpuSocket = "LGA1156"
	SOCKET_LGA1155   CpuSocket = "LGA1155"
	SOCKET_LGA1150   CpuSocket = "LGA1150"
	SOCKET_LGA1151   CpuSocket = "LGA1151"
	SOCKET_LGA1151v2 CpuSocket = "LGA1151v2"
	SOCKET_LGA1200   CpuSocket = "LGA1200"
	SOCKET_LGA1700   CpuSocket = "LGA1700"
	SOCKET_LGA1851   CpuSocket = "LGA1851"
)

type CpuChars struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	PCores      uint64    `json:"pcores"`
	ECores      uint64    `json:"ecores"`
	Threads     uint64    `json:"threads"`
	BaseFreqMHz uint64    `json:"base_freq_mhz"`
	MaxFreqMHz  uint64    `json:"max_freq_mhz"`
	Socket      CpuSocket `json:"socket"`
	L1KB        uint64    `json:"l1_kb"`
	L2KB        uint64    `json:"l2_kb"`
	L3KB        uint64    `json:"l3_kb"`
	TecProcNM   uint64    `json:"tecproc_nm"`
	TDPWatt     uint64    `json:"tdp_watt"`
	ReleaseYear uint64    `json:"release_year"`
}

func NewCpuChars(id uint64, name string, pcores, ecores, threads, bfmhz, mfmhz uint64, socket CpuSocket, l1, l2, l3, tp, tdp, ry uint64) *CpuChars {
	return &CpuChars{
		id, name, pcores, ecores, threads, bfmhz, mfmhz, socket, l1, l2, l3, tp, tdp, ry,
	}
}
