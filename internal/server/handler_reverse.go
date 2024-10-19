package server

import (
	"bspfp/bsrevproxy/internal/config"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func handleReverseProxy(cfg *config.ReverseProxyConfig, w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, cfg.RequestPathPrefix) {
		return false
	}

	trimmedPath := strings.TrimPrefix(r.URL.Path, cfg.RequestPathPrefix)
	targetUrl, err := url.Parse(cfg.TargetUrl + trimmedPath)
	if err != nil {
		HttpErr(w, http.StatusBadRequest)
		return true
	}

	proxyReq, err := http.NewRequest(r.Method, targetUrl.String(), r.Body)
	if err != nil {
		HttpErr(w, http.StatusBadRequest)
		return true
	}

	proxyReq.Header.Add("X-Forwarded-For", r.RemoteAddr)
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	client := &http.Client{Timeout: time.Duration(cfg.TimeoutSec) * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		HttpErr(w, http.StatusBadGateway)
		return true
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		HttpErr(w, http.StatusInternalServerError)
	}

	return true
}
