package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/auth"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/email"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/logger"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	userRepo     *repository.UserRepository
	emailService *email.EmailService
}

func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		emailService: email.NewEmailService(),
	}
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger("auth-handler")
	
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("Invalid signup request body", map[string]interface{}{"error": err.Error()})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		http.Error(w, "Email, password, and name are required", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = "developer"
	}

	if req.Role != "developer" && req.Role != "consumer" {
		http.Error(w, "Role must be 'developer' or 'consumer'", http.StatusBadRequest)
		return
	}

	// Hash password
	passwordHash, err := models.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Generate verification token
	verificationToken, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		Email:             req.Email,
		PasswordHash:      passwordHash,
		Name:              req.Name,
		Role:              req.Role,
		EmailVerified:     false,
		VerificationToken: verificationToken,
	}

	if err := h.userRepo.Create(user); err != nil {
		log.Error("Failed to create user", map[string]interface{}{"email": req.Email, "error": err.Error()})
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Info("User created successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email, "role": user.Role})

	// Send verification email
	if err := h.emailService.SendVerificationEmail(user.Email, user.Name, verificationToken); err != nil {
		log.Error("Failed to send verification email", map[string]interface{}{"error": err.Error()})
		// Don't fail signup if email fails
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger("auth-handler")
	
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("Invalid login request body", map[string]interface{}{"error": err.Error()})
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		log.Warn("Login attempt with invalid email", map[string]interface{}{"email": req.Email})
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if !models.CheckPassword(req.Password, user.PasswordHash) {
		log.Warn("Login attempt with invalid password", map[string]interface{}{"user_id": user.ID, "email": user.Email})
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	log.Info("User logged in successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email})

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger("auth-handler")
	vars := mux.Vars(r)
	token := vars["token"]

	if token == "" {
		http.Error(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	// Get user by verification token
	user, err := h.userRepo.GetByVerificationToken(token)
	if err != nil {
		log.Warn("Invalid verification token", map[string]interface{}{"token": token})
		http.Error(w, "Invalid or expired verification token", http.StatusBadRequest)
		return
	}

	// Check if already verified
	if user.EmailVerified {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Email already verified",
		})
		return
	}

	// Verify email
	if err := h.userRepo.VerifyEmail(user.ID); err != nil {
		log.Error("Failed to verify email", map[string]interface{}{"user_id": user.ID, "error": err.Error()})
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}

	log.Info("Email verified successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email})

	// Send welcome email
	if err := h.emailService.SendWelcomeEmail(user.Email, user.Name); err != nil {
		log.Error("Failed to send welcome email", map[string]interface{}{"error": err.Error()})
	}

	// Generate JWT token
	token, err = auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger("auth-handler")

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "If your email is registered, you will receive a password reset link",
		})
		return
	}

	// Generate reset token
	resetToken, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}

	// Set expiry to 1 hour from now
	expires := time.Now().Add(1 * time.Hour)

	// Save reset token
	if err := h.userRepo.SetPasswordResetToken(user.ID, resetToken, &expires); err != nil {
		log.Error("Failed to set password reset token", map[string]interface{}{"user_id": user.ID, "error": err.Error()})
		http.Error(w, "Failed to process password reset", http.StatusInternalServerError)
		return
	}

	log.Info("Password reset requested", map[string]interface{}{"user_id": user.ID, "email": user.Email})

	// Send reset email
	if err := h.emailService.SendPasswordResetEmail(user.Email, user.Name, resetToken); err != nil {
		log.Error("Failed to send password reset email", map[string]interface{}{"error": err.Error()})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If your email is registered, you will receive a password reset link",
	})
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger("auth-handler")

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		http.Error(w, "Token and new password are required", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if len(req.NewPassword) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Get user by reset token
	user, err := h.userRepo.GetByPasswordResetToken(req.Token)
	if err != nil {
		log.Warn("Invalid password reset token", map[string]interface{}{"token": req.Token})
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}

	// Check if token is expired
	if user.PasswordResetExpires != nil && user.PasswordResetExpires.Before(time.Now()) {
		log.Warn("Expired password reset token", map[string]interface{}{"user_id": user.ID})
		http.Error(w, "Reset token has expired", http.StatusBadRequest)
		return
	}

	// Hash new password
	passwordHash, err := models.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Update password
	if err := h.userRepo.UpdatePassword(user.ID, passwordHash); err != nil {
		log.Error("Failed to update password", map[string]interface{}{"user_id": user.ID, "error": err.Error()})
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	log.Info("Password reset successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successfully",
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger("auth-handler")

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "Old password and new password are required", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify old password
	if !models.CheckPassword(req.OldPassword, user.PasswordHash) {
		log.Warn("Invalid old password", map[string]interface{}{"user_id": userID})
		http.Error(w, "Invalid old password", http.StatusUnauthorized)
		return
	}

	// Hash new password
	passwordHash, err := models.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Update password
	if err := h.userRepo.UpdatePassword(user.ID, passwordHash); err != nil {
		log.Error("Failed to change password", map[string]interface{}{"user_id": user.ID, "error": err.Error()})
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}

	log.Info("Password changed successfully", map[string]interface{}{"user_id": user.ID})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}

// generateToken creates a secure random token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
