package main

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func main() {
    // Set Gin mode
    gin.SetMode(gin.ReleaseMode)
    
    // Create router
    r := gin.Default()
    
    // Basic route untuk testing
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Rekber API is running!",
            "status":  "success",
        })
    })
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "service": "rekber-api",
        })
    })
    
    // Start server
    port := ":8080"
    log.Printf("Server starting on port %s", port)
    if err := r.Run(port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}