package main

import (
    "fmt"
    "log"
    // "os"
    "time"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "rekber/internal/models"
    "rekber/scripts/utils"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
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

    // Check if tables exist
    if !tablesExist(db) {
        log.Fatal("Tables don't exist. Please run migrations first: make migrate-up")
    }

    // Seed data
    seedData(db)
    log.Println("✅ Database seeded successfully!")
}

func tablesExist(db *gorm.DB) bool {
    var count int64
    db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
    return count > 0
}

func seedData(db *gorm.DB) {
    // Check if data already exists
    var existingUserCount int64
    db.Model(&models.User{}).Count(&existingUserCount)
    
    if existingUserCount > 0 {
        log.Println("⚠️  Data already exists, skipping seed")
        return
    }

    // Create users
    users := []models.User{
        {
            Username:     "admin",
            Email:        "admin@rekber.com",
            PasswordHash: "$2a$10$YourHashedPasswordHere", // In production, use bcrypt
            FullName:     "Administrator",
            Phone:        "081234567890",
            Role:         "admin",
            IsVerified:   true,
            IsActive:     true,
            CreatedAt:    time.Now(),
            UpdatedAt:    time.Now(),
        },
        {
            Username:     "seller1",
            Email:        "seller1@example.com",
            PasswordHash: "$2a$10$YourHashedPasswordHere",
            FullName:     "John Seller",
            Phone:        "081234567891",
            Role:         "seller",
            IsVerified:   true,
            IsActive:     true,
            BankName:     "Bank ABC",
            AccountNumber: "1234567890",
            AccountHolder: "John Seller",
            CreatedAt:    time.Now(),
            UpdatedAt:    time.Now(),
        },
        {
            Username:     "buyer1",
            Email:        "buyer1@example.com",
            PasswordHash: "$2a$10$YourHashedPasswordHere",
            FullName:     "Jane Buyer",
            Phone:        "081234567892",
            Role:         "buyer",
            IsVerified:   true,
            IsActive:     true,
            CreatedAt:    time.Now(),
            UpdatedAt:    time.Now(),
        },
    }

    for i := range users {
        result := db.Create(&users[i])
        if result.Error != nil {
            log.Printf("Error creating user %s: %v", users[i].Username, result.Error)
        } else {
            log.Printf("Created user: %s (ID: %d)", users[i].Username, users[i].ID)
        }
    }

    // Create transactions
    transactions := []models.Transaction{
        {
            TransactionCode: "REKBER-001",
            Title:           "Laptop Gaming",
            Description:     "Laptop gaming dengan spek tinggi",
            Amount:          15000000,
            Fee:             75000,
            Status:          models.StatusCompleted,
            BuyerID:         3, // buyer1
            SellerID:        2, // seller1
            PaidAt:          timePtr(time.Now().Add(-7 * 24 * time.Hour)),
            ValidatedAt:     timePtr(time.Now().Add(-6 * 24 * time.Hour)),
            ProcessingAt:    timePtr(time.Now().Add(-5 * 24 * time.Hour)),
            ShippedAt:       timePtr(time.Now().Add(-4 * 24 * time.Hour)),
            CompletedAt:     timePtr(time.Now().Add(-3 * 24 * time.Hour)),
            CreatedAt:       time.Now().Add(-8 * 24 * time.Hour),
            UpdatedAt:       time.Now().Add(-3 * 24 * time.Hour),
        },
        {
            TransactionCode: "REKBER-002",
            Title:           "Smartphone Flagship",
            Description:     "Smartphone terbaru dengan kamera 108MP",
            Amount:          12000000,
            Fee:             60000,
            Status:          models.StatusShipped,
            BuyerID:         3,
            SellerID:        2,
            ShippingData:    "JNE Express - Tracking No: 1234567890",
            PaidAt:          timePtr(time.Now().Add(-2 * 24 * time.Hour)),
            ValidatedAt:     timePtr(time.Now().Add(-2 * 24 * time.Hour)),
            ProcessingAt:    timePtr(time.Now().Add(-1 * 24 * time.Hour)),
            ShippedAt:       timePtr(time.Now().Add(-12 * time.Hour)),
            AutoCompleteAt:  timePtr(time.Now().Add(36 * time.Hour)), // 2x24 jam dari shipped
            CreatedAt:       time.Now().Add(-3 * 24 * time.Hour),
            UpdatedAt:       time.Now().Add(-12 * time.Hour),
        },
        {
            TransactionCode: "REKBER-003",
            Title:           "Monitor 4K 27 inch",
            Description:     "Monitor gaming 144Hz",
            Amount:          5000000,
            Fee:             25000,
            Status:          models.StatusWaitingPayment,
            BuyerID:         3,
            SellerID:        2,
            CreatedAt:       time.Now().Add(-2 * time.Hour),
            UpdatedAt:       time.Now().Add(-2 * time.Hour),
        },
    }

    // Recalculate total amount
    for i := range transactions {
        transactions[i].TotalAmount = transactions[i].Amount + transactions[i].Fee
    }

    for i := range transactions {
        result := db.Create(&transactions[i])
        if result.Error != nil {
            log.Printf("Error creating transaction %s: %v", transactions[i].TransactionCode, result.Error)
        } else {
            log.Printf("Created transaction: %s (ID: %d)", transactions[i].TransactionCode, transactions[i].ID)
        }
    }

    // Create testimonies
    testimonies := []models.Testimony{
        {
            TransactionID: 1,
            UserID:        3,
            Rating:        5,
            Comment:       "Sangat puas dengan transaksi ini! Penjual responsif dan barang sesuai deskripsi.",
            CreatedAt:     time.Now().Add(-2 * 24 * time.Hour),
            UpdatedAt:     time.Now().Add(-2 * 24 * time.Hour),
        },
    }

    for i := range testimonies {
        result := db.Create(&testimonies[i])
        if result.Error != nil {
            log.Printf("Error creating testimony for transaction %d: %v", testimonies[i].TransactionID, result.Error)
        } else {
            log.Printf("Created testimony for transaction: %d", testimonies[i].TransactionID)
        }
    }
}

func timePtr(t time.Time) *time.Time {
    return &t
}