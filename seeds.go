package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	users := []models.User{
		*models.NewUser(0, "yellowpeacock117", "jennie.nichols@example.com", "Default", Sha256("bibi")),
		*models.NewUser(0, "sadmouse784", "ievfimiya.dibrova@example.com", "Default", Sha256("around")),
		*models.NewUser(0, "browntiger738", "mariana.garica@example.com", "Admin", Sha256("killer1")),
		*models.NewUser(0, "whitebear910", "diane.fontai@example.com", "Default", Sha256("nancy1")),
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO Users (Name, Email, Role, PasswordHash) VALUES ($1, $2, $3, $4)", user.Name, user.Email, user.Role, user.PasswdHash)

		if err != nil {
			fmt.Println("Error while inserting users: ", err)
		}
	}
}

var media_ids []uint64

func InsertMedias(db *sql.DB, ids []uint64) {
	medias := []models.Media{
		{0, "C:\\2.png", "Image", ids[0]},
		{0, "C:\\3.mp4", "Video", ids[0]},
		{0, "C:\\4.png", "Image", ids[2]},
		{0, "C:\\5.mp4", "Video", ids[2]},
		{0, "C:\\6.png", "Image", ids[3]},
		{0, "C:\\7.png", "Image", ids[3]},
		{0, "C:\\8.png", "Image", ids[4]},
		{0, "C:\\9.png", "Image", ids[4]},
		{0, "C:\\10.mp4", "Video", ids[5]},
		{0, "C:\\1.png", "Image", ids[5]},
		{0, "C:\\11.png", "Image", ids[6]},
		{0, "C:\\12.mp4", "Video", ids[6]},
		{0, "C:\\13.png", "Image", ids[6]},
		{0, "C:\\14.mp4", "Video", ids[6]},
		{0, "C:\\15.mp4", "Video", ids[7]},
		{0, "C:\\16.png", "Image", ids[8]},
		{0, "C:\\17.png", "Image", ids[8]},
		{0, "C:\\18.mp4", "Video", ids[8]},
		{0, "C:\\19.mp4", "Video", ids[9]},
		{0, "C:\\20.png", "Image", ids[10]},
	}

	for _, media := range medias {
		var id uint64
		err := db.QueryRow("INSERT INTO Medias (url, type, product_id) VALUES ($1, $2, $3) returning id", media.Url, media.Type, media.ProductID).Scan(&id)

		media_ids = append(media_ids, id)

		if err != nil {
			log.Fatalf("Error while adding medias: %s", err.Error())
		}
	}
}

func InsertLaptops(db *sql.DB) {
	cpus := []models.CpuChars{
		{
			Name:         "i9-14900HX",
			PCores:       8,
			ECores:       16,
			Threads:      32,
			BasePFreqMHz: 2200,
			MaxPFreqMHz:  5800,
			BaseEFreqMHz: 1600,
			MaxEFreqMHz:  4100,
			Socket:       models.SOCKET_BGA1964,
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
			Socket:       models.SOCKET_UNKNOWN,
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
			Socket:       models.SOCKET_FL1,
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
			Socket:       models.SOCKET_BGA2049,
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
			Socket:       models.SOCKET_BGA1964,
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
			Socket:       models.SOCKET_BGA2833,
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
			Socket:       models.SOCKET_BGA2049,
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
			Socket:       models.SOCKET_UNKNOWN,
			L1KB:         112,
			L2KB:         4000,
			L3KB:         0,
			TecProcNM:    4,
			TDPWatt:      0,
			ReleaseYear:  2023,
		},
	}

	ids := make([]uint64, 0, len(cpus))
	laptop_ids := make([]uint64, 0, len(cpus))

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
		CpuID  uint64
		Ram    int16
		Gpu    string
	}{
		{"MSI Titan 18", 614999, 0, 13, ids[0], 32, "RTX 4090"},
		{"Lenovo Legion Y9000P", 425999, 32, 28, ids[0], 32, "RTX 4090"},
		{"Apple Macbook Pro", 421599, 3, 83, ids[1], 48, "Apple M3 Max"},
		{"MSI Raider 18", 417999, 10, 38, ids[0], 32, "RTX 4090"},
		{"ASUS ROG Zephyrus Duo 16", 389999, 5, 47, ids[2], 32, "RTX 4090"},
		{"MSI Vector 17", 374999, 5, 123, ids[0], 32, "RTX 4080"},
		{"ASUS VivoBook Pro 15", 179999, 42, 82, ids[3], 24, "RTX 4060"},
		{"MSI Sword 17 HX", 179999, 6, 35, ids[4], 16, "RTX 4070"},
		{"MSI Summit 13 AI+ Evo", 179999, 0, 4, ids[5], 32, "Intel Arc Graphics"},
		{"Honor MagicBook Art 14", 149999, 2, 64, ids[6], 32, "Intel Arc Graphics"},
		{"Apple MacBook Air", 179499, 30, 32, ids[7], 16, "Apple M3"},
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

		err = tx.QueryRow("INSERT INTO Products (name, price, selled, stock, chars_table_name, chars_id) VALUES ($1, $2, $3, $4, $5, $6) returning id", laptop.Name, laptop.Price, laptop.Selled, laptop.Stock, "LaptopChars", charId).Scan(&productId)

		if err != nil {
			fmt.Println("Error while inserting laptops: ", err)
			return
		}

		if err := tx.Commit(); err != nil {
			fmt.Println("Error on commiting laptops: ", err)
			return
		}

		laptop_ids = append(laptop_ids, productId)
	}

	InsertMedias(db, laptop_ids)
}

func InsertCategories(db *sql.DB) {
	cats := []models.Category{
		{0, "Процессоры", "Сердце вашего компьютера! В нашем ассортименте процессоры для любых нужд — от бюджетных моделей до высококлассных чипов для гейминга и работы с тяжелыми приложениями", "", "cpu"},
		{0, "Ноутбуки", "Идеальный выбор для тех, кто ценит мобильность и производительность. У нас представлены ноутбуки для работы, учебы, развлечений и гейминга с различными характеристиками и дизайнами.", "", "laptop"},
		{0, "Видеокарты", "Для тех, кто ценит графику и производительность в играх или профессиональной работе. Мы предлагаем видеокарты от лидеров отрасли с отличными характеристиками для любого бюджета.", "", "gpu"},
		{0, "ОЗУ", "Увеличьте быстродействие вашего ПК с помощью высококачественной оперативной памяти. У нас есть ОЗУ для любых нужд — от стандартных моделей до сверхбыстрых для энтузиастов и профессионалов.", "", "ram"},
		{0, "ПК", "Готовые решения для работы, учёбы и гейминга. В нашем ассортименте — как стандартные офисные ПК, так и мощные игровые системы с топовыми комплектующими для самых требовательных пользователей.", "", "pc"},
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
