package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/dataset-depot/migration-service/internal/config"
	"github.com/dataset-depot/migration-service/internal/httpserver"
	"github.com/dataset-depot/migration-service/internal/storage"
	"github.com/dataset-depot/migration-service/internal/migrate"
	"github.com/dataset-depot/migration-service/internal/handlers"
)

func main() {
	cfg := config.Load()

	db := storage.MustOpenPostgres(cfg.Database)
	defer db.Close()

	migrator := migrate.NewGooseMigrator(db)

	h := handlers.New(db, migrator, cfg.Security.AdminToken, handlers.NewCloudSQL(cfg.CloudSQL))

	srv := httpserver.New(httpserver.Opts{
		Addr:								cfg.HTTP.Addr,
		ReadTimeout:				cfg.HTTP.ReadTimeout,
		WriteTimeout:				cfg.HTTP.WriteTimeout,
	}, h.Routes(cfg.HTTP.MaxBodyBytes))

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("http server start failed: %v", err)
		}
	}()
	log.Printf("server listening on %s", cfg.Http.Addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<- stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	} else {
		log.Println("server stopped cleanly")
	}
}

