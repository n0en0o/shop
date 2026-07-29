package main

import (
	"context"
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
	"github.com/n0en0o/marketplace/internal/basket/applications/interfaces"
	"github.com/n0en0o/marketplace/internal/basket/applications/queries"
	"github.com/n0en0o/marketplace/internal/basket/config"
	"github.com/n0en0o/marketplace/internal/basket/infrastructure/cache"
	"github.com/n0en0o/marketplace/internal/basket/infrastructure/persistence"
	"github.com/n0en0o/marketplace/internal/shared"
	"github.com/redis/go-redis/v9"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("redis ping: ", err)
	}

	log.Println("redis connected")
	defer redisClient.Close()

	pgRepo := persistence.NewCartRepository(db)
	var repo interfaces.CartRepository = cache.NewRedisCartRepository(
		pgRepo,
		redisClient,
	)

	saveCartHandler := commands.NewSaveCartHandler(repo)
	getCartHandler := queries.NewGetCartHandler(repo)
	removeCartHandler := commands.NewRemoveCartHandler(repo)
	cartHandler := handlers.NewCartHandler(
		saveCartHandler,
		getCartHandler,
		removeCartHandler,
	)

	r := gin.Default()

	r.Use(shared.ErrorHandlerMiddleware())

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.RegisterRoutes(r, cartHandler)

	log.Printf("starting basket server on port: %s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
