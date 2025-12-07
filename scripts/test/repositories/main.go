package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "rekber/internal/repositories"
    "rekber/internal/models"
    "rekber/scripts/utils"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Database connection
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        utils.GetEnv("DB_HOST", "localhost"),
        utils.GetEnv("DB_PORT", "5433"),
        utils.GetEnv("DB_USER", "postgres"),
        utils.GetEnv("DB_PASSWORD", "password"),
        utils.GetEnv("DB_NAME", "db_rekber"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize repository
    repos := &repositories.Repository{
        User:        repositories.NewUserRepository(db),
        Transaction: repositories.NewTransactionRepository(db),
        Testimony:   repositories.NewTestimonyRepository(db),
    }

    // Test user repository
    fmt.Println("=== Testing User Repository ===")
    testUserRepository(repos.User, db)

    fmt.Println("\n✅ All repository tests completed!")
}

func testUserRepository(repo repositories.UserRepository, db *gorm.DB) {
    // Generate unique username and email
    timestamp := time.Now().Unix()
    uniqueUsername := fmt.Sprintf("testuser%d", timestamp)
    uniqueEmail := fmt.Sprintf("test%d@example.com", timestamp)
    
    // Create user with unique data
    user := &models.User{
        Username:     uniqueUsername,
        Email:        uniqueEmail,
        PasswordHash: "$2a$10$hashedpasswordfortesting",
        FullName:     "Test User",
        Phone:        "081234567890",
        Role:         "buyer",
    }
    
    if err := repo.Create(user); err != nil {
        log.Printf("Error creating user: %v", err)
    } else {
        fmt.Printf("✅ Created user with ID: %d, Username: %s\n", user.ID, user.Username)
    }

    // Find by email
    foundUser, err := repo.FindByEmail(uniqueEmail)
    if err != nil {
        log.Printf("Error finding user: %v", err)
    } else if foundUser != nil {
        fmt.Printf("✅ Found user by email: %s (ID: %d)\n", foundUser.Username, foundUser.ID)
    }

    // Count users
    count, err := repo.Count()
    if err != nil {
        log.Printf("Error counting users: %v", err)
    } else {
        fmt.Printf("📊 Total users: %d\n", count)
    }
}