# Stumpfworks Hub

Central registry server for Stumpfworks NAS templates and applications.

## Overview

Stumpfworks Hub is a standalone microservice that provides:
- **Template Registry**: Docker Compose templates for one-click deployments
- **App Store**: Addon/plugin registry for NAS extensions
- **Version Management**: Template and app versioning
- **Metadata API**: Search, categories, ratings (future)

## Architecture

```
┌─────────────────┐         ┌─────────────────┐
│  Stumpfworks    │  HTTP   │  Stumpfworks    │
│      NAS        ├────────>│      Hub        │
│                 │         │                 │
└─────────────────┘         └────────┬────────┘
                                     │
                            ┌────────▼────────┐
                            │  Templates/Apps │
                            │   (JSON Files)  │
                            └─────────────────┘
```

## Features

### Template Registry
- JSON-based template storage
- Categories: media, automation, download, monitoring
- Variable substitution
- Requirement specifications (RAM, disk, ports)

### App Store (Future)
- Addon registry
- Dependency management
- Version updates
- Community submissions

## API Endpoints

### Templates
```
GET  /api/v1/templates              - List all templates
GET  /api/v1/templates/categories   - Get categories
GET  /api/v1/templates/{id}         - Get specific template
GET  /api/v1/templates/search?q=... - Search templates
```

### Apps (Future)
```
GET  /api/v1/apps                   - List all apps
GET  /api/v1/apps/{id}              - Get specific app
```

## Running

```bash
# Development
go run cmd/hub/main.go

# Production
go build -o hub cmd/hub/main.go
./hub --port 8090
```

## Configuration

Environment variables:
- `HUB_PORT`: Server port (default: 8090)
- `HUB_TEMPLATES_DIR`: Templates directory (default: ./templates)
- `HUB_APPS_DIR`: Apps directory (default: ./apps)
- `HUB_CACHE_TTL`: Cache TTL in minutes (default: 60)

## Development

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Format code
go fmt ./...

# Build
go build -o hub cmd/hub/main.go
```

## Template Format

Templates are stored as JSON files in `templates/{category}/{name}.json`:

```json
{
  "id": "plex",
  "name": "Plex Media Server",
  "description": "Stream your media collection",
  "icon": "🎬",
  "category": "media",
  "author": "StumpfWorks",
  "version": "1.0.0",
  "variables": {
    "MEDIA_PATH": "/mnt/media",
    "PUID": "1000"
  },
  "compose": "version: '3.8'\nservices:\n  plex:\n    ...",
  "requirements": {
    "min_memory_mb": 2048,
    "min_disk_gb": 10,
    "ports": [32400]
  }
}
```

## License

Proprietary - Stumpf.Works
