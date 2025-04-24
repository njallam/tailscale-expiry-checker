package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/godbus/dbus/v5"
	"tailscale.com/client/local"
)

const NotificationIDPath = "/tmp/tailscale-expiry-checker-notification-id"

func main() {
	hours := flag.Int("hours", 24, "hours")
	flag.Parse()

	ctx := context.Background()

	var ts local.Client
	status, err := ts.Status(ctx)
	if err != nil {
		fmt.Printf("Failed to get Tailscale status: %v", err)
		os.Exit(1)
	}

	timeToExpiry := time.Until(*status.Self.KeyExpiry).Truncate(time.Second)
	fmt.Printf("Tailscale node key expiring in %s\n", timeToExpiry)

	dbusConn, err := dbus.SessionBus()
	if err != nil {
		fmt.Printf("Failed to connect to DBus session bus: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbusConn.Close(); err != nil {
			fmt.Printf("Failed to close DBus connection: %v", err)
		}
	}()

	notificationsObject := dbusConn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	var notificationID uint32 = 0
	if data, err := os.ReadFile(NotificationIDPath); err == nil {
		if id, err := strconv.ParseUint(string(data), 10, 32); err == nil {
			notificationID = uint32(id)
		}
	}

	if timeToExpiry > time.Duration(*hours)*time.Hour {
		call := notificationsObject.Call("org.freedesktop.Notifications.CloseNotification", 0, notificationID)
		if call.Err != nil {
			fmt.Printf("Failed to close notification: %v", call.Err)
			os.Exit(1)
		}
		return
	}

	var title, body string
	if timeToExpiry <= 0 {
		title = "Node Key Expired"
		body = "Your node key has expired."
	} else {
		title = "Node Key Expiring Soon"
		body = fmt.Sprintf("Your node key will expire in %s.\n", timeToExpiry)
	}

	call := notificationsObject.Call("org.freedesktop.Notifications.Notify", 0,
		"Tailscale Expiry Checker", // Application Name
		notificationID,             // Notification ID
		"tailscale",                // Icon
		title,                      // Title
		body,                       // Body
		[]string{},                 // No actions
		map[string]dbus.Variant{},  // Hints
		int32(-1),                  // No timeout
	)
	if call.Err != nil {
		fmt.Printf("Failed to send notification: %v", call.Err)
		os.Exit(1)
	}

	err = os.WriteFile(NotificationIDPath, []byte(strconv.FormatUint(uint64(call.Body[0].(uint32)), 10)), 0644)
	if err != nil {
		fmt.Printf("Failed to save notification ID: %v", err)
		os.Exit(1)
	}
}
