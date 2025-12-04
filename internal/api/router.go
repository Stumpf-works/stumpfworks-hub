// Package api provides HTTP API handlers
package api

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/Stumpf-works/stumpfworks-hub/internal/registry"
	"github.com/go-chi/chi/v5"
)

//go:embed landing.html
var landingPageHTML string

// NewRouter creates the API router
func NewRouter(reg *registry.Registry) http.Handler {
	r := chi.NewRouter()

	// Template routes
	r.Get("/templates", handleListTemplates(reg))
	r.Get("/templates/categories", handleGetTemplateCategories(reg))
	r.Get("/templates/search", handleSearchTemplates(reg))
	r.Get("/templates/{id}", handleGetTemplate(reg))

	// App routes
	r.Get("/apps", handleListApps(reg))
	r.Get("/apps/categories", handleGetAppCategories(reg))
	r.Get("/apps/search", handleSearchApps(reg))
	r.Get("/apps/{id}", handleGetApp(reg))

	return r
}

// Response wrappers
type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondSuccess(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, successResponse{Success: true, Data: data})
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, errorResponse{Success: false, Error: message})
}

// Template handlers
func handleListTemplates(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates, err := reg.ListTemplates()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondSuccess(w, templates)
	}
}

func handleGetTemplateCategories(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := reg.GetTemplateCategories()
		respondSuccess(w, categories)
	}
}

func handleSearchTemplates(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			respondError(w, http.StatusBadRequest, "query parameter 'q' is required")
			return
		}

		templates, err := reg.SearchTemplates(query)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondSuccess(w, templates)
	}
}

func handleGetTemplate(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		template, err := reg.GetTemplate(id)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondSuccess(w, template)
	}
}

// App handlers
func handleListApps(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apps, err := reg.ListApps()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondSuccess(w, apps)
	}
}

func handleGetAppCategories(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := reg.GetAppCategories()
		respondSuccess(w, categories)
	}
}

func handleSearchApps(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			respondError(w, http.StatusBadRequest, "query parameter 'q' is required")
			return
		}

		apps, err := reg.SearchApps(query)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondSuccess(w, apps)
	}
}

func handleGetApp(reg *registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		app, err := reg.GetApp(id)
		if err != nil {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondSuccess(w, app)
	}
}

// ServeLandingPage returns a handler for the landing page
func ServeLandingPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(landingPageHTML))
	}
}
