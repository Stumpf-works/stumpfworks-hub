# Deployment Guide - Stumpfworks Hub

## Server Setup (apt.stumpfworks.de)

### 1. Build the Hub

```bash
# On your development machine
cd /path/to/stumpfworks-hub
GOOS=linux GOARCH=amd64 go build -o hub cmd/hub/main.go

# Or build on server directly
ssh apt.stumpfworks.de
cd /opt/stumpfworks-hub
go build -o hub cmd/hub/main.go
```

### 2. Directory Structure

```bash
# Create directories on server
sudo mkdir -p /opt/stumpfworks-hub
sudo mkdir -p /var/lib/stumpfworks-hub/templates
sudo mkdir -p /var/lib/stumpfworks-hub/apps
sudo mkdir -p /var/log/stumpfworks-hub

# Set permissions
sudo chown -R www-data:www-data /opt/stumpfworks-hub
sudo chown -R www-data:www-data /var/lib/stumpfworks-hub
sudo chown -R www-data:www-data /var/log/stumpfworks-hub
```

### 3. Copy Files

```bash
# Copy binary
sudo cp hub /opt/stumpfworks-hub/

# Copy templates
sudo cp -r templates/* /var/lib/stumpfworks-hub/templates/

# Copy service file
sudo cp deployment/stumpfworks-hub.service /etc/systemd/system/
```

### 4. Configure Systemd Service

```bash
# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable stumpfworks-hub
sudo systemctl start stumpfworks-hub

# Check status
sudo systemctl status stumpfworks-hub

# View logs
sudo journalctl -u stumpfworks-hub -f
```

### 5. Nginx Reverse Proxy

Add to nginx config:

```nginx
server {
    listen 80;
    server_name hub.stumpfworks.de;

    location / {
        proxy_pass http://localhost:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 6. SSL Certificate

```bash
sudo certbot --nginx -d hub.stumpfworks.de
```

## NAS Configuration

On Stumpfworks NAS instances, set the Hub URL:

```bash
# In /etc/environment or systemd service file
STUMPFWORKS_HUB_URL=https://hub.stumpfworks.de
```

If not set, the NAS will default to `https://hub.stumpfworks.de`.

For development/testing, you can override:

```bash
STUMPFWORKS_HUB_URL=http://localhost:8090
```

## Updating Templates

### Add New Template

1. Create JSON file in `/var/lib/stumpfworks-hub/templates/{category}/{name}.json`
2. Hub will automatically load it on next request (cache TTL: 60 minutes)
3. Or restart Hub: `sudo systemctl restart stumpfworks-hub`

Example:

```bash
sudo nano /var/lib/stumpfworks-hub/templates/media/emby.json
```

```json
{
  "id": "emby",
  "name": "Emby Server",
  "description": "Media server alternative",
  "icon": "📺",
  "category": "media",
  "author": "StumpfWorks",
  "version": "1.0.0",
  "variables": {
    "MEDIA_PATH": "/mnt/media"
  },
  "compose": "...",
  "requirements": {},
  "tags": ["media"],
  "created_at": "2025-12-04T22:00:00Z",
  "updated_at": "2025-12-04T22:00:00Z"
}
```

### Force Cache Refresh

```bash
# Restart Hub
sudo systemctl restart stumpfworks-hub

# Or wait for cache TTL (60 minutes)
```

## Monitoring

### Health Check

```bash
curl https://hub.stumpfworks.de/health
```

### Logs

```bash
# Real-time logs
sudo journalctl -u stumpfworks-hub -f

# Last 100 lines
sudo journalctl -u stumpfworks-hub -n 100

# Today's logs
sudo journalctl -u stumpfworks-hub --since today
```

## Backup

### Templates Backup

```bash
# Backup templates
tar -czf hub-templates-$(date +%Y%m%d).tar.gz /var/lib/stumpfworks-hub/templates/

# Restore
tar -xzf hub-templates-20251204.tar.gz -C /
```

## Troubleshooting

### Hub Not Starting

```bash
# Check service status
sudo systemctl status stumpfworks-hub

# Check logs
sudo journalctl -u stumpfworks-hub -n 50

# Check permissions
ls -la /opt/stumpfworks-hub
ls -la /var/lib/stumpfworks-hub
```

### NAS Can't Reach Hub

```bash
# Test from NAS
curl https://hub.stumpfworks.de/health
curl https://hub.stumpfworks.de/api/v1/templates

# Check NAS environment
systemctl show stumpfworks-nas | grep STUMPFWORKS_HUB_URL
```

### Templates Not Loading

```bash
# Validate JSON
cd /var/lib/stumpfworks-hub/templates
find . -name "*.json" -exec json_pp < {} \; > /dev/null

# Check Hub logs for parsing errors
sudo journalctl -u stumpfworks-hub | grep -i error
```
