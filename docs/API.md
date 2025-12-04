# Stumpfworks Hub API Documentation

## Base URL

Production: `https://hub.stumpfworks.de`
Development: `http://localhost:8090`

## Authentication

Currently, all endpoints are public (read-only). Authentication will be added for template/app submission in the future.

## Response Format

All responses follow this format:

```json
{
  "success": true,
  "data": { ... }
}
```

Or for errors:

```json
{
  "success": false,
  "error": "Error message"
}
```

## Endpoints

### Templates

#### List All Templates

```
GET /api/v1/templates
```

Returns a list of all available templates (metadata only).

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "jellyfin",
      "name": "Jellyfin Media Server",
      "description": "Open-source media server...",
      "icon": "🍿",
      "category": "media",
      "author": "StumpfWorks",
      "version": "1.0.0",
      "updated_at": "2025-12-04T22:00:00Z"
    }
  ]
}
```

#### Get Specific Template

```
GET /api/v1/templates/{id}
```

Returns full template details including Docker Compose file.

**Parameters:**
- `id` (path) - Template ID

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "jellyfin",
    "name": "Jellyfin Media Server",
    "description": "Open-source media server...",
    "icon": "🍿",
    "category": "media",
    "author": "StumpfWorks",
    "version": "1.0.0",
    "compose": "version: '3.8'\nservices:\n...",
    "variables": {
      "MEDIA_PATH": "/mnt/media",
      "CONFIG_PATH": "/var/lib/stumpfworks/jellyfin/config"
    },
    "requirements": {
      "min_memory_mb": 1024,
      "min_disk_gb": 5,
      "ports": [8096, 8920],
      "notes": ["Hardware transcoding requires..."]
    },
    "tags": ["media", "streaming"],
    "created_at": "2025-12-04T22:00:00Z",
    "updated_at": "2025-12-04T22:00:00Z"
  }
}
```

#### Get Template Categories

```
GET /api/v1/templates/categories
```

Returns list of all template categories.

**Response:**
```json
{
  "success": true,
  "data": ["media", "automation", "download", "monitoring"]
}
```

#### Search Templates

```
GET /api/v1/templates/search?q={query}
```

Search templates by name, description, or category.

**Parameters:**
- `q` (query) - Search query

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "jellyfin",
      "name": "Jellyfin Media Server",
      ...
    }
  ]
}
```

### Apps (Future)

#### List All Apps

```
GET /api/v1/apps
```

Returns a list of all available apps for the NAS.

#### Get Specific App

```
GET /api/v1/apps/{id}
```

Returns full app details.

#### Get App Categories

```
GET /api/v1/apps/categories
```

Returns list of all app categories.

#### Search Apps

```
GET /api/v1/apps/search?q={query}
```

Search apps by name, description, or category.

## Rate Limiting

Currently no rate limiting. Will be added in production.

## Caching

The Hub caches templates and apps in memory with a configurable TTL (default: 60 minutes).

## CORS

All origins are allowed for GET requests.

## Error Codes

- `200` - Success
- `400` - Bad Request (missing parameters)
- `404` - Not Found (template/app doesn't exist)
- `500` - Internal Server Error
