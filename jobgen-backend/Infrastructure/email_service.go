package infrastructure

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	from     string
	host     string
	port     int
	username string
	password string
	dialer   *gomail.Dialer
}

func NewEmailService() domain.IEmailService {
	port, _ := strconv.Atoi(Env.EmailPort)
    dialer := gomail.NewDialer(Env.EmailHost, port, Env.EmailUsername, Env.EmailPassword)
    dialer.SSL = false // Use StartTLS (default behavior)

	return &EmailService{
		from:     Env.EmailFrom,
		host:     Env.EmailHost,
		port:     port,
		username: Env.EmailUsername,
		password: Env.EmailPassword,
		dialer:   dialer,
	}
}

func (e *EmailService) SendWelcomeEmail(ctx context.Context, user *domain.User, otp string) error {
	subject := "Welcome to JobGen - Verify Your Email"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .otp-code { font-size: 32px; font-weight: bold; text-align: center; background: white; padding: 20px; margin: 20px 0; border: 2px dashed #667eea; letter-spacing: 5px; border-radius: 8px; }
        .button { display: inline-block; padding: 12px 24px; background: #667eea; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to JobGen!</h1>
            <p>Your AI-Powered Career Assistant</p>
        </div>
        <div class="content">
            <h2>Hi %s!</h2>
            <p>Thank you for joining JobGen. We're excited to help you find your dream remote job!</p>
            
            <p>To get started, please verify your email address using the code below:</p>
            
            <div class="otp-code">%s</div>
            
            <p><strong>This code expires in 15 minutes.</strong></p>
            
            <p>Once verified, you'll be able to:</p>
            <ul>
                <li>Upload and optimize your CV with AI</li>
                <li>Get personalized job recommendations</li>
                <li>Chat with our AI career assistant</li>
                <li>Apply to remote jobs worldwide</li>
            </ul>
            
            <p>If you didn't create this account, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>This is an automated email from JobGen. Please do not reply.</p>
            <p>&copy; 2024 JobGen. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, user.FullName, otp)

	return e.sendEmail(user.Email, subject, body)
}

func (e *EmailService) SendPasswordResetEmail(ctx context.Context, user *domain.User, resetToken string) error {
	subject := "Reset Your JobGen Password"
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", Env.FrontendURL, resetToken)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; padding: 12px 24px; background: #667eea; color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .warning { background: #fff3cd; border: 1px solid #ffeeba; padding: 15px; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>We received a request to reset your JobGen account password.</p>
            
            <p>Click the button below to reset your password:</p>
            <p style="text-align: center;">
                <a href="%s" class="button">Reset Password</a>
            </p>
            
            <div class="warning">
                <strong>Security Notice:</strong>
                <ul>
                    <li>This link expires in 1 hour</li>
                    <li>If you didn't request this reset, please ignore this email</li>
                    <li>For security, this link can only be used once</li>
                </ul>
            </div>
            
            <p>If the button doesn't work, copy and paste this link:</p>
            <p><a href="%s">%s</a></p>
        </div>
        <div class="footer">
            <p>This is an automated email from JobGen. Please do not reply.</p>
            <p>&copy; 2024 JobGen. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, user.FullName, resetLink, resetLink, resetLink)

	return e.sendEmail(user.Email, subject, body)
}

func (e *EmailService) SendAccountDeactivationEmail(ctx context.Context, user *domain.User) error {
	subject := "JobGen Account Deactivated"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #dc3545; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Account Deactivated</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>Your JobGen account has been deactivated by an administrator.</p>
            <p>If you believe this is an error, please contact our support team.</p>
        </div>
        <div class="footer">
            <p>This is an automated email from JobGen. Please do not reply.</p>
        </div>
    </div>
</body>
</html>`, user.FullName)

	return e.sendEmail(user.Email, subject, body)
}

func (e *EmailService) SendRoleChangeNotification(ctx context.Context, user *domain.User, newRole domain.Role) error {
	subject := "JobGen Account Role Updated"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #28a745; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Role Updated</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>Your JobGen account role has been updated to: <strong>%s</strong></p>
            <p>This change takes effect immediately. Please log out and log back in to see the changes.</p>
        </div>
        <div class="footer">
            <p>This is an automated email from JobGen. Please do not reply.</p>
        </div>
    </div>
</body>
</html>`, user.FullName, string(newRole))

	return e.sendEmail(user.Email, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return e.dialer.DialAndSend(m)
}
