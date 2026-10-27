package payment

import (
	"context"
	"database/sql"
	"log"
	"sort"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
	StatusExpired = "EXPIRED"
)

type Order struct {
	ID            string    `json:"id"`
	FoodID        string    `json:"food_id"`
	FoodName      string    `json:"food_name"`
	Store         string    `json:"store"`
	Qty           int       `json:"qty"`
	Amount        int       `json:"amount"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	Reference     string    `json:"reference,omitempty"`
	PaymentURL    string    `json:"payment_url,omitempty"`
	VANumber      string    `json:"va_number,omitempty"`
	QRString      string    `json:"qr_string,omitempty"`
	QRImage       string    `json:"qr_image,omitempty"`
	ResultCode    string    `json:"result_code,omitempty"`
	StatusMessage string    `json:"status_message,omitempty"`
	IsSimulation  bool      `json:"is_simulation,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type OrderStore interface {
	Save(order Order)
	Get(id string) (Order, bool)
	Update(id string, fn func(*Order)) (Order, bool)
	All() []Order
}

// ==========================================
// MemoryStore (In-Memory Fallback)
// ==========================================

type MemoryStore struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{orders: make(map[string]Order)}
}

func (s *MemoryStore) Save(order Order) {
	order.UpdatedAt = time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = order
}

func (s *MemoryStore) Get(id string) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[id]
	return order, ok
}

func (s *MemoryStore) Update(id string, fn func(*Order)) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[id]
	if !ok {
		return Order{}, false
	}

	fn(&order)
	order.UpdatedAt = time.Now()
	s.orders[id] = order
	return order, true
}

func (s *MemoryStore) All() []Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.After(orders[j].CreatedAt)
	})
	return orders
}

// ==========================================
// PGStore (PostgreSQL: Supabase / Neon)
// ==========================================

type PGStore struct {
	db *sql.DB
}

func NewPGStore(db *sql.DB) *PGStore {
	return &PGStore{db: db}
}

func (s *PGStore) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(64) PRIMARY KEY,
		food_id VARCHAR(64) NOT NULL DEFAULT '',
		food_name VARCHAR(255) NOT NULL DEFAULT '',
		store_name VARCHAR(255) NOT NULL DEFAULT '',
		qty INT NOT NULL DEFAULT 1,
		amount INT NOT NULL DEFAULT 0,
		customer_name VARCHAR(255) NOT NULL DEFAULT '',
		customer_email VARCHAR(255) NOT NULL DEFAULT '',
		customer_phone VARCHAR(64) NOT NULL DEFAULT '',
		method VARCHAR(32) NOT NULL DEFAULT '',
		status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
		reference TEXT NOT NULL DEFAULT '',
		payment_url TEXT NOT NULL DEFAULT '',
		va_number VARCHAR(128) NOT NULL DEFAULT '',
		qr_string TEXT NOT NULL DEFAULT '',
		qr_image TEXT NOT NULL DEFAULT '',
		result_code VARCHAR(32) NOT NULL DEFAULT '',
		status_message TEXT NOT NULL DEFAULT '',
		is_simulation BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *PGStore) Save(order Order) {
	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	order.UpdatedAt = now

	query := `
	INSERT INTO orders (
		id, food_id, food_name, store_name, qty, amount,
		customer_name, customer_email, customer_phone, method,
		status, reference, payment_url, va_number, qr_string,
		qr_image, result_code, status_message, is_simulation,
		created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16, $17, $18, $19,
		$20, $21
	)
	ON CONFLICT (id) DO UPDATE SET
		food_id = EXCLUDED.food_id,
		food_name = EXCLUDED.food_name,
		store_name = EXCLUDED.store_name,
		qty = EXCLUDED.qty,
		amount = EXCLUDED.amount,
		customer_name = EXCLUDED.customer_name,
		customer_email = EXCLUDED.customer_email,
		customer_phone = EXCLUDED.customer_phone,
		method = EXCLUDED.method,
		status = EXCLUDED.status,
		reference = EXCLUDED.reference,
		payment_url = EXCLUDED.payment_url,
		va_number = EXCLUDED.va_number,
		qr_string = EXCLUDED.qr_string,
		qr_image = EXCLUDED.qr_image,
		result_code = EXCLUDED.result_code,
		status_message = EXCLUDED.status_message,
		is_simulation = EXCLUDED.is_simulation,
		updated_at = EXCLUDED.updated_at;
	`

	_, err := s.db.Exec(
		query,
		order.ID, order.FoodID, order.FoodName, order.Store, order.Qty, order.Amount,
		order.CustomerName, order.CustomerEmail, order.CustomerPhone, order.Method,
		order.Status, order.Reference, order.PaymentURL, order.VANumber, order.QRString,
		order.QRImage, order.ResultCode, order.StatusMessage, order.IsSimulation,
		order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		log.Printf("[DB ERROR] Gagal menyimpan order %s: %v", order.ID, err)
	}
}

