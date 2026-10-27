package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"refood-backend/backend/config"
	"refood-backend/backend/payment"
	"refood-backend/backend/routes"
)

//go:embed all:frontend
var embeddedFrontend embed.FS

// getFrontendFS mengembalikan file system frontend:
// 1. Dari disk lokal jika ada (untuk live edit saat development lokal)
// 2. Dari embedded fs jika running di environment standalone / Vercel
func getFrontendFS() (http.FileSystem, fs.FS, string) {
	candidates := []string{"./frontend", "../frontend"}
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "frontend"),
			filepath.Join(exeDir, "..", "frontend"),
		)
	}

	for _, dir := range candidates {
		info, err := os.Stat(filepath.Join(dir, "index.html"))
		if err == nil && !info.IsDir() {
			abs, err := filepath.Abs(dir)
			if err == nil {
				log.Println("[FRONTEND] Menyajikan dari disk lokal:", abs)
				return http.Dir(abs), os.DirFS(abs), abs
			}
			log.Println("[FRONTEND] Menyajikan dari disk lokal:", dir)
			return http.Dir(dir), os.DirFS(dir), dir
		}
	}

	log.Println("[FRONTEND] Menyajikan dari embedded filesystem (Vercel/Self-contained)")
	sub, err := fs.Sub(embeddedFrontend, "frontend")
	if err != nil {
		log.Fatal("Gagal memuat embedded frontend:", err)
	}
	return http.FS(sub), sub, ""
}

func serveFile(diskPath string, subFS fs.FS, filename string, contentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data []byte
		var err error

		if diskPath != "" {
			data, err = os.ReadFile(filepath.Join(diskPath, filename))
		}
		if err != nil || diskPath == "" {
			data, err = fs.ReadFile(subFS, filename)
		}

		if err != nil {
			c.String(http.StatusNotFound, "File tidak ditemukan: "+filename)
			return
		}

		if contentType != "" {
			c.Header("Content-Type", contentType)
		}
		c.Data(http.StatusOK, contentType, data)
	}
}

func main() {
	cfg := config.LoadConfig()

	paymentStore := payment.NewStore(cfg.DSN())
	paymentGateway := payment.NewGateway(cfg)
	paymentHandler := payment.NewHandler(cfg, paymentStore, paymentGateway)

	r := gin.Default()

	httpFS, subFS, diskPath := getFrontendFS()
	fileServer := http.FileServer(httpFS)

	// Static assets
	r.GET("/landing.css", serveFile(diskPath, subFS, "landing.css", "text/css; charset=utf-8"))
	r.GET("/landing.js", serveFile(diskPath, subFS, "landing.js", "application/javascript; charset=utf-8"))
	r.GET("/app.css", serveFile(diskPath, subFS, "app.css", "text/css; charset=utf-8"))
	r.GET("/app.js", serveFile(diskPath, subFS, "app.js", "application/javascript; charset=utf-8"))
	r.GET("/pembayaran.css", serveFile(diskPath, subFS, "pembayaran.css", "text/css; charset=utf-8"))
	r.GET("/pembayaran.js", serveFile(diskPath, subFS, "pembayaran.js", "application/javascript; charset=utf-8"))

	// Web Pages
	pages := map[string]string{
		"/":                "index.html",
		"/login":           "login.html",
		"/pembayaran":      "pembayaran.html",
		"/dashboard":       "dashboard.html",
		"/verifikasi":      "verifikasi.html",
		"/tinjau":          "tinjau.html",
		"/pengguna":        "pengguna.html",
		"/kategori":        "kategori.html",
		"/tambah-kategori": "tambah-kategori.html",
		"/pelaporan":       "laporan.html",
		"/laporan":         "laporan.html",
		"/detail-laporan":  "detail-laporan.html",
		"/akun":            "akun.html",
		"/pengaturan":      "coming-soon.html",
	}

	for routePath, filename := range pages {
		r.GET(routePath, serveFile(diskPath, subFS, filename, "text/html; charset=utf-8"))
	}

	// Static directory: /assets/*
	r.GET("/assets/*filepath", func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	// Shared API routes
	routes.RegisterAPIRoutes(r, cfg, paymentHandler)

	// Fallback untuk file statis lain di root frontend (misal favicon.ico)
	r.NoRoute(func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path != "" {
			var exists bool
			if diskPath != "" {
				if _, err := os.Stat(filepath.Join(diskPath, path)); err == nil {
					exists = true
				}
			}
			if !exists {
				if _, err := fs.Stat(subFS, path); err == nil {
					exists = true
				}
			}
			if exists {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Halaman tidak ditemukan"})
	})

	log.Printf("======================================================")
	log.Printf("🚀 Server ReFood aktif di %s (Port :%s)", cfg.BaseURL, cfg.Port)
	log.Printf("💳 Merchant Code: %s | Environment: %s", cfg.MerchantCode, cfg.Environment)
	log.Printf("🔗 Webhook Callback: %s", cfg.CallbackURL())
	log.Printf("🧪 Mode Simulasi Bayar: %v", cfg.EnableSimulation)
	log.Printf("======================================================")

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
