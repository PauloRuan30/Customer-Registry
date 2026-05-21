package main

import (
	"context"
	"customer-registry-api/internal/handler"
	"customer-registry-api/internal/repository"
	"customer-registry-api/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Database Connection
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database didn't respond to ping: %v\n", err)
	}
	log.Println("Connected to PostgreSQL successfully.")

	// 3. Dependency Injection (Wiring it together)
	repo := repository.NewPostgresCustomerRepository(pool)
	svc := service.NewCustomerService(repo)
	custHandler := handler.NewCustomerHandler(svc)

	// 4. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register API Routes
	custHandler.RegisterRoutes(r)

	// 5. Start Server
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
