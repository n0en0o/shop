package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/n0en0o/marketplace/internal/catalog/api"
	"github.com/n0en0o/marketplace/internal/catalog/api/handlers"
	"github.com/n0en0o/marketplace/internal/catalog/applications/commands"
	"github.com/n0en0o/marketplace/internal/catalog/applications/queries"
	"github.com/n0en0o/marketplace/internal/catalog/config"
	"github.com/n0en0o/marketplace/internal/catalog/infrastructure/persistence"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("invalid config: ", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("can't connect to db: ", err)
	}

	log.Println("connected to db successfully")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("postgres.WithInstance: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MigrationsPath, "postgres", driver)
	if err != nil {
		log.Fatal("migrate.NewWithDatabaseInstance: ", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("migrate.Up(): ", err)
	}

	log.Println("migrations completed successfully")

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

	log.Printf("starting server on port: %s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
