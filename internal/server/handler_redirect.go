package server

import (
	"bspfp/bsrevproxy/internal/config"
	"net/http"
	"strings"
)

func handleRedirect(cfg *config.RedirectConfig, w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, cfg.RequestPathPrefix) {
		return false
	}

	targetUrl := cfg.TargetUrl

	if cfg.PassSubPath {
		subPath := strings.TrimPrefix(r.URL.Path, cfg.RequestPathPrefix)
		targetUrl = targetUrl + subPath
	}

	if cfg.PassQuery && r.URL.RawQuery != "" {
		targetUrl = targetUrl + "?" + r.URL.RawQuery
	}

	http.Redirect(w, r, targetUrl, http.StatusFound)

	return true
}
