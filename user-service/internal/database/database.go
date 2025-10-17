package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
    "time"
    "userservice/internal/models" 
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
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email VARCHAR(255) UNIQUE NOT NULL,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    `

    _, err := db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to initialize database schema: %w", err)
    }

    log.Println("✅ Database schema initialized successfully")
    return nil
}

func (db *DB) CreateUser(user *models.User, passwordHash string) error {
    query := `
    INSERT INTO users (email, first_name, last_name, password_hash)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, updated_at`

    err := db.QueryRow(query, user.Email, user.FirstName, user.LastName, passwordHash).Scan(
        &user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

func (db *DB) GetUserByEmail(email string) (*models.User, string, error) {
    var user models.User
    var passwordHash string
    
    query := `SELECT id, email, first_name, last_name, password_hash, created_at, updated_at 
              FROM users WHERE email = $1`
    
    err := db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.FirstName, &user.LastName, 
        &passwordHash, &user.CreatedAt, &user.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, "", fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, "", fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, passwordHash, nil
}

func (db *DB) GetUserByID(id string) (*models.User, error) {
    var user models.User
    
    query := `SELECT id, email, first_name, last_name, created_at, updated_at 
              FROM users WHERE id = $1`
    
    err := db.QueryRow(query, id).Scan(
        &user.ID, &user.Email, &user.FirstName, &user.LastName, 
        &user.CreatedAt, &user.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

func (db *DB) Close() error {
    return db.DB.Close()
}