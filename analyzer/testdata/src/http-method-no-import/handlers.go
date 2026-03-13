package httpmethodnoimport

import "net/http"

func handleOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
