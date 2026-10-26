package payment

import (
	"sort"
	"sync"
	"time"
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

type Store struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewStore() *Store {
	return &Store{orders: make(map[string]Order)}
}

func (s *Store) Save(order Order) {
	order.UpdatedAt = time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = order
}

func (s *Store) Get(id string) (Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[id]
	return order, ok
}

func (s *Store) Update(id string, fn func(*Order)) (Order, bool) {
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

func (s *Store) All() []Order {
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
