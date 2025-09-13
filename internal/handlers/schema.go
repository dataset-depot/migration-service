package handlers

import (
	"net/http"
	"strconv"
)

func (h *Handlers) schemaVersion(w http.ResponseWriter, r *http.Request) {
	v, err := h.migrator.Version()
	if err != nil { http.Error(w, err.Error(), 500); return }
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.FormatInt(v, 10)))
}
