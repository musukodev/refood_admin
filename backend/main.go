package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	r := gin.Default()

	// Static files
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/dashboard", "./frontend/dashboard.html")
	r.StaticFile("/verifikasi", "./frontend/verifikasi.html")
	r.StaticFile("/tinjau", "./frontend/tinjau.html")
	r.StaticFile("/pengguna", "./frontend/pengguna.html")
	r.StaticFile("/kategori", "./frontend/kategori.html")
	r.StaticFile("/tambah-kategori", "./frontend/tambah-kategori.html")
	r.StaticFile("/pelaporan", "./frontend/laporan.html")
	r.StaticFile("/laporan", "./frontend/laporan.html")
	r.StaticFile("/pengaturan", "./frontend/coming-soon.html")

	// API route for login
	r.POST("/api/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username/Email dan Password harus diisi"})
			return
		}

		// TODO: Implement actual database authentication here
		// For now, we use a dummy check
		if req.Username == "contoh@gmail.com" && req.Password == "password123" {
			c.JSON(http.StatusOK, gin.H{
				"message": "Login berhasil",
				"token":   "dummy-jwt-token-12345",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Username atau Password salah",
			})
		}
	})

	// API route for dashboard data
	r.GET("/api/dashboard", func(c *gin.Context) {
		// Mock data for the dashboard
		c.JSON(http.StatusOK, gin.H{
			"stats": gin.H{
				"activeUsers":   1428,
				"savedPortions": 215,
				"newPartners":   8,
				"incidents":     2,
			},
			"transactions": gin.H{
				"completed":    142,
				"surplusValue": "Rp 3.850.000",
				"peakInfo":     "Puncak: 12:30 (Makan Siang)",
				// Dummy data transaksi per 30 menit (06:00 s.d 17:30)
				"chartData": []int{
					15, 18, 22, // 06:00 - 07:00 (Pagi mulai buka)
					20, 25, 30, // 07:30 - 08:30 (Sarapan)
					38, 42, 36, // 09:00 - 10:00
					48, 65, 88, // 10:30 - 11:30 (Mulai ramai jelang siang)
					125, 95, 60, // 12:00 - 13:00 (Puncak jam makan siang!)
					45, 38, 30, // 13:30 - 14:30 (Mulai reda)
					42, 68, 85, // 15:00 - 16:00 (Surplus sore toko mulai di-listing)
					98, 105, 110, // 16:30 - 17:30 (Now: Penyelamatan makanan sore hari)
				},
				"chartLabels": []string{
					"06.00", "", "",
					"", "", "",
					"09.00", "", "",
					"", "", "",
					"12.00", "", "",
					"", "", "",
					"15.00", "", "",
					"", "", "17.30 (Now)",
				},
			},
			"pendingPartners": []gin.H{
				{
					"name": "Rotikita Bakery Mer...",
					"type": "Surplus roti sore • Jakp",
					"icon": "croissant",
				},
				{
					"name": "Kopi Pagi Senopati",
					"type": "Pastry & cold brew • Jal...",
					"icon": "coffee",
				},
				{
					"name": "Dapoer Nusantara M...",
					"type": "Makanan siap saji • Jaks...",
					"icon": "utensils",
				},
			},
		})
	})

	// API route for verifikasi toko
	r.GET("/api/verifikasi", func(c *gin.Context) {
		// Mock data
		menunggu := []gin.H{}
		for i := 0; i < 3; i++ {
			menunggu = append(menunggu, gin.H{
				"id":        i + 1,
				"name":      "Artisan Bakery",
				"owner":     "Sri Wahyuni",
				"location":  "Batam Center",
				"time":      "Diajukan hari ini, 09:15 WIB",
				"status":    "Menunggu",
				"image":     "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop",
			})
		}

		terverifikasi := []gin.H{}
		for i := 0; i < 6; i++ {
			terverifikasi = append(terverifikasi, gin.H{
				"id":        i + 101,
				"name":      "Artisan Bakery",
				"owner":     "Sri Wahyuni",
				"location":  "Batam Center",
				"time":      "Disetujui hari ini, 09:15 WIB",
				"status":    "Terverifikasi",
				"image":     "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"menunggu":      menunggu,
			"terverifikasi": terverifikasi,
		})
	})

	// API route for manajemen pengguna
	r.GET("/api/pengguna", func(c *gin.Context) {
		// Mock data for all users
		var semua []gin.H
		
		// 5 Toko Terverifikasi
		for i := 0; i < 5; i++ {
			semua = append(semua, gin.H{
				"id":        i + 1,
				"name":      "Morning Bakery",
				"subtitle":  "Budi Santoso (Pemilik) • morningbakery@gmail.com",
				"type":      "Usaha",
				"location":  "Batam Center, Kepri",
				"status":    "Aktif",
				"avatar":    "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop",
				"isStore":   true,
			})
		}

		// 5 Individu (Aktif)
		for i := 5; i < 10; i++ {
			semua = append(semua, gin.H{
				"id":        i + 1,
				"name":      "Budi Santoso",
				"subtitle":  "budisantoso@gmail.com",
				"type":      "Individu",
				"location":  "Batam Center, Kepri",
				"status":    "Aktif",
				"avatar":    "BS",
				"isStore":   false,
			})
		}
		
		// 5 Individu (Diblokir)
		for i := 10; i < 15; i++ {
			semua = append(semua, gin.H{
				"id":        i + 1,
				"name":      "Budi Santoso",
				"subtitle":  "budisantoso@gmail.com",
				"type":      "Individu",
				"location":  "Batam Center, Kepri",
				"status":    "Diblokir",
				"avatar":    "BS",
				"isStore":   false,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"total": 1428,
			"toko": 174,
			"diblokir": 5,
			"users": semua,
		})
	})

	log.Println("Server backend berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
