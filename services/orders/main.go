package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dauvanmuoi/shopxpress-app/pkg/httpx"
	"github.com/dauvanmuoi/shopxpress-app/pkg/logger"
)

type Order struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

var orders sync.Map

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "ord-" + hex.EncodeToString(b)
}

func handleOrders(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var o Order
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		o.ID = newID()
		o.CreatedAt = time.Now().UTC()
		orders.Store(o.ID, o)
		log.Info("create order", "id", o.ID, "product_id", o.ProductID, "qty", o.Quantity)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(o)
	}
}

func handleOrderByID(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/orders/")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		v, ok := orders.Load(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		log.Info("get order", "id", id)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
}

func main() {
	log := logger.New("orders")
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", httpx.Healthz)
	mux.HandleFunc("/version", httpx.Version("orders"))
	mux.HandleFunc("/api/orders", handleOrders(log))
	mux.HandleFunc("/api/orders/", handleOrderByID(log))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Info("starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}
