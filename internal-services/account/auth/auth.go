package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHelpers struct {
	DB *pgx.Conn
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

// For now i'm using this to store the hashed password got from the auth.users table
type AuthUser struct {
	EnryptedPass string `json:"encrypted_password"`
}

func (s *AuthHelpers) GetUserAuth(ctx context.Context, authID string, db *pgx.Conn) (*AuthUser, error) {
	var authUser AuthUser
	query := `SELECT encrypted_password FROM auth.users WHERE id = $1;`
	err := s.DB.QueryRow(ctx, query, authID).Scan(&authUser.EnryptedPass)
	if err != nil {
		return nil, fmt.Errorf("error querying auth user: %v", err)
	}
	return &authUser, nil
}
func (a *AuthHelpers) VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
func (a *AuthHelpers) HashPassword(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
