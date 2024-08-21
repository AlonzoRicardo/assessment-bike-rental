package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/bike-rental/src/controllers"
	"github.com/yourusername/bike-rental/src/cronjobs"
	"github.com/yourusername/bike-rental/src/database"
	"github.com/yourusername/bike-rental/src/logger"
)

func main() {
	logger.Init()

	log.Info().Msg("Initializing db connection...")

	// Load the configuration
	config, err := database.LoadConfig("config.toml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize the database connection
	db, err := database.InitDB(&config.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	// Migrate the schema
	database.Migrate(&config.Database)
	// Clean the database
	// database.CleanDatabase(db)
	// Seed the database with fixtures
	database.SeedDatabase(db)

	// Initialize the HTTP server and routes...
	r := chi.NewRouter()
	// Installing logger middleware for debugging...
	r.Use(logger.LoggerMiddleware)

	r.Get("/assignments", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllAssignments(w, r, db)
	})
	r.Post("/bikes/assign", func(w http.ResponseWriter, r *http.Request) {
		controllers.AssignBike(w, r, db)
	})
	r.Post("/bikes/unassign", func(w http.ResponseWriter, r *http.Request) {
		controllers.UnassignBike(w, r, db)
	})
	r.Get("/bikes/available", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAvailableBikes(w, r, db)
	})
	r.Get("/bikes", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllBikes(w, r, db)
	})

	// Set up the cron job to run the function every hour
	log.Info().Msg("Setting up cronjobs...")
	c := cron.New()
	c.AddFunc("@hourly", func() { cronjobs.AutoUnassignOverdueBikes(db) })
	c.Start()

	log.Info().Msg("Starting server...")
	http.ListenAndServe(":8080", r)
}
