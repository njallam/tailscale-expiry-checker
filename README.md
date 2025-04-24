# Tailscale Expiry Checker
Oneshot Go systemd user service & timer to periodically check for Tailscale node key expiry, sending a notification if the key is due to expire in the next 24 hours (number of hours adjustable with `-hours` flag)

## Installation
* Clone this repository
* Run [/install.sh](install.sh), or use as a reference for installing manually

## Disclaimer
Tailscale is a registered trademark of Tailscale Inc. This project is not affiliated with, endorsed by, or sponsored by Tailscale Inc.