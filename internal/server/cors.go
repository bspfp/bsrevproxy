package server

import (
	"bspfp/bsrevproxy/internal/config"
	"net/http"
	"strings"
)

func handleCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", config.Value.CORS.AllowOrigin)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", config.Value.CORS.AllowHeaders)

		if !strings.Contains("GET, PUT, DELETE, OPTIONS", r.Header.Get("Access-Control-Request-Method")) {
			HttpErr(w, http.StatusMethodNotAllowed)
			return true
		}

		w.WriteHeader(http.StatusNoContent)
		return true
	}

	return false
}
