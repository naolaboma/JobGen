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
	adminEmail string // New field for admin email
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
		adminEmail: Env.AdminEmail, // Initialize admin email
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
	subject := "Your JobGen Password Reset Code"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; text-align: center; }
        .otp-box { display: inline-block; font-size: 32px; font-weight: bold; letter-spacing: 5px; padding: 15px 25px; background: #fff; border: 2px dashed #667eea; border-radius: 8px; margin: 20px 0; }
        .warning { background: #fff3cd; border: 1px solid #ffeeba; padding: 15px; border-radius: 5px; margin: 20px 0; text-align: left; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Code</h1>
        </div>
        <div class="content">
            <h2>Hi %s,</h2>
            <p>We received a request to reset your JobGen account password.</p>
            <p>Use the code below to reset your password:</p>

            <div class="otp-box">%s</div>

            <div class="warning">
                <strong>Security Notice:</strong>
                <ul>
                    <li>This code expires in 15 minutes</li>
                    <li>If you didn't request this reset, please ignore this email</li>
                    <li>For security, this code can only be used once</li>
                </ul>
            </div>
        </div>
        <div class="footer">
            <p>This is an automated email from JobGen. Please do not reply.</p>
            <p>&copy; 2024 JobGen. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, user.FullName, resetToken)

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
func (e *EmailService) SendContactFormToAdmin(ctx context.Context, contact *domain.Contact) error {
	if e.adminEmail == "" {
		return fmt.Errorf("admin email is not configured")
	}

	subject := fmt.Sprintf("JobGen Contact Form: %s (From: %s)", contact.Subject, contact.Name)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>New Contact Form Submission</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f4f6f8;
            color: #333;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 650px;
            margin: 40px auto;
            background-color: #ffffff;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 4px 15px rgba(0,0,0,0.1);
            border-top: 6px solid #667eea;
        }
        .header {
            background-color: #667eea;
            color: #ffffff;
            text-align: center;
            padding: 25px 20px;
        }
        .header h2 {
            margin: 0;
        }
        .content {
            padding: 25px 30px;
        }
        .field {
            margin-bottom: 15px;
        }
        .field strong {
            display: inline-block;
            width: 90px;
            color: #555;
        }
        .message-box {
            background: #f1f5fb;
            border: 1px solid #cbd5e1;
            box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
            padding: 18px 22px;
            border-radius: 8px;
            margin-top: 20px;
            white-space: pre-wrap;
            word-wrap: break-word;
            font-size: 14px;
            line-height: 1.6;
            color: #1f2937;
        }
        .button {
            display: inline-block;
            background-color: #667eea;
            color: #ffffff;
            text-decoration: none;
            padding: 12px 25px;
            border-radius: 5px;
            margin-top: 20px;
            transition: background 0.3s ease;
        }
        .button:hover {
            background-color: #5562c1;
        }
        .footer {
            text-align: center;
            padding: 15px;
            font-size: 12px;
            color: #999;
            border-top: 1px solid #eee;
        }
        @media screen and (max-width: 680px) {
            .container { margin: 20px; }
            .content { padding: 20px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>New Contact Form Submission</h2>
        </div>
        <div class="content">
            <div class="field"><strong>Name:</strong> %s</div>
            <div class="field"><strong>Email:</strong> %s</div>
            <div class="field"><strong>Subject:</strong> %s</div>
            <div class="message-box">
                %s
            </div>
            <p style="margin-top: 20px; font-size: 14px;">Received on: %s</p>
            <a href="mailto:%s" class="button">Reply to User</a>
        </div>
        <div class="footer">
            <p>This is an automated notification from JobGen contact form.</p>
        </div>
    </div>
</body>
</html>`,
		contact.Name,
		contact.Email,
		contact.Subject,
		contact.Message,
		contact.CreatedAt.Format("2006-01-02 15:04:05 MST"),
		contact.Email,
	)

	return e.sendEmail(e.adminEmail, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return e.dialer.DialAndSend(m)
}
