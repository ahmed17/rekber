package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "rekber/internal/handlers/middleware"
    "rekber/internal/services"
)

type TransactionHandler struct {
    service *services.Service
}

func NewTransactionHandler(service *services.Service) *TransactionHandler {
    return &TransactionHandler{service: service}
}

// CreateTransaction handler
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    var req services.CreateTransactionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    transaction, err := h.service.Transaction.CreateTransaction(req, userID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusCreated, transaction, "Transaction created successfully")
}

// GetTransactionByID handler
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    role, _ := c.Get("user_role")
    userRole := role.(string)

    // Get transaction ID from URL parameter
    transactionIDStr := c.Param("id")
    transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    transaction, err := h.service.Transaction.GetTransactionByID(uint(transactionID), userID, userRole)
    if err != nil {
        ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, transaction, "Transaction retrieved successfully")
}

// GetTransactionByCode handler
func (h *TransactionHandler) GetTransactionByCode(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    role, _ := c.Get("user_role")
    userRole := role.(string)

    code := c.Param("code")
    transaction, err := h.service.Transaction.GetTransactionByCode(code, userID, userRole)
    if err != nil {
        ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, transaction, "Transaction retrieved successfully")
}

// ListUserTransactions handler
func (h *TransactionHandler) ListUserTransactions(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    role, _ := c.Get("user_role")
    userRole := role.(string)

    // Get pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    transactions, total, err := h.service.Transaction.ListUserTransactions(userID, userRole, page, limit)
    if err != nil {
        ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessPaginationResponse(c, http.StatusOK, transactions, page, limit, total)
}

// ListAllTransactions handler (admin only)
func (h *TransactionHandler) ListAllTransactions(c *gin.Context) {
    // Get pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    transactions, total, err := h.service.Transaction.ListAllTransactions(page, limit)
    if err != nil {
        ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    SuccessPaginationResponse(c, http.StatusOK, transactions, page, limit, total)
}

// UploadPaymentProof handler
func (h *TransactionHandler) UploadPaymentProof(c *gin.Context) {
    userID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    transactionIDStr := c.Param("id")
    transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    var req services.UploadPaymentProofRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    transaction, err := h.service.Transaction.UploadPaymentProof(uint(transactionID), req, userID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, transaction, "Payment proof uploaded successfully")
}

// UpdateShippingData handler (seller only)
func (h *TransactionHandler) UpdateShippingData(c *gin.Context) {
    sellerID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    transactionIDStr := c.Param("id")
    transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    var req services.UpdateShippingDataRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
        return
    }

    transaction, err := h.service.Transaction.UpdateShippingData(uint(transactionID), req, sellerID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, transaction, "Shipping data updated successfully")
}

// CompleteTransaction handler (buyer only)
func (h *TransactionHandler) CompleteTransaction(c *gin.Context) {
    buyerID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    transactionIDStr := c.Param("id")
    transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    transaction, err := h.service.Transaction.CompleteTransaction(uint(transactionID), buyerID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, transaction, "Transaction completed successfully")
}

// CancelTransaction handler (buyer only)
func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
    buyerID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    transactionIDStr := c.Param("id")
    transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    transaction, err := h.service.Transaction.CancelTransaction(uint(transactionID), buyerID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, transaction, "Transaction cancelled successfully")
}

// ValidatePayment handler (admin only)
func (h *TransactionHandler) ValidatePayment(c *gin.Context) {
    adminID, ok := middleware.GetUserIDFromContext(c)
    if !ok {
        ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    transactionIDStr := c.Param("id")
    transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
        return
    }

    transaction, err := h.service.Transaction.ValidatePayment(uint(transactionID), adminID)
    if err != nil {
        ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    SuccessResponse(c, http.StatusOK, transaction, "Payment validated successfully")
}

// SetupRoutes mengatur routing untuk transaction handler
func (h *TransactionHandler) SetupRoutes(router *gin.RouterGroup) {
    transactions := router.Group("/transactions")
    transactions.Use(middleware.AuthMiddleware())
    {
        // Buyer routes
        transactions.POST("/", h.CreateTransaction)
        transactions.GET("/my", h.ListUserTransactions)
        transactions.PUT("/:id/upload-payment", h.UploadPaymentProof)
        transactions.PUT("/:id/complete", h.CompleteTransaction)
        transactions.PUT("/:id/cancel", h.CancelTransaction)
        
        // Seller routes
        transactions.PUT("/:id/shipping-data", middleware.SellerMiddleware(), h.UpdateShippingData)
        
        // Admin routes
        transactions.GET("/", middleware.AdminMiddleware(), h.ListAllTransactions)
        transactions.PUT("/:id/validate-payment", middleware.AdminMiddleware(), h.ValidatePayment)
        
        // Common routes
        transactions.GET("/:id", h.GetTransactionByID)
        transactions.GET("/code/:code", h.GetTransactionByCode)
    }
}