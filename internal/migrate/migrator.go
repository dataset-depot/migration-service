package migrate

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"
	"github.com/pressly/goose/v3"
)

type Runner interface {
	Version() (int64, error)
	UpFromDir(ctx context.Context, dir string) error
}

type gooseRunner struct { db * sql.DB }

func NewGooseMigrator(db *sql.DB) Runner {
	goose.SetDialect("postgres")
	return &gooseRunner{db: db}
}

func (g *gooseRunner) Version() (int64, error) {
	return goose.GetDBVersion(g.db)
}

func (g *gooseRunner) UpFromDir(ctx context.Context, dir string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute); defer cancel()
	return goose.UpContext(ctx, g.db, filepath.Clean(dir))
}
