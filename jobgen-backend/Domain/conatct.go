package domain

import (
	"context"
	"time"
)

// Contact represents a contact form submission
type Contact struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Subject   string    `json:"subject" bson:"subject"`
	Message   string    `json:"message" bson:"message"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Replied   bool      `json:"replied" bson:"replied"` // To track if a response has been sent
}

// IContactRepository provides methods for interacting with Contact storage
type IContactRepository interface {
	Create(ctx context.Context, contact *Contact) error
}

// IContactUsecase provides methods for Contact business logic
type IContactUsecase interface {
	SubmitContactForm(ctx context.Context, contact *Contact) error
}

// ContactFormRequest is a DTO for incoming contact form data
type ContactFormRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=100"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required,min=3,max=200"`
	Message string `json:"message" binding:"required,min=10,max=2000"`
}
