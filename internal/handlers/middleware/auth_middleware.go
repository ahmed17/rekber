package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "rekber/internal/services"
)

// AuthMiddleware memvalidasi JWT token dari Authorization header
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header is required",
            })
            c.Abort()
            return
        }

        // Format: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization format. Expected: Bearer <token>",
            })
            c.Abort()
            return
        }

        token := parts[1]
        
        // Get service from context atau buat baru
        svc := services.GetService()
        if svc == nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Service not available",
            })
            c.Abort()
            return
        }

        userID, role, err := svc.Auth.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        // Set user info ke context
        c.Set("user_id", userID)
        c.Set("user_role", role)
        c.Next()
    }
}

// AdminMiddleware memastikan user adalah admin
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("user_role")
        if !exists || role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Admin access required",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// SellerMiddleware memastikan user adalah seller
func SellerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("user_role")
        if !exists || role != "seller" {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Seller access required",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// BuyerMiddleware memastikan user adalah buyer
func BuyerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("user_role")
        if !exists || role != "buyer" {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Buyer access required",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}

// GetUserIDFromContext mendapatkan user ID dari context
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, false
    }
    
    // Convert ke uint
    switch v := userID.(type) {
    case uint:
        return v, true
    case int:
        return uint(v), true
    case int64:
        return uint(v), true
    case float64:
        return uint(v), true
    default:
        return 0, false
    }
}