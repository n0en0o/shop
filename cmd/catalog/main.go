package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/api"
	"github.com/n0en0o/marketplace/internal/catalog/api/handlers"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
	"github.com/n0en0o/marketplace/internal/catalog/infrastructure/persistence"
)

func main() {

	const (
		appPort = "9001"
		pgHost  = "localhost"
		pgPort  = "9101"
		pgDB    = "catalog_db_dev"
		pgUser  = "postgres"
		pgPass  = "123456789"
		pgSSL   = "disable"
	)

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		pgHost, pgPort, pgDB, pgUser, pgPass, pgSSL)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	brandRepo := persistence.NewBrandRepository(db)
	listBrandsHandler := queries.NewBrandsHandler(brandRepo)
	brandsHandler := handlers.NewBrandsHandler(listBrandsHandler)

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	api.RegisterRoutes(r, brandsHandler)

	if err := r.Run(":" + appPort); err != nil {
		log.Fatal(err)
	}
}
