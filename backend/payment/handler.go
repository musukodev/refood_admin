package payment

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/duitku-go/common"

	"refood-backend/catalog"
	"refood-backend/config"
)

type CheckoutRequest struct {
	FoodID        string `json:"food_id"`
	Qty           int    `json:"qty"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	PaymentMethod string `json:"payment_method"`
}

type Handler struct {
	Cfg     *config.Config
	Store   *Store
	Gateway *Gateway
}

func NewHandler(cfg *config.Config, store *Store, gateway *Gateway) *Handler {
	return &Handler{Cfg: cfg, Store: store, Gateway: gateway}
}

func (h *Handler) HandleFoods(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"foods": catalog.Foods,
	})
}

func (h *Handler) HandleCheckout(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format permintaan tidak valid."})
		return
	}

	req.FoodID = strings.TrimSpace(req.FoodID)
	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.CustomerEmail = strings.TrimSpace(req.CustomerEmail)
	req.CustomerPhone = strings.TrimSpace(req.CustomerPhone)
	req.PaymentMethod = strings.ToLower(strings.TrimSpace(req.PaymentMethod))

	if err := validateCheckout(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food, ok := catalog.Find(req.FoodID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Makanan tidak ditemukan."})
		return
	}
	if req.Qty > food.Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok tersisa %d porsi.", food.Stock)})
		return
	}

	method := req.PaymentMethod
	if method == "" {
		method = "pop"
	}

	order := Order{
		ID:            newOrderID(),
		FoodID:        food.ID,
		FoodName:      food.Name,
		Store:         food.Store,
		Qty:           req.Qty,
		Amount:        food.Price * req.Qty,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		Method:        method,
		Status:        StatusPending,
		CreatedAt:     time.Now(),
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 25*time.Second)
	defer cancel()

	if method == "pop" {
		if err := h.Gateway.CreateInvoice(ctx, &order); err != nil {
			log.Printf("[CHECKOUT] [ERROR] POP order %s: %v", order.ID, err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "Gagal terhubung ke Duitku. " + err.Error()})
			return
		}
	} else {
		duitkuMethod, ok := toDuitkuMethod(method)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Metode pembayaran tidak dikenal."})
			return
		}
		if err := h.Gateway.CreateDirect(ctx, &order, duitkuMethod); err != nil {
			log.Printf("[CHECKOUT] [ERROR] direct order %s: %v", order.ID, err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "Gagal membuat transaksi Duitku. " + err.Error()})
			return
		}
	}

	h.Store.Save(order)
	log.Printf("[CHECKOUT] Order %s dibuat: %s x%d = Rp%d (%s)", order.ID, food.Name, order.Qty, order.Amount, method)

	c.JSON(http.StatusCreated, gin.H{
		"order":       order,
		"order_id":    order.ID,
		"payment_url": order.PaymentURL,
		"va_number":   order.VANumber,
		"qr_string":   order.QRString,
	})
}

func (h *Handler) HandleStatus(c *gin.Context) {
	orderID := strings.TrimSpace(c.Query("order_id"))
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id wajib diisi."})
		return
	}

	order, ok := h.Store.Get(orderID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan. Mungkin server sudah direstart."})
		return
	}

	warning := ""
	if order.Status == StatusPending && !order.IsSimulation {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		if err := h.Gateway.SyncStatus(ctx, &order); err != nil {
			warning = "Gagal sinkron status ke Duitku: " + err.Error()
		} else {
			h.Store.Save(order)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"order":              order,
		"simulation_enabled": h.Gateway.SimulationEnabled(),
		"warning":            warning,
	})
}

// HandleCallback menerima webhook dari Duitku setelah pembayaran selesai.
func (h *Handler) HandleCallback(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	if err := c.Request.ParseForm(); err != nil {
		log.Printf("[CALLBACK] Gagal parsing form: %v", err)
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	merchantCode := c.PostForm("merchantCode")
	amount := c.PostForm("amount")
	orderID := c.PostForm("merchantOrderId")
	signature := c.PostForm("signature")
	resultCode := c.PostForm("resultCode")
	reference := c.PostForm("reference")

	log.Printf("[CALLBACK] Diterima Callback dari Duitku: OrderID=%s, Amount=%s, ResultCode=%s, Ref=%s",
		orderID, amount, resultCode, reference)

	if merchantCode != h.Cfg.MerchantCode {
		log.Printf("[CALLBACK] merchantCode tidak cocok: %s", merchantCode)
		c.String(http.StatusBadRequest, "Invalid merchant")
		return
	}
	if orderID == "" {
		c.String(http.StatusBadRequest, "Invalid order")
		return
	}

	order, exists := h.Store.Get(orderID)
	if exists && amount != "" {
		paidAmount, err := strconv.Atoi(amount)
		if err != nil || paidAmount != order.Amount {
			log.Printf("[CALLBACK] nominal tidak cocok untuk %s: callback=%s order=%d", orderID, amount, order.Amount)
			c.String(http.StatusBadRequest, "Invalid amount")
			return
		}
	}

	validSignature := h.Gateway.VerifyCallbackSignature(merchantCode, amount, orderID, signature)
	if !validSignature {
		log.Printf("[CALLBACK] [WARNING] Signature tidak cocok! Received: %s", signature)
		if h.Cfg.Environment == common.ProductionEnv {
			c.String(http.StatusUnauthorized, "Invalid signature")
			return
		}
	} else {
		log.Printf("[CALLBACK] [OK] Signature tervalidasi dengan sukses!")
	}

	if !exists {
		paidAmount, _ := strconv.Atoi(amount)
		order = Order{
			ID:            orderID,
			Amount:        paidAmount,
			Method:        "DUITKU",
			Status:        StatusPending,
			Reference:     reference,
			CustomerName:  "Pelanggan",
			CustomerEmail: "-",
			CreatedAt:     time.Now(),
		}
	}

	order.Reference = reference
	order.ResultCode = resultCode
	if resultCode == "00" {
		order.Status = StatusSuccess
		order.StatusMessage = "Pembayaran berhasil"
	} else {
		order.Status = StatusFailed
		order.StatusMessage = "Pembayaran gagal atau kedaluwarsa"
	}
	h.Store.Save(order)

	// Duitku mengharapkan response HTTP 200 dengan body "OK".
	c.String(http.StatusOK, "OK")
}

func (h *Handler) HandleSimulate(c *gin.Context) {
	if !h.Gateway.SimulationEnabled() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Simulasi pembayaran tidak aktif."})
		return
	}

	var body struct {
		OrderID string `json:"order_id"`
	}
	_ = c.ShouldBindJSON(&body)

	orderID := strings.TrimSpace(body.OrderID)
	if orderID == "" {
		orderID = strings.TrimSpace(c.Query("order_id"))
	}
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id wajib diisi."})
		return
	}

	order, ok := h.Store.Update(orderID, func(o *Order) {
		o.Status = StatusSuccess
		o.ResultCode = "00"
		o.StatusMessage = "Pembayaran disimulasikan berhasil"
		o.IsSimulation = true
	})
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan."})
		return
	}

	log.Printf("[SIMULASI] Order %s berhasil diset ke SUCCESS via simulasi", orderID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pembayaran berhasil disimulasikan.",
		"order":   order,
	})
}

func (h *Handler) HandlePaymentMethods(c *gin.Context) {
	if h.Cfg.IsDummy() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Kredensial Duitku belum diisi di backend/.env."})
		return
	}

	amount := 10000
	if val, err := strconv.Atoi(c.Query("amount")); err == nil && val > 0 {
		amount = val
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	res, err := h.Gateway.PaymentMethods(ctx, amount, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"amount": amount,
		"data":   res,
	})
}

func toDuitkuMethod(method string) (common.PaymentMethod, bool) {
	switch method {
	case "qris":
		return common.SPAYQR, true
	case "va_bca":
		return common.BCAVA, true
	case "va_mandiri":
		return common.MANDIRIVA, true
	case "va_bni":
		return common.BNIVA, true
	default:
		return "", false
	}
}

func validateCheckout(req CheckoutRequest) error {
	if req.FoodID == "" {
		return fmt.Errorf("Makanan wajib dipilih.")
	}
	if req.Qty < 1 || req.Qty > 10 {
		return fmt.Errorf("Jumlah harus antara 1 sampai 10.")
	}
	if len(req.CustomerName) < 2 {
		return fmt.Errorf("Nama wajib diisi minimal 2 karakter.")
	}
	if _, err := mail.ParseAddress(req.CustomerEmail); err != nil {
		return fmt.Errorf("Format email tidak valid.")
	}
	phoneDigits := strings.NewReplacer("+", "", " ", "", "-", "", "(", "", ")", "").Replace(req.CustomerPhone)
	if len(phoneDigits) < 9 || len(phoneDigits) > 15 {
		return fmt.Errorf("Nomor WhatsApp tidak valid.")
	}
	return nil
}

func newOrderID() string {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("RF-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("RF-%s-%s", time.Now().Format("20060102-150405"), strings.ToUpper(hex.EncodeToString(buf)))
}
