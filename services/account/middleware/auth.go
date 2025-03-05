package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// These middleware code can be reused for later

// func FetchUserUUID() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		encodedEmail := c.Param("email")
// 		email, err := url.QueryUnescape(encodedEmail)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
// 			c.Abort()
// 			return
// 		}

// 		uuid, err := helpers.GetUserUUIDByEmail(email)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 			c.Abort()
// 			return
// 		}

// 		// Storing UUID in Gin context for later use
// 		c.Set("email", email)
// 		c.Set("uuid", uuid)

// 		c.Next()
// 	}
// }

// func DeleteAuthUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Retrieving UUID from context
// 		uuid, exists := c.Get("uuid")
// 		if !exists {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found in context"})
// 			c.Abort()
// 			return
// 		}

// 		err := helpers.DeleteUserByUUID(uuid.(string))
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user from Auth: " + err.Error()})
// 			c.Abort()
// 			return
// 		}
// 		c.Next()
// 	}
// }

type AuthService struct {
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

// FUTURE: Implement this for better error handling
// type AuthError struct {
// 	Code      int    `json:"code"`
// 	ErrorCode string `json:"error_code"`
// 	Message   string `json:"msg"`
// }

func (s *AuthService) SignUpUserToDb(email, password string) (*AuthResponse, error) {
	baseURL := os.Getenv("LOCAL_AUTH_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("missing URL environment variable")
	}
	url := fmt.Sprintf("%s/auth/v1/signup", baseURL)
	reqBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	apiKey := os.Getenv("LOCAL_DB_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing KEY environment variable")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	var token AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &token, nil
}

func (s *AuthService) LoginUserToDb(email, password string) (*AuthResponse, error) {
	baseURL := os.Getenv("LOCAL_AUTH_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("missing LOCAL_AUTH_URL environment variable")
	}
	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", baseURL)
	reqBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	apiKey := os.Getenv("LOCAL_DB_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing LOCAL_DB_KEY environment variable")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	// Decoding response into AuthResponse
	var token AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &token, nil
}

// json.NewDecoder is so finnicky kay if maka decode naka ka-isa di na siya pwede ma decode again even if it's for error handling
// will definitely have to improve this in the future, but for now this works

// Adding user in auth.users table
// authQuery := `
// INSERT INTO auth.users (id, email, encrypted_password, aud, instance_id, role, created_at, email_confirmed_at,
// raw_app_meta_data, raw_user_meta_data)
// VALUES (gen_random_uuid(), $1, crypt($2, gen_salt('bf')), 'authenticated', gen_random_uuid(),'authenticated',
// NOW(), NOW(), '{"provider": "email", "providers": ["email"]}'::jsonb,
// '{"email_verified": true}'::jsonb)  RETURNING id;`
// var authID uuid.UUID
// err = tx.QueryRow(ctx, authQuery, req.Email, req.Password).Scan(&authID)
// if err != nil {
// 	s.Logger.Error("Failed to Create User in auth: %v", err)
// 	return nil, status.Errorf(codes.Internal, "Error creating user in auth DB: %v", err)
// }
// // Adding user in auth.identities table
// identityQuery := `
// INSERT INTO auth.identities (provider_id, user_id, identity_data, provider, last_sign_in_at, created_at, updated_at)
// VALUES ($1, $2, $3, 'email', NOW(), NOW(), NOW());`

// identityData := fmt.Sprintf(`{"sub": "%s", "email": "%s"}`, authID, req.Email)
// _, err = tx.Exec(ctx, identityQuery, req.Email, authID, identityData)
// if err != nil {
// 	s.Logger.Error("Failed to Insert User Identity: %v", err)
// 	return nil, status.Errorf(codes.Internal, "Error inserting user identity: %v", err)
// }
// Generating the user id for the account table
