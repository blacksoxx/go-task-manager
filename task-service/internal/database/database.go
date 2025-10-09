package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
    "time"
    "taskservice/internal/models"
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
    CREATE TABLE IF NOT EXISTS tasks (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title VARCHAR(255) NOT NULL,
        description TEXT,
        status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed')),
        user_id UUID NOT NULL,
        due_date TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
    CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
    CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
    `

    _, err := db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to initialize tasks schema: %w", err)
    }

    log.Println("✅ Task Service: Database schema initialized successfully")
    return nil
}

func (db *DB) CreateTask(task *models.Task) error {
    query := `
    INSERT INTO tasks (title, description, user_id, due_date)
    VALUES ($1, $2, $3, $4)
    RETURNING id, status, created_at, updated_at`

    err := db.QueryRow(query, task.Title, task.Description, task.UserID, task.DueDate).Scan(
        &task.ID, &task.Status, &task.CreatedAt, &task.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create task: %w", err)
    }
    
    return nil
}

func (db *DB) GetTaskByID(id string) (*models.Task, error) {
    var task models.Task
    
    query := `SELECT id, title, description, status, user_id, due_date, created_at, updated_at 
              FROM tasks WHERE id = $1`
    
    err := db.QueryRow(query, id).Scan(
        &task.ID, &task.Title, &task.Description, &task.Status, &task.UserID,
        &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("task not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get task: %w", err)
    }
    
    return &task, nil
}

func (db *DB) GetTasksByUserID(userID string) ([]*models.Task, error) {
    query := `SELECT id, title, description, status, user_id, due_date, created_at, updated_at 
              FROM tasks WHERE user_id = $1 ORDER BY created_at DESC`
    
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query tasks: %w", err)
    }
    defer rows.Close()

    var tasks []*models.Task
    for rows.Next() {
        var task models.Task
        err := rows.Scan(
            &task.ID, &task.Title, &task.Description, &task.Status, &task.UserID,
            &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan task: %w", err)
        }
        tasks = append(tasks, &task)
    }

    return tasks, nil
}

func (db *DB) UpdateTask(id string, req *models.UpdateTaskRequest) (*models.Task, error) {
    query := `
    UPDATE tasks 
    SET title = COALESCE($1, title),
        description = COALESCE($2, description),
        status = COALESCE($3, status),
        due_date = $4,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = $5
    RETURNING id, title, description, status, user_id, due_date, created_at, updated_at`

    var task models.Task
    err := db.QueryRow(query, req.Title, req.Description, req.Status, req.DueDate, id).Scan(
        &task.ID, &task.Title, &task.Description, &task.Status, &task.UserID,
        &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("task not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to update task: %w", err)
    }
    
    return &task, nil
}

func (db *DB) DeleteTask(id string) error {
    query := `DELETE FROM tasks WHERE id = $1`
    
    result, err := db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}

func (db *DB) Close() error {
    return db.DB.Close()
}
