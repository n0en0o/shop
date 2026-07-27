package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/api"
	"github.com/n0en0o/marketplace/internal/catalog/api/handlers"
	"github.com/n0en0o/marketplace/internal/catalog/applications/commands"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
	"github.com/n0en0o/marketplace/internal/catalog/infrastructure/persistence"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	appPort := os.Getenv("CATALOG_APP_PORT")
	pgHost := os.Getenv("CATALOG_PG_HOST")
	pgPort := os.Getenv("CATALOG_PG_PORT")
	pgDB := os.Getenv("CATALOG_PG_DATABASE")
	pgUser := os.Getenv("CATALOG_PG_USER")
	pgPass := os.Getenv("CATALOG_PG_PASSWORD")
	pgSSL := os.Getenv("CATALOG_PG_SSLMODE")
	migrationsPath := os.Getenv("CATALOG_MIGRATION_PATH")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		pgHost, pgPort, pgDB, pgUser, pgPass, pgSSL)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("can't connect to db: ", err)
	}

	log.Println("connected to db successfuly")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("postgres.WithInstance: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatal("migrate.NewWithDatabaseInstance: ", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("migrate.Up(): ", err)
	}

	log.Println("miggrations completed successfuly")

	brandRepo := persistence.NewBrandRepository(db)
	categoryRepo := persistence.NewCategoryRepository(db)
	itemRepo := persistence.NewItemRepository(db)

	listBrandsHandler := queries.NewBrandsHandler(brandRepo)
	listCategoriesHandler := queries.NewCategoriesHandler(categoryRepo)
	listCatalogItemsHandler := queries.NewCatalogItemsHandler(itemRepo)
	itemByIDHandler := queries.NewCatalogItemByIDHandler(itemRepo)
	itemsByTitleHandler := queries.NewCatalogItemsByTitleHandler(itemRepo)
	createCatalogItemHandler := commands.NewCreateCatalogItemHandler(itemRepo)
	updateCatalogItemHandler := commands.NewUpdateCatalogItemHandler(itemRepo)
	deleteCatalogItemHandler := commands.NewDeleteCatalogItemHandler(itemRepo)
	itemsByBrandHandler := queries.NewCatalogItemsByBrandHandler(itemRepo)

	brandsHandler := handlers.NewBrandsHandler(listBrandsHandler)
	categoriesHandler := handlers.NewCategoriesHandler(listCategoriesHandler)
	catalogItemsHandler := handlers.NewCatalogItemsHandler(
		listCatalogItemsHandler,
		itemByIDHandler,
		itemsByTitleHandler,
		itemsByBrandHandler,
		createCatalogItemHandler,
		updateCatalogItemHandler,
		deleteCatalogItemHandler,
	)

	listItemsV2 := queries.NewCatalogItemsV2Handler(itemRepo)
	itemsHandlerV2 := handlers.NewCatalogItemsHandlerV2(listItemsV2)

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.RegisterRoutes(
		r,
		brandsHandler,
		categoriesHandler,
		catalogItemsHandler,
		itemsHandlerV2,
	)

	log.Printf("starting server on port: %s\n", appPort)
	if err := r.Run(":" + appPort); err != nil {
		log.Fatal(err)
	}
}
