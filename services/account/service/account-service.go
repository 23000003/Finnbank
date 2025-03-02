package service

import (
	"context"
	"finnbank/services/account/helpers"
	"finnbank/services/common/grpc/account"
	"finnbank/services/common/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountService struct {
	DB     *pgx.Conn
	Logger *utils.Logger
	Grpc   account.AccountServiceServer
	account.UnimplementedAccountServiceServer
}

// Create New Account
// PARAMS:  email, fullname, password, address, account type
func (s *AccountService) AddAccount(ctx context.Context, req *account.AddAccountRequest) (*account.AddAccountResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)

	authQuery := `
	INSERT INTO auth.users (id, email, encrypted_password, aud, instance_id)
	VALUES (gen_random_uuid(),$1, crypt($2, gen_salt('bf')), 'authenticated', gen_random_uuid())
	RETURNING id;`
	var authID uuid.UUID
	err = tx.QueryRow(ctx, authQuery, req.Email, req.Password).Scan(&authID)
	if err != nil {
		s.Logger.Error("Failed to Create User in auth: %v", err)
		return nil, status.Errorf(codes.Internal, "Error creating user in auth DB: %v", err)
	}
	userID, err := helpers.GenAccNum()
	if err != nil {
		s.Logger.Error("Failed to Generate Account Number: %v", err)
		return nil, status.Errorf(codes.Internal, "Error generating account number: %v", err)
	}
	accQuery := `
	INSERT INTO account (email, full_name, phone_number, account_number, address, account_type, auth_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, accQuery, req.Email, req.FullName, req.PhoneNumber, userID, req.Address, req.AccountType, authID)
	if err != nil {
		s.Logger.Error("Failed to Create User in table: %v", err)
		return nil, status.Error(codes.Internal, "Error creating user")
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.Error("Transaction commit failed: %v", err)
		return nil, status.Errorf(codes.Internal, "Error finalizing account creation")
	}

	return &account.AddAccountResponse{
		Email:         req.Email,
		FullName:      req.FullName,
		PhoneNumber:   req.PhoneNumber,
		Address:       req.Address,
		AccountType:   req.AccountType,
		AccountNumber: userID,
	}, nil
}

// Fetch Account
// PARAMS:  account number
func (s *AccountService) GetAccountById(ctx context.Context, req *account.AccountRequest) (*account.AccountResponse, error) {
	var (
		email, fullName, phoneNumber, address, accountType, accountNumber string
		hasCard                                                           bool
		dateCreated                                                       time.Time
	)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)

	accQuery := `
	SELECT email, full_name, phone_number, address, account_type, account_number, has_card, date_created
	FROM account WHERE account_number = $1;
	`

	err = tx.QueryRow(ctx, accQuery, req.AccountNumber).Scan(
		&email, &fullName, &phoneNumber, &address, &accountType, &accountNumber, &hasCard, &dateCreated,
	)
	if err != nil {
		s.Logger.Error("Failed to Fetch Account: %v", err)
		return nil, status.Errorf(codes.Internal, "Error fetching account from DB: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.Error("Transaction commit failed: %v", err)
		return nil, status.Errorf(codes.Internal, "Error finalizing account creation")
	}

	gotAcc := &account.Account{
		Email:         email,
		FullName:      fullName,
		PhoneNumber:   phoneNumber,
		Address:       address,
		AccountType:   accountType,
		AccountNumber: accountNumber,
		HasCard:       hasCard,
		DateCreated:   timestamppb.New(dateCreated),
	}

	return &account.AccountResponse{
		Account: gotAcc,
	}, nil
}

// Update Account
// PARAMS: account number, fullname, phone number, address
func (s *AccountService) UpdateAccount(ctx context.Context, req *account.UpdateRequest) (*account.AccountResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE account
		SET full_name = $1, phone_number = $2, address = $3
		WHERE account_number = $4;
	`
	result, err := tx.Exec(ctx, updateQuery, req.FullName, req.PhoneNumber, req.Address, req.AccountNumber)
	if err != nil {
		s.Logger.Error("Could not Update account: %v", err)
		return nil, status.Errorf(codes.Internal, "Error updating account: %v", err)
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		s.Logger.Error("No account found with account_number: %s", req.AccountNumber)
		return nil, status.Errorf(codes.NotFound, "No account found with the given account number")
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.Error("Transaction commit failed: %v", err)
		return nil, status.Errorf(codes.Internal, "Error finalizing account creation")
	}
	accReq := &account.AccountRequest{
		AccountNumber: req.AccountNumber,
	}

	res, err := s.GetAccountById(ctx, accReq)
	if err != nil {
		s.Logger.Error("Could not Fetch newly updated account: %v", err)
		return nil, status.Errorf(codes.Internal, "Error fetching updated account: %v", err)
	}
	acc := &account.Account{
		Email:         res.Account.Email,
		FullName:      res.Account.FullName,
		PhoneNumber:   res.Account.PhoneNumber,
		Address:       res.Account.Address,
		AccountType:   res.Account.AccountType,
		AccountNumber: res.Account.AccountNumber,
		HasCard:       res.Account.HasCard,
		Balance:       res.Account.Balance,
		DateCreated:   res.Account.DateCreated,
	}
	return &account.AccountResponse{
		Account: acc,
	}, nil
}

func (s *AccountService) UpdateCardStatus(ctx context.Context, req *account.CardUpdateRequest) (*account.CardUpdateResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE account
		SET has_card = 'TRUE'
		WHERE account_number = $1;
	`
	res, err := tx.Exec(ctx, updateQuery, req.AccountNumber)
	if err != nil {
		s.Logger.Error("Could not Update card status in DB: %v", err)
		return nil, status.Errorf(codes.Internal, "Error updating card status: %v", err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		s.Logger.Error("No account found with account_number: %s", req.AccountNumber)
		return nil, status.Errorf(codes.NotFound, "No account found with the given account number")
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.Error("Transaction commit failed: %v", err)
		return nil, status.Errorf(codes.Internal, "Error finalizing account creation")
	}

	return &account.CardUpdateResponse{
		Status: "Sucessfully Updated Card Status",
	}, nil
}
