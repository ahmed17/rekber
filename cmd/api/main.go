package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "rekber/internal/config"
    // "rekber/internal/models"
    "rekber/pkg/database"
)

func main() {
    // Load configuration
    if err := config.LoadConfig(); err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Connect to database
    if err := config.ConnectDatabase(); err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer config.CloseDatabase()

    // Run migrations
    if err := runMigrations(); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Auto migrate GORM models untuk development
    // if config.AppConfig.AppEnv == "development" {
    //     if err := models.AutoMigrate(); err != nil {
    //         log.Printf("Warning: Auto migration failed: %v", err)
    //     }
    // }

    // Set Gin mode
    if config.AppConfig.AppEnv == "production" {
        gin.SetMode(gin.ReleaseMode)
    } else {
        gin.SetMode(gin.DebugMode)
    }

    // Create router
    router := gin.Default()

    // Add middlewares
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    // Basic routes
    router.GET("/", healthCheck)
    router.GET("/health", healthCheck)

    // API v1 routes will be added here later
    setupRoutes(router)

    // Start server with graceful shutdown
    startServer(router)
}

func healthCheck(c *gin.Context) {
    // Check database connection
    db, err := config.DB.DB()
    dbStatus := "connected"
    
    if err != nil {
        dbStatus = "disconnected"
    } else {
        if err := db.Ping(); err != nil {
            dbStatus = "disconnected"
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "service":   "rekber-api",
        "version":   "1.0.0",
        "database":  dbStatus,
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

func runMigrations() error {
    log.Println("Running database migrations...")
    
    // Get generic database object from GORM
    sqlDB, err := config.DB.DB()
    if err != nil {
        return err
    }

    // Run migrations
    if err := database.RunMigrations(sqlDB, config.AppConfig.DBName); err != nil {
        return err
    }

    return nil
}

func setupRoutes(router *gin.Engine) {
    // API v1 routes will be defined here
    api := router.Group("/api/v1")
    {
        api.GET("/status", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "Rekber API v1 is running",
            })
        })
    }
}

func startServer(router *gin.Engine) {
    port := ":" + config.AppConfig.AppPort
    
    srv := &http.Server{
        Addr:    port,
        Handler: router,
    }

    // Graceful shutdown
    go func() {
        log.Printf("Server starting on port %s in %s mode", port, config.AppConfig.AppEnv)
        
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited properly")
}