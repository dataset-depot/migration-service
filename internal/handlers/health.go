package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"github.com/dataset-depot/migration-service/internal/migrate"
)

type Handlers struct {
	db					*sql.DB
	migrator		migrate.Migrator
	cloudSQL		CloudSQL
	adminToken	string
}

type CloudSQL interface {
	ListInstances(ctx context.Context) (any, error)
	ListDatabases(ctx context.Context, instance, project string) (any, error)
	CreateDatabase(ctx context.Context, project, instance, name string) (any, error)
}

func New(db *sql.DB, m migrate.Migrator, token string, cloud CloudSQL) *Handlers {
	return &Handlers{db: db, migrator: m, adminToken: token, cloudSQL: cloud}
}

func (h *Handlers) health(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(); err != nil {
		http.Error(w, "database unhealthy", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}
