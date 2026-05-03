package main

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/dauvanmuoi/shopxpress-app/pkg/httpx"
	"github.com/dauvanmuoi/shopxpress-app/pkg/logger"
)

func newProxy(target string, log *slog.Logger) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		log.Error("invalid target url", "target", target, "err", err)
		os.Exit(1)
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error("upstream error", "target", target, "path", r.URL.Path, "err", err)
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}
	return rp
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	log := logger.New("gateway")

	productsURL := getEnv("PRODUCTS_URL", "http://localhost:8081")
	ordersURL := getEnv("ORDERS_URL", "http://localhost:8082")
	log.Info("upstream config", "products", productsURL, "orders", ordersURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", httpx.Healthz)
	mux.HandleFunc("/version", httpx.Version("gateway"))
	mux.Handle("/api/products", newProxy(productsURL, log))
	mux.Handle("/api/products/", newProxy(productsURL, log))
	mux.Handle("/api/orders", newProxy(ordersURL, log))
	mux.Handle("/api/orders/", newProxy(ordersURL, log))

	port := getEnv("PORT", "8080")
	log.Info("starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Error("server failed", "err", err)
		os.Exit(1)
	}
}
