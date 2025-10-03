package database

import (
	"database/sql"
	"fmt"
	"log"
	"auth-service/internal/models"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sql.DB
}

func NewPostgresDB() (*DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=password dbname=taskmanager sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Auth Service: Connected to PostgreSQL successfully")
	return &DB{db}, nil
}

func (db *DB) Init() error {
    query := `
    -- Create users table if it doesn't exist
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email VARCHAR(255) UNIQUE NOT NULL,
        first_name VARCHAR(100) NOT NULL,
        last_name VARCHAR(100) NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    -- Migration: Ensure first_name exists
    DO $$ 
    BEGIN 
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                      WHERE table_name='users' AND column_name='first_name') THEN
            ALTER TABLE users ADD COLUMN first_name VARCHAR(100);
        END IF;
    END $$;

    -- Migration: Ensure last_name exists  
    DO $$ 
    BEGIN 
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                      WHERE table_name='users' AND column_name='last_name') THEN
            ALTER TABLE users ADD COLUMN last_name VARCHAR(100);
        END IF;
    END $$;

    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    `

    _, err := db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to initialize users schema: %w", err)
    }

    log.Println("✅ Auth Service: Database schema initialized successfully")
    return nil
}

func (db *DB) CreateUser(user *models.User, password string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
	INSERT INTO users (email, password_hash, first_name, last_name)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at`

	err = db.QueryRow(query, user.Email, string(hashedPassword), user.FirstName, user.LastName).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var passwordHash string

	query := `SELECT id, email, password_hash, first_name, last_name, created_at, updated_at 
			  FROM users WHERE email = $1`

	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FirstName, &user.LastName, 
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Password = passwordHash // Store hash for verification
	return &user, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}