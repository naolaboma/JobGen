package controllers

import (
	"net/http"
	"strconv"

	domain "jobgen-backend/Domain"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatUsecase domain.IChatUsecase
}

func NewChatController(chatUsecase domain.IChatUsecase) *ChatController {
	return &ChatController{chatUsecase: chatUsecase}
}

// @Summary Send a message to the AI chatbot
// @Description Send a message to the AI career assistant and receive a response
// @Tags AI Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.ChatRequest true "Chat message"
// @Success 200 {object} StandardResponse "Message sent successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /chat/message [post]
func (c *ChatController) SendMessage(ctx *gin.Context) {
	var req domain.ChatRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}
	
	// Get user ID from authentication context
	userID := ctx.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(ctx, "User not authenticated")
		return
	}
	req.UserID = userID
	
	response, err := c.chatUsecase.SendMessage(ctx, &req)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to process message: "+err.Error())
		return
	}
	
	SuccessResponse(ctx, http.StatusOK, "Message processed successfully", response)
}

// @Summary Get chat session history
// @Description Retrieve the message history for a specific chat session
// @Tags AI Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path string true "Session ID"
// @Success 200 {object} StandardResponse "History retrieved successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 404 {object} StandardResponse "Session not found"
// @Router /chat/session/{session_id} [get]
func (c *ChatController) GetSessionHistory(ctx *gin.Context) {
	sessionID := ctx.Param("session_id")
	userID := ctx.GetString("user_id")
	
	if userID == "" {
		UnauthorizedResponse(ctx, "User not authenticated")
		return
	}
	
	history, err := c.chatUsecase.GetSessionHistory(ctx, sessionID, userID)
	if err != nil {
		NotFoundResponse(ctx, "Session not found")
		return
	}
	
	SuccessResponse(ctx, http.StatusOK, "History retrieved successfully", history)
}

// @Summary Get user's chat sessions
// @Description Retrieve a list of the user's chat sessions
// @Tags AI Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of results to return (default 10)"
// @Param offset query int false "Number of results to skip (default 0)"
// @Success 200 {object} StandardResponse "Sessions retrieved successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Router /chat/sessions [get]
func (c *ChatController) GetUserSessions(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	
	if userID == "" {
		UnauthorizedResponse(ctx, "User not authenticated")
		return
	}
	
	limit := parseIntQueryParam(ctx, "limit", 10)
	offset := parseIntQueryParam(ctx, "offset", 0)
	
	sessions, err := c.chatUsecase.GetUserSessions(ctx, userID, limit, offset)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to retrieve sessions: "+err.Error())
		return
	}
	
	SuccessResponse(ctx, http.StatusOK, "Sessions retrieved successfully", sessions)
}

// @Summary Delete a chat session
// @Description Delete a specific chat session and all its messages
// @Tags AI Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path string true "Session ID"
// @Success 200 {object} StandardResponse "Session deleted successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 404 {object} StandardResponse "Session not found"
// @Router /chat/session/{session_id} [delete]
func (c *ChatController) DeleteSession(ctx *gin.Context) {
	sessionID := ctx.Param("session_id")
	userID := ctx.GetString("user_id")
	
	if userID == "" {
		UnauthorizedResponse(ctx, "User not authenticated")
		return
	}
	
	err := c.chatUsecase.DeleteSession(ctx, sessionID, userID)
	if err != nil {
		NotFoundResponse(ctx, "Session not found")
		return
	}
	
	SuccessResponse(ctx, http.StatusOK, "Session deleted successfully", nil)
}

// Helper function to parse integer query parameters
func parseIntQueryParam(ctx *gin.Context, param string, defaultValue int) int {
	value := ctx.Query(param)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}