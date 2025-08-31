package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse represents the standardized API response format
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo provides detailed error information
type ErrorInfo struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginatedResponse represents paginated data response
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// Response helper functions
func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	response := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	ctx.JSON(statusCode, response)
}

func ErrorResponse(ctx *gin.Context, statusCode int, code, message string, details interface{}) {
	response := StandardResponse{
		Success: false,
		Message: "Request failed",
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	ctx.JSON(statusCode, response)
}

func ValidationErrorResponse(ctx *gin.Context, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input data", err.Error())
}

func UnauthorizedResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func ForbiddenResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func NotFoundResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func ConflictResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusConflict, "CONFLICT", message, nil)
}

func InternalErrorResponse(ctx *gin.Context, message string) {
	ErrorResponse(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

func PaginatedSuccessResponse(ctx *gin.Context, statusCode int, message string, paginatedData *PaginatedResponse) {
	response := StandardResponse{
		Success: true,
		Message: message,
		Data:    paginatedData,
	}
	ctx.JSON(statusCode, response)
}
