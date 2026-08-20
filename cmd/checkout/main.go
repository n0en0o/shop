package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/n0en0o/shop/internal/checkout/config"
	"github.com/n0en0o/shop/internal/shared"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("invalid config: ", err)
	}

	db, err := openDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("openDB: %v", err)
	}
	defer db.Close()
	log.Println("postgres connected")

	if err := runMigrations(db, cfg.MigrationsPath); err != nil {
		log.Fatalf("runMigrations: %v", err)
	}
	log.Println("checkout migrations completed successfully")

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("starting checkout server on port: %s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := shared.WaitForDB(db, 30, 2*time.Second); err != nil {
		db.Close()
		return nil, fmt.Errorf("WaitForDB: %w", err)
	}

	return db, nil
}

func runMigrations(migrationDB *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(migrationDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
