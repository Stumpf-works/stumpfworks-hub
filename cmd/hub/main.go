// Stumpfworks Hub - Central registry for templates and apps
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stumpf-works/stumpfworks-hub/internal/api"
	"github.com/Stumpf-works/stumpfworks-hub/internal/registry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	defaultPort         = 8090
	defaultTemplatesDir = "./templates"
	defaultAppsDir      = "./apps"
	defaultCacheTTL     = 60 // minutes
)

func main() {
	// Parse flags
	port := flag.Int("port", defaultPort, "Server port")
	templatesDir := flag.String("templates", defaultTemplatesDir, "Templates directory")
	appsDir := flag.String("apps", defaultAppsDir, "Apps directory")
	cacheTTL := flag.Int("cache-ttl", defaultCacheTTL, "Cache TTL in minutes")
	flag.Parse()

	// Override with environment variables if set
	if envPort := os.Getenv("HUB_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envTemplates := os.Getenv("HUB_TEMPLATES_DIR"); envTemplates != "" {
		*templatesDir = envTemplates
	}
	if envApps := os.Getenv("HUB_APPS_DIR"); envApps != "" {
		*appsDir = envApps
	}
	if envCacheTTL := os.Getenv("HUB_CACHE_TTL"); envCacheTTL != "" {
		fmt.Sscanf(envCacheTTL, "%d", cacheTTL)
	}

	fmt.Printf("🚀 Starting Stumpfworks Hub\n")
	fmt.Printf("   Port:          %d\n", *port)
	fmt.Printf("   Templates:     %s\n", *templatesDir)
	fmt.Printf("   Apps:          %s\n", *appsDir)
	fmt.Printf("   Cache TTL:     %d minutes\n", *cacheTTL)

	// Initialize registry
	reg := registry.NewRegistry(*templatesDir, *appsDir, time.Duration(*cacheTTL)*time.Minute)
	if err := reg.Initialize(); err != nil {
		fmt.Printf("❌ Failed to initialize registry: %v\n", err)
		os.Exit(1)
	}

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// CORS - allow all origins for now (hub is public)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Landing page
	r.Get("/", api.ServeLandingPage())

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes
	r.Mount("/api/v1", api.NewRouter(reg))

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		fmt.Printf("✅ Hub listening on http://localhost:%d\n", *port)
		fmt.Printf("📡 API: http://localhost:%d/api/v1\n", *port)
		fmt.Printf("💚 Health: http://localhost:%d/health\n\n", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down Hub...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("❌ Server forced to shutdown: %v\n", err)
	}

	fmt.Println("✅ Hub stopped gracefully")
}
