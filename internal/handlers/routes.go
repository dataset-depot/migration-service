package handlers

import (
	"net/http"
	"github.com/dataset-depot/migration-service/internal/httpserver"
)

func (h *Handlers) Routes(maxBodyBytes int64) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/schema/version", h.schemaVersion)
	mux.HandleFunc("/migrate/upload", h.migrateUpload)
	mux.HandleFunc("/admin/instances", h.listInstances)
	mux.HandleFunc("/admin/databases", h.listDatabases)
	mux.HandleFunc("/admin/create-database", h.createDatabase)

	protected := http.NewServeMux()
	protected.Handle("/", mux)

	root := http.NewServeMux()
	root.Handle("/health", http.HandlerFunc(h.health))
	root.Handle("/", httpserver.HeaderAuth(h.adminToken, protected))
	return root
}
