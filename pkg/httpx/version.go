package httpx

import (
	"encoding/json"
	"net/http"
	"os"
)

type VersionInfo struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

func Version(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v := VersionInfo{
			Service: service,
			Version: getEnv("IMAGE_TAG", "dev"),
			Commit:  getEnv("GIT_COMMIT", "unknown"),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
