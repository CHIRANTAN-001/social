package main

import (
	"github.com/CHIRANTAN-001/social/internal/db"
	"github.com/CHIRANTAN-001/social/internal/env"
	"github.com/CHIRANTAN-001/social/internal/store"
	"go.uber.org/zap"
)

var version = "0.0.1"

func main() {
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Load environment variables from a .env file if present
	// err := godotenv.Load()
	// if err != nil {
	// 	logger.Fatal("No .env file found")
	// }

	// Initialize application configuration using environment variables with fallbacks
	cfg := config{
		addr: ":" + env.GetString("ADDR", "8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Establish a new database connection using the configuration
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err) // Panic if the database connection cannot be established
	}

	defer db.Close() // Ensure the database connection is closed when main exits

	logger.Info("Database connection established")

	// Initialize the storage layer with the database connection
	store := store.NewStorage(db)

	// Create a new application instance with the configuration and storage
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	// Mount the application's routes and middleware
	mux := app.mount()

	// Start the HTTP server and log any errors that occur
	logger.Fatal(app.serve(mux))
}

