#!/bin/bash
if [ "$(id -u)" -eq 0 ]; then
    echo "This script should only be run as a regular user, not root."
    exit 1
fi

go install
sed "s|\$HOME|$HOME|g" "tailscale-expiry-checker.service" > ~/.config/systemd/user/tailscale-expiry-checker.service
cp tailscale-expiry-checker.timer ~/.config/systemd/user/tailscale-expiry-checker.timer
systemctl --user daemon-reload
systemctl --user enable --now tailscale-expiry-checker.timer
