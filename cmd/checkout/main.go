package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/n0en0o/shop/internal/checkout/api"
	"github.com/n0en0o/shop/internal/checkout/api/handlers"
	"github.com/n0en0o/shop/internal/checkout/applications/commands"
	"github.com/n0en0o/shop/internal/checkout/applications/queries"
	"github.com/n0en0o/shop/internal/checkout/config"
	"github.com/n0en0o/shop/internal/checkout/infrastructure/persistence"
	checkoutMsg "github.com/n0en0o/shop/internal/checkout/messaging"
	"github.com/n0en0o/shop/internal/shared"
	sharedMsg "github.com/n0en0o/shop/internal/shared/messaging"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	repo := persistence.NewOrderRepository(db)
	orderByIDHandler := queries.NewOrderByIDQueryHandler(repo)
	ordersByAccountNameHandler := queries.NewOrdersByAccountNameQueryHandler(repo)
	processOrderHandler := commands.NewProcessOrderSubmissionHandler(repo)

	orderHandler := handlers.NewOrderHandler(
		orderByIDHandler,
		ordersByAccountNameHandler,
	)

	rabbitConfig := sharedMsg.RabbitMQConfig{
		Host:     cfg.RabbitMQHost,
		Port:     cfg.RabbitMQPort,
		Username: cfg.RabbitMQUser,
		Password: cfg.RabbitMQPassword,
	}

	consumer, err := sharedMsg.NewRabbitMQConsumer(rabbitConfig)
	if err != nil {
		log.Printf("unable to connect to RabbitMQ: %v", err)
		log.Print("service will start without consumer")
	} else {
		defer consumer.Close()

		err := consumer.SetupQueue(
			sharedMsg.OrderSubmittedExchange,
			"direct",
			sharedMsg.OrderSubmittedQueue,
			sharedMsg.OrderSubmittedEventType,
		)
		if err != nil {
			log.Printf("queue configuration error: %v", err)
		} else {
			orderConsumer := checkoutMsg.NewOrderSubmittedConsumer(processOrderHandler)
			go func() {
				log.Print("launch consumer for OrderSubmittedEvent...")
				err := consumer.Consume(
					ctx,
					sharedMsg.OrderSubmittedQueue,
					orderConsumer.HandleMessage,
				)
				if err != nil && !errors.Is(err, context.Canceled) {
					log.Printf("consumer stopped: %v", err)
				}
			}()
		}
	}

	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.RegisterRoutes(r, orderHandler)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("starting checkout server on port: %s\n", cfg.AppPort)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case <-sigChan:
		log.Print("services stopped...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("checkout server shutdown error: %v", err)
		}
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		return
	}

	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
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
