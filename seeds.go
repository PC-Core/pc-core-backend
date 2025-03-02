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

	laptops := []struct {
		Name   string
		Price  float64
		Selled uint64
		Stock  uint64
		medias []uint64
		Cpu    string
		Ram    int16
		Gpu    string
	}{
		{"MSI Titan 18", 614999, 0, 13, media_ids[0:2], "i9-14900HX", 32, "RTX 4090"},
		{"Lenovo Legion Y9000P", 425999, 32, 28, media_ids[2:4], "i9-14900HX", 32, "RTX 4090"},
		{"Apple Macbook Pro", 421599, 3, 83, media_ids[4:6], "Apple M3 Max", 48, "Apple M3 Max"},
		{"MSI Raider 18", 417999, 10, 38, media_ids[6:8], "i9-14900HX", 32, "RTX 4090"},
		{"ASUS ROG Zephyrus Duo 16", 389999, 5, 47, media_ids[8:10], "Ryzen 9 7945HX", 32, "RTX 4090"},
		{"MSI Vector 17", 374999, 5, 123, media_ids[10:12], "i9-14900HX", 32, "RTX 4080"},
		{"ASUS VivoBook Pro 15", 179999, 42, 82, []uint64{}, "Core Ultra 9 185H", 24, "RTX 4060"},
		{"MSI Sword 17 HX", 179999, 6, 35, media_ids[12:14], "i7-13700HX", 16, "RTX 4070"},
		{"MSI Summit 13 AI+ Evo", 179999, 0, 4, media_ids[14:16], "Core Ultra 7 258V", 32, "Intel Arc Graphics"},
		{"Honor MagicBook Art 14", 149999, 2, 64, media_ids[16:18], "Core Ultra 7 155H", 32, "Intel Arc Graphics"},
		{"Apple MacBook Air", 179499, 30, 32, media_ids[18:20], "Apple M3", 16, "Apple M3"},
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

func Clear(db *sql.DB) {
	tables := []string{
		"Cart",
		"Users",
		"LaptopChars",
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
