package service

import (
	"context"
	"finnbank/common/grpc/account"
	"finnbank/common/utils"
	"finnbank/internal-services/account/auth"
	"finnbank/internal-services/account/helpers"
	middleware "finnbank/internal-services/account/account_service"
	"sync"
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
	Auth   *auth.AuthService
	Grpc   account.AccountServiceServer
	account.UnimplementedAccountServiceServer
}

// Create New Account
// PARAMS:  email, fullname, password, address, account type
// This shit is getting out of hand so i will have to move some of this code somewhere else
// FUTURE: add concurrency if some of these queries will be moved
func (s *AccountService) AddAccount(ctx context.Context, req *account.AddAccountRequest) (*account.AddAccountResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)
	res, err := s.Auth.SignUpUserToDb(req.Email, req.Password)
	if err != nil {
		s.Logger.Error("Authentication failed: %v", err)
		return nil, status.Errorf(codes.Canceled, "Invalid email or password")
	}
	var authID string = res.User.ID
	if authID == "" {
		s.Logger.Error("Empty authID: %v", authID)
		return nil, status.Errorf(codes.Internal, "Auth Error")
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
// FUTURE: I think i'll have to make another version of this where i use the email instead of the account number
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

// Fetch Account
// PARAMS: email
func (s *AccountService) GetAccountByEmail(ctx context.Context, req *account.EmailRequest) (*account.AccountResponse, error) {
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

	accQuery := `SELECT email, full_name, phone_number, address, account_type, account_number, has_card, date_created
	FROM account WHERE email = $1;
	`
	err = tx.QueryRow(ctx, accQuery, req.Email).Scan(
		&email, &fullName, &phoneNumber, &address, &accountType, &accountNumber, &hasCard, &dateCreated,
	)
	if err != nil {
		s.Logger.Error("Could not Fetch account from Db: %v", err)
		return nil, status.Errorf(codes.Internal, "Error Fetching from DB: %v", err)
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
func (s *AccountService) UpdatePassword(ctx context.Context, req *account.UpdatePasswordRequest) (*account.UpdatePasswordResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)

	authUser, err := s.Auth.GetUserAuth(ctx, req.AuthId)
	if err != nil {
		s.Logger.Error("Failed to fetch user: %v", err)
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	isValid, err := s.Auth.VerifyPassword(authUser.EnryptedPass, req.OldPassword)
	if err != nil || !isValid {
		s.Logger.Error("Invalid old password: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "Invalid old password")
	}

	newPasswordHash, err := s.Auth.HashPassword(req.NewPassword)
	if err != nil {
		s.Logger.Error("Failed to hash new password: %v", err)
		return nil, status.Errorf(codes.Internal, "Error hashing new password")
	}

	updateQuery := `UPDATE auth.users SET encrypted_password = $1, updated_at = NOW() WHERE id = $2`
	_, err = s.DB.Exec(ctx, updateQuery, newPasswordHash, req.AuthId)
	if err != nil {
		s.Logger.Error("Failed to update password: %v", err)
		return nil, status.Errorf(codes.Internal, "Error updating password")
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.Error("Transaction commit failed: %v", err)
		return nil, status.Errorf(codes.Internal, "Error finalizing password update")
	}
	// Return Message here
	return &account.UpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully",
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
	// TODO: Account Validation if User already has a card

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

// TODO: DELETE ROUTES

// FUTURE: Will have to move this somewhere else
func (s *AccountService) DeleteFromAuth(ctx context.Context, req uuid.UUID) (string, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return "Could not Connect to Db", status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)

	deleteQuery := `
		DELETE FROM auth.users
		WHERE id = $1;
	`
	res, err := tx.Exec(ctx, deleteQuery, req)
	if err != nil {
		s.Logger.Error("Could not delete from auth : %v", err)
		return "Error deleting row in auth", status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		s.Logger.Error("No account found with uuid: %s", req)
		return "No Account Found", status.Errorf(codes.NotFound, "No account found with the given account number")
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.Logger.Error("Transaction commit failed: %v", err)
		return "Error Committing Db transaction", status.Errorf(codes.Internal, "Error finalizing account creation")
	}

	return "Successfully Deleted account from Auth", nil

}

func (s *AccountService) DeleteAccount(ctx context.Context, req *account.DeleteUserRequest) (*account.DeleteUserResponse, error) {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.Logger.Error("Could not start DB Transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "Error starting DB: %v", err)
	}
	defer tx.Rollback(ctx)
	var UUID uuid.UUID
	uuidQuery := `SELECT auth_id FROM account WHERE account_number = $1;`
	err = tx.QueryRow(ctx, uuidQuery, req.AccountNumber).Scan(&UUID)
	if err != nil {
		s.Logger.Error("Could not fetch UUID for account_number: %s", req.AccountNumber)
		return nil, status.Errorf(codes.NotFound, "No account found with the given account number")
	}

	deleteQuery := `
		DELETE FROM account
		WHERE account_number = $1;
	`
	res, err := tx.Exec(ctx, deleteQuery, req.AccountNumber)
	if err != nil {
		s.Logger.Error("Could not Delete account: %v", err)
		return nil, status.Errorf(codes.Internal, "Error deleting account in DB: %v", err)
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
	var wg sync.WaitGroup
	var authDeleteErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, authDeleteErr = s.DeleteFromAuth(ctx, UUID)
	}()

	wg.Wait()

	if authDeleteErr != nil {
		s.Logger.Error("Could not delete from Auth: %v", authDeleteErr)
		return nil, status.Errorf(codes.Internal, "Error deleting from Auth: %v", authDeleteErr)
	}

	return &account.DeleteUserResponse{
		Message: "Successfully Deleted Account",
	}, nil
}

// AUTH SERVICES
func (s *AccountService) LoginUser(ctx context.Context, req *account.LoginRequest) (*account.LoginResponse, error) {
	var (
		res *account.AccountResponse
		err error
		wg  sync.WaitGroup
	)

	errCh := make(chan error, 1)

	tok, authErr := s.Auth.LoginUserToDb(req.Email, req.Password)
	if authErr != nil {
		s.Logger.Error("Authentication failed: %v", authErr)
		return nil, status.Errorf(codes.Unauthenticated, "Invalid email or password")
	}
	if tok == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Authentication failed: empty response")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		emailReq := &account.EmailRequest{Email: req.Email}
		res, err = s.GetAccountByEmail(ctx, emailReq)
		if err != nil {
			errCh <- status.Errorf(codes.Internal, "Failed to fetch Account: %v", err)
		}
	}()

	wg.Wait()
	close(errCh)

	if err = <-errCh; err != nil {
		return nil, err
	}

	// Construct account data
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

	// Return the response
	return &account.LoginResponse{
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		ExpiresIn:    int32(tok.ExpiresIn),
		AuthId:       tok.User.ID,
		Email:        tok.User.Email,
		Account:      acc,
	}, nil
}
