package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/godbus/dbus/v5"
)

const NotificationIDPath = "/tmp/tailscale-expiry-checker-notification-id"

func getNotificationID() uint32 {
	if data, err := os.ReadFile(NotificationIDPath); err == nil {
		if id, err := strconv.ParseUint(string(data), 10, 32); err == nil {
			return uint32(id)
		}
	}
	return 0
}

func sendNotification(notificationsObject dbus.BusObject, notificationID uint32, title, body string) {
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
		fmt.Printf("Failed to send notification: %v\n", call.Err)
		os.Exit(1)
	}

	err := os.WriteFile(NotificationIDPath, []byte(strconv.FormatUint(uint64(call.Body[0].(uint32)), 10)), 0644)
	if err != nil {
		fmt.Printf("Failed to save notification ID: %v\n", err)
		os.Exit(1)
	}
}

func clearNotification(notificationsObject dbus.BusObject, notificationID uint32) {
	call := notificationsObject.Call("org.freedesktop.Notifications.CloseNotification", 0, notificationID)
	if call.Err != nil {
		fmt.Printf("Failed to close notification: %v\n", call.Err)
		os.Exit(1)
	}
}
