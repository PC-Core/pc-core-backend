package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	SEED_TYPE_USER     = "user"
	SEED_TYPE_LAPTOP   = "laptop"
	SEED_TYPE_CATEGORY = "category"
)

var SEED_TYPE_NAMES = []string{SEED_TYPE_USER, SEED_TYPE_LAPTOP, SEED_TYPE_CATEGORY}

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

func InsertLaptops(db *sql.DB) {
	laptops := []struct {
		Name   string
		Price  float64
		Selled uint64
		Stock  uint64
		Cpu    string
		Ram    int16
		Gpu    string
	}{
		{"MSI Titan 18", 614999, 0, 13, "i9-14900HX", 32, "RTX 4090"},
		{"Lenovo Legion Y9000P", 425999, 32, 28, "i9-14900HX", 32, "RTX 4090"},
		{"Apple Macbook Pro", 421599, 3, 83, "Apple M3 Max", 48, "Apple M3 Max"},
		{"MSI Raider 18", 417999, 10, 38, "i9-14900HX", 32, "RTX 4090"},
		{"ASUS ROG Zephyrus Duo 16", 389999, 5, 47, "Ryzen 9 7945HX", 32, "RTX 4090"},
		{"MSI Vector 17", 374999, 5, 123, "i9-14900HX", 32, "RTX 4080"},
		{"ASUS VivoBook Pro 15", 179999, 42, 82, "Core Ultra 9 185H", 24, "RTX 4060"},
		{"MSI Sword 17 HX", 179999, 6, 35, "i7-13700HX", 16, "RTX 4070"},
		{"MSI Summit 13 AI+ Evo", 179999, 0, 4, "Core Ultra 7 258V", 32, "Intel Arc Graphics"},
		{"Honor MagicBook Art 14", 149999, 2, 64, "Core Ultra 7 155H", 32, "Intel Arc Graphics"},
		{"Apple MacBook Air", 179499, 30, 32, "Apple M3", 16, "Apple M3"},
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

		err = tx.QueryRow("INSERT INTO LaptopChars (cpu, ram, gpu) VALUES ($1, $2, $3) returning id", laptop.Cpu, laptop.Ram, laptop.Gpu).Scan(&charId)

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

var SEED_TYPES = map[string]func(*sql.DB){
	SEED_TYPE_USER:     InsertUsers,
	SEED_TYPE_LAPTOP:   InsertLaptops,
	SEED_TYPE_CATEGORY: InsertCategories,
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

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_IBYTE_CONN"))

	if err != nil {
		fmt.Println("Error while connecting the db: ", err)
		return
	}

	handleCliArgs(os.Args, db)
}
