package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

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
	SOCKET_BGA1964   CpuSocket = "BGA1964"
	SOCKET_FL1       CpuSocket = "FL1"
	SOCKET_UNKNOWN   CpuSocket = "UNKNOWN"
	SOCKET_BGA2049   CpuSocket = "BGA2049"
	SOCKET_BGA2833   CpuSocket = "BGA2833"
)

const (
	SEED_TYPE_USER     = "user"
	SEED_TYPE_LAPTOP   = "laptop"
	SEED_TYPE_CATEGORY = "category"
	SEED_ALL           = "all"
	CLEAR_ALL          = "clear"
)

var SEED_TYPE_NAMES = []string{SEED_TYPE_USER, SEED_TYPE_LAPTOP, SEED_TYPE_CATEGORY, SEED_ALL, CLEAR_ALL}

func Sha256(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InsertUsers(db *sql.DB) {
	users := []struct {
		Name     string
		Email    string
		Role     string
		Password string
	}{
		{"yellowpeacock117", "jennie.nichols@example.com", "Default", "bibi"},
		{"sadmouse784", "ievfimiya.dibrova@example.com", "Default", "around"},
		{"browntiger738", "mariana.garica@example.com", "Admin", "killer1"},
		{"whitebear910", "diane.fontai@example.com", "Default", "nancy1"},
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO Users (Name, Email, Role, PasswordHash) VALUES ($1, $2, $3, $4)", user.Name, user.Email, user.Role, Sha256(user.Password))

		if err != nil {
			fmt.Println("Error while inserting users: ", err)
		}
	}
}

var media_ids []uint64

func InsertMedias(db *sql.DB) {
	medias := []struct {
		Url  string
		Type string
	}{
		{"C:\\2.png", "Image"},
		{"C:\\3.mp4", "Video"},
		{"C:\\4.png", "Image"},
		{"C:\\5.mp4", "Video"},
		{"C:\\6.png", "Image"},
		{"C:\\7.png", "Image"},
		{"C:\\8.png", "Image"},
		{"C:\\9.png", "Image"},
		{"C:\\10.mp4", "Video"},
		{"C:\\1.png", "Image"},
		{"C:\\11.png", "Image"},
		{"C:\\12.mp4", "Video"},
		{"C:\\13.png", "Image"},
		{"C:\\14.mp4", "Video"},
		{"C:\\15.mp4", "Video"},
		{"C:\\16.png", "Image"},
		{"C:\\17.png", "Image"},
		{"C:\\18.mp4", "Video"},
		{"C:\\19.mp4", "Video"},
		{"C:\\20.png", "Image"},
	}

	for _, media := range medias {
		var id uint64
		err := db.QueryRow("INSERT INTO Medias (url, type) VALUES ($1, $2) returning id", media.Url, media.Type).Scan(&id)

		media_ids = append(media_ids, id)

		if err != nil {
			log.Fatalf("Error while adding medias: %s", err.Error())
		}
	}
}

func InsertLaptops(db *sql.DB) {
	InsertMedias(db)

	cpus := []struct {
		Name         string    `json:"name"`
		PCores       uint64    `json:"pcores"`
		ECores       uint64    `json:"ecores"`
		Threads      uint64    `json:"threads"`
		BasePFreqMHz uint64    `json:"base_p_freq_mhz"`
		MaxPFreqMHz  uint64    `json:"max_p_freq_mhz"`
		BaseEFreqMHz uint64    `json:"base_e_freq_mhz"`
		MaxEFreqMHz  uint64    `json:"max_e_freq_mhz"`
		Socket       CpuSocket `json:"socket"`
		L1KB         uint64    `json:"l1_kb"`
		L2KB         uint64    `json:"l2_kb"`
		L3KB         uint64    `json:"l3_kb"`
		TecProcNM    uint64    `json:"tecproc_nm"`
		TDPWatt      uint64    `json:"tdp_watt"`
		ReleaseYear  uint64    `json:"release_year"`
	}{
		{
			Name:         "i9-14900HX",
			PCores:       8,
			ECores:       16,
			Threads:      32,
			BasePFreqMHz: 2200,
			MaxPFreqMHz:  5800,
			BaseEFreqMHz: 1600,
			MaxEFreqMHz:  4100,
			Socket:       SOCKET_BGA1964,
			L1KB:         80,
			L2KB:         2000,
			L3KB:         36000,
			TecProcNM:    10,
			TDPWatt:      55,
			ReleaseYear:  2024,
		},
		{
			Name:         "Apple M3 Max",
			PCores:       12,
			ECores:       4,
			Threads:      16,
			BasePFreqMHz: 4100,
			MaxPFreqMHz:  4100,
			BaseEFreqMHz: 2700,
			MaxEFreqMHz:  2700,
			Socket:       SOCKET_UNKNOWN,
			L1KB:         192,
			L2KB:         32000,
			L3KB:         0,
			TecProcNM:    3,
			TDPWatt:      78,
			ReleaseYear:  2023,
		},
		{
			Name:         "Ryzen 9 7945HX",
			PCores:       16,
			ECores:       0,
			Threads:      32,
			BasePFreqMHz: 2500,
			MaxPFreqMHz:  5400,
			BaseEFreqMHz: 0,
			MaxEFreqMHz:  0,
			Socket:       SOCKET_FL1,
			L1KB:         64,
			L2KB:         1000,
			L3KB:         64000,
			TecProcNM:    5,
			TDPWatt:      55,
			ReleaseYear:  2023,
		},
		{
			Name:         "Core Ultra 9 185H",
			PCores:       6,
			ECores:       10,
			Threads:      22,
			BasePFreqMHz: 3900,
			MaxPFreqMHz:  5100,
			BaseEFreqMHz: 1900,
			MaxEFreqMHz:  3800,
			Socket:       SOCKET_BGA2049,
			L1KB:         112,
			L2KB:         2000,
			L3KB:         24000,
			TecProcNM:    7,
			TDPWatt:      45,
			ReleaseYear:  2023,
		},
		{
			Name:         "i7-13700HX",
			PCores:       8,
			ECores:       8,
			Threads:      24,
			BasePFreqMHz: 2100,
			MaxPFreqMHz:  5000,
			BaseEFreqMHz: 1500,
			MaxEFreqMHz:  3700,
			Socket:       SOCKET_BGA1964,
			L1KB:         80,
			L2KB:         2000,
			L3KB:         30000,
			TecProcNM:    10,
			TDPWatt:      55,
			ReleaseYear:  2023,
		},
		{
			Name:         "Core Ultra 7 258V",
			PCores:       4,
			ECores:       4,
			Threads:      8,
			BasePFreqMHz: 2200,
			MaxPFreqMHz:  4800,
			BaseEFreqMHz: 2200,
			MaxEFreqMHz:  3700,
			Socket:       SOCKET_BGA2833,
			L1KB:         48,
			L2KB:         192,
			L3KB:         2500,
			TecProcNM:    3,
			TDPWatt:      17,
			ReleaseYear:  2024,
		},
		{
			Name:         "Core Ultra 7 155H",
			PCores:       6,
			ECores:       10,
			Threads:      22,
			BasePFreqMHz: 3800,
			MaxPFreqMHz:  4800,
			BaseEFreqMHz: 1800,
			MaxEFreqMHz:  3800,
			Socket:       SOCKET_BGA2049,
			L1KB:         112,
			L2KB:         2000,
			L3KB:         24000,
			TecProcNM:    7,
			TDPWatt:      28,
			ReleaseYear:  2023,
		},
		{
			Name:         "Apple M3",
			PCores:       4,
			ECores:       4,
			Threads:      8,
			BasePFreqMHz: 4100,
			MaxPFreqMHz:  4100,
			BaseEFreqMHz: 2700,
			MaxEFreqMHz:  2700,
			Socket:       SOCKET_UNKNOWN,
			L1KB:         112,
			L2KB:         4000,
			L3KB:         0,
			TecProcNM:    4,
			TDPWatt:      0,
			ReleaseYear:  2023,
		},
	}

	ids := make([]uint64, 0, len(cpus))

	for _, cpu := range cpus {
		var (
			charId uint64
		)

		err := db.QueryRow(fmt.Sprintf("INSERT INTO %s (name, pcores, ecores, threads, base_p_freq_mhz, max_p_freq_mhz, base_e_freq_mhz, max_e_freq_mhz, socket, l1_kb, l2_kb, l3_kb, tecproc_nm, tdp_watt, release_year) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning id", "CpuChars"), cpu.Name, cpu.PCores, cpu.ECores, cpu.Threads, cpu.BasePFreqMHz, cpu.MaxPFreqMHz, cpu.BaseEFreqMHz, cpu.MaxEFreqMHz, cpu.Socket, cpu.L1KB, cpu.L2KB, cpu.L3KB, cpu.TecProcNM, cpu.TDPWatt, cpu.ReleaseYear).Scan(&charId)

		if err != nil {
			panic(err)
		}

		ids = append(ids, charId)
	}

	laptops := []struct {
		Name   string
		Price  float64
		Selled uint64
		Stock  uint64
		medias []uint64
		CpuID  uint64
		Ram    int16
		Gpu    string
	}{
		{"MSI Titan 18", 614999, 0, 13, media_ids[0:2], ids[0], 32, "RTX 4090"},
		{"Lenovo Legion Y9000P", 425999, 32, 28, media_ids[2:4], ids[0], 32, "RTX 4090"},
		{"Apple Macbook Pro", 421599, 3, 83, media_ids[4:6], ids[1], 48, "Apple M3 Max"},
		{"MSI Raider 18", 417999, 10, 38, media_ids[6:8], ids[0], 32, "RTX 4090"},
		{"ASUS ROG Zephyrus Duo 16", 389999, 5, 47, media_ids[8:10], ids[2], 32, "RTX 4090"},
		{"MSI Vector 17", 374999, 5, 123, media_ids[10:12], ids[0], 32, "RTX 4080"},
		{"ASUS VivoBook Pro 15", 179999, 42, 82, []uint64{}, ids[3], 24, "RTX 4060"},
		{"MSI Sword 17 HX", 179999, 6, 35, media_ids[12:14], ids[4], 16, "RTX 4070"},
		{"MSI Summit 13 AI+ Evo", 179999, 0, 4, media_ids[14:16], ids[5], 32, "Intel Arc Graphics"},
		{"Honor MagicBook Art 14", 149999, 2, 64, media_ids[16:18], ids[6], 32, "Intel Arc Graphics"},
		{"Apple MacBook Air", 179499, 30, 32, media_ids[18:20], ids[7], 16, "Apple M3"},
	}

	for _, laptop := range laptops {
		var (
			charId    uint64
			productId uint64
		)

		tx, err := db.Begin()

		if err != nil {
			fmt.Println("Error while starting transaction: ", err)
			return
		}

		defer tx.Rollback()

		err = tx.QueryRow("INSERT INTO LaptopChars (cpu_id, ram, gpu) VALUES ($1, $2, $3) returning id", laptop.CpuID, laptop.Ram, laptop.Gpu).Scan(&charId)

		if err != nil {
			fmt.Println("Error while inserting products: ", err)
			return
		}

		err = tx.QueryRow("INSERT INTO Products (name, price, selled, stock, chars_table_name, chars_id, medias) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id", laptop.Name, laptop.Price, laptop.Selled, laptop.Stock, "LaptopChars", charId, pq.Array(laptop.medias)).Scan(&productId)

		if err != nil {
			fmt.Println("Error while inserting laptops: ", err)
			return
		}

		if err := tx.Commit(); err != nil {
			fmt.Println("Error on commiting laptops: ", err)
			return
		}
	}
}

func InsertCategories(db *sql.DB) {
	cats := []struct {
		Title       string
		Description string
		Icon        string
		Slug        string
	}{
		{"Процессоры", "Сердце вашего компьютера! В нашем ассортименте процессоры для любых нужд — от бюджетных моделей до высококлассных чипов для гейминга и работы с тяжелыми приложениями", "", "cpu"},
		{"Ноутбуки", "Идеальный выбор для тех, кто ценит мобильность и производительность. У нас представлены ноутбуки для работы, учебы, развлечений и гейминга с различными характеристиками и дизайнами.", "", "laptop"},
		{"Видеокарты", "Для тех, кто ценит графику и производительность в играх или профессиональной работе. Мы предлагаем видеокарты от лидеров отрасли с отличными характеристиками для любого бюджета.", "", "gpu"},
		{"ОЗУ", "Увеличьте быстродействие вашего ПК с помощью высококачественной оперативной памяти. У нас есть ОЗУ для любых нужд — от стандартных моделей до сверхбыстрых для энтузиастов и профессионалов.", "", "ram"},
		{"ПК", "Готовые решения для работы, учёбы и гейминга. В нашем ассортименте — как стандартные офисные ПК, так и мощные игровые системы с топовыми комплектующими для самых требовательных пользователей.", "", "pc"},
	}

	for _, cat := range cats {
		_, err := db.Exec("INSERT INTO Categories (title, description, icon, slug) VALUES ($1, $2, $3, $4)", cat.Title, cat.Description, cat.Icon, cat.Slug)

		if err != nil {
			fmt.Println("Error while inserting users: ", err)
		}
	}
}

// func InsertCpus(db *sql.DB) {
// 	cpus := []struct {
// 	}
// }

func Clear(db *sql.DB) {
	tables := []string{
		"Cart",
		"Users",
		"LaptopChars",
		"CpuChars",
		"Categories",
		"Medias",
		"Products",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))

		if err != nil {
			log.Fatalf("Error while deleting from %s: %s", table, err.Error())
		}
	}
}

