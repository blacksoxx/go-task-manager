package handlers

import (
	"encoding/json"
	"net/http"
	"notification-service/internal/database"
	"notification-service/internal/models"
	"strconv"

	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	db *database.DB
}

func NewNotificationHandler(db *database.DB) *NotificationHandler {
	return &NotificationHandler{db: db}
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.UserID == "" || req.Title == "" || req.Message == "" {
		http.Error(w, `{"error": "user_id, title, and message are required"}`, http.StatusBadRequest)
		return
	}

	// Validate notification type
	if req.Type != models.EmailNotification && req.Type != models.InAppNotification && req.Type != models.PushNotification {
		http.Error(w, `{"error": "Invalid notification type. Must be email, in_app, or push"}`, http.StatusBadRequest)
		return
	}

	// Ensure Data is not nil
	if req.Data == nil {
		req.Data = make(map[string]interface{})
	}

	notification := &models.Notification{
		UserID:  req.UserID,
		Title:   req.Title,
		Message: req.Message,
		Type:    req.Type,
		Status:  models.StatusPending,
		Data:    req.Data, // This will now be {} instead of nil
	}

	if err := h.db.CreateNotification(notification); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(models.NotificationResponse{Notification: notification}); err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
	}
}

func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]

	notification, err := h.db.GetNotificationByID(notificationID)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(models.NotificationResponse{Notification: notification}); err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
	}
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	notifications, err := h.db.GetNotificationsByUserID(userID, limit, offset)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := models.NotificationsResponse{
		Notifications: notifications,
		Total:         len(notifications),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
	}
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]

	if err := h.db.MarkNotificationAsRead(notificationID); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]

	var req struct {
		Status models.NotificationStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate status
	if req.Status != models.StatusSent && req.Status != models.StatusFailed && req.Status != models.StatusPending {
		http.Error(w, `{"error": "Invalid status. Must be sent, failed, or pending"}`, http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateNotificationStatus(notificationID, req.Status); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["id"]

	if err := h.db.DeleteNotification(notificationID); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":    "healthy",
		"service":   "notificationservice",
		"timestamp": "2024-01-01T00:00:00Z", // We'll make this dynamic in main.go
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
	}
}