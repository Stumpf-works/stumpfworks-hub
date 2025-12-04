// Package registry provides template and app registry functionality
package registry

import "time"

// Template represents a Docker Compose template
type Template struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Icon        string                 `json:"icon"`
	Category    string                 `json:"category"`
	Author      string                 `json:"author"`
	Version     string                 `json:"version"`
	Compose     string                 `json:"compose"`
	Variables   map[string]string      `json:"variables"`
	Requirements TemplateRequirements  `json:"requirements,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Tags        []string               `json:"tags,omitempty"`
	Screenshots []string               `json:"screenshots,omitempty"`
}

// TemplateRequirements specifies system requirements
type TemplateRequirements struct {
	MinMemoryMB int      `json:"min_memory_mb,omitempty"`
	MinDiskGB   int      `json:"min_disk_gb,omitempty"`
	Ports       []int    `json:"ports,omitempty"`
	Notes       []string `json:"notes,omitempty"`
}

// App represents an addon/plugin for the NAS
type App struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Category    string            `json:"category"`
	Author      string            `json:"author"`
	Version     string            `json:"version"`
	Dependencies []string         `json:"dependencies,omitempty"`
	Packages    []string          `json:"packages,omitempty"`
	Services    []string          `json:"services,omitempty"`
	MinNASVersion string          `json:"min_nas_version,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Tags        []string          `json:"tags,omitempty"`
	Screenshots []string          `json:"screenshots,omitempty"`
	InstallScript string          `json:"install_script,omitempty"`
	UninstallScript string        `json:"uninstall_script,omitempty"`
}

// TemplateMetadata contains summary info for template listing
type TemplateMetadata struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AppMetadata contains summary info for app listing
type AppMetadata struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	UpdatedAt   time.Time `json:"updated_at"`
}
