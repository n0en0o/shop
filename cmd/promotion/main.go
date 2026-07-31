package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/n0en0o/marketplace/internal/promotion/config"
	promotiongrpc "github.com/n0en0o/marketplace/internal/promotion/grpc"
	"github.com/n0en0o/marketplace/internal/promotion/grpc/pb"
	"google.golang.org/grpc"
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
	log.Println("mysql connected")

	if err := runMigrations(db, cfg.MigrationsPath); err != nil {
		log.Fatalf("runMigrations: %v", err)
	}
	log.Println("promotion migrations completed successfully")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := runGRPCServer(ctx, cfg.GRPCPort); err != nil {
		log.Fatalf("runGRPCServer: %v", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sql.Ping: %w", err)
	}

	return db, nil
}

func runMigrations(migrationDB *sql.DB, migrationsPath string) error {
	driver, err := mysql.WithInstance(migrationDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func runGRPCServer(ctx context.Context, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	greeterService := promotiongrpc.NewGreeterService()

	pb.RegisterGreeterServer(grpcServer, greeterService)

	go func() {
		<-ctx.Done()
		log.Println("shutdown signal received, stopped gRPC server...")

		timer := time.AfterFunc(10*time.Second, func() {
			log.Println("timeout exceeded, server stop")
			grpcServer.Stop()
		})

		defer timer.Stop()
		grpcServer.GracefulStop()
	}()

	log.Printf("gRPC server listening on: %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("grpcServer.Serve: %w", err)
	}

	return nil
}
