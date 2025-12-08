package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rekber/internal/handlers/middleware"
    "rekber/internal/services"
)

type TestimonyHandler struct {
    service *services.Service
}

func NewTestimonyHandler(service *services.Service) *TestimonyHandler {
    return &TestimonyHandler{service: service}
}

// CreateTestimony handler
func (h *TestimonyHandler) CreateTestimony(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    var req services.CreateTestimonyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    testimony, err := h.service.Testimony.CreateTestimony(req, userID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusCreated, testimony, "Testimony created successfully")
}

// GetTestimonyByID handler
func (h *TestimonyHandler) GetTestimonyByID(c *gin.Context) {
    testimonyIDStr := c.Param("id")
    testimonyID, err := strconv.ParseUint(testimonyIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid testimony ID")
        return
    }

    testimony, err := h.service.Testimony.GetTestimonyByID(uint(testimonyID))
    if err != nil {
        ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, testimony, "Testimony retrieved successfully")
}

// UpdateTestimony handler
func (h *TestimonyHandler) UpdateTestimony(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    testimonyIDStr := c.Param("id")
    testimonyID, err := strconv.ParseUint(testimonyIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid testimony ID")
        return
    }

    var req services.UpdateTestimonyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    testimony, err := h.service.Testimony.UpdateTestimony(uint(testimonyID), req, userID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, testimony, "Testimony updated successfully")
}

// DeleteTestimony handler
func (h *TestimonyHandler) DeleteTestimony(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    testimonyIDStr := c.Param("id")
    testimonyID, err := strconv.ParseUint(testimonyIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid testimony ID")
        return
    }

    err = h.service.Testimony.DeleteTestimony(uint(testimonyID), userID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, nil, "Testimony deleted successfully")
}

// ListTestimonies handler
func (h *TestimonyHandler) ListTestimonies(c *gin.Context) {
    // Get pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    testimonies, total, err := h.service.Testimony.ListTestimonies(page, limit)
    if err != nil {
        ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessPaginationResponse(c, http.StatusOK, testimonies, page, limit, total)
}

// ListUserTestimonies handler
func (h *TestimonyHandler) ListUserTestimonies(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    // Get pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    testimonies, total, err := h.service.Testimony.ListUserTestimonies(userID, page, limit)
    if err != nil {
        ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessPaginationResponse(c, http.StatusOK, testimonies, page, limit, total)
}

// SetupRoutes mengatur routing untuk testimony handler
func (h *TestimonyHandler) SetupRoutes(router *gin.RouterGroup) {
    testimonies := router.Group("/testimonies")
    {
        // Public routes
        testimonies.GET("/", h.ListTestimonies)
        testimonies.GET("/:id", h.GetTestimonyByID)
        
        // Protected routes
        protected := testimonies.Group("")
        protected.Use(middleware.AuthMiddleware())
        {
            protected.POST("/", h.CreateTestimony)
            protected.PUT("/:id", h.UpdateTestimony)
            protected.DELETE("/:id", h.DeleteTestimony)
            protected.GET("/my", h.ListUserTestimonies)
        }
    }
}