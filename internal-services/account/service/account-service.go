package service

import (
	"bytes"
	"context"
	"encoding/json"
	pb "finnbank/common/grpc/auth"
	"finnbank/common/utils"
	"finnbank/internal-services/account/auth"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	Logger *utils.Logger
	DB     *pgx.Conn
	Helper *auth.AuthHelpers
	Grpc   pb.AuthServiceServer
	pb.UnimplementedAuthServiceServer
}

func (s *AuthService) SignUpUser(ctx context.Context, in *pb.SignUpRequest) (*pb.AuthResponse, error) {
	baseURL := os.Getenv("DB_URL")
	if baseURL == "" {
		s.Logger.Error("missing URL environment variable")
		return nil, fmt.Errorf("missing URL environment variable")
	}
	url := fmt.Sprintf("%s/auth/v1/signup", baseURL)
	reqBody, err := json.Marshal(map[string]string{
		"email":    in.Email,
		"password": in.Password,
	})
	if err != nil {
		s.Logger.Error("failed to create request body: %v", err)
		return nil, fmt.Errorf("failed to create request body: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		s.Logger.Error("failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	apiKey := os.Getenv("DB_ANON_KEY")
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
	var token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		User        struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	userInfo := &pb.UserInfo{
		Id:    token.User.ID,
		Email: token.User.Email,
	}
	authResponse := &pb.AuthResponse{
		User:      userInfo,
		TokenType: token.TokenType,
	}
	return authResponse, nil
}

func (s *AuthService) LoginUser(c context.Context, in *pb.LoginRequest) (*pb.AuthResponse, error) {
	baseURL := os.Getenv("DB_URL")
	if baseURL == "" {
		s.Logger.Error("missing DB_URL environment variable")
		return nil, fmt.Errorf("missing DB_URL environment variable")
	}
	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", baseURL)
	reqBody, err := json.Marshal(map[string]string{
		"email":    in.Email,
		"password": in.Password,
	})
	if err != nil {
		s.Logger.Error("failed to create request body: %v", err)
		return nil, fmt.Errorf("failed to create request body: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		s.Logger.Error("failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	apiKey := os.Getenv("DB_ANON_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing KEY environment variable")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		s.Logger.Error("login failed, status: %d, response: %s", resp.StatusCode, bodyBytes)
		return nil, fmt.Errorf("login failed: %s", bodyBytes)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	var token struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	userInfo := &pb.UserInfo{
		Id:    token.User.ID,
		Email: token.User.Email,
	}
	authResponse := &pb.AuthResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		ExpiresIn:    int32(token.ExpiresIn),
		RefreshToken: token.RefreshToken,
		User:         userInfo,
	}
	return authResponse, nil
}

// TODO: change this
func (s *AuthService) HashAndEncryptPassowrd(c context.Context, in *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	oldEncryptedPassword, err := s.Helper.GetUserAuth(c, in.AuthID, s.DB)
	if err != nil {
		s.Logger.Error("Failed to get user auth: %v", err)
		return nil, err
	}
	ok := s.Helper.VerifyPassword(oldEncryptedPassword.EnryptedPass, in.OldPassword)
	if !ok {
		s.Logger.Error("Old password is incorrect")
		return nil, fmt.Errorf("old password is incorrect")
	}
	newEncryptedPassword, err := s.Helper.HashPassword(in.NewPassword)
	if err != nil {
		s.Logger.Error("Failed to hash new password: %v", err)
		return nil, err
	}

	authResponse := &pb.UpdatePasswordResponse{
		EncryptedPassword: newEncryptedPassword,
	}
	return authResponse, nil

}
