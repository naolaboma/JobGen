package controllers

import (
	domain "jobgen-backend/Domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ContactController struct {
	contactUsecase domain.IContactUsecase
}
// todo what if a single user make contact request multiple time so consider DDOS attack and implement it a correct manner for next time
func NewContactController(contactUsecase domain.IContactUsecase) *ContactController {
	return &ContactController{
		contactUsecase: contactUsecase,
	}
}

// @Summary Submit a contact form
// @Description Allows users to submit general inquiries or feedback
// @Tags Public
// @Accept json
// @Produce json
// @Param request body domain.ContactFormRequest true "Contact form details"
// @Success 200 {object} StandardResponse "Contact form submitted successfully"
// @Failure 400 {object} StandardResponse "Bad request (validation error)"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /contact [post]
func (c *ContactController) SubmitContactForm(ctx *gin.Context) {
	var req domain.ContactFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	contact := &domain.Contact{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
		CreatedAt: time.Now(),
	}

	if err := c.contactUsecase.SubmitContactForm(ctx, contact); err != nil {
		InternalErrorResponse(ctx, "Failed to submit contact form")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Your message has been received. We will get back to you shortly.", nil)
}
