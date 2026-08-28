package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/n0en0o/shop/internal/promotion/applications/commands"
	"github.com/n0en0o/shop/internal/promotion/applications/queries"
	"github.com/n0en0o/shop/internal/promotion/config"
	"github.com/n0en0o/shop/internal/promotion/grpc"
	"github.com/n0en0o/shop/internal/promotion/grpc/pb"
	"github.com/n0en0o/shop/internal/promotion/infrastructure/persistence"
	"github.com/n0en0o/shop/internal/shared"
	"github.com/n0en0o/shop/internal/shared/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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

	const serviceName = "promotion"
	metrics.AppInfo.WithLabelValues(serviceName, "1.0.0", runtime.Version()).Set(1)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		log.Printf("metrics server listening on: %s", cfg.MetricsPort)
		if err := http.ListenAndServe(":"+cfg.MetricsPort, mux); err != nil {
			log.Fatalf("metrics server error: %v", err)
		}
	}()

	if err := runGRPCServer(ctx, cfg.GRPCPort, db, serviceName); err != nil {
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

	if err := shared.WaitForDB(db, 30, 2*time.Second); err != nil {
		db.Close()
		return nil, fmt.Errorf("WaitForDB: %w", err)
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

func runGRPCServer(ctx context.Context, port string, db *sql.DB, serviceName string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.GRPCUnaryServerInterceptor(serviceName)),
		grpc.StreamInterceptor(metrics.GRPCStreamServerInterceptor(serviceName)),
	)
	greeterService := promotiongrpc.NewGreeterService()

	repo := persistence.NewPromoRepository(db)
	queryHandler := queries.NewGetByCatalogItemHandler(repo)
	commandHandler := commands.NewCreatePromoHandler(repo)
	updateHandler := commands.NewUpdatePromoHandler(repo)
	deleteHandler := commands.NewDeletePromoHandler(repo)
	promoService := promotiongrpc.NewPromotionService(
		queryHandler,
		commandHandler,
		updateHandler,
		deleteHandler,
	)

	pb.RegisterGreeterServer(grpcServer, greeterService)
	pb.RegisterPromotionServiceServer(grpcServer, promoService)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	go func() {
		<-ctx.Done()
		log.Println("shutdown signal received, stopping gRPC server...")

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
