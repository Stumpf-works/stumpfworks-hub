#!/bin/bash
# Stumpfworks Hub Deployment Script
# Usage: ./deploy.sh [server]
# Example: ./deploy.sh apt.stumpfworks.de

set -e

SERVER=${1:-apt.stumpfworks.de}
HUB_DIR="/opt/stumpfworks-hub"
DATA_DIR="/var/lib/stumpfworks-hub"
LOG_DIR="/var/log/stumpfworks-hub"

echo "🚀 Deploying Stumpfworks Hub to $SERVER..."

# Build for Linux
echo "📦 Building Hub binary..."
GOOS=linux GOARCH=amd64 go build -o hub cmd/hub/main.go

# Create tar package
echo "📦 Creating deployment package..."
tar -czf hub-deploy.tar.gz hub templates/ deployment/stumpfworks-hub.service

# Upload to server
echo "📤 Uploading to server..."
scp hub-deploy.tar.gz root@$SERVER:/tmp/

# Deploy on server
echo "🔧 Installing on server..."
ssh root@$SERVER << 'ENDSSH'
set -e

# Stop service if running
systemctl stop stumpfworks-hub || true

# Create directories
mkdir -p /opt/stumpfworks-hub
mkdir -p /var/lib/stumpfworks-hub/templates
mkdir -p /var/lib/stumpfworks-hub/apps
mkdir -p /var/log/stumpfworks-hub

# Extract package
cd /opt/stumpfworks-hub
tar -xzf /tmp/hub-deploy.tar.gz

# Copy templates to data directory
cp -r templates/* /var/lib/stumpfworks-hub/templates/ || true

# Install systemd service
cp deployment/stumpfworks-hub.service /etc/systemd/system/

# Set permissions
chown -R www-data:www-data /opt/stumpfworks-hub
chown -R www-data:www-data /var/lib/stumpfworks-hub
chown -R www-data:www-data /var/log/stumpfworks-hub
chmod +x /opt/stumpfworks-hub/hub

# Reload systemd and start service
systemctl daemon-reload
systemctl enable stumpfworks-hub
systemctl restart stumpfworks-hub

# Cleanup
rm /tmp/hub-deploy.tar.gz

echo "✅ Hub deployed successfully!"
systemctl status stumpfworks-hub
ENDSSH

# Cleanup local package
rm hub-deploy.tar.gz hub

echo ""
echo "✅ Deployment complete!"
echo "📡 Hub should be running at: http://$SERVER:8090"
echo ""
echo "Next steps:"
echo "1. Configure nginx reverse proxy for https://hub.stumpfworks.de"
echo "2. Run: certbot --nginx -d hub.stumpfworks.de"
echo "3. Test: curl https://hub.stumpfworks.de/health"
echo ""
echo "Useful commands:"
echo "  sudo systemctl status stumpfworks-hub"
echo "  sudo journalctl -u stumpfworks-hub -f"
echo "  curl http://$SERVER:8090/api/v1/templates"
