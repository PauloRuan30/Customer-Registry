package main

import (
	"context"
	"customer-registry-api/internal/handler"
	"customer-registry-api/internal/repository"
	"customer-registry-api/internal/service"
	"log/slog"
	"net/http"
	"os"

	_ "customer-registry-api/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	// Import of Datadog Tracer and Chi middleware
	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// @title Customer Registry API
// @version 1.0
// @description Microservice for managing and validating customer records.
// @host localhost:8080
// @BasePath /
func main() {
	// 3. Configure slog to output JSON (Datadog parses this automatically)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 4. Start the Datadog Tracer
	tracer.Start(
		tracer.WithService("customer-registry-api"),
		tracer.WithEnv("dev"),
		tracer.WithServiceVersion("1.0.0"),
	)
	defer tracer.Stop() // Ensure traces flush when the app exits

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL environment variable is required")
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	slog.Info("Connected to PostgreSQL successfully")

	repo := repository.NewPostgresCustomerRepository(pool)
	svc := service.NewCustomerService(repo)
	custHandler := handler.NewCustomerHandler(svc)

	r := chi.NewRouter()

	r.Use(chitrace.Middleware(chitrace.WithServiceName("customer-registry-api")))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	})

	custHandler.RegisterRoutes(r)

	slog.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
