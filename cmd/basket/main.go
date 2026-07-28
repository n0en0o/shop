package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/n0en0o/marketplace/internal/basket/api"
	"github.com/n0en0o/marketplace/internal/basket/api/handlers"
	"github.com/n0en0o/marketplace/internal/basket/applications/commands"
	"github.com/n0en0o/marketplace/internal/basket/config"
	"github.com/n0en0o/marketplace/internal/basket/infrastructure/persistence"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("invalid config: ", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("can't connect to db: ", err)
	}

	log.Println("connected to basket db successfully")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("migrate driver: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsPath, "postgres", driver)
	if err != nil {
		log.Fatal("migrate init: ", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("migrate up: ", err)
	}

	log.Println("basket migrations completed successfully")

	repo := persistence.NewCartRepository(db)
	saveCartHandler := commands.NewSaveCartHandler(repo)
	cartHandler := handlers.NewCartHandler(
		saveCartHandler,
	)

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.RegisterRoutes(r, cartHandler)

	log.Printf("starting basket server on port: %s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
