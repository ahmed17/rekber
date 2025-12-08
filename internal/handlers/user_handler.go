package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rekber/internal/handlers/middleware"
    "rekber/internal/services"
)

type UserHandler struct {
    service *services.Service
}

func NewUserHandler(service *services.Service) *UserHandler {
    return &UserHandler{service: service}
}

// ListUsers handler (admin only)
func (h *UserHandler) ListUsers(c *gin.Context) {
    // Get pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    users, total, err := h.service.User.ListUsers(page, limit)
    if err != nil {
        ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessPaginationResponse(c, http.StatusOK, users, page, limit, total)
}

// GetUserByID handler (admin only)
func (h *UserHandler) GetUserByID(c *gin.Context) {
    userIDStr := c.Param("id")
    userID, err := strconv.ParseUint(userIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
        return
    }

    user, err := h.service.User.GetUserByID(uint(userID))
    if err != nil {
        ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, user, "User retrieved successfully")
}

// SetupRoutes mengatur routing untuk user handler
func (h *UserHandler) SetupRoutes(router *gin.RouterGroup) {
    users := router.Group("/users")
    users.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
    {
        users.GET("/", h.ListUsers)
        users.GET("/:id", h.GetUserByID)
    }
}