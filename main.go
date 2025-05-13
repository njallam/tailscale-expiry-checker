package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/godbus/dbus/v5"
	"tailscale.com/client/local"
)

func main() {
	hours := flag.Int("hours", 24, "hours")
	flag.Parse()

	ctx := context.Background()

	dbusConn, err := dbus.SessionBus()
	if err != nil {
		fmt.Printf("Failed to connect to DBus session bus: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbusConn.Close(); err != nil {
			fmt.Printf("Failed to close DBus connection: %v\n", err)
		}
	}()

	notificationsObject := dbusConn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	notificationID := getNotificationID()

	var ts local.Client
	status, err := ts.Status(ctx)
	if err != nil {
		fmt.Printf("Failed to get Tailscale status: %v\n", err)
		os.Exit(1)
	}

	if status.BackendState == "NeedsLogin" {
		fmt.Println("Logged out.")
		var body string
		if status.AuthURL != "" {
			body = fmt.Sprintf("Log in at: %s", status.AuthURL)
		} else {
			body = "Tailscale is logged out."
		}
		sendNotification(notificationsObject, notificationID, "Needs Login", body)
		return
	}

	if status.Self.KeyExpiry == nil {
		fmt.Printf("Unable to get key expiry. BackendState is %s\n", status.BackendState)
		os.Exit(1)
	}

	timeToExpiry := time.Until(*status.Self.KeyExpiry).Truncate(time.Second)
	fmt.Printf("Tailscale node key expiring in %s\n", timeToExpiry)

	if timeToExpiry > time.Duration(*hours)*time.Hour {
		clearNotification(notificationsObject, notificationID)
		return
	}

	var title, body string
	if timeToExpiry <= 0 {
		title = "Node Key Expired"
		body = "Your node key has expired."
	} else {
		title = "Node Key Expiring Soon"
		body = fmt.Sprintf("Your node key will expire in %s.", timeToExpiry)
	}

	sendNotification(notificationsObject, notificationID, title, body)
}
