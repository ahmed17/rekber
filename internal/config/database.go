package config

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
    SSLMode  string
}

// GetDBConfig returns database configuration
func GetDBConfig() DatabaseConfig {
    return DatabaseConfig{
        Host:     AppConfig.DBHost,
        Port:     AppConfig.DBPort,
        User:     AppConfig.DBUser,
        Password: AppConfig.DBPassword,
        Name:     AppConfig.DBName,
        SSLMode:  AppConfig.DBSSLMode,
    }
}

// ConnectDatabase initializes database connection
func ConnectDatabase() error {
    dbConfig := GetDBConfig()
    
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        dbConfig.Host,
        dbConfig.Port,
        dbConfig.User,
        dbConfig.Password,
        dbConfig.Name,
        dbConfig.SSLMode,
    )

    // Configure GORM
    gormConfig := &gorm.Config{}
    
    if AppConfig.AppEnv == "development" {
        gormConfig.Logger = logger.Default.LogMode(logger.Info)
    }

    // Open connection
    database, err := gorm.Open(postgres.Open(dsn), gormConfig)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Get generic database object sql.DB to use its functions
    sqlDB, err := database.DB()
    if err != nil {
        return fmt.Errorf("failed to get generic database object: %w", err)
    }

    // Set connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    DB = database
    log.Println("Database connected successfully")
    
    return nil
}

// CloseDatabase closes database connection
func CloseDatabase() error {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err != nil {
            return err
        }
        return sqlDB.Close()
    }
    return nil
}