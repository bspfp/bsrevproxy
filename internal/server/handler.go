package server

import (
	"bspfp/bsrevproxy/internal/config"
	"net/http"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	appCfg := config.V()

	if handleCORS(appCfg.CORS, w, r) {
		return
	}

	for _, cfg := range appCfg.StaticDirs {
		if handleStatic(cfg, w, r) {
			return
		}
	}

	for _, cfg := range appCfg.ReverseProxies {
		if handleReverseProxy(cfg, w, r) {
			return
		}
	}

	for _, cfg := range appCfg.Redirects {
		if handleRedirect(cfg, w, r) {
			return
		}
	}

	http.Redirect(w, r, appCfg.DefaultRedirectUrl, http.StatusFound)
}
