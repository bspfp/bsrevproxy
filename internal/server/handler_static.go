package server

import (
	"bspfp/bsrevproxy/internal/config"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func handleStatic(cfg *config.StaticDirConfig, w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, cfg.RequestPathPrefix) {
		return false
	}

	localFilePath := filepath.Join(cfg.LocalPath, strings.TrimPrefix(r.URL.Path, cfg.RequestPathPrefix))

	file, err := os.Open(localFilePath)
	if err != nil {
		HttpErr(w, http.StatusNotFound)
		return true
	}
	defer file.Close()

	if cfg.MimeType != "" {
		w.Header().Set("Content-Type", cfg.MimeType)
	} else {
		ext := filepath.Ext(localFilePath)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", mimeType)
	}

	if _, err = io.Copy(w, file); err != nil {
		HttpErr(w, http.StatusNotFound)
	}

	return true
}
