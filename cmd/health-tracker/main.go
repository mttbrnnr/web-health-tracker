package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"web-health-tracker/internal/db"
	"web-health-tracker/internal/handler"
	"web-health-tracker/web"
)

func main() {
	// Initialize database
	database, err := db.Open()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if err := db.InitSchema(database); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	if err := db.SeedFoods(database); err != nil {
		log.Fatalf("Failed to seed foods: %v", err)
	}

	// Initialize templates
	if err := handler.InitTemplates(web.Templates); err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	// Set database for handlers
	handler.SetDB(database)

	// Create router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files (embedded)
	staticFS, _ := fs.Sub(web.Static, "static")
	fileServer := http.FileServer(http.FS(staticFS))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Routes
	r.Get("/", handler.Home)
	r.Get("/day/{date}", handler.Day)

	// Day endpoints
	r.Post("/day/{date}/weight", handler.UpdateWeight)
	r.Post("/day/{date}/type", handler.UpdateDayType)

	// Meal endpoints
	r.Post("/meals/{meal}/save", handler.SaveMeal)
	r.Post("/meals/{meal}/yesterday", handler.SameAsYesterday)

	// Food endpoints
	r.Post("/foods/custom", handler.AddCustomFood)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
