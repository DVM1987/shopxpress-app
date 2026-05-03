package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/dauvanmuoi/shopxpress-app/pkg/httpx"
	"github.com/dauvanmuoi/shopxpress-app/pkg/logger"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

var catalog = []Product{
	{ID: "p-001", Name: "Ao thun ShopXpress", Price: 199000, Stock: 42},
	{ID: "p-002", Name: "Quan jean basic", Price: 449000, Stock: 18},
	{ID: "p-003", Name: "Giay sneaker SX-Run", Price: 899000, Stock: 7},
	{ID: "p-004", Name: "Balo laptop 15in", Price: 599000, Stock: 25},
	{ID: "p-005", Name: "Mu luoi trai", Price: 149000, Stock: 60},
}

func listProducts(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("list products", "count", len(catalog), "remote", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(catalog)
	}
}

func main() {
	log := logger.New("products")
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", httpx.Healthz)
	mux.HandleFunc("/version", httpx.Version("products"))
	mux.HandleFunc("/api/products", listProducts(log))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Info("starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}
