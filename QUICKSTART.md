# Stumpfworks Hub - Quick Start Guide

## ✅ Hub ist Ready!

### Was ist komplett:
- ✅ 10 Templates (Plex, Jellyfin, Sonarr, Radarr, Prowlarr, Transmission, qBittorrent, Portainer, Uptime Kuma, Complete Media Stack)
- ✅ 4 Kategorien (media, automation, download, monitoring)
- ✅ REST API komplett funktional
- ✅ Systemd Service
- ✅ Nginx Config
- ✅ Deployment Script
- ✅ NAS Backend Integration

### Template Übersicht:

**Media (3):**
- Plex Media Server
- Jellyfin Media Server
- Complete Media Stack (all-in-one)

**Automation (3):**
- Sonarr (TV)
- Radarr (Movies)
- Prowlarr (Indexer)

**Download (2):**
- Transmission
- qBittorrent

**Monitoring (2):**
- Portainer (Docker UI)
- Uptime Kuma (Uptime monitoring)

## Deployment auf 46.4.25.15

### 1. Hub Deployen

```bash
cd /Users/sebastianstumpf/Documents/GitHub/stumpfworks-hub
./deployment/deploy.sh 46.4.25.15
```

Das Script macht automatisch:
- Baut Hub Binary für Linux
- Uploaded alles zum Server
- Installiert Systemd Service
- Kopiert alle 10 Templates
- Startet den Hub

### 2. Nginx Konfigurieren

```bash
ssh root@46.4.25.15

# Nginx Config kopieren
cp /opt/stumpfworks-hub/deployment/nginx-hub.conf /etc/nginx/sites-available/hub.stumpfworks.de

# Symlink erstellen
ln -s /etc/nginx/sites-available/hub.stumpfworks.de /etc/nginx/sites-enabled/

# Nginx testen und reloaden
nginx -t
systemctl reload nginx
```

### 3. Cloudflare

Da du Cloudflare nutzt (Subdomain schon hinzugefügt):
- Cloudflare managed automatisch SSL
- A Record zeigt auf 46.4.25.15
- Proxy aktiviert (orange cloud)

### 4. Testen

```bash
# Health check
curl http://hub.stumpfworks.de/health

# Templates abrufen
curl http://hub.stumpfworks.de/api/v1/templates

# Kategorien
curl http://hub.stumpfworks.de/api/v1/templates/categories
```

## NAS Konfiguration

Das NAS Backend ist schon vorbereitet! Es nutzt automatisch `https://hub.stumpfworks.de`.

Falls du eine andere URL testen willst:
```bash
export STUMPFWORKS_HUB_URL=http://hub.stumpfworks.de
```

## Monitoring

### Hub Status prüfen:
```bash
ssh root@46.4.25.15
systemctl status stumpfworks-hub
```

### Logs anschauen:
```bash
journalctl -u stumpfworks-hub -f
```

### Templates aktualisieren:
```bash
# Neue JSON Datei hinzufügen
nano /var/lib/stumpfworks-hub/templates/media/neue-app.json

# Hub neustarten (lädt alle Templates neu)
systemctl restart stumpfworks-hub
```

## Nächste Schritte

Nach erfolgreichem Deployment:

1. **NAS testen:**
   - Docker Manager → Templates Tab
   - Sollte jetzt 10 Templates vom Hub zeigen

2. **Timeout-Fix testen:**
   - Template deployen (z.B. Jellyfin)
   - Sollte jetzt nicht mehr timeout nach 60 Sekunden
   - 10 Minuten Timeout für Deployments

3. **Weitere Templates hinzufügen:**
   - JSON Dateien in `/var/lib/stumpfworks-hub/templates/{kategorie}/`
   - Hub restart oder warte 60 Minuten (Cache TTL)

## Troubleshooting

### Hub startet nicht:
```bash
journalctl -u stumpfworks-hub -n 50
```

### NAS kann Hub nicht erreichen:
```bash
# Vom NAS aus testen
curl http://hub.stumpfworks.de/health
```

### Template wird nicht geladen:
```bash
# JSON validieren
cd /var/lib/stumpfworks-hub/templates
find . -name "*.json" -exec python3 -m json.tool {} \; > /dev/null
```

## API Endpoints

Alle Endpoints unter `http://hub.stumpfworks.de/api/v1`:

- `GET /templates` - Liste aller Templates
- `GET /templates/{id}` - Einzelnes Template
- `GET /templates/categories` - Alle Kategorien
- `GET /templates/search?q=plex` - Suche

## Files auf dem Server

```
/opt/stumpfworks-hub/
├── hub                           # Binary
└── deployment/                   # Config files

/var/lib/stumpfworks-hub/
├── templates/
│   ├── media/                    # 3 templates
│   ├── automation/               # 3 templates
│   ├── download/                 # 2 templates
│   └── monitoring/               # 2 templates
└── apps/                         # Future: App Store

/etc/systemd/system/
└── stumpfworks-hub.service       # Service file

/etc/nginx/sites-available/
└── hub.stumpfworks.de            # Nginx config
```

## Ready to Deploy! 🚀

Alles ist vorbereitet. Einfach `./deployment/deploy.sh 46.4.25.15` ausführen!
