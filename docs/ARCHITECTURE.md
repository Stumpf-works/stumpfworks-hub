# Stumpfworks Hub Architecture

## Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Stumpfworks Ecosystem                     │
└─────────────────────────────────────────────────────────────┘

┌──────────────────┐         ┌──────────────────┐
│  Stumpfworks NAS │         │  Stumpfworks Hub │
│  (192.168.x.x)   │◄────────┤ (hub.stumpfworks │
│                  │  HTTPS  │      .de:443)    │
└──────────────────┘         └──────────────────┘
         │                            │
         │                            │
         ▼                            ▼
  ┌─────────────┐           ┌──────────────────┐
  │   Docker    │           │   Templates/Apps │
  │  Containers │           │   (JSON Files)   │
  └─────────────┘           └──────────────────┘
```

## Components

### Stumpfworks Hub

**Location:** `apt.stumpfworks.de:8090` → `hub.stumpfworks.de` (via nginx)

**Purpose:**
- Central registry for Docker Compose templates
- App Store for NAS addons/plugins
- Version management
- Metadata API (search, categories, tags)

**Technology:**
- Go 1.24+
- Chi router
- JSON file storage
- In-memory caching (60 min TTL)

**Endpoints:**
- `/api/v1/templates` - Template registry
- `/api/v1/apps` - App store (future)
- `/health` - Health check

### Stumpfworks NAS

**Purpose:**
- Consumer of Hub API
- Docker stack deployment
- Template rendering
- User management

**Integration:**
- Environment variable: `STUMPFWORKS_HUB_URL=https://hub.stumpfworks.de`
- HTTP client with fallback to builtin templates
- Automatic retry logic

## Data Flow

### Template Deployment Flow

```
1. User selects template in NAS UI
   └─> Frontend: GET /api/v1/docker/templates

2. NAS Backend fetches from Hub
   └─> Hub Client: GET https://hub.stumpfworks.de/api/v1/templates/{id}

3. Hub returns template JSON
   └─> Includes: compose file, variables, requirements

4. NAS renders template with user variables
   └─> Variable substitution: {{MEDIA_PATH}} → /mnt/media

5. NAS deploys to Docker
   └─> docker compose up -d

6. Stack created and running
   └─> Containers, networks, volumes
```

### Cache Strategy

**Hub Side:**
- In-memory cache with 60-minute TTL
- Reload on cache miss or TTL expiry
- Full scan of templates directory on reload

**NAS Side:**
- No caching (always fetches latest)
- Fallback to builtin templates on network error
- 30-second HTTP timeout

## File Structure

### Hub

```
stumpfworks-hub/
├── cmd/
│   └── hub/
│       └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   └── router.go            # HTTP handlers
│   ├── registry/
│   │   ├── types.go             # Data models
│   │   └── registry.go          # Business logic
│   └── storage/                 # Future: DB support
├── templates/
│   ├── media/
│   │   ├── plex.json
│   │   ├── jellyfin.json
│   │   └── emby.json
│   ├── automation/
│   │   ├── sonarr.json
│   │   ├── radarr.json
│   │   └── prowlarr.json
│   ├── download/
│   │   ├── transmission.json
│   │   └── qbittorrent.json
│   └── monitoring/
│       └── ...
├── apps/                        # Future: App Store
├── deployment/
│   ├── stumpfworks-hub.service
│   └── deploy.sh
└── docs/
    ├── API.md
    ├── DEPLOYMENT.md
    └── ARCHITECTURE.md
```

### NAS

```
stumpfworks-nas/
└── backend/
    └── internal/
        └── docker/
            ├── hub_client.go    # Hub API client
            ├── templates.go     # Template functions (Hub-aware)
            └── compose.go       # Docker Compose operations
```

## Security

### Current (v1.0)
- Public read-only API
- No authentication required
- CORS: Allow all origins (read-only endpoints)
- HTTPS via Let's Encrypt (nginx reverse proxy)

### Future (v2.0)
- API keys for template submission
- User accounts and ratings
- Verified authors/publishers
- Template validation and scanning
- Rate limiting

## Scalability

### Current Load
- Expected: 10-50 NAS instances
- Templates: ~20-50 total
- Request rate: ~10 req/min
- Cache hit rate: >95%

### Future Scaling
- Add Redis for distributed caching
- Add PostgreSQL for metadata
- Add CDN for template files
- Horizontal scaling with load balancer

## Deployment

### Production Servers

**hub.stumpfworks.de:**
- Server: apt.stumpfworks.de (shared with APT repo)
- User: www-data
- Service: systemd
- Logs: journalctl -u stumpfworks-hub
- Data: /var/lib/stumpfworks-hub/

**Backup:**
- Daily backup of /var/lib/stumpfworks-hub/templates/
- Git repository as source of truth

### Development

**Local Hub:**
```bash
cd stumpfworks-hub
go run cmd/hub/main.go --port 8090
```

**Local NAS (pointing to local Hub):**
```bash
export STUMPFWORKS_HUB_URL=http://localhost:8090
cd stumpfworks-nas/backend
go run cmd/nas/main.go
```

## Monitoring

### Health Checks

```bash
# Hub health
curl https://hub.stumpfworks.de/health

# Template count
curl https://hub.stumpfworks.de/api/v1/templates | jq '.data | length'

# Categories
curl https://hub.stumpfworks.de/api/v1/templates/categories
```

### Logs

```bash
# Real-time Hub logs
ssh apt.stumpfworks.de
sudo journalctl -u stumpfworks-hub -f

# NAS template fetch logs
sudo journalctl -u stumpfworks-nas | grep -i "hub\|template"
```

## Future Enhancements

### Phase 2: App Store
- Addon/plugin registry
- Dependency management
- Installation scripts
- Version updates

### Phase 3: Community
- User submissions
- Template ratings and reviews
- Screenshots and demos
- Usage statistics

### Phase 4: Advanced
- Template versioning (v1, v2, etc.)
- Breaking change notifications
- Automatic updates
- Template marketplace
