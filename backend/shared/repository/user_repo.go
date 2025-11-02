package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	
	query := `
		INSERT INTO users (id, email, password_hash, name, role, email_verified, verification_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`
	
	return r.db.QueryRow(
		query, user.ID, user.Email, user.PasswordHash, user.Name, user.Role, user.EmailVerified, user.VerificationToken,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	
	query := `
		SELECT id, email, password_hash, name, role, 
		       COALESCE(email_verified, false) as email_verified,
		       verification_token, password_reset_token, password_reset_expires,
		       created_at, updated_at
		FROM users WHERE email = $1
	`
	
	var verificationToken, resetToken sql.NullString
	var resetExpires sql.NullTime
	
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.EmailVerified, &verificationToken, &resetToken,
		&resetExpires, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	if err != nil {
		return nil, err
	}
	
	// Handle NULL values
	if verificationToken.Valid {
		user.VerificationToken = verificationToken.String
	}
	if resetToken.Valid {
		user.PasswordResetToken = resetToken.String
	}
	if resetExpires.Valid {
		user.PasswordResetExpires = &resetExpires.Time
	}
	
	return user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	
	query := `
		SELECT id, email, password_hash, name, role,
		       COALESCE(email_verified, false) as email_verified,
		       verification_token, password_reset_token, password_reset_expires,
		       created_at, updated_at
		FROM users WHERE id = $1
	`
	
	var verificationToken, resetToken sql.NullString
	var resetExpires sql.NullTime
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.EmailVerified, &verificationToken, &resetToken,
		&resetExpires, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	if err != nil {
		return nil, err
	}
	
	// Handle NULL values
	if verificationToken.Valid {
		user.VerificationToken = verificationToken.String
	}
	if resetToken.Valid {
		user.PasswordResetToken = resetToken.String
	}
	if resetExpires.Valid {
		user.PasswordResetExpires = &resetExpires.Time
	}
	
	return user, nil
}

func (r *UserRepository) GetByVerificationToken(token string) (*models.User, error) {
	user := &models.User{}
	
	query := `
		SELECT id, email, password_hash, name, role, email_verified, verification_token,
		       password_reset_token, password_reset_expires, created_at, updated_at
		FROM users WHERE verification_token = $1
	`
	
	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.EmailVerified, &user.VerificationToken, &user.PasswordResetToken,
		&user.PasswordResetExpires, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	return user, err
}

func (r *UserRepository) GetByPasswordResetToken(token string) (*models.User, error) {
	user := &models.User{}
	
	query := `
		SELECT id, email, password_hash, name, role, email_verified, verification_token,
		       password_reset_token, password_reset_expires, created_at, updated_at
		FROM users WHERE password_reset_token = $1
	`
	
	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.EmailVerified, &user.VerificationToken, &user.PasswordResetToken,
		&user.PasswordResetExpires, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	return user, err
}

func (r *UserRepository) VerifyEmail(userID string) error {
	query := `
		UPDATE users 
		SET email_verified = TRUE, verification_token = NULL, updated_at = NOW()
		WHERE id = $1
	`
	
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserRepository) SetPasswordResetToken(userID, token string, expires *time.Time) error {
	query := `
		UPDATE users 
		SET password_reset_token = $1, password_reset_expires = $2, updated_at = NOW()
		WHERE id = $3
	`
	
	_, err := r.db.Exec(query, token, expires, userID)
	return err
}

func (r *UserRepository) UpdatePassword(userID, passwordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = $1, password_reset_token = NULL, password_reset_expires = NULL, updated_at = NOW()
		WHERE id = $2
	`
	
	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}
