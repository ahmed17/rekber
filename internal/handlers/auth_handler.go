package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "rekber/internal/handlers/middleware"
    "rekber/internal/services"
)

type AuthHandler struct {
    service *services.Service
}

func NewAuthHandler(service *services.Service) *AuthHandler {
    return &AuthHandler{service: service}
}

// RegisterRequest untuk validasi input register
type RegisterRequest struct {
    Username    string `json:"username" binding:"required,min=3,max=50"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=6"`
    FullName    string `json:"full_name" binding:"required,min=2,max=100"`
    Phone       string `json:"phone" binding:"omitempty"`
    Role        string `json:"role" binding:"omitempty,oneof=buyer seller"`
    
    BankName      string `json:"bank_name" binding:"omitempty"`
    AccountNumber string `json:"account_number" binding:"omitempty"`
    AccountHolder string `json:"account_holder" binding:"omitempty"`
}

// LoginRequest untuk validasi input login
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// Register handler
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    // Convert to service request
    serviceReq := services.CreateUserRequest{
        Username:     req.Username,
        Email:        req.Email,
        Password:     req.Password,
        FullName:     req.FullName,
        Phone:        req.Phone,
        Role:         req.Role,
        BankName:     req.BankName,
        AccountNumber: req.AccountNumber,
        AccountHolder: req.AccountHolder,
    }

    user, err := h.service.User.Register(serviceReq)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusCreated, user, "User registered successfully")
}

// Login handler
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    // Convert to service request
    serviceReq := services.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    }

    token, user, err := h.service.User.Login(serviceReq)
    if err != nil {
        ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    response := gin.H{
        "token": token,
        "user":  user,
    }

    SuccessResponse(c, http.StatusOK, response, "Login successful")
}

// Profile handler untuk mendapatkan profile user saat ini
func (h *AuthHandler) Profile(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    user, err := h.service.User.GetProfile(userID)
    if err != nil {
        ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, user, "Profile retrieved successfully")
}

// UpdateProfile handler untuk update profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    var req services.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    user, err := h.service.User.UpdateProfile(userID, req)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, user, "Profile updated successfully")
}

// ChangePassword handler untuk ganti password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    var req struct {
        OldPassword string `json:"old_password" binding:"required"`
        NewPassword string `json:"new_password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    err := h.service.User.ChangePassword(userID, req.OldPassword, req.NewPassword)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, nil, "Password changed successfully")
}

// SetupRoutes mengatur routing untuk auth handler
func (h *AuthHandler) SetupRoutes(router *gin.RouterGroup) {
    auth := router.Group("/auth")
    {
        // Public routes
        auth.POST("/register", h.Register)
        auth.POST("/login", h.Login)

        // Protected routes
        auth.GET("/profile", middleware.AuthMiddleware(), h.Profile)
        auth.PUT("/profile", middleware.AuthMiddleware(), h.UpdateProfile)
        auth.PUT("/change-password", middleware.AuthMiddleware(), h.ChangePassword)
    }
}