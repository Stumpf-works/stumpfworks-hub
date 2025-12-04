// Package registry manages templates and apps
package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Registry manages templates and apps with caching
type Registry struct {
	templatesDir string
	appsDir      string
	cacheTTL     time.Duration

	// In-memory cache
	templates     map[string]*Template
	apps          map[string]*App
	categories    []string
	appCategories []string
	lastLoad      time.Time
	mu            sync.RWMutex
}

// NewRegistry creates a new registry instance
func NewRegistry(templatesDir, appsDir string, cacheTTL time.Duration) *Registry {
	return &Registry{
		templatesDir: templatesDir,
		appsDir:      appsDir,
		cacheTTL:     cacheTTL,
		templates:    make(map[string]*Template),
		apps:         make(map[string]*App),
	}
}

// Initialize loads all templates and apps into memory
func (r *Registry) Initialize() error {
	fmt.Println("📚 Loading templates and apps...")

	if err := r.loadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	if err := r.loadApps(); err != nil {
		// Don't fail if apps dir doesn't exist yet
		fmt.Printf("⚠️  Apps directory not found or empty: %v\n", err)
	}

	fmt.Printf("✅ Loaded %d templates and %d apps\n", len(r.templates), len(r.apps))
	return nil
}

// loadTemplates loads all templates from disk
func (r *Registry) loadTemplates() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.templates = make(map[string]*Template)
	categoryMap := make(map[string]bool)

	err := filepath.Walk(r.templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		// Read template file
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("⚠️  Failed to read %s: %v\n", path, err)
			return nil
		}

		// Parse template
		var tpl Template
		if err := json.Unmarshal(data, &tpl); err != nil {
			fmt.Printf("⚠️  Failed to parse %s: %v\n", path, err)
			return nil
		}

		// Validate required fields
		if tpl.ID == "" || tpl.Name == "" {
			fmt.Printf("⚠️  Invalid template in %s: missing ID or Name\n", path)
			return nil
		}

		// Store template
		r.templates[tpl.ID] = &tpl
		categoryMap[tpl.Category] = true

		return nil
	})

	if err != nil {
		return err
	}

	// Extract categories
	r.categories = make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		r.categories = append(r.categories, cat)
	}

	r.lastLoad = time.Now()
	return nil
}

// loadApps loads all apps from disk
func (r *Registry) loadApps() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if apps directory exists
	if _, err := os.Stat(r.appsDir); os.IsNotExist(err) {
		return nil // Apps directory doesn't exist yet, that's OK
	}

	r.apps = make(map[string]*App)
	categoryMap := make(map[string]bool)

	err := filepath.Walk(r.appsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-JSON files
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		// Read app file
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("⚠️  Failed to read %s: %v\n", path, err)
			return nil
		}

		// Parse app
		var app App
		if err := json.Unmarshal(data, &app); err != nil {
			fmt.Printf("⚠️  Failed to parse %s: %v\n", path, err)
			return nil
		}

		// Validate required fields
		if app.ID == "" || app.Name == "" {
			fmt.Printf("⚠️  Invalid app in %s: missing ID or Name\n", path)
			return nil
		}

		// Store app
		r.apps[app.ID] = &app
		categoryMap[app.Category] = true

		return nil
	})

	if err != nil {
		return err
	}

	// Extract categories
	r.appCategories = make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		r.appCategories = append(r.appCategories, cat)
	}

	return nil
}

// Reload reloads templates and apps if cache expired
func (r *Registry) Reload() error {
	r.mu.RLock()
	needsReload := time.Since(r.lastLoad) > r.cacheTTL
	r.mu.RUnlock()

	if needsReload {
		return r.Initialize()
	}
	return nil
}

// GetTemplate returns a template by ID
func (r *Registry) GetTemplate(id string) (*Template, error) {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()

	tpl, ok := r.templates[id]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", id)
	}
	return tpl, nil
}

// ListTemplates returns all templates
func (r *Registry) ListTemplates() ([]*TemplateMetadata, error) {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*TemplateMetadata, 0, len(r.templates))
	for _, tpl := range r.templates {
		result = append(result, &TemplateMetadata{
			ID:          tpl.ID,
			Name:        tpl.Name,
			Description: tpl.Description,
			Icon:        tpl.Icon,
			Category:    tpl.Category,
			Author:      tpl.Author,
			Version:     tpl.Version,
			UpdatedAt:   tpl.UpdatedAt,
		})
	}
	return result, nil
}

// GetTemplateCategories returns all template categories
func (r *Registry) GetTemplateCategories() []string {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.categories
}

// SearchTemplates searches templates by query
func (r *Registry) SearchTemplates(query string) ([]*TemplateMetadata, error) {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(query)
	result := make([]*TemplateMetadata, 0)

	for _, tpl := range r.templates {
		if strings.Contains(strings.ToLower(tpl.Name), query) ||
			strings.Contains(strings.ToLower(tpl.Description), query) ||
			strings.Contains(strings.ToLower(tpl.Category), query) {
			result = append(result, &TemplateMetadata{
				ID:          tpl.ID,
				Name:        tpl.Name,
				Description: tpl.Description,
				Icon:        tpl.Icon,
				Category:    tpl.Category,
				Author:      tpl.Author,
				Version:     tpl.Version,
				UpdatedAt:   tpl.UpdatedAt,
			})
		}
	}
	return result, nil
}

// GetApp returns an app by ID
func (r *Registry) GetApp(id string) (*App, error) {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()

	app, ok := r.apps[id]
	if !ok {
		return nil, fmt.Errorf("app not found: %s", id)
	}
	return app, nil
}

// ListApps returns all apps
func (r *Registry) ListApps() ([]*AppMetadata, error) {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*AppMetadata, 0, len(r.apps))
	for _, app := range r.apps {
		result = append(result, &AppMetadata{
			ID:          app.ID,
			Name:        app.Name,
			Description: app.Description,
			Icon:        app.Icon,
			Category:    app.Category,
			Author:      app.Author,
			Version:     app.Version,
			UpdatedAt:   app.UpdatedAt,
		})
	}
	return result, nil
}

// GetAppCategories returns all app categories
func (r *Registry) GetAppCategories() []string {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.appCategories
}

// SearchApps searches apps by query
func (r *Registry) SearchApps(query string) ([]*AppMetadata, error) {
	r.Reload()
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(query)
	result := make([]*AppMetadata, 0)

	for _, app := range r.apps {
		if strings.Contains(strings.ToLower(app.Name), query) ||
			strings.Contains(strings.ToLower(app.Description), query) ||
			strings.Contains(strings.ToLower(app.Category), query) {
			result = append(result, &AppMetadata{
				ID:          app.ID,
				Name:        app.Name,
				Description: app.Description,
				Icon:        app.Icon,
				Category:    app.Category,
				Author:      app.Author,
				Version:     app.Version,
				UpdatedAt:   app.UpdatedAt,
			})
		}
	}
	return result, nil
}
