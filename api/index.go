package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"refood-backend/backend/config"
	"refood-backend/backend/payment"
	"refood-backend/backend/routes"
)

var (
	app  *gin.Engine
	once sync.Once
)

func initApp() *gin.Engine {
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(gin.Recovery())

		cfg := config.LoadConfig()
		paymentStore := payment.NewStore()
		paymentGateway := payment.NewGateway(cfg)
		paymentHandler := payment.NewHandler(cfg, paymentStore, paymentGateway)

		routes.RegisterAPIRoutes(r, cfg, paymentHandler)

		app = r
	})
	return app
}

// Handler adalah entrypoint standar Serverless Function untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	initApp().ServeHTTP(w, r)
}
