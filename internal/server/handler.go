package server

import (
	"bspfp/bsrevproxy/internal/config"
	"net/http"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}

	for _, cfg := range config.Value.StaticDirs {
		if handleStatic(cfg, w, r) {
			return
		}
	}

	for _, cfg := range config.Value.ReverseProxies {
		if handleReverseProxy(cfg, w, r) {
			return
		}
	}

	for _, cfg := range config.Value.Redirects {
		if handleRedirect(cfg, w, r) {
			return
		}
	}

	http.Redirect(w, r, config.Value.DefaultRedirectUrl, http.StatusFound)
}
