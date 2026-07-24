package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/api"
	"github.com/n0en0o/marketplace/internal/catalog/api/handlers"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
	"github.com/n0en0o/marketplace/internal/catalog/infrastructure/persistence"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := os.Getenv("CATALOG_APP_PORT")
	pgHost := os.Getenv("CATALOG_PG_HOST")
	pgPort := os.Getenv("CATALOG_PG_PORT")
	pgDB := os.Getenv("CATALOG_PG_DATABASE")
	pgUser := os.Getenv("CATALOG_PG_USER")
	pgPass := os.Getenv("CATALOG_PG_PASSWORD")
	pgSSL := os.Getenv("CATALOG_PG_SSLMODE")

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
	categoryRepo := persistence.NewCategoryRepository(db)
	itemRepo := persistence.NewItemRepository(db)

	listBrandsHandler := queries.NewBrandsHandler(brandRepo)
	listCategoriesHandler := queries.NewCategoriesHandler(categoryRepo)
	listCatalogItemsHandler := queries.NewCatalogItemsHandler(itemRepo)

	brandsHandler := handlers.NewBrandsHandler(listBrandsHandler)
	categoriesHandler := handlers.NewCategoriesHandler(listCategoriesHandler)
	catalogItemsHandler := handlers.NewCatalogItemsHandler(listCatalogItemsHandler)

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	api.RegisterRoutes(
		r,
		brandsHandler,
		categoriesHandler,
		catalogItemsHandler,
	)

	if err := r.Run(":" + appPort); err != nil {
		log.Fatal(err)
	}
}
