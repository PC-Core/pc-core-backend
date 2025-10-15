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
	SEED_TYPE_MEDIA    = "media"
	SEED_ALL           = "all"
	CLEAR_ALL          = "clear"
	SEED_TYPE_GPU      = "gpu"
	SEED_TYPE_KEYBOARD = "keyboard"
)

var SEED_TYPE_NAMES = []string{SEED_TYPE_USER, SEED_TYPE_LAPTOP, SEED_TYPE_CATEGORY, SEED_TYPE_MEDIA, SEED_ALL, CLEAR_ALL, SEED_TYPE_GPU}

type MediaDownloadTask struct {
	URL        string
	ObjectName string
	ProductID  uint64
	Type       string
}

func Sha256(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
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

func InsertMedias(db *sql.DB, laptopIDs []uint64) {
	dowloadTask := []MediaDownloadTask{{"https://cdn1.ozone.ru/s3/multimedia-2/6714764426.jpg", "laptops/msi-titan-1.jpg", 1, "Image"},
		{"https://www.ixbt.com/img/r30/00/02/56/91/msi-titan-gt77-12uhs-big.jpg", "laptops/msi-titan-2.jpg", 1, "Image"},
		{"https://rutube.ru/video/24f70de9c47c9df1be60adfbd6b3d3db/?ysclid=mf7837kqg3115263121", "laptops/msi-titan-review.mp4", 1, "Video"},
		{"https://i.ebayimg.com/images/g/IUoAAeSwPG5oMKqS/s-l500.jpg", "laptops/lenovo-legion-1.jpg", 2, "Image"},
		{"https://cdn1.youla.io/files/images/780_780/65/df/65df0a0dc3773b77740cf344-2.jpg", "laptops/lenovo-legion-2.jpg", 2, "Image"},
		{"https://n.cdn.cdek.shopping/images/shopping/9257a5ea990e4f7fb298cf575b4dd336.jpg?v=1", "laptops/macbook-pro-1.jpg", 3, "Image"},
		{"https://appgreatstore.ru/upload/iblock/1e9/ndvccnc38q1a2l6i3yzdn45jdkkysxe5.jpg", "laptops/macbook-pro-2.jpg", 3, "Image"},
		{"https://static.1k.by/images/productsimages/ip/big/pp6/4/4931862/i128c38fef.jpg", "laptops/MSI-Raider-18", 4, "Image"},
		{"https://overclockers.ru/st/legacy/blog/413830/606041_O.jpg", "laptops/MSI-Raider-18", 4, "Image"},
		{"https://thedigitaltech.com/wp-content/uploads/2022/07/Asus-ROG-Zephyrus-Duo-16.jpg", "laptops/ASUS-ROG-Zephyrus-Duo-16", 5, "Image"},
		{"https://wit.ru/images/1/10625/10627/10640/90NR0D71-M000X0-421567.png", "laptops/ASUS-ROG-Zephyrus-Duo-16", 5, "Image"},
		{"https://static.1k.by/images/productsimages/ip/big/pp6/4/4931862/i128c38fef.jpg", "laptops/MSI-Vector-17", 6, "Image"},
		{"https://img.mvideo.ru/Pdb/400293516b4.jpg", "laptops/MSI-Vector-17", 6, "Image"},
		{"https://comparema.ru/image/cache/catalog/products/3441782_3-1200x800.jpg", "laptops/ASUS-VivoBook-Pro-15", 7, "Image"},
		{"https://cdn1.ozone.ru/s3/multimedia-u/6479990442.jpg", "laptops/ASUS-VivoBook-Pro-15", 7, "Image"},
		{"https://tm.by/sites/default/files/styles/uc_product_full/public/content/product/2024/11/2/i835792-20241126_0.jpg", "laptops/MSI-Sword-17-HX", 8, "Image"},
		{"https://main-cdn.sbermegamarket.ru/big2/hlr-system/166/078/535/511/615/44/100071445455b5.png", "laptops/MSI-Sword-17-HX", 8, "Image"},
		{"https://m.media-amazon.com/images/I/714G8bC-F1L.jpg", "laptops/MSI-Summit-13-AI+-Evo", 9, "Image"},
		{"https://avatars.mds.yandex.net/i?id=7487728a923587b878ff5e9c0bbfb7d0780eb10a-16485021-images-thumbs&n=13", "laptops/MSI-Summit-13-AI+-Evo", 9, "Image"},
		{"https://mews.biggeek.ru/wp-content/uploads/2022/07/bfarsace_190101_5333_0005.jpg", "laptops/Apple-MacBook-Air", 10, "Image"},
		{"https://overclockers.ru/st/legacy/blog/422120/358217_O.jpg", "laptops/Apple-MacBook-Air", 10, "Image"},
		{"https://m.media-amazon.com/images/I/71RuOUPT+4L.jpg", "laptops/Dell-XPS-17", 11, "Image"},
		{"https://m.media-amazon.com/images/I/71Kn5i5fsrL.jpg", "laptops/Dell-XPS-17", 11, "Image"},
		{"https://m.media-amazon.com/images/I/71RuOUPT+4L.jpg", "laptops/Dell-XPS-15", 12, "Image"},
		{"https://m.media-amazon.com/images/I/71RuOUPT+4L.jpg", "laptops/Dell-XPS-15", 12, "Image"},
		{"https://cache3.youla.io/files/images/780_780/5b/03/5b0306f0e7696a72e95fb742.jpg", "laptops/HP-Omen-17", 13, "Image"},
		{"https://i.pinimg.com/736x/08/18/d4/0818d4903ca522622de32b48ccdcd932.jpg", "laptops/HP-Omen-17", 13, "Image"},
		{"https://pigmentarius.ru/upload/medialibrary/3ac/0xxtl71pz9280bvq0pzcz2t9wmvyu6mr.jpg", "laptops/HP-Spectre-x360-14", 14, "Image"},
		{"https://main-cdn.sbermegamarket.ru/big2/hlr-system/1628142416/100023661450b5.jpg", "laptops/HP-Spectre-x360-14", 14, "Image"},
		{"https://static01.servicebox.ru/img/uploads/news/6614e1b3cf74f_1712644531.jpeg", "laptops/Razer-Blade-18", 15, "Image"},
		{"https://www.notebookcheck-ru.com/uploads/tx_nbc2/blade_01.jpg", "laptops/Razer-Blade-18", 15, "Image"},
		{"https://www.notebookcheck-ru.com/uploads/tx_nbc2/blade_01.jpg", "laptops/Razer-Blade-18", 16, "Image"},
		{"https://www.notebookcheck-ru.com/uploads/tx_nbc2/blade_01.jpg", "laptops/Razer-Blade-18", 16, "Image"},
		{"https://www.cdrinfo.com/d7/system/files/styles/siteberty_image_770x484/private/new_site_image/2020/gigabyte_aorus17g.jpg?itok=FSbVJYlW", "laptops/Gigabyte-Aorus-17", 17, "Image"},
		{"https://www.cdrinfo.com/d7/system/files/styles/siteberty_image_770x484/private/new_site_image/2020/gigabyte_aorus17g.jpg?itok=FSbVJYlW", "laptops/Gigabyte-Aorus-17", 17, "Image"},
		{"https://c.dns-shop.ru/thumb/st4/fit/wm/0/0/549a0b1bc8b1b5e200fd28edaa260295/84560d0d011ac958db8e55cf001f4ae1f76900dab0409de57fa676540a3967aa.jpg.webp", "laptops/Gigabyte-Aero-16-OLED", 18, "Image"},
		{"https://c.dns-shop.ru/thumb/st4/fit/0/0/794102e30c43ada217a69495e290e1ad/946fdc59e840364576d5e71525cc5d35e2438d8622914c753d36959a225d1b4e.jpg.webp", "laptops/Gigabyte-Aero-16-OLED", 18, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/ccadb1f98aa909e6ad96a9818607fdd3/f5cd838f524900de7a6151f32173465652ca351aa69cfd5386c284c9728c7b7d.jpg.webp", "laptops/Samsung-Galaxy-Book4-Ultra", 19, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/2d42b2861daa886897c3241892c3fa89/90ebe9e40b50a636bb77d99c8200f3222c1c21d58bf223c216305ceedbe03c6f.jpg.webp", "laptops/Samsung-Galaxy-Book4-Ultra", 19, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/2d6bcec68df7c383e7a07dcab9515adf/148e2980db1a3a10d4777d448262fab966000e090948d2d33642298a1bce8fb8.jpg.webp", "laptops/Samsung-Galaxy-Book4-Pro", 20, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/addf544e5b369a9fdfedec1be4ae8319/88ada83908a3b640c579f699459808ac885a4f573c7ced103644563f59ade509.jpg.webp", "laptops/Samsung-Galaxy-Book4-Pro", 20, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/ecd0d4a753caa2e3458411c0740cc800/6f11a7a1c608fcc6908402d94df217bd326a975d4c9e136d47123dc1db3e401a.jpg.webp", "laptops/Acer-Predator-Helios-18", 21, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/badaf84f3e48552f609f49b9f50bdfad/dee83ba14cf5b284f468e596a2e14e2b77a044099fe9ea697b74af3e553f8764.jpg.webp", "laptops/Acer-Predator-Helios-18", 21, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/ecd0d4a753caa2e3458411c0740cc800/6f11a7a1c608fcc6908402d94df217bd326a975d4c9e136d47123dc1db3e401a.jpg.webp", "laptops/Acer-Predator-Helios-16", 22, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/badaf84f3e48552f609f49b9f50bdfad/dee83ba14cf5b284f468e596a2e14e2b77a044099fe9ea697b74af3e553f8764.jpg.webp", "laptops/Acer-Predator-Helios-16", 22, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/3036587_1.jpg", "laptops/Acer-Swift-X-14", 23, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/3036587_2.jpg", "laptops/Acer-Swift-X-14", 23, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/noutbuk_asus_zenbook_duo_ux8406ca_ql221w_duo_touch_14_14_fhd_fhd_oled_60hz_400nits_touch_ultra_7_255h_2.0ghz_16gb_lpddr5x_ssd_1tb_arc_graphics_win11_inkwell_gray_90nb14x1_m00c70_5.jpg", "laptops/ASUS-Zenbook-Pro-Duo-14", 24, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/noutbuk_asus_zenbook_duo_ux8406ca_ql221w_duo_touch_14_14_fhd_fhd_oled_60hz_400nits_touch_ultra_7_255h_2.0ghz_16gb_lpddr5x_ssd_1tb_arc_graphics_win11_inkwell_gray_90nb14x1_m00c70_1.jpg", "laptops/ASUS-Zenbook-Pro-Duo-14", 24, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/241630471ce2bbc8a55673f6c7ac2ed4/5b5a9a787bdb5e970bea99f60f829bd57b53cd02179d2e0c0bd1a68eb9d6a8ce.jpg.webp", "laptops/ASUS-TUF-Gaming-F15", 25, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/644fecaa0ca3b884d037badf707f9365/dd1e4cbb67650916f2d89b1a6f0c9ca357fe23d99772bc5f628b700a769fe49e.jpg.webp", "laptops/ASUS-TUF-Gaming-F15", 25, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/noutbuk_asus_tuf_gaming_a17_fa706nfr_hx007_17.3_fhd_ips_144hz_250nits_ryzen_7_7435hs_3.1ghz_16gb_ddr5_ssd_512gb_rtx_2050_4gb_noos_graphite_black_90nr0jw5_m00080__3276410_1.jpg", "laptops/ASUS-TUF-Gaming-A17", 26, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/noutbuk_asus_tuf_gaming_a17_fa706nfr_hx007_17.3_fhd_ips_144hz_250nits_ryzen_7_7435hs_3.1ghz_16gb_ddr5_ssd_512gb_rtx_2050_4gb_noos_graphite_black_90nr0jw5_m00080__3276410_4.jpg", "laptops/ASUS-TUF-Gaming-A17", 26, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/db7664bd13f5ac044f5aba55b3a375c1/be6c48142d7cd590b9ba88b6686ea1dc676773e97040233c356b76a117fc1cb8.jpg.webp", "laptops/Lenovo-Yoga-Slim-7", 27, "Image"},
		{"https://c.dns-shop.ru/thumb/st4/fit/wm/0/0/ebc20c3e48d27a12c2ca199f5aa71bab/b85df134be0d2db537543b57048fdb16509ea80e069b0aa78672e608764047c4.jpg.webp", "laptops/Lenovo-Yoga-Slim-7", 27, "Image"},
		{"https://c.dns-shop.ru/thumb/st4/fit/0/0/817b8391b8332efec99cb572dc4aabc6/63add0435dc2ca1f4082d79c725ff635bf24115760a1825efa496a176cf2d807.jpg.webp", "laptops/Lenovo-IdeaPad-Gaming-3", 28, "Image"},
		{"https://c.dns-shop.ru/thumb/st4/fit/wm/0/0/24960407824a70fd350b53853a4f1ac1/9eb549b1f29c510c7198a4c60b90552bf3bb880965c8b3abfca1e65d77502158.jpg.webp", "laptops/Lenovo-IdeaPad-Gaming-3", 28, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/f9ef35fd9ec68a0601a890e5d7b0e3f1/a003d1c674eeb891feea71ae6e2aa305247a90817ce780a2b8601e2ffecce9e3.png.webp", "laptops/Lenovo-Legion-Slim-5", 29, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/b080a238b6102fcaebcddfbc0c4a8faf/0ee6a10b1fc35684de0b028f195d098f7dd18ff036234ec79f6940ce925eba0a.jpg.webp", "laptops/Lenovo-Legion-Slim-5", 29, "Image"},
		{"https://cdn.citilink.ru/PvviYFG7xLrMc1Cg-iICyATPKKzDq79G4Cw5ddM-b5M/resizing_type:fit/gravity:sm/width:1200/height:1200/plain/product-images/4c734b8b-2ef5-413d-9f1a-cab6f65593e1.jpg", "laptops/Lenovo-ThinkPad-X1-Carbon", 30, "Image"},
		{"https://cdn.citilink.ru/IgQMw6ci2q8NkyvQ4QyBgHamSC9AdRoH5rQi3QYSalU/resizing_type:fit/gravity:sm/width:1200/height:1200/plain/product-images/1210e825-ca43-4142-8aef-972270d89dcc.jpg", "laptops/Lenovo-ThinkPad-X1-Carbon", 30, "Image"},
		{"https://microsoft-surface.ru/wp-content/uploads/2023/12/Surface-Laptop-Studio-2-Store-1200.png", "laptops/Microsoft-Surface-Laptop-Studio-2", 31, "Image"},
		{"https://microsoft-surface.ru/wp-content/uploads/2023/12/Surface-Laptop-Studio-2-Store-3.png", "laptops/Microsoft-Surface-Laptop-Studio-2", 31, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/noutbuk_microsoft_surface_laptop_6_13.5_2.2k_ips_60hz_400nits_ultra_7_165h_32gb_lpddr5x_ssd_512gb_arc_graphics_win11pro_platinum_zjz_00026__3249886_1.jpg", "laptops/Microsoft-Surface-Laptop-6", 32, "Image"},
		{"https://static.onlinetrade.ru/img/items/b/noutbuk_microsoft_surface_laptop_6_13.5_2.2k_ips_60hz_400nits_ultra_7_165h_32gb_lpddr5x_ssd_512gb_arc_graphics_win11pro_platinum_zjz_00026__3249886_5.jpg", "laptops/Microsoft-Surface-Laptop-6", 32, "Image"},
		{"https://cdn.citilink.ru/kkq4UUlKVQib0-2V5XmdZsue_ACniIB6cgruX06rV-4/resizing_type:fit/gravity:sm/width:1200/height:1200/plain/product-images/1fd8acf5-1c44-4c7d-a0c8-2f998d65ed64.jpg", "laptops/MSI-Prestige-16", 33, "Image"},
		{"https://cdn.citilink.ru/g3UdE4Ui0mkg7dL3S1OYwPaeTTFj92Pfujpx3eJ0z9Q/resizing_type:fit/gravity:sm/width:1200/height:1200/plain/product-images/1c78f760-55f8-4585-9efd-d7f048871b5a.jpg", "laptops/MSI-Prestige-16", 33, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/d3f47e0e9458fb5c347874db77a306b6/f63ee020c5873fec66c00dfa2f47a7a60902fd443c001736826fad6a3d16b150.jpg.webp", "laptops/MSI-Stealth-17-Studio", 34, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/3f5142dcd60d2504af694ee6b97dfdee/ca07bcbef732933a6e2a6ae752564983178a83a4e01f5ad5082af1a24b1837a0.jpg.webp", "laptops/MSI-Stealth-17-Studio", 34, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/8b2fa800deb6557c4b215ce0b121dc33/64cad8497f91c7b926d47580ba6e67c39534e707153591a1b161302c07634135.jpg.webp", "laptops/MSI-Katana-15", 35, "Image"},
		{"https://c.dns-shop.ru/thumb/st4/fit/wm/0/0/7f9a8843d29232957e81e9f665a5fc37/b0729b24605d841a23aa7f829bdb4d5af6f9078d63e57c328bdc9b79823f3e20.jpg.webp", "laptops/MSI-Katana-15", 35, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/0/0/2c7521d1a86370d26ddbee231015ff54/e15f3d8f62b6418efc287bd8221c21c4477f9b79564d57b9193e0225dc8eff16.jpg.webp", "laptops/MSI-Cyborg-14", 36, "Image"},
		{"https://c.dns-shop.ru/thumb/st1/fit/wm/0/0/92b4ebcc14400e849f603aef35e7d122/ea5d898b4c4bad39b6e1d92f13869c6314a24e88e8dc14d809b6b1581a55f5e1.jpg.webp", "laptops/MSI-Cyborg-14", 36, "Image"},
		{"https://n.cdn.cdek.shopping/images/shopping/34e1c9f001e446c9bdf7d6f8271deb46.jpg?v=1", "laptops/Alienware-m18", 37, "Image"},
		{"https://n.cdn.cdek.shopping/images/shopping/ddf82cbbb2c94ddfa887e0e6412f0274.jpg?v=1", "laptops/Alienware-m18", 37, "Image"},
		{"https://n.cdn.cdek.shopping/images/shopping/e4b7aac10c89442e9266a203e5c81881.jpg?v=1", "laptops/Alienware-x16", 38, "Image"},
		{"https://n.cdn.cdek.shopping/images/shopping/3d188eacbed045f78ff1b4e2366d446e.jpg?v=1", "laptops/Alienware-x16", 38, "Image"},
		{"https://optim.tildacdn.com/stor3637-3135-4666-b735-393465306166/-/format/webp/17462120.jpg.webp", "laptops/Alienware-m16", 39, "Image"},
		{"https://optim.tildacdn.com/stor3061-6233-4639-b032-323939333236/-/format/webp/83179001.jpg.webp", "laptops/Alienware-m16", 39, "Image"},
		{"https://n.cdn.cdek.shopping/images/shopping/2d8673f5f3aa4219976606f8ffbfcc05.jpg?v=1", "laptops/Alienware-x14", 40, "Image"},
		{"https://n.cdn.cdek.shopping/images/shopping/c410e4f8a5f24ceeab09537c72321624.jpg?v=1", "laptops/Alienware-x14", 40, "Image"},
		{"https://ir.ozone.ru/s3/multimedia-1-9/wc1000/7663191237.jpg", "laptops/Huawei-MateBook-X-Pro", 41, "Image"},
		{"https://ir.ozone.ru/s3/multimedia-1-9/wc1000/7663191273.jpg", "laptops/Huawei-MateBook-X-Pro", 41, "Image"}}

	for _, task := range dowloadTask {
		id := laptopIDs[task.ProductID-1]

		_, err := db.Exec("INSERT INTO Medias (url, type, product_id) VALUES ($1, $2, $3)",
			task.URL, task.Type, id)

		if err != nil {
			log.Printf("Ошибка сохранения медиа в базу: %v", err)
		} else {
			log.Printf("Успешно загружено")
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

func InsertMouse(db *sql.DB) []uint64 {
	mouses := []models.MouseChars{
		{
			Name:        "MCHOSE",
			TypeMouses:  "mouse",
			Dpi:         26000,
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

func InsertAll(db *sql.DB) {
	InsertUsers(db)
	InsertLaptops(db)
	InsertCategories(db)
	InsertGpus(db)
}

func formatHelpMessage() string {
	return fmt.Sprintf("Setup seeds with MinIO media download.\n\nUsage:\n\tgo run seeds.go [SEED TYPES]\n\nSeed Types:\n\t%s", strings.Join(SEED_TYPE_NAMES, "\n\t"))
}

func handleCliArgs(args []string, db *sql.DB, seedTypes map[string]func(*sql.DB)) {
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
		fn(db)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("PCCORE_POSTGRES_CONN"))
	if err != nil {
		log.Fatalf("Error while connecting the db: %v", err)
	}
	defer db.Close()

	SEED_TYPES := map[string]func(*sql.DB){
		SEED_TYPE_USER:     func(db *sql.DB) { InsertUsers(db) },
		SEED_TYPE_LAPTOP:   InsertLaptops,
		SEED_TYPE_CATEGORY: func(db *sql.DB) { InsertCategories(db) },
		SEED_TYPE_MEDIA: func(db *sql.DB) {

		},
		SEED_ALL:      InsertAll,
		CLEAR_ALL:     func(db *sql.DB) { Clear(db) },
		SEED_TYPE_GPU: func(d *sql.DB) { InsertGpus(db) },
	}

	handleCliArgs(os.Args, db, SEED_TYPES)
}
