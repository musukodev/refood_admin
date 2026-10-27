package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"refood-backend/backend/config"
	"refood-backend/backend/payment"
	"refood-backend/backend/routes"
)

// findFrontendDir mencari folder frontend, baik saat server dijalankan dari
// dalam folder backend (go run main.go) maupun dari root project (go run ./backend).
func findFrontendDir() string {
	candidates := []string{
		"./frontend",
		"../frontend",
	}
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
			if err != nil {
				return dir
			}
			return abs
		}
	}

	log.Fatal("Folder frontend tidak ditemukan. Jalankan server dari folder backend atau root project.")
	return ""
}

func main() {
	cfg := config.LoadConfig()

	paymentStore := payment.NewStore()
	paymentGateway := payment.NewGateway(cfg)
	paymentHandler := payment.NewHandler(cfg, paymentStore, paymentGateway)

	r := gin.Default()

	frontendDir := findFrontendDir()
	log.Println("Menyajikan frontend dari:", frontendDir)

	// Static files
	r.Static("/static", frontendDir)
	r.Static("/assets", filepath.Join(frontendDir, "assets"))
	r.StaticFile("/app.css", filepath.Join(frontendDir, "app.css"))
	r.StaticFile("/app.js", filepath.Join(frontendDir, "app.js"))
	r.StaticFile("/landing.css", filepath.Join(frontendDir, "landing.css"))
	r.StaticFile("/landing.js", filepath.Join(frontendDir, "landing.js"))
	r.StaticFile("/pembayaran.css", filepath.Join(frontendDir, "pembayaran.css"))
	r.StaticFile("/pembayaran.js", filepath.Join(frontendDir, "pembayaran.js"))

	// Web Pages
	r.StaticFile("/", filepath.Join(frontendDir, "index.html"))
	r.StaticFile("/login", filepath.Join(frontendDir, "login.html"))
	r.StaticFile("/pembayaran", filepath.Join(frontendDir, "pembayaran.html"))
	r.StaticFile("/dashboard", filepath.Join(frontendDir, "dashboard.html"))
	r.StaticFile("/verifikasi", filepath.Join(frontendDir, "verifikasi.html"))
	r.StaticFile("/tinjau", filepath.Join(frontendDir, "tinjau.html"))
	r.StaticFile("/pengguna", filepath.Join(frontendDir, "pengguna.html"))
	r.StaticFile("/kategori", filepath.Join(frontendDir, "kategori.html"))
	r.StaticFile("/tambah-kategori", filepath.Join(frontendDir, "tambah-kategori.html"))
	r.StaticFile("/pelaporan", filepath.Join(frontendDir, "laporan.html"))
	r.StaticFile("/laporan", filepath.Join(frontendDir, "laporan.html"))
	r.StaticFile("/detail-laporan", filepath.Join(frontendDir, "detail-laporan.html"))
	r.StaticFile("/akun", filepath.Join(frontendDir, "akun.html"))
	r.StaticFile("/pengaturan", filepath.Join(frontendDir, "coming-soon.html"))

	// Register shared API routes
	routes.RegisterAPIRoutes(r, cfg, paymentHandler)

	log.Printf("======================================================")
	log.Printf("🚀 Server ReFood aktif di %s (Port :%s)", cfg.BaseURL, cfg.Port)
	log.Printf("💳 Merchant Code: %s | Environment: %s", cfg.MerchantCode, cfg.Environment)
	log.Printf("🔗 Webhook Callback: %s", cfg.CallbackURL())
	log.Printf("🧪 Mode Simulasi Bayar: %v", cfg.EnableSimulation)
	log.Printf("======================================================")

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
