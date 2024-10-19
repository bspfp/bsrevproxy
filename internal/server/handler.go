package server

import (
	"bspfp/bsrevproxy/internal/config"
	"fmt"
	"log"
	"net/http"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	appCfg := config.V()

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

	log.Printf("No matching route found: %q\n", []string{
		"method", r.Method,
		"url", r.URL.String(),
		"remote", r.RemoteAddr,
		"user-agent", r.UserAgent(),
		"host", r.Host,
		"headers", fmt.Sprintf("%q", r.Header),
	})

	http.Redirect(w, r, appCfg.DefaultRedirectUrl, http.StatusFound)
}
