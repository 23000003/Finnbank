package service

import (
	"context"
	"finnbank/services/account/helpers"
	"finnbank/services/common/grpc/account"
	"finnbank/services/common/utils"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountService struct {
	DB     *pgx.Conn
	Logger *utils.Logger
	Grpc   account.AccountServiceServer
	account.UnimplementedAccountServiceServer
}

// func (s* AccountService) GetAccountbyId(ctx * context.Context, req* account.AccountRequest) (*account.AccountResponse, error) {

// 	return
// // }

func (s *AccountService) AddAccount(ctx context.Context, req *account.AddAccountRequest) (*account.AddAccountResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)

	authQuery := `
	INSERT INTO auth.users (id, email, encrypted_password, aud, instance_id)
	VALUES (gen_random_uuid(),$1, crypt($2, gen_salt('bf')), 'authenticated', gen_random_uuid())`

	_, err = tx.Exec(ctx, authQuery, req.Email, req.Password)
	if err != nil {
		s.Logger.Error("Failed to Create User in auth: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	// email            VARCHAR(255) UNIQUE NOT NULL,
	// full_name        VARCHAR(255) NOT NULL,
	// phone_number     VARCHAR(11) UNIQUE NOT NULL,
	// account_number   varchar(20) UNIQUE NOT NULL,
	// address          TEXT,
	// account_type     VARCHAR(20) CHECK (account_type IN ('Business', 'Personal')) NOT NULL,
	userID, err := helpers.GenAccNum(req.FullName)
	if err != nil {
		s.Logger.Error("Failed to Generate Account Number: %v", err)
		return nil, status.Errorf(codes.Internal, "Error generating account number: %v", err)
	}
	accQuery := `
	INSERT INTO account (email, full_name, phone_number, account_number, address, account_type)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, accQuery, req.Email, req.FullName, req.PhoneNumber, userID, req.Address, req.AccountType)
	if err != nil {
		s.Logger.Error("Failed to Create User in table: %v", err)
		return nil, status.Error(codes.Internal, "Error creating user")
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.Error("Transaction commit failed: %v", err)
		return nil, status.Errorf(codes.Internal, "Error finalizing account creation")
	}

	// message AddAccountResponse {
	// 	string email = 1;
	// 	string full_name = 2;
	// 	string phone_number = 3;
	// 	string password = 4;
	// 	string address = 5;
	// 	string account_type = 6;
	// 	string account_number = 7;
	//   }

	return &account.AddAccountResponse{
		Email:         req.Email,
		FullName:      req.FullName,
		PhoneNumber:   req.PhoneNumber,
		Address:       req.Address,
		AccountType:   req.AccountType,
		AccountNumber: userID,
	}, nil
}
