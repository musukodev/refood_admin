package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/idoyudha/duitku-go/common"
	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	BaseURL          string
	MerchantCode     string
	APIKey           string
	Environment      common.Environment
	EnableSimulation bool
}

func LoadConfig() *Config {
	// Load .env if present
	loaded := false
	for _, path := range []string{".env", "backend/.env"} {
		if err := godotenv.Load(path); err == nil {
			loaded = true
			break
		}
	}
	if !loaded {
		log.Println("[INFO] File .env tidak ditemukan, menggunakan environment variable bawaan / fallback.")
	}

	port := getEnv("PORT", "8080")
	baseURL := strings.TrimRight(getEnv("BASE_URL", "http://localhost:"+port), "/")
	merchantCode := getEnv("MERCHANT_CODE", "D12345")
	apiKey := getEnv("API_KEY", "your_sandbox_api_key_here")
	envStr := strings.ToUpper(getEnv("ENVIRONMENT", "SANDBOX"))

	env := common.SandboxEnv
	if envStr == "PRODUCTION" {
		env = common.ProductionEnv
	}

	enableSimulation, err := strconv.ParseBool(getEnv("ENABLE_PAYMENT_SIMULATION", "true"))
	if err != nil {
		enableSimulation = true
	}

	cfg := &Config{
		Port:             port,
		BaseURL:          baseURL,
		MerchantCode:     merchantCode,
		APIKey:           apiKey,
		Environment:      env,
		EnableSimulation: enableSimulation,
	}

	if merchantCode == "D12345" || apiKey == "your_sandbox_api_key_here" {
		log.Println("[WARNING] Kredensial Duitku masih dummy. Isi MERCHANT_CODE dan API_KEY di backend/.env.")
	}

	return cfg
}

func (c *Config) IsDummy() bool {
	return c.MerchantCode == "D12345" || c.APIKey == "your_sandbox_api_key_here"
}

func (c *Config) CallbackURL() string {
	return c.BaseURL + "/api/callback"
}

func (c *Config) ReturnURL(orderID string) string {
	return c.BaseURL + "/pembayaran?order_id=" + orderID
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
