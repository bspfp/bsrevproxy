package server

import "net/http"

func HttpErr(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