func InsertAll(db *sql.DB) {
	InsertUsers(db)
	InsertLaptops(db)
	InsertCategories(db)
}

var SEED_TYPES = map[string]func(*sql.DB){
	SEED_TYPE_USER:     InsertUsers,
	SEED_TYPE_LAPTOP:   InsertLaptops,
	SEED_TYPE_CATEGORY: InsertCategories,
	SEED_ALL:           InsertAll,
	CLEAR_ALL:          Clear,
}

func formatHelpMessage() string {
	return fmt.Sprintf("Setup seeds.\n\nUsage:\n\tgo run seeds.go [SEED TYPES]\n\nSeed Types:\n\t%s", strings.Join(SEED_TYPE_NAMES, "\n\t"))
}

func handleCliArgs(args []string, db *sql.DB) {
	if len(args) == 1 {
		fmt.Println(formatHelpMessage())
		return
	}

	for _, arg := range args[1:] {
		fn, ok := SEED_TYPES[arg]

		if !ok {
			fmt.Printf("Error: unknown parameter: %s\n", arg)
			os.Exit(-1)
		}

		fn(db)
	}
}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", os.Getenv("PCCORE_POSTGRES_CONN"))

	if err != nil {
		fmt.Println("Error while connecting the db: ", err)
		return
	}

	handleCliArgs(os.Args, db)
}
