package types

import (
	"time"

	"github.com/graphql-go/graphql"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response.
type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int32   `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	AuthID       string  `json:"auth_id"`
	Email        string  `json:"email"`
	Account      Account `json:"account"`
}

// Account represents an account.
type Account struct {
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	AccountType   string    `json:"account_type"`
	AccountNumber string    `json:"account_number"`
	HasCard       bool      `json:"has_card"`
	Balance       float64   `json:"balance"`
	DateCreated   time.Time `json:"date_created"`
}

// AccountsResponse represents a response containing multiple accounts.
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// AccountRequest represents a request for a specific account.
type AccountRequest struct {
	AccountNumber string `json:"account_number"`
}

// AccountResponse represents a response containing a single account.
type AccountResponse struct {
	Account Account `json:"account"`
}

// UpdateRequest represents a request to update account details.
type UpdateRequest struct {
	AccountNumber string `json:"account_number"`
	FullName      string `json:"full_name"`
	PhoneNumber   string `json:"phone_number"`
	Address       string `json:"address"`
}

// CardUpdateRequest represents a request to update card details.
type CardUpdateRequest struct {
	AccountNumber string `json:"account_number"`
}

// CardUpdateResponse represents a response for a card update request.
type CardUpdateResponse struct {
	Status string `json:"status"`
}

// UpdatePasswordRequest represents a request to update a password.
type UpdatePasswordRequest struct {
	AuthID      string `json:"auth_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UpdatePasswordResponse represents a response for a password update request.
type UpdatePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DeleteUserRequest represents a request to delete a user.
type DeleteUserRequest struct {
	AccountNumber string `json:"account_number"`
}

// DeleteUserResponse represents a response for a user deletion request.
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// AddAccountRequest represents a request to add a new account.
type AddAccountRequest struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	AccountType string `json:"account_type"`
	Nationality string `json:"nationality"`
}

// AddAccountResponse represents a response for adding a new account.
type AddAccountResponse struct {
	ID            string    `json:"account_id"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	AccountType   string    `json:"account_type"`
	AccountNumber string    `json:"account_number"`
	AuthID        string    `json:"auth_id"`
	HasCard       bool      `json:"has_card"`
	Balance       float64   `json:"balance"`
	DateCreated   time.Time `json:"date_created"`
	Nationality   string    `json:"nationality"`
}

// For signup input
var AccountInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "AccountInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"email": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"password": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"first_name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"last_name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"phone_number": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"address": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"account_type": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"nationality": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	},
)
