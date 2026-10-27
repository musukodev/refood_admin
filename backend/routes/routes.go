package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"refood-backend/backend/config"
	"refood-backend/backend/payment"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterAPIRoutes(r *gin.Engine, cfg *config.Config, paymentHandler *payment.Handler) {
	// Health check
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"app":         "ReFood API",
			"environment": cfg.Environment,
		})
	})

	// Payment API routes (Simulasi & Integrasi Duitku)
	r.GET("/api/foods", paymentHandler.HandleFoods)
	r.POST("/api/checkout", paymentHandler.HandleCheckout)
	r.GET("/api/checkout/status", paymentHandler.HandleStatus)
	r.POST("/api/callback", paymentHandler.HandleCallback)
	r.POST("/api/simulate-payment", paymentHandler.HandleSimulate)
	r.GET("/api/payment-methods", paymentHandler.HandlePaymentMethods)

	// API route for login
	r.POST("/api/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username/Email dan Password harus diisi"})
			return
		}

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
				"chartData": []int{
					15, 18, 22,
					20, 25, 30,
					38, 42, 36,
					48, 65, 88,
					125, 95, 60,
					45, 38, 30,
					42, 68, 85,
					98, 105, 110,
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
		menunggu := []gin.H{}
		for i := 0; i < 3; i++ {
			menunggu = append(menunggu, gin.H{
				"id":       i + 1,
				"name":     "Artisan Bakery",
				"owner":    "Sri Wahyuni",
				"location": "Batam Center",
				"time":     "Diajukan hari ini, 09:15 WIB",
				"status":   "Menunggu",
				"image":    "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop",
			})
		}

		terverifikasi := []gin.H{}
		for i := 0; i < 6; i++ {
			terverifikasi = append(terverifikasi, gin.H{
				"id":       i + 101,
				"name":     "Artisan Bakery",
				"owner":    "Sri Wahyuni",
				"location": "Batam Center",
				"time":     "Disetujui hari ini, 09:15 WIB",
				"status":   "Terverifikasi",
				"image":    "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"menunggu":      menunggu,
			"terverifikasi": terverifikasi,
		})
	})

	// API route for manajemen pengguna
	r.GET("/api/pengguna", func(c *gin.Context) {
		var semua []gin.H

		// 5 Toko Terverifikasi
		for i := 0; i < 5; i++ {
			semua = append(semua, gin.H{
				"id":       i + 1,
				"name":     "Morning Bakery",
				"subtitle": "Budi Santoso (Pemilik) • morningbakery@gmail.com",
				"type":     "Usaha",
				"location": "Batam Center, Kepri",
				"status":   "Aktif",
				"avatar":   "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=150&h=150&fit=crop",
				"isStore":  true,
			})
		}

		// 5 Individu (Aktif)
		for i := 5; i < 10; i++ {
			semua = append(semua, gin.H{
				"id":       i + 1,
				"name":     "Budi Santoso",
				"subtitle": "budisantoso@gmail.com",
				"type":     "Individu",
				"location": "Batam Center, Kepri",
				"status":   "Aktif",
				"avatar":   "BS",
				"isStore":  false,
			})
		}

		// 5 Individu (Diblokir)
		for i := 10; i < 15; i++ {
			semua = append(semua, gin.H{
				"id":       i + 1,
				"name":     "Budi Santoso",
				"subtitle": "budisantoso@gmail.com",
				"type":     "Individu",
				"location": "Batam Center, Kepri",
				"status":   "Diblokir",
				"avatar":   "BS",
				"isStore":  false,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"total":    1428,
			"toko":     174,
			"diblokir": 5,
			"users":    semua,
		})
	})
}
