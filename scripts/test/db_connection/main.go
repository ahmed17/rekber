package main

import (
    "database/sql"
    "fmt"
    "log"
    
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "rekber/scripts/utils"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Build connection string
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        utils.GetEnv("DB_HOST", "localhost"),
        utils.GetEnv("DB_PORT", "5433"),
        utils.GetEnv("DB_USER", "postgres"),
        utils.GetEnv("DB_PASSWORD", "password"),
        utils.GetEnv("DB_NAME", "db_rekber"),
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