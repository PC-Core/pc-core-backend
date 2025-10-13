package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PC-Core/pc-core-backend/pkg/config"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	SEED_TYPE_USER     = "user"
	SEED_TYPE_LAPTOP   = "laptop"
	SEED_TYPE_CATEGORY = "category"
	SEED_TYPE_MEDIA    = "media"
	SEED_ALL           = "all"
	CLEAR_ALL          = "clear"
	SEED_TYPE_GPU      = "gpu"
	SEED_TYPE_KEYBOARD = "keyboard"
)

var SEED_TYPE_NAMES = []string{SEED_TYPE_USER, SEED_TYPE_LAPTOP, SEED_TYPE_CATEGORY, SEED_TYPE_MEDIA, SEED_ALL, CLEAR_ALL, SEED_TYPE_GPU}

// MinIOConfig конфигурация MinIO
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// MediaDownloadTask задача на загрузку медиа
type MediaDownloadTask struct {
	URL        string
	ObjectName string
	ProductID  uint64
	Type       string
}

// BatchMediaDownloader для пакетной загрузки медиа
type BatchMediaDownloader struct {
	minioClient *minio.Client
	config      MinIOConfig
	httpClient  *http.Client
	semaphore   chan struct{}
	wg          sync.WaitGroup
	mu          sync.Mutex
	results     []MediaDownloadResult
}

// MediaDownloadResult результат загрузки
type MediaDownloadResult struct {
	URL        string
	ObjectName string
	Success    bool
	Error      error
	Duration   time.Duration
}

func Sha256(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InitMinIOClient(config MinIOConfig) (*minio.Client, error) {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MinIO: %w", err)
	}

	// Проверяем и создаем бакет если нужно
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки бакета: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("ошибка создания бакета: %w", err)
		}
	}

	return minioClient, nil
}

func NewBatchMediaDownloader(config MinIOConfig, maxConcurrent int) (*BatchMediaDownloader, error) {
	minioClient, err := InitMinIOClient(config)
	if err != nil {
		return nil, err
	}

	return &BatchMediaDownloader{
		minioClient: minioClient,
		config:      config,
		httpClient: &http.Client{
			Timeout: 30 * time.Minute,
			Transport: &http.Transport{
				MaxIdleConns:        maxConcurrent,
				MaxIdleConnsPerHost: maxConcurrent,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		semaphore: make(chan struct{}, maxConcurrent),
		results:   make([]MediaDownloadResult, 0),
	}, nil
}

func (b *BatchMediaDownloader) DownloadAndUploadMedia(ctx context.Context, tasks []MediaDownloadTask) []MediaDownloadResult {
	for _, task := range tasks {
		fmt.Println(task)
		b.wg.Add(1)
		go b.processMediaTask(ctx, task)
	}

	b.wg.Wait()
	return b.results
}

func (b *BatchMediaDownloader) processMediaTask(ctx context.Context, task MediaDownloadTask) {
	defer b.wg.Done()

	// Захватываем слот семафора
	b.semaphore <- struct{}{}
	defer func() { <-b.semaphore }()

	startTime := time.Now()
	result := MediaDownloadResult{
		URL:        task.URL,
		ObjectName: task.ObjectName,
	}

	// Скачиваем и загружаем файл
	err := b.downloadAndUploadSingleMedia(ctx, task)
	if err != nil {
		result.Error = err
		result.Success = false
	} else {
		result.Success = true
	}

	result.Duration = time.Since(startTime)

	b.mu.Lock()
	b.results = append(b.results, result)
	b.mu.Unlock()
}

func (b *BatchMediaDownloader) downloadAndUploadSingleMedia(ctx context.Context, task MediaDownloadTask) error {
	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, "GET", task.URL, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Выполняем запрос
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка скачивания: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP ошибка: %s", resp.Status)
	}

	// Определяем content type
	contentType := "application/octet-stream"
	if task.Type == "Image" {
		contentType = getImageContentType(task.URL)
	} else if task.Type == "Video" {
		contentType = "video/mp4"
	}

	// Загружаем в MinIO
	_, err = b.minioClient.PutObject(ctx, b.config.Bucket, task.ObjectName, resp.Body, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("ошибка загрузки в MinIO: %w", err)
	}

	return nil
}

func getImageContentType(url string) string {
	ext := strings.ToLower(url[strings.LastIndex(url, "."):])
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// ReadMediaTasksFromFile читает задачи на загрузку медиа из файла
func ReadMediaTasksFromFile(filename string) ([]MediaDownloadTask, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	tasks := make([]MediaDownloadTask, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}

		url := strings.TrimSpace(parts[0])
		objectName := strings.TrimSpace(parts[1])
		productID := strings.TrimSpace(parts[2])
		mediaType := strings.TrimSpace(parts[3])

		var pid uint64
		fmt.Sscanf(productID, "%d", &pid)

		tasks = append(tasks, MediaDownloadTask{
			URL:        url,
			ObjectName: objectName,
			ProductID:  pid,
			Type:       mediaType,
		})
	}

	return tasks, nil
}

// getDefaultMediaTasks возвращает дефолтные задачи если файл не найден
func getDefaultMediaTasks(laptopIDs []uint64) []MediaDownloadTask {
	if len(laptopIDs) < 11 {
		return []MediaDownloadTask{}
	}

	return []MediaDownloadTask{
		{"https://example.com/images/2.png", "laptops/2.png", laptopIDs[0], "Image"},
		{"https://example.com/videos/3.mp4", "laptops/3.mp4", laptopIDs[0], "Video"},
		{"https://example.com/images/4.png", "laptops/4.png", laptopIDs[2], "Image"},
		{"https://example.com/videos/5.mp4", "laptops/5.mp4", laptopIDs[2], "Video"},
		{"https://example.com/images/6.png", "laptops/6.png", laptopIDs[3], "Image"},
		{"https://example.com/images/7.png", "laptops/7.png", laptopIDs[3], "Image"},
		{"https://example.com/images/8.png", "laptops/8.png", laptopIDs[4], "Image"},
		{"https://example.com/images/9.png", "laptops/9.png", laptopIDs[4], "Image"},
		{"https://example.com/videos/10.mp4", "laptops/10.mp4", laptopIDs[5], "Video"},
		{"https://example.com/images/1.png", "laptops/1.png", laptopIDs[5], "Image"},
		{"https://example.com/images/11.png", "laptops/11.png", laptopIDs[6], "Image"},
		{"https://example.com/videos/12.mp4", "laptops/12.mp4", laptopIDs[6], "Video"},
		{"https://example.com/images/13.png", "laptops/13.png", laptopIDs[6], "Image"},
		{"https://example.com/videos/14.mp4", "laptops/14.mp4", laptopIDs[6], "Video"},
		{"https://example.com/videos/15.mp4", "laptops/15.mp4", laptopIDs[7], "Video"},
		{"https://example.com/images/16.png", "laptops/16.png", laptopIDs[8], "Image"},
		{"https://example.com/images/17.png", "laptops/17.png", laptopIDs[8], "Image"},
		{"https://example.com/videos/18.mp4", "laptops/18.mp4", laptopIDs[8], "Video"},
		{"https://example.com/videos/19.mp4", "laptops/19.mp4", laptopIDs[9], "Video"},
		{"https://example.com/images/20.png", "laptops/20.png", laptopIDs[10], "Image"},
	}
}

// getExistingLaptopIDs получает ID существующих ноутбуков из базы
func getExistingLaptopIDs(db *sql.DB) []uint64 {
	rows, err := db.Query("SELECT id FROM Products WHERE chars_table_name = 'LaptopChars'")
	if err != nil {
		log.Printf("Ошибка получения ID ноутбуков: %v", err)
		return []uint64{}
	}
	defer rows.Close()

	var ids []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			log.Printf("Ошибка сканирования ID: %v", err)
			continue
		}
		ids = append(ids, id)
	}

	return ids
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

func InsertMedias(db *sql.DB, minioConfig MinIOConfig, laptopIDs []uint64) {
	var mediaTasks []MediaDownloadTask
	var err error

	// Пытаемся прочитать из файла, если он существует
	if _, err := os.Stat("media_urls.txt"); err == nil {
		mediaTasks, err = ReadMediaTasksFromFile("media_urls.txt")
		if err != nil {
			log.Printf("Ошибка чтения файла с ссылками: %v", err)
			// Fallback на дефолтные ссылки
			mediaTasks = getDefaultMediaTasks(laptopIDs)
		}
	} else {
		// Используем дефолтные ссылки если файла нет
		mediaTasks = getDefaultMediaTasks(laptopIDs)
	}

	// Создаем загрузчик
	downloader, err := NewBatchMediaDownloader(minioConfig, 5) // 5 одновременных загрузок
	if err != nil {
		log.Fatalf("Ошибка создания загрузчика: %v", err)
	}

	// Загружаем медиафайлы
	ctx := context.Background()
	results := downloader.DownloadAndUploadMedia(ctx, mediaTasks)

	// Сохраняем информацию о медиа в базу
	for i, result := range results {
		task := mediaTasks[i]
		if result.Success {
			// Сохраняем MinIO URL в базу
			path := fmt.Sprintf("%s/%s", minioConfig.Bucket, task.ObjectName)
			minioURL := fmt.Sprintf("https://%s/%s", minioConfig.Endpoint, path)

			id := laptopIDs[task.ProductID-1]

			_, err := db.Exec("INSERT INTO Medias (url, type, product_id) VALUES ($1, $2, $3)",
				path, task.Type, id)

			if err != nil {
				log.Printf("Ошибка сохранения медиа в базу: %v", err)
			} else {
				log.Printf("Успешно загружено: %s -> %s", task.URL, minioURL)
			}
		} else {
			log.Printf("Ошибка загрузки %s: %v", task.URL, result.Error)
		}
	}
}

func InsertKeyboards(db *sql.DB) []uint64 {
	keyboards := []models.KeyboardChars{
		{
		Name:          "Red Square",
		TypeKeyBoards: "Mechanic",
		Switches:      "Yellow",
		ReleaseYear:   2023,
		},
		{
		Name:          "Red Dragon",
		TypeKeyBoards: "Mechanic",
		Switches:      "Blue",
		ReleaseYear:   2020,
		},
		{
		Name:          "Razer",
		TypeKeyBoards: "Mechanic",
		Switches:      "Red",
		ReleaseYear:   2017,
		},
	}

	idk := make([]uint64, 0, len(keyboards))

	for _, keyboard := range keyboards {
		var (
			charkId uint64
		)
		err := db.QueryRow(fmt.Sprintf("INSERT INTO %s (name, type_key_boards, switches, release_year) VALUES ($1, $2, $3, $4) returning id", "KeyboardChars"), keyboard.Name, keyboard.TypeKeyBoards, keyboard.Switches, keyboard.ReleaseYear).Scan(&charkId)

		if err != nil {
			panic(err)
		}

		idk = append(idk, charkId)
	}

	return idk
}

func InsertMouse(db *sql.DB) []uint64{
	mouses := []models.MouseChars{
		{
			Name: "MCHOSE",
			TypeMouses: "mouse",
			Dpi: 26000,
			ReleaseYear: 2025,
		},
	}

		idm := make([]uint64, 0, len(mouses))

	for _, mouse := range mouses {
		var (
			charmId uint64
		)
		err := db.QueryRow(fmt.Sprintf("INSERT INTO %s (name, type_mouses, dpi, release_year) VALUES ($1, $2, $3, $4) returning id", "MouseChars"), mouse.Name, mouse.TypeMouses, mouse.Dpi, mouse.ReleaseYear).Scan(&charmId)

		if err != nil {
			panic(err)
		}

		idm = append(idm, charmId)
	}

	return idm
}

func InsertGpus(db *sql.DB) []uint64 {
	gpus := []models.GpuChars{
		{
			Name:         "RTX 4090",
			MemoryGB:     24,
			MemoryType:   "GDDR6X",
			BusWidthBit:  384,
			BaseFreqMHz:  2235,
			BoostFreqMHz: 2520,
			TecprocNm:    5,
			TDPWatt:      450,
			ReleaseYear:  2022,
		},
		{
			Name:         "Apple M3 Max",
			MemoryGB:     128,
			MemoryType:   "LPDDR5",
			BusWidthBit:  512,
			BaseFreqMHz:  400,
			BoostFreqMHz: 1400,
			TecprocNm:    3,
			TDPWatt:      60,
			ReleaseYear:  2023,
		},
		{
			Name:         "Apple M3",
			MemoryGB:     24,
			MemoryType:   "LPDDR5",
			BusWidthBit:  128,
			BaseFreqMHz:  400,
			BoostFreqMHz: 1100,
			TecprocNm:    3,
			TDPWatt:      30,
			ReleaseYear:  2023,
		},
		{
			Name:         "RTX 4080",
			MemoryGB:     16,
			MemoryType:   "GDDR6X",
			BusWidthBit:  256,
			BaseFreqMHz:  2210,
			BoostFreqMHz: 2510,
			TecprocNm:    5,
			TDPWatt:      320,
			ReleaseYear:  2022,
		},
		{
			Name:         "RTX 4070",
			MemoryGB:     12,
			MemoryType:   "GDDR6X",
			BusWidthBit:  192,
			BaseFreqMHz:  1920,
			BoostFreqMHz: 2475,
			TecprocNm:    5,
			TDPWatt:      200,
			ReleaseYear:  2023,
		},
		{
			Name:         "RTX 4060",
			MemoryGB:     8,
			MemoryType:   "GDDR6",
			BusWidthBit:  128,
			BaseFreqMHz:  1830,
			BoostFreqMHz: 2460,
			TecprocNm:    5,
			TDPWatt:      115,
			ReleaseYear:  2023,
		},
		{
			Name:         "RTX 4050",
			MemoryGB:     6,
			MemoryType:   "GDDR6",
			BusWidthBit:  96,
			BaseFreqMHz:  1600,
			BoostFreqMHz: 2200,
			TecprocNm:    5,
			TDPWatt:      75,
			ReleaseYear:  2023,
		},
		{
			Name:         "Intel Arc A770",
			MemoryGB:     16,
			MemoryType:   "GDDR6",
			BusWidthBit:  256,
			BaseFreqMHz:  2100,
			BoostFreqMHz: 2400,
			TecprocNm:    6,
			TDPWatt:      225,
			ReleaseYear:  2022,
		},
	}

	idg := make([]uint64, 0, len(gpus))

	for _, gpu := range gpus {
		var (
			chargId uint64
		)
		err := db.QueryRow(fmt.Sprintf("INSERT INTO %s (name, memory_gb, memory_type, bus_width_bit, base_freq_mhz, boost_freq_mhz, tecproc_nm, tdp_watt, release_year) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id", "GpuChars"), gpu.Name, gpu.MemoryGB, gpu.MemoryType, gpu.BusWidthBit, gpu.BaseFreqMHz, gpu.BoostFreqMHz, gpu.TecprocNm, gpu.TDPWatt, gpu.ReleaseYear).Scan(&chargId)

		if err != nil {
			panic(err)
		}

		idg = append(idg, chargId)
	}

	return idg
}

func InsertLaptops(db *sql.DB, minioConfig MinIOConfig) {
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

	gpuIds := InsertGpus(db)

	laptops := []struct {
		Name   string
		Price  float64
		Selled uint64
		Stock  uint64
		CpuID  uint64
		Ram    int16
		Gpu    string
		GpuId  uint64
	}{
		{"MSI Titan 18", 614999, 0, 13, ids[0], 32, "RTX 4090", gpuIds[0]},
		{"Lenovo Legion Y9000P", 425999, 32, 28, ids[0], 32, "RTX 4090", gpuIds[0]},
		{"Apple Macbook Pro", 421599, 3, 83, ids[1], 48, "Apple M3 Max", gpuIds[1]},
		{"MSI Raider 18", 417999, 10, 38, ids[0], 32, "RTX 4090", gpuIds[0]},
		{"ASUS ROG Zephyrus Duo 16", 389999, 5, 47, ids[2], 32, "RTX 4090", gpuIds[0]},
		{"MSI Vector 17", 374999, 5, 123, ids[0], 32, "RTX 4080", gpuIds[3]},
		{"ASUS VivoBook Pro 15", 179999, 42, 82, ids[3], 24, "RTX 4060", gpuIds[5]},
		{"MSI Sword 17 HX", 179999, 6, 35, ids[4], 16, "RTX 4070", gpuIds[4]},
		{"MSI Summit 13 AI+ Evo", 179999, 0, 4, ids[5], 32, "Intel Arc Graphics", gpuIds[7]},
		{"Honor MagicBook Art 14", 149999, 2, 64, ids[6], 32, "Intel Arc Graphics", gpuIds[7]},
		{"Apple MacBook Air", 179499, 30, 32, ids[7], 16, "Apple M3", gpuIds[2]},
		{"Dell XPS 17", 359999, 15, 50, ids[0], 32, "RTX 4080", gpuIds[3]},
		{"Dell XPS 15", 289999, 20, 40, ids[3], 16, "RTX 4070", gpuIds[4]},
		{"HP Omen 17", 299999, 12, 33, ids[4], 16, "RTX 4070", gpuIds[4]},
		{"HP Spectre x360 14", 189999, 8, 60, ids[6], 16, "Intel Arc Graphics", gpuIds[7]},
		{"Razer Blade 18", 449999, 5, 20, ids[0], 32, "RTX 4090", gpuIds[0]},
		{"Razer Blade 16", 379999, 7, 22, ids[2], 32, "RTX 4080", gpuIds[3]},
		{"Gigabyte Aorus 17", 319999, 6, 15, ids[4], 32, "RTX 4070", gpuIds[4]},
		{"Gigabyte Aero 16 OLED", 299999, 10, 18, ids[3], 32, "RTX 4060", gpuIds[5]},
		{"Samsung Galaxy Book4 Ultra", 279999, 4, 25, ids[6], 32, "RTX 4070", gpuIds[4]},
		{"Samsung Galaxy Book4 Pro", 219999, 12, 30, ids[6], 16, "Intel Arc Graphics", gpuIds[7]},
		{"Acer Predator Helios 18", 349999, 9, 27, ids[0], 32, "RTX 4080", gpuIds[3]},
		{"Acer Predator Helios 16", 309999, 10, 21, ids[4], 32, "RTX 4070", gpuIds[4]},
		{"Acer Swift X 14", 159999, 11, 42, ids[6], 16, "RTX 4050", gpuIds[6]},
		{"ASUS Zenbook Pro Duo 14", 279999, 6, 34, ids[3], 32, "RTX 4060", gpuIds[5]},
		{"ASUS TUF Gaming F15", 199999, 40, 50, ids[4], 16, "RTX 4060", gpuIds[5]},
		{"ASUS TUF Gaming A17", 219999, 22, 44, ids[2], 16, "RTX 4060", gpuIds[5]},
		{"Lenovo Yoga Slim 7", 189999, 30, 38, ids[6], 16, "Intel Arc Graphics", gpuIds[7]},
		{"Lenovo IdeaPad Gaming 3", 159999, 28, 45, ids[4], 16, "RTX 4050", gpuIds[6]},
		{"Lenovo Legion Slim 5", 279999, 18, 32, ids[2], 32, "RTX 4070", gpuIds[4]},
		{"Lenovo ThinkPad X1 Carbon", 259999, 12, 40, ids[6], 16, "Intel Arc Graphics", gpuIds[7]},
		{"Microsoft Surface Laptop Studio 2", 299999, 8, 23, ids[3], 32, "RTX 4060", gpuIds[5]},
		{"Microsoft Surface Laptop 6", 189999, 5, 30, ids[6], 16, "Intel Arc Graphics", gpuIds[7]},
		{"MSI Prestige 16", 209999, 3, 28, ids[6], 16, "RTX 4050", gpuIds[6]},
		{"MSI Stealth 17 Studio", 329999, 9, 20, ids[0], 32, "RTX 4080", gpuIds[3]},
		{"MSI Katana 15", 179999, 17, 45, ids[4], 16, "RTX 4060", gpuIds[5]},
		{"MSI Cyborg 14", 149999, 12, 55, ids[6], 16, "RTX 4050", gpuIds[6]},
		{"Alienware m18", 469999, 4, 14, ids[0], 32, "RTX 4090", gpuIds[0]},
		{"Alienware x16", 419999, 5, 19, ids[2], 32, "RTX 4080", gpuIds[3]},
		{"Alienware m16", 339999, 6, 26, ids[4], 32, "RTX 4070", gpuIds[4]},
		{"Alienware x14", 259999, 8, 33, ids[6], 16, "RTX 4060", gpuIds[5]},
		{"Huawei MateBook X Pro", 179999, 10, 40, ids[6], 16, "Intel Arc Graphics", gpuIds[7]},
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

	InsertMedias(db, minioConfig, laptop_ids)
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

func Clear(db *sql.DB) {
	tables := []string{
		"Cart",
		"Users",
		"LaptopChars",
		"CpuChars",
		"Categories",
		"Medias",
		"Products",
		"comments",
		"commentreactions",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))

		if err != nil {
			log.Fatalf("Error while deleting from %s: %s", table, err.Error())
		}
	}
}

func InsertAll(db *sql.DB, minioConfig MinIOConfig) {
	InsertUsers(db)
	InsertLaptops(db, minioConfig)
	InsertCategories(db)
	InsertGpus(db)
}

func GetMinIOConfig(cfg *config.Config) MinIOConfig {
	bucket := cfg.MinIOConn.Bucket

	if bucket == "" {
		bucket = "pccore"
	}

	return MinIOConfig{
		Endpoint:  cfg.MinIOConn.Ep,
		AccessKey: os.Getenv("MINIO_ACCESS"),
		SecretKey: os.Getenv("MINIO_SECRET"),
		Bucket:    bucket,
		UseSSL:    cfg.MinIOConn.Secure,
	}
}

func formatHelpMessage() string {
	return fmt.Sprintf("Setup seeds with MinIO media download.\n\nUsage:\n\tgo run seeds.go [SEED TYPES]\n\nSeed Types:\n\t%s", strings.Join(SEED_TYPE_NAMES, "\n\t"))
}

func handleCliArgs(args []string, db *sql.DB, config MinIOConfig, seedTypes map[string]func(*sql.DB, MinIOConfig)) {
	if len(args) == 1 {
		fmt.Println(formatHelpMessage())
		return
	}

	for _, arg := range args[1:] {
		fn, ok := seedTypes[arg]
		if !ok {
			fmt.Printf("Error: unknown parameter: %s\n", arg)
			os.Exit(1)
		}
		fn(db, config)
	}
}

func main() {
	cfg, err := config.ParseConfig("cfg.yml")

	if err != nil {
		panic(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	minioConfig := GetMinIOConfig(cfg)

	db, err := sql.Open("postgres", os.Getenv("PCCORE_POSTGRES_CONN"))
	if err != nil {
		log.Fatalf("Error while connecting the db: %v", err)
	}
	defer db.Close()

	SEED_TYPES := map[string]func(*sql.DB, MinIOConfig){
		SEED_TYPE_USER:     func(db *sql.DB, config MinIOConfig) { InsertUsers(db) },
		SEED_TYPE_LAPTOP:   InsertLaptops,
		SEED_TYPE_CATEGORY: func(db *sql.DB, config MinIOConfig) { InsertCategories(db) },
		SEED_TYPE_MEDIA: func(db *sql.DB, config MinIOConfig) {

		},
		SEED_ALL:      InsertAll,
		CLEAR_ALL:     func(db *sql.DB, config MinIOConfig) { Clear(db) },
		SEED_TYPE_GPU: func(d *sql.DB, mi MinIOConfig) { InsertGpus(db) },
	}

	handleCliArgs(os.Args, db, minioConfig, SEED_TYPES)
}
