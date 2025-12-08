package handlers

import (
    "github.com/gin-gonic/gin"
)

// Response structure untuk standarisasi response
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
    Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponse mengembalikan response sukses
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message ...string) {
    resp := Response{
        Success: true,
        Data:    data,
    }
    
    if len(message) > 0 {
        resp.Message = message[0]
    }
    
    c.JSON(statusCode, resp)
}

// ErrorResponse mengembalikan response error
func ErrorResponse(c *gin.Context, statusCode int, errorMessage string) {
    c.JSON(statusCode, Response{
        Success: false,
        Error:   errorMessage,
    })
}

// PaginationResponse untuk response dengan pagination
type PaginationResponse struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    Total      int64       `json:"total"`
    TotalPages int         `json:"total_pages"`
}

// SuccessPaginationResponse mengembalikan response dengan pagination
func SuccessPaginationResponse(c *gin.Context, statusCode int, data interface{}, page, limit int, total int64) {
    totalPages := 0
    if limit > 0 {
        totalPages = int((total + int64(limit) - 1) / int64(limit))
    }
    
    resp := PaginationResponse{
        Data:       data,
        Page:       page,
        Limit:      limit,
        Total:      total,
        TotalPages: totalPages,
    }
    
    c.JSON(statusCode, Response{
        Success: true,
        Data:    resp,
    })
}