package handlers

import (
    "github.com/gin-gonic/gin"
    "rekber/internal/services"
)

// Handler container untuk semua handlers
type Handler struct {
    auth        *AuthHandler
    user        *UserHandler
    transaction *TransactionHandler
    testimony   *TestimonyHandler
}

// NewHandler membuat instance baru Handler
func NewHandler(service *services.Service) *Handler {
    return &Handler{
        auth:        NewAuthHandler(service),
        user:        NewUserHandler(service),
        transaction: NewTransactionHandler(service),
        testimony:   NewTestimonyHandler(service),
    }
}

// SetupRoutes mengatur semua routes
func (h *Handler) SetupRoutes(router *gin.Engine) {
    // API v1 routes
    api := router.Group("/api/v1")
    {
        // Setup masing-masing handler routes
        h.auth.SetupRoutes(api)
        h.user.SetupRoutes(api)
        h.transaction.SetupRoutes(api)
        h.testimony.SetupRoutes(api)
        
        // Health check
        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "status":  "healthy",
                "version": "1.0",
            })
        })
    }
}