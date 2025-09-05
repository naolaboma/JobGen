package usecases

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"time"
)

type contactUsecase struct {
	contactRepo    domain.IContactRepository
	emailService   domain.IEmailService
	contextTimeout time.Duration
}

func NewContactUsecase(
	contactRepo domain.IContactRepository,
	emailService domain.IEmailService,
	timeout time.Duration,
) domain.IContactUsecase {
	return &contactUsecase{
		contactRepo:    contactRepo,
		emailService:   emailService,
		contextTimeout: timeout,
	}
}

func (u *contactUsecase) SubmitContactForm(ctx context.Context, contact *domain.Contact) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Store the contact message in the database
	if err := u.contactRepo.Create(ctx, contact); err != nil {
		return fmt.Errorf("failed to store contact message: %w", err)
	}

	// Send the contact message to the admin email
	// This is a fire-and-forget operation; we log errors but don't fail the request
	// if email sending fails after the message is stored.
	if err := u.emailService.SendContactFormToAdmin(ctx, contact); err != nil {
		fmt.Printf("Warning: Failed to send contact form email to admin: %v\n", err)
	}

	return nil
}
