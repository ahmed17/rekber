//go:build ignore
// +build ignore

package main

import (
    "fmt"
    "log"
    
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "rekber/internal/config"
    "rekber/internal/repositories"
    "rekber/internal/services"
)

func main() {
    // Load environment
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Load configuration FIRST
    if err := config.LoadConfig(); err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Database connection
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.AppConfig.DBHost,
        config.AppConfig.DBPort,
        config.AppConfig.DBUser,
        config.AppConfig.DBPassword,
        config.AppConfig.DBName,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
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
        // Take only first 20 chars for display
        tokenPreview := token
        if len(token) > 20 {
            tokenPreview = token[:20]
        }
        fmt.Printf("✅ Generated token: %s...\n", tokenPreview)
        
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
    
    // Coba cari user yang sudah ada
    userResponse, err := svc.User.GetUserByID(1)
    if err != nil {
        fmt.Printf("⚠️  Cannot get user by ID: %v\n", err)
    } else {
        fmt.Printf("✅ Found user: %s (%s)\n", userResponse.Username, userResponse.Email)
    }
    
    // Test transaction service
    fmt.Println("\n=== Testing Transaction Service ===")
    testTransactionService(svc)
    
    fmt.Println("\n✅ Service layer test completed!")
}

func testTransactionService(svc *services.Service) {
    // Coba list transactions
    transactions, total, err := svc.Transaction.ListAllTransactions(1, 10)
    if err != nil {
        fmt.Printf("⚠️  Cannot list transactions: %v\n", err)
        return
    }
    
    fmt.Printf("📊 Total transactions: %d\n", total)
    if len(transactions) > 0 {
        fmt.Printf("📋 Found %d transactions\n", len(transactions))
        for i, tx := range transactions {
            fmt.Printf("  %d. %s - %s (%.2f)\n", 
                i+1, tx.TransactionCode, tx.Title, tx.Amount)
        }
    }
}