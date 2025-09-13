package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/dataset-depot/internal/migration-service/config"
)

func MustOpenPostgres(c config.Database) *sql.DB {
	var dsn string
	if c.UseSocket {
		dsn = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable",
			c.InstanceConnectionName, c.User, c.Password, c.Name)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.Name)
	}
	
	db, err := sql.Open("pgx", dsn); if err != nil { panic(err) }

	db.SetMaxOpenConns(c.MaxOpen)
	db.SetMaxIdleConns(c.MaxIdle)
	db.SetConnMaxIdleTime(c.IdleTime)
	db.SetConnMaxLifetime(c.LifeTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		panic(fmt.Errorf("ping database: %w", err))
	}
	return db
}
