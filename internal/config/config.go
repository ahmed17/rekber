package config

import (
    "log"
    "os"
    "strconv"
    
    "github.com/joho/godotenv"
)

type Config struct {
    AppName      string
    AppEnv       string
    AppPort      string
    AppURL       string
    
    DBHost       string
    DBPort       string
    DBUser       string
    DBPassword   string
    DBName       string
    DBSSLMode    string
    
    JWTSecret    string
    JWTExpire    int
    
    AutoCompleteHours int
    MaxUploadSize     int64
    UploadPath        string
}

var AppConfig *Config

func LoadConfig() error {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    AppConfig = &Config{
        AppName:      getEnv("APP_NAME", "Rekber"),
        AppEnv:       getEnv("APP_ENV", "development"),
        AppPort:      getEnv("APP_PORT", "8080"),
        AppURL:       getEnv("APP_URL", "http://localhost:8080"),
        
        DBHost:       getEnv("DB_HOST", "localhost"),
        DBPort:       getEnv("DB_PORT", "5432"),
        DBUser:       getEnv("DB_USER", "postgres"),
        DBPassword:   getEnv("DB_PASSWORD", ""),
        DBName:       getEnv("DB_NAME", "rekber_db"),
        DBSSLMode:    getEnv("DB_SSLMODE", "disable"),
        
        JWTSecret:    getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
        JWTExpire:    getEnvAsInt("JWT_EXPIRE_HOURS", 72),
        
        AutoCompleteHours: getEnvAsInt("AUTO_COMPLETE_HOURS", 48),
        MaxUploadSize:     getEnvAsInt64("MAX_UPLOAD_SIZE", 5*1024*1024), // 5MB
        UploadPath:        getEnv("UPLOAD_PATH", "./uploads"),
    }
    
    return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
        return value
    }
    return defaultValue
}