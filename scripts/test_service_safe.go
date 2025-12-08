//go:build ignore
// +build ignore

package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "rekber/internal/config"
    "rekber/internal/repositories"
    "rekber/internal/services"
)

func main() {
    // Ensure .env exists
    if _, err := os.Stat(".env"); os.IsNotExist(err) {
        createDefaultEnv()
    }
    
    // Load environment
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found, using defaults")
    }
    
    // Load configuration
    if err := config.LoadConfig(); err != nil {
        log.Printf("Warning: Failed to load config: %v", err)
        // Set default config
        setupDefaultConfig()
    }

    // Test database connection
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        getConfigValue("DB_HOST", "localhost"),
        getConfigValue("DB_PORT", "5433"),
        getConfigValue("DB_USER", "postgres"),
        getConfigValue("DB_PASSWORD", "password"),
        getConfigValue("DB_NAME", "db_rekber"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("⚠️  Cannot connect to database: %v", err)
        log.Println("⚠️  Testing auth service only...")
        testAuthServiceOnly()
        return
    }

    // Initialize repository
    repo := &repositories.Repository{
        User:        repositories.NewUserRepository(db),
        Transaction: repositories.NewTransactionRepository(db),
        Testimony:   repositories.NewTestimonyRepository(db),
    }
    
    // Initialize service
    services.InitService(repo)
    svc := services.GetService()
    
    fmt.Println("✅ Service layer initialized successfully!")
    testAllServices(svc)
}

func testAuthServiceOnly() {
    // Manual auth service creation for testing
    authSvc := services.NewAuthService()
    
    token, err := authSvc.GenerateToken(1, "admin")
    if err != nil {
        fmt.Printf("❌ Error generating token: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Generated token: %s...\n", safeSubstring(token, 20))
    
    userID, role, err := authSvc.ValidateToken(token)
    if err != nil {
        fmt.Printf("❌ Error validating token: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Validated token - UserID: %d, Role: %s\n", userID, role)
    fmt.Println("\n⚠️  Note: Database tests skipped (connection failed)")
}

func testAllServices(svc *services.Service) {
    fmt.Println("Available services:")
    fmt.Println("  - User Service")
    fmt.Println("  - Auth Service")
    fmt.Println("  - Transaction Service")
    fmt.Println("  - Testimony Service")
    
    // Test auth service
    token, err := svc.Auth.GenerateToken(1, "admin")
    if err != nil {
        fmt.Printf("❌ Error generating token: %v\n", err)
    } else {
        fmt.Printf("\n✅ Generated token: %s...\n", safeSubstring(token, 20))
        
        // Validate token
        userID, role, err := svc.Auth.ValidateToken(token)
        if err != nil {
            fmt.Printf("❌ Error validating token: %v\n", err)
        } else {
            fmt.Printf("✅ Validated token - UserID: %d, Role: %s\n", userID, role)
        }
    }
    
    // Test user service
    fmt.Println("\n=== Testing User Service ===")
    if userResp, err := svc.User.GetUserByID(1); err != nil {
        fmt.Printf("⚠️  Cannot get user by ID: %v\n", err)
    } else {
        fmt.Printf("✅ Found user: %s (%s)\n", userResp.Username, userResp.Email)
    }
    
    // Test transaction service
    fmt.Println("\n=== Testing Transaction Service ===")
    if transactions, total, err := svc.Transaction.ListAllTransactions(1, 10); err != nil {
        fmt.Printf("⚠️  Cannot list transactions: %v\n", err)
    } else {
        fmt.Printf("📊 Total transactions: %d\n", total)
        if len(transactions) > 0 {
            fmt.Printf("📋 Found %d transactions\n", len(transactions))
        }
    }
    
    fmt.Println("\n✅ Service layer test completed!")
}

func createDefaultEnv() {
    fmt.Println("Creating default .env file...")
    content := `APP_NAME=Rekber
APP_ENV=development
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=db_rekber
DB_SSLMODE=disable
JWT_SECRET=test-jwt-secret-change-in-production
JWT_EXPIRE_HOURS=72`
    
    if err := os.WriteFile(".env", []byte(content), 0644); err != nil {
        log.Printf("Error creating .env: %v", err)
    }
}

func setupDefaultConfig() {
    // Set default config if LoadConfig fails
    config.AppConfig = &config.Config{
        JWTSecret: "test-jwt-secret-change-in-production",
        JWTExpire: 72,
    }
}

func getConfigValue(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func safeSubstring(s string, n int) string {
    if len(s) > n {
        return s[:n]
    }
    return s
}