package tests

import (
	"net/http"

	controllers "jobgen-backend/Delivery/Controllers"
	domain "jobgen-backend/Domain"

	"github.com/stretchr/testify/mock"
)

func (suite *APITestSuite) TestRegister() {
	suite.Run("successful registration", func() {
		reqBody := controllers.RegisterRequest{
			Email:           "newuser@example.com",
			Username:        "newuser",
			Password:        "Password123!",
			FullName:        "New User",
			PhoneNumber:     "+1234567890",
			Location:        "New York",
			Skills:          []string{"Go", "JavaScript"},
			ExperienceYears: 2,
			Bio:             "Software Developer",
		}

		suite.userUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

		w := suite.makeRequest("POST", "/api/v1/auth/register", reqBody, nil)

		suite.Equal(http.StatusCreated, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "registered successfully")
	})

	suite.Run("registration with missing full name", func() {
		reqBody := controllers.RegisterRequest{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "Password123!",
			FullName: "", // Missing full name
		}

		w := suite.makeRequest("POST", "/api/v1/auth/register", reqBody, nil)

		suite.Equal(http.StatusBadRequest, w.Code)
		response := suite.parseResponse(w)
		suite.False(response.Success)
		suite.Equal("VALIDATION_ERROR", response.Error.Code)
	})

	// todo check the tests
	// suite.Run("registration with duplicate email", func() {
	// 	reqBody := controllers.RegisterRequest{
	// 		Email:    "existing@example.com",
	// 		Username: "newuser",
	// 		Password: "Password123!",
	// 		FullName: "New User",
	// 	}

	// 	suite.userUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(domain.ErrEmailTaken)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/register", reqBody, nil)

	// 	suite.Equal(http.StatusConflict, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "Email already registered")
	// })

	// suite.Run("registration with duplicate username", func() {
	// 	reqBody := controllers.RegisterRequest{
	// 		Email:    "newuser@example.com",
	// 		Username: "existinguser",
	// 		Password: "Password123!",
	// 		FullName: "New User",
	// 	}

	// 	suite.userUsecase.On("Register", mock.Anything, mock.AnythingOfType("*domain.User")).Return(domain.ErrUsernameTaken)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/register", reqBody, nil)

	// 	suite.Equal(http.StatusConflict, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "Username already taken")
	// })
}

func (suite *APITestSuite) TestLogin() {
	suite.Run("successful login", func() {
		reqBody := controllers.LoginRequest{
			Email:    "test@example.com",
			Password: "Password123!",
		}

		tokenResponse := &domain.TokenResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
		}

		suite.userUsecase.On("Login", mock.Anything, "test@example.com", "Password123!").Return(tokenResponse, nil)

		w := suite.makeRequest("POST", "/api/v1/auth/login", reqBody, nil)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "Login successful")
		
		// Check if refresh token cookie is set
		cookies := w.Result().Cookies()
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "refresh_token" {
				refreshTokenCookie = cookie
				break
			}
		}
		suite.NotNil(refreshTokenCookie)
		suite.Equal("refresh-token", refreshTokenCookie.Value)
		suite.True(refreshTokenCookie.HttpOnly)
	})

	// todo check this test
	// suite.Run("login with invalid credentials", func() {
	// 	reqBody := controllers.LoginRequest{
	// 		Email:    "test@example.com",
	// 		Password: "wrongpassword",
	// 	}

	// 	suite.userUsecase.On("Login", mock.Anything, "test@example.com", "wrongpassword").Return(nil, domain.ErrInvalidCredentials)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/login", reqBody, nil)

	// 	suite.Equal(http.StatusUnauthorized, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "Invalid email or password")
	// })

	// suite.Run("login with unverified email", func() {
	// 	reqBody := controllers.LoginRequest{
	// 		Email:    "unverified@example.com",
	// 		Password: "Password123!",
	// 	}

	// 	suite.userUsecase.On("Login", mock.Anything, "unverified@example.com", "Password123!").Return(nil, domain.ErrUserNotVerified)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/login", reqBody, nil)

	// 	suite.Equal(http.StatusUnauthorized, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "verify your email")
	// })

	// suite.Run("login with deactivated account", func() {
	// 	reqBody := controllers.LoginRequest{
	// 		Email:    "deactivated@example.com",
	// 		Password: "Password123!",
	// 	}

	// 	suite.userUsecase.On("Login", mock.Anything, "deactivated@example.com", "Password123!").Return(nil, domain.ErrUserDeactivated)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/login", reqBody, nil)

	// 	suite.Equal(http.StatusUnauthorized, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "deactivated")
	// })
}

func (suite *APITestSuite) TestVerifyEmail() {
	suite.Run("successful email verification", func() {
		reqBody := controllers.VerifyEmailRequest{
			Email: "test@example.com",
			OTP:   "123456",
		}

		suite.userUsecase.On("VerifyEmail", mock.Anything, domain.VerifyEmailInput{
			Email: "test@example.com",
			OTP:   "123456",
		}).Return(nil)

		w := suite.makeRequest("POST", "/api/v1/auth/verify-email", reqBody, nil)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "verified successfully")
	})

	// todo check this test
	// suite.Run("verification with invalid OTP", func() {
	// 	reqBody := controllers.VerifyEmailRequest{
	// 		Email: "test@example.com",
	// 		OTP:   "invalid",
	// 	}

	// 	suite.userUsecase.On("VerifyEmail", mock.Anything, domain.VerifyEmailInput{
	// 		Email: "test@example.com",
	// 		OTP:   "invalid",
	// 	}).Return(domain.ErrInvalidOTP)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/verify-email", reqBody, nil)

	// 	suite.Equal(http.StatusBadRequest, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Equal("INVALID_OTP", response.Error.Code)
	// })
}

func (suite *APITestSuite) TestRefreshToken() {
	suite.Run("successful token refresh", func() {
		reqBody := controllers.RefreshTokenRequest{
			RefreshToken: "valid-refresh-token",
		}

		tokenResponse := &domain.TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
		}

		suite.authUsecase.On("RefreshToken", mock.Anything, "valid-refresh-token").Return(tokenResponse, nil)

		w := suite.makeRequest("POST", "/api/v1/auth/refresh", reqBody, nil)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "refreshed successfully")
	})

	// suite.Run("refresh with invalid token", func() {
	// 	reqBody := controllers.RefreshTokenRequest{
	// 		RefreshToken: "invalid-refresh-token",
	// 	}

	// 	suite.authUsecase.On("RefreshToken", mock.Anything, "invalid-refresh-token").Return(nil, domain.ErrInvalidToken)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/refresh", reqBody, nil)

	// 	suite.Equal(http.StatusUnauthorized, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "Invalid or expired")
	// })

	suite.Run("refresh with expired token", func() {
		reqBody := controllers.RefreshTokenRequest{
			RefreshToken: "expired-refresh-token",
		}

		suite.authUsecase.On("RefreshToken", mock.Anything, "expired-refresh-token").Return(nil, domain.ErrTokenExpired)

		w := suite.makeRequest("POST", "/api/v1/auth/refresh", reqBody, nil)

		suite.Equal(http.StatusUnauthorized, w.Code)
		response := suite.parseResponse(w)
		suite.False(response.Success)
	})
}

func (suite *APITestSuite) TestLogout() {
	suite.Run("successful logout", func() {
		suite.authUsecase.On("Logout", mock.Anything, "test-user-id").Return(nil)

		w := suite.makeAuthenticatedRequest("POST", "/api/v1/auth/logout", nil, "test-user-id", domain.RoleUser)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "Logged out successfully")

		// Check if refresh token cookie is cleared
		cookies := w.Result().Cookies()
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "refresh_token" {
				refreshTokenCookie = cookie
				break
			}
		}
		suite.NotNil(refreshTokenCookie)
		suite.Equal("", refreshTokenCookie.Value)
		suite.Equal(-1, refreshTokenCookie.MaxAge)
	})

	suite.Run("logout without authentication", func() {
		w := suite.makeRequest("POST", "/api/v1/auth/logout", nil, nil)

		suite.Equal(http.StatusUnauthorized, w.Code)
		response := suite.parseResponse(w)
		suite.False(response.Success)
		suite.Equal("MISSING_AUTH_HEADER", response.Error.Code)
	})
}

func (suite *APITestSuite) TestRequestPasswordReset() {
	suite.Run("successful password reset request", func() {
		reqBody := controllers.RequestPasswordResetRequest{
			Email: "test@example.com",
		}

		suite.userUsecase.On("RequestPasswordReset", mock.Anything, "test@example.com").Return("reset-token", nil)

		w := suite.makeRequest("POST", "/api/v1/auth/forgot-password", reqBody, nil)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "reset link has been sent")
	})

	suite.Run("password reset request with non-existent email", func() {
		reqBody := controllers.RequestPasswordResetRequest{
			Email: "nonexistent@example.com",
		}

		suite.userUsecase.On("RequestPasswordReset", mock.Anything, "nonexistent@example.com").Return("", domain.ErrUserNotFound)

		w := suite.makeRequest("POST", "/api/v1/auth/forgot-password", reqBody, nil)

		// Should still return success for security reasons
		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
	})
}

func (suite *APITestSuite) TestResetPassword() {
	suite.Run("successful password reset", func() {
		reqBody := controllers.ResetPasswordRequest{
			Token:       "valid-reset-token",
			NewPassword: "NewPassword123!",
		}

		suite.userUsecase.On("ResetPassword", mock.Anything, domain.ResetPasswordInput{
			Token:       "valid-reset-token",
			NewPassword: "NewPassword123!",
		}).Return(nil)

		w := suite.makeRequest("POST", "/api/v1/auth/reset-password", reqBody, nil)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "Password reset successful")
	})

	suite.Run("password reset with invalid token", func() {
		reqBody := controllers.ResetPasswordRequest{
			Token:       "invalid-token",
			NewPassword: "NewPassword123!",
		}

		suite.userUsecase.On("ResetPassword", mock.Anything, domain.ResetPasswordInput{
			Token:       "invalid-token",
			NewPassword: "NewPassword123!",
		}).Return(domain.ErrInvalidResetToken)

		w := suite.makeRequest("POST", "/api/v1/auth/reset-password", reqBody, nil)

		suite.Equal(http.StatusBadRequest, w.Code)
		response := suite.parseResponse(w)
		suite.False(response.Success)
		suite.Equal("INVALID_TOKEN", response.Error.Code)
	})
}

func (suite *APITestSuite) TestChangePassword() {
	suite.Run("successful password change", func() {
		reqBody := controllers.ChangePasswordRequest{
			OldPassword: "OldPassword123!",
			NewPassword: "NewPassword123!",
		}

		suite.userUsecase.On("ChangePassword", mock.Anything, "test-user-id", "OldPassword123!", "NewPassword123!").Return(nil)

		w := suite.makeAuthenticatedRequest("POST", "/api/v1/auth/change-password", reqBody, "test-user-id", domain.RoleUser)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "Password changed successfully")
	})

	// todo check this test
	// suite.Run("password change with wrong old password", func() {
	// 	reqBody := controllers.ChangePasswordRequest{
	// 		OldPassword: "WrongPassword",
	// 		NewPassword: "NewPassword123!",
	// 	}

	// 	suite.userUsecase.On("ChangePassword", mock.Anything, "test-user-id", "WrongPassword", "NewPassword123!").Return(domain.ErrInvalidCredentials)

	// 	w := suite.makeAuthenticatedRequest("POST", "/api/v1/auth/change-password", reqBody, "test-user-id", domain.RoleUser)

	// 	suite.Equal(http.StatusUnauthorized, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "Invalid old password")
	// })

	suite.Run("password change without authentication", func() {
		reqBody := controllers.ChangePasswordRequest{
			OldPassword: "OldPassword123!",
			NewPassword: "NewPassword123!",
		}

		w := suite.makeRequest("POST", "/api/v1/auth/change-password", reqBody, nil)

		suite.Equal(http.StatusUnauthorized, w.Code)
		response := suite.parseResponse(w)
		suite.False(response.Success)
	})
}

func (suite *APITestSuite) TestResendOTP() {
	suite.Run("successful OTP resend", func() {
		reqBody := controllers.ResendOTPRequest{
			Email: "test@example.com",
		}

		suite.userUsecase.On("ResendOTP", mock.Anything, "test@example.com").Return(nil)

		w := suite.makeRequest("POST", "/api/v1/auth/resend-otp", reqBody, nil)

		suite.Equal(http.StatusOK, w.Code)
		response := suite.parseResponse(w)
		suite.True(response.Success)
		suite.Contains(response.Message, "Verification code resent")
	})

	// todo check this test
	// suite.Run("OTP resend for non-existent user", func() {
	// 	reqBody := controllers.ResendOTPRequest{
	// 		Email: "nonexistent@example.com",
	// 	}

	// 	suite.userUsecase.On("ResendOTP", mock.Anything, "nonexistent@example.com").Return(domain.ErrUserNotFound)

	// 	w := suite.makeRequest("POST", "/api/v1/auth/resend-otp", reqBody, nil)

	// 	suite.Equal(http.StatusNotFound, w.Code)
	// 	response := suite.parseResponse(w)
	// 	suite.False(response.Success)
	// 	suite.Contains(response.Message, "User not found")
	// })
}
