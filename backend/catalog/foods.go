package catalog

// Food adalah data menu surplus yang ditampilkan di landing page.
// Harga di sini adalah sumber kebenaran; harga dari browser tidak pernah dipakai.
type Food struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Store         string `json:"store"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	Category      string `json:"category"`
	OriginalPrice int    `json:"original_price"`
	Price         int    `json:"price"`
	Stock         int    `json:"stock"`
}

var Foods = []Food{
	{
		ID:            "roti-artisan",
		Name:          "Surplus Roti Artisan",
		Store:         "Rotikita Bakery",
		Description:   "Sourdough, croissant, dan roti manis panggang hari ini. Masih lembut, rasa tetap juara.",
		Image:         "/assets/foods/roti-artisan.jpg",
		Category:      "Bakery",
		OriginalPrice: 50000,
		Price:         25000,
		Stock:         12,
	},
	{
		ID:            "nasi-box-ayam",
		Name:          "Nasi Box Ayam Bakar",
		Store:         "Dapoer Nusantara",
		Description:   "Nasi hangat dengan ayam bakar bumbu kecap, lalapan, dan sambal rumahan.",
		Image:         "/assets/foods/nasi-box.jpg",
		Category:      "Makanan Berat",
		OriginalPrice: 35000,
		Price:         18000,
		Stock:         8,
	},
	{
		ID:            "pastry-coffee",
		Name:          "Paket Pastry & Cold Brew",
		Store:         "Kopi Pagi Senopati",
		Description:   "Dua pastry butter dan satu botol cold brew 250 ml. Cocok untuk teman kerja sore.",
		Image:         "/assets/foods/pastry-coffee.jpg",
		Category:      "Minuman & Pastry",
		OriginalPrice: 60000,
		Price:         30000,
		Stock:         10,
	},
	{
		ID:            "sayur-buah",
		Name:          "Paket Sayur & Buah Segar",
		Store:         "Pasar Segar Batam",
		Description:   "Sayuran hijau dan buah musiman pilihan yang masih segar, diselamatkan dari sisa display.",
		Image:         "/assets/foods/sayur-buah.jpg",
		Category:      "Segar",
		OriginalPrice: 30000,
		Price:         15000,
		Stock:         15,
	},
	{
		ID:            "sushi-set",
		Name:          "Sushi Set Salmon",
		Store:         "Sakura Kitchen",
		Description:   "Delapan potong sushi salmon dan maki sayur, dibuat pagi ini dan siap disantap.",
		Image:         "/assets/foods/sushi-set.jpg",
		Category:      "Makanan Berat",
		OriginalPrice: 70000,
		Price:         35000,
		Stock:         6,
	},
	{
		ID:            "brownies-fudgy",
		Name:          "Brownies Fudgy",
		Store:         "Sweet Lab",
		Description:   "Brownies cokelat pekat dengan tekstur fudgy, dipotong besar dan siap dibawa pulang.",
		Image:         "/assets/foods/brownies.jpg",
		Category:      "Dessert",
		OriginalPrice: 40000,
		Price:         20000,
		Stock:         9,
	},
}

func Find(id string) (Food, bool) {
	for _, food := range Foods {
		if food.ID == id {
			return food, true
		}
	}
	return Food{}, false
}
