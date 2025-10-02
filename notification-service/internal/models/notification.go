package models

import (
	"encoding/json"
	"time"
)

type NotificationType string

const (
	EmailNotification NotificationType = "email"
	InAppNotification NotificationType = "in_app"
	PushNotification  NotificationType = "push"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
	StatusRead    NotificationStatus = "read"
)

type Notification struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Type      NotificationType       `json:"type"`
	Status    NotificationStatus     `json:"status"`
	Data      map[string]interface{} `json:"data,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	ReadAt    *time.Time             `json:"read_at,omitempty"`
}

type CreateNotificationRequest struct {
	UserID  string                 `json:"user_id"`
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Type    NotificationType       `json:"type"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type NotificationResponse struct {
	Notification *Notification `json:"notification"`
}

type NotificationsResponse struct {
	Notifications []*Notification `json:"notifications"`
	Total         int             `json:"total"`
}

type MarkAsReadRequest struct {
	ReadAt time.Time `json:"read_at,omitempty"`
}

// Scan and Value methods for JSONB data field
func (n *Notification) Scan(value interface{}) error {
	if value == nil {
		n.Data = nil
		return nil
	}
	
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil
	}
	
	return json.Unmarshal(data, &n.Data)
}