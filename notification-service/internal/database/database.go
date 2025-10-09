package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"notification-service/internal/models"
	"time"
	"os"
    "regexp"
    "strings"
	_ "github.com/lib/pq"
)


type DB struct {
    *sql.DB
}

func NewPostgresDB() (*DB, error) {
    // Get database URL from environment variable
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        // Fallback to localhost (for development outside Docker)
        connStr = "host=localhost port=5432 user=postgres password=password dbname=taskmanager sslmode=disable"
    } else {
        // Log that we're using the environment variable (for debugging)
        log.Printf("🔗 Using DATABASE_URL from environment: %s", maskPassword(connStr))
    }
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Add connection retry logic
    var pingErr error
    for i := 0; i < 10; i++ {
        pingErr = db.Ping()
        if pingErr == nil {
            break
        }
        log.Printf("⏳ Database not ready yet (attempt %d/10): %v", i+1, pingErr)
        time.Sleep(2 * time.Second)
    }

    if pingErr != nil {
        return nil, fmt.Errorf("failed to ping database after retries: %w", pingErr)
    }

    log.Println("✅ Service: Connected to PostgreSQL successfully")
    return &DB{db}, nil
}

// Helper function to mask password in logs
func maskPassword(connStr string) string {
    if strings.Contains(connStr, "password=") {
        return regexp.MustCompile(`password=[^&]+`).ReplaceAllString(connStr, "password=***")
    }
    return connStr
}

func (db *DB) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		type VARCHAR(20) NOT NULL CHECK (type IN ('email', 'in_app', 'push')),
		status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'read')),
		data JSONB DEFAULT '{}'::jsonb,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		read_at TIMESTAMP NULL
	);

	CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
	CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to initialize notifications schema: %w", err)
	}

	log.Println("✅ Notification Service: Database schema initialized successfully")
	return nil
}

func (db *DB) CreateNotification(notification *models.Notification) error {
	query := `
	INSERT INTO notifications (user_id, title, message, type, status, data)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`

	// Handle JSON data properly - use json.RawMessage or empty object
	var jsonData interface{}
	if notification.Data != nil && len(notification.Data) > 0 {
		// Convert map to JSON bytes
		dataBytes, err := json.Marshal(notification.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		jsonData = dataBytes
	} else {
		// Use empty JSON object
		jsonData = []byte("{}")
	}

	err := db.QueryRow(query,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.Type,
		notification.Status,
		jsonData,
	).Scan(&notification.ID, &notification.CreatedAt, &notification.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}
func (db *DB) GetNotificationByID(id string) (*models.Notification, error) {
	var notification models.Notification
	var dataBytes []byte

	query := `SELECT id, user_id, title, message, type, status, data, created_at, updated_at, read_at 
			  FROM notifications WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Title,
		&notification.Message,
		&notification.Type,
		&notification.Status,
		&dataBytes,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&notification.ReadAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("notification not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	// Unmarshal JSON data - handle empty case
	if len(dataBytes) > 0 && string(dataBytes) != "null" {
		if err := json.Unmarshal(dataBytes, &notification.Data); err != nil {
			// If unmarshal fails, set empty map
			notification.Data = make(map[string]interface{})
		}
	} else {
		notification.Data = make(map[string]interface{})
	}

	return &notification, nil
}

func (db *DB) GetNotificationsByUserID(userID string, limit, offset int) ([]*models.Notification, error) {
	query := `SELECT id, user_id, title, message, type, status, data, created_at, updated_at, read_at 
			  FROM notifications WHERE user_id = $1 
			  ORDER BY created_at DESC 
			  LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var notification models.Notification
		var dataBytes []byte

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Message,
			&notification.Type,
			&notification.Status,
			&dataBytes,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&notification.ReadAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		// Unmarshal JSON data - handle empty case
		if len(dataBytes) > 0 && string(dataBytes) != "null" {
			if err := json.Unmarshal(dataBytes, &notification.Data); err != nil {
				notification.Data = make(map[string]interface{})
			}
		} else {
			notification.Data = make(map[string]interface{})
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (db *DB) MarkNotificationAsRead(id string) error {
	query := `UPDATE notifications 
			  SET status = $1, read_at = $2, updated_at = $2 
			  WHERE id = $3 AND read_at IS NULL`

	result, err := db.Exec(query, models.StatusRead, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or already read")
	}

	return nil
}

func (db *DB) UpdateNotificationStatus(id string, status models.NotificationStatus) error {
	query := `UPDATE notifications 
			  SET status = $1, updated_at = $2 
			  WHERE id = $3`

	result, err := db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (db *DB) DeleteNotification(id string) error {
	query := `DELETE FROM notifications WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}