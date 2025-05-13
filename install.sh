#!/bin/bash
set -euo pipefail

if [ "$(id -u)" -eq 0 ]; then
    echo "This script should only be run as a regular user, not root."
    exit 1
fi

go install
mkdir -p ~/.config/systemd/user
sed "s|\$HOME|$HOME|g" "tailscale-expiry-checker.service" > ~/.config/systemd/user/tailscale-expiry-checker.service
cp tailscale-expiry-checker.timer ~/.config/systemd/user/tailscale-expiry-checker.timer
systemctl --user daemon-reload
systemctl --user enable --now tailscale-expiry-checker.timer
