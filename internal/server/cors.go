package server

import (
	"bspfp/bsrevproxy/internal/config"
	"net/http"
	"strings"
)

const allowedMethods = "GET, POST, PUT, DELETE, OPTIONS"

func handleCORS(cfg *config.CORSConfig, w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", cfg.AllowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == http.MethodOptions {
		if !strings.Contains(allowedMethods, r.Header.Get("Access-Control-Request-Method")) {
			HttpErr(w, http.StatusMethodNotAllowed)
			return true
		}

		w.WriteHeader(http.StatusNoContent)
		return true
	}

	return false
}
