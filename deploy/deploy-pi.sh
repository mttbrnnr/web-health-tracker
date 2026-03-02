#!/bin/bash
set -e

# Configuration
PI_HOST="${PI_HOST:-raspberrypi.local}"
PI_USER="${PI_USER:-pi}"
INSTALL_DIR="${INSTALL_DIR:-/opt/health-tracker}"
BINARY_NAME="health-tracker"

echo "==> Building for Raspberry Pi (linux/arm64)..."
make build-pi

echo "==> Creating directories on Pi..."
ssh "${PI_USER}@${PI_HOST}" "sudo mkdir -p ${INSTALL_DIR}/data && sudo chown -R ${PI_USER}:${PI_USER} ${INSTALL_DIR}"

echo "==> Copying binary to Pi..."
scp "./build/${BINARY_NAME}-arm64" "${PI_USER}@${PI_HOST}:${INSTALL_DIR}/${BINARY_NAME}"

echo "==> Copying systemd service file..."
scp "./deploy/health-tracker.service" "${PI_USER}@${PI_HOST}:/tmp/health-tracker.service"
ssh "${PI_USER}@${PI_HOST}" "sudo mv /tmp/health-tracker.service /etc/systemd/system/"

echo "==> Reloading systemd and restarting service..."
ssh "${PI_USER}@${PI_HOST}" "sudo systemctl daemon-reload && sudo systemctl enable health-tracker && sudo systemctl restart health-tracker"

echo "==> Checking service status..."
ssh "${PI_USER}@${PI_HOST}" "sudo systemctl status health-tracker --no-pager" || true

echo ""
echo "==> Deployment complete!"
echo "    Access at: http://${PI_HOST}:3000"
