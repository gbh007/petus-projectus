package client

import (
	notificationClient "app/services/notification/client"
)

const (
	ButtonKind   = notificationClient.ButtonKind
	SuccessLevel = notificationClient.SuccessLevel
	ErrorLevel   = notificationClient.ErrorLevel
)

type Notification notificationClient.Notification
