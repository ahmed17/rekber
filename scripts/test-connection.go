package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Build connection string
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_PORT", "5433"),
        getEnv("DB_USER", "postgres"),
        getEnv("DB_PASSWORD", "postgres"),
        getEnv("DB_NAME", "db_rekber"),
    )
    
    fmt.Println("Connection string:", connStr)
    
    // Connect
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer db.Close()
    
    // Test connection
    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    fmt.Println("✅ Database connected successfully!")
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}