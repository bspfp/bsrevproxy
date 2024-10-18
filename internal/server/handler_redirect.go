package server

import (
	"bspfp/bsrevproxy/internal/config"
	"net/http"
	"strings"
)

func handleRedirect(cfg *config.RedirectConfig, w http.ResponseWriter, r *http.Request) bool {
	// 조건) 요청 경로가 RequestUrl로 시작하는지 확인
	if !strings.HasPrefix(r.URL.String(), cfg.RequestUrl) {
		return false
	}

	// TargetUrl 생성
	targetUrl := cfg.TargetUrl

	// PassSubPath가 true일 경우 RequestUrl을 제거한 경로를 TargetUrl에 추가
	if cfg.PassSubPath {
		subPath := strings.TrimPrefix(r.URL.Path, cfg.RequestUrl)
		targetUrl = targetUrl + subPath
	}

	// PassQuery가 true일 경우 쿼리 파라미터를 TargetUrl에 추가
	if cfg.PassQuery && r.URL.RawQuery != "" {
		targetUrl = targetUrl + "?" + r.URL.RawQuery
	}

	// 리다이렉트
	http.Redirect(w, r, targetUrl, http.StatusFound) // 302 리다이렉트

	return true
}