func (s *PGStore) Get(id string) (Order, bool) {
	query := `
	SELECT id, food_id, food_name, store_name, qty, amount,
	       customer_name, customer_email, customer_phone, method,
	       status, reference, payment_url, va_number, qr_string,
	       qr_image, result_code, status_message, is_simulation,
	       created_at, updated_at
	FROM orders WHERE id = $1;
	`
	row := s.db.QueryRow(query, id)

	var o Order
	err := row.Scan(
		&o.ID, &o.FoodID, &o.FoodName, &o.Store, &o.Qty, &o.Amount,
		&o.CustomerName, &o.CustomerEmail, &o.CustomerPhone, &o.Method,
		&o.Status, &o.Reference, &o.PaymentURL, &o.VANumber, &o.QRString,
		&o.QRImage, &o.ResultCode, &o.StatusMessage, &o.IsSimulation,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("[DB ERROR] Gagal mengambil order %s: %v", id, err)
		}
		return Order{}, false
	}
	return o, true
}

func (s *PGStore) Update(id string, fn func(*Order)) (Order, bool) {
	order, ok := s.Get(id)
	if !ok {
		return Order{}, false
	}
	fn(&order)
	s.Save(order)
	return order, true
}

func (s *PGStore) All() []Order {
	query := `
	SELECT id, food_id, food_name, store_name, qty, amount,
	       customer_name, customer_email, customer_phone, method,
	       status, reference, payment_url, va_number, qr_string,
	       qr_image, result_code, status_message, is_simulation,
	       created_at, updated_at
	FROM orders ORDER BY created_at DESC;
	`
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("[DB ERROR] Gagal mengambil daftar orders: %v", err)
		return nil
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID, &o.FoodID, &o.FoodName, &o.Store, &o.Qty, &o.Amount,
			&o.CustomerName, &o.CustomerEmail, &o.CustomerPhone, &o.Method,
			&o.Status, &o.Reference, &o.PaymentURL, &o.VANumber, &o.QRString,
			&o.QRImage, &o.ResultCode, &o.StatusMessage, &o.IsSimulation,
			&o.CreatedAt, &o.UpdatedAt,
		); err == nil {
			orders = append(orders, o)
		}
	}
	return orders
}

// ==========================================
// Factory: NewStore
// ==========================================

func NewStore(connStr string) OrderStore {
	if connStr == "" {
		log.Println("------------------------------------------------------")
		log.Println("[DATABASE] ⚠️ STATUS: TIDAK TERHUBUNG (Kredensial DB kosong)")
		log.Println("[DATABASE] ℹ️ MODE: Menggunakan In-Memory Store (Penyimpanan Sementara di RAM)")
		log.Println("------------------------------------------------------")
		return NewMemoryStore()
	}

	log.Println("------------------------------------------------------")
	log.Println("[DATABASE] ⏳ Menghubungkan ke PostgreSQL database...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("[DATABASE] ❌ STATUS: GAGAL TERHUBUNG (Error Open: %v)", err)
		log.Println("[DATABASE] ℹ️ MODE: Fallback otomatis ke In-Memory Store")
		log.Println("------------------------------------------------------")
		return NewMemoryStore()
	}

	// Atur connection pool agar hemat connection di serverless
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("[DATABASE] ❌ STATUS: GAGAL TERHUBUNG (Ping error: %v)", err)
		log.Println("[DATABASE] ℹ️ MODE: Fallback otomatis ke In-Memory Store")
		log.Println("------------------------------------------------------")
		return NewMemoryStore()
	}

	pgStore := NewPGStore(db)
	if err := pgStore.migrate(); err != nil {
		log.Printf("[DATABASE] ❌ STATUS: TERHUBUNG TAPI MIGRASI GAGAL (%v)", err)
		log.Println("[DATABASE] ℹ️ MODE: Fallback otomatis ke In-Memory Store")
		log.Println("------------------------------------------------------")
		return NewMemoryStore()
	}

	log.Println("[DATABASE] ✅ STATUS: BERHASIL TERHUBUNG KE POSTGRESQL!")
	log.Println("[DATABASE] 📦 TABEL: 'orders' siap digunakan untuk menyimpan data")
	log.Println("------------------------------------------------------")
	return pgStore
}
