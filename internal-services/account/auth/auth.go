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

// FUTURE: Implement this for better error handling
// type AuthError struct {
// 	Code      int    `json:"code"`
// 	ErrorCode string `json:"error_code"`
// 	Message   string `json:"msg"`
// }

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

// func (s *AuthService) SignUpUserToDb(email, password string) (*AuthResponse, error) {
// 	// LOCAL_AUTH_URL <-- LOCAL DB
// 	// DB_URL <-- PROD DB
// 	baseURL := os.Getenv("LOCAL_AUTH_URL")
// 	if baseURL == "" {
// 		return nil, fmt.Errorf("missing URL environment variable")
// 	}
// 	url := fmt.Sprintf("%s/auth/v1/signup", baseURL)
// 	reqBody, err := json.Marshal(map[string]string{
// 		"email":    email,
// 		"password": password,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request body: %v", err)
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %v", err)
// 	}
// 	// LOCAL_DB_KEY <-- LOCAL ANON KEY
// 	// DB_ANON_KEY <-- PROD ANON KEY
// 	apiKey := os.Getenv("LOCAL_DB_KEY")
// 	if apiKey == "" {
// 		return nil, fmt.Errorf("missing KEY environment variable")
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("apikey", apiKey)
// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()
// 	var token AuthResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %v", err)
// 	}

// 	return &token, nil
// }

// func (s *AuthService) LoginUserToDb(email, password string) (*AuthResponse, error) {
// 	// LOCAL_AUTH_URL <-- LOCAL DB
// 	// DB_URL <-- PROD DB
// 	baseURL := os.Getenv("LOCAL_AUTH_URL")
// 	if baseURL == "" {
// 		return nil, fmt.Errorf("missing DB_URL environment variable")
// 	}
// 	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", baseURL)
// 	reqBody, err := json.Marshal(map[string]string{
// 		"email":    email,
// 		"password": password,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request body: %v", err)
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %v", err)
// 	}
// 	// LOCAL_DB_KEY <-- LOCAL ANON KEY
// 	// DB_ANON_KEY <-- PROD ANON KEY
// 	apiKey := os.Getenv("LOCAL_DB_KEY")
// 	if apiKey == "" {
// 		return nil, fmt.Errorf("missing DB_ANON_KEY environment variable")
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("apikey", apiKey)
// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("authentication failed: %s", string(body))
// 	}
// 	// Decoding response into AuthResponse
// 	var token AuthResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %v", err)
// 	}
// 	if token.AccessToken == "" {
// 		return nil, fmt.Errorf("authentication failed: empty access token")
// 	}

// 	return &token, nil
// }
