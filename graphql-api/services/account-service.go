package services

import (
	"context"
	pb "finnbank/common/grpc/auth"
	"finnbank/common/utils"
	"finnbank/graphql-api/types"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountService struct {
	db          *pgxpool.Pool
	authService pb.AuthServiceClient
	l           *utils.Logger
}

func NewAccountService(db *pgxpool.Pool, logger *utils.Logger, pb pb.AuthServiceClient) *AccountService {
	return &AccountService{
		db:          db,
		l:           logger,
		authService: pb,
	}
}

func (s *AccountService) CreateUser(ctx *context.Context, in *types.AddAccountRequest) (*types.AddAccountResponse, error) {
	req := &pb.SignUpRequest{
		Email:    in.Email,
		Password: in.Password,
	}

	// i created this validation to avoid conflict if grpc is success and insert is not
	// no email validation since it is done in the auth service
	// no account_number validation since its too impossible for it to be the same
	// Check phone number
	phoneQuery := `SELECT EXISTS (SELECT 1 FROM account WHERE phone_number = $1)`
	var phoneExists bool
	if err := s.db.QueryRow(*ctx, phoneQuery, in.PhoneNumber).Scan(&phoneExists); err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, fmt.Errorf("account with phone number already exists")
	}

	accNum, err := generateRandomNumber(16)

	if err != nil {
		return nil, fmt.Errorf("failed to generate account number: %v", err)
	}

	// TODO: This seems really bad, will have to find a better way for this somehow
	authRes, err := s.authService.SignUpUser(*ctx, req)
	if err != nil {
		return nil, err
	}
	if in.AccountType != "Personal" && in.AccountType != "Business" {
		return nil, fmt.Errorf("account type must be either Personal or Business")
	}
	var acc types.Account
	var res types.AddAccountResponse
	createQuery := `
	INSERT INTO account (
		email, first_name, middle_name, last_name, phone_number, address, nationality,
		account_type, account_number, has_card, auth_id, birthdate, national_id
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING 
		id, email, first_name, middle_name, last_name, phone_number, address, nationality,
		account_type, account_number, has_card, date_created, date_updated, auth_id, birthdate, national_id, account_status
	`

	s.l.Info("Creating account for user: %s", authRes)
	s.l.Info("Account number: %s", accNum)

	err = s.db.QueryRow(*ctx,
		createQuery,
		in.Email,
		in.FirstName,
		in.MiddleName,
		in.LastName,
		in.PhoneNumber,
		in.Address,
		in.Nationality,
		in.AccountType,
		accNum,
		false,
		authRes.User.Id,
		in.BirthDate,
		in.NationalID).
		Scan(
			&acc.ID,
			&acc.Email,
			&acc.FirstName,
			&acc.MiddleName,
			&acc.LastName,
			&acc.PhoneNumber,
			&acc.Address,
			&acc.Nationality,
			&acc.AccountType,
			&acc.AccountNumber,
			&acc.HasCard,
			&acc.DateCreated,
			&acc.DateUpdated,
			&acc.AuthID,
			&acc.BirthDate,
			&acc.NationalID,
			&acc.AccountStatus,
		)
	if err != nil {
		return nil, err
	}
	res.Account = acc

	return &res, nil
}

func (s *AccountService) FetchUserByAccountNumber(ctx *context.Context, req string) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE account_number = $1
	`
	err := s.db.QueryRow(*ctx, query, req).Scan(
		&acc.ID,
		&acc.AuthID,
		&acc.Email,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.AccountStatus,
		&acc.NationalID,
		&acc.BirthDate,
		&acc.FirstName,
		&acc.MiddleName,
		&acc.LastName,
	)
	if err != nil {
		return nil, err
	}
	return &types.AccountResponse{Account: acc}, nil
}

func (s *AccountService) FetchUserByEmail(ctx *context.Context, req string) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE email = $1
	`
	err := s.db.QueryRow(*ctx, query, req).Scan(
		&acc.ID,
		&acc.AuthID,
		&acc.Email,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.AccountStatus,
		&acc.NationalID,
		&acc.BirthDate,
		&acc.FirstName,
		&acc.MiddleName,
		&acc.LastName,
	)
	if err != nil {
		return nil, err
	}
	return &types.AccountResponse{Account: acc}, nil
}

func (s *AccountService) FetchUserById(ctx *context.Context, req string) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE id = $1
	`
	err := s.db.QueryRow(*ctx, query, req).Scan(
		&acc.ID,
		&acc.AuthID,
		&acc.Email,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.AccountStatus,
		&acc.NationalID,
		&acc.BirthDate,
		&acc.FirstName,
		&acc.MiddleName,
		&acc.LastName,
	)
	if err != nil {
		return nil, err
	}
	return &types.AccountResponse{Account: acc}, nil
}

// Potentially build the user dashboard with this, since auth_id is passed during login
func (s *AccountService) FetchUserByAuthID(ctx *context.Context, req string) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE auth_id = $1
	`
	err := s.db.QueryRow(*ctx, query, req).Scan(
		&acc.ID,
		&acc.AuthID,
		&acc.Email,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.AccountStatus,
		&acc.NationalID,
		&acc.BirthDate,
		&acc.FirstName,
		&acc.MiddleName,
		&acc.LastName,
	)
	if err != nil {
		return nil, err
	}
	return &types.AccountResponse{Account: acc}, nil
}

func (s *AccountService) Login(ctx *context.Context, in *types.LoginRequest) (*types.LoginResponse, error) {
	req := &pb.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	}
	authRes, err := s.authService.LoginUser(*ctx, req)
	if err != nil {
		return nil, err
	}
	var res types.LoginResponse
	res.AccessToken = authRes.AccessToken
	res.TokenType = authRes.TokenType
	res.ExpiresIn = authRes.ExpiresIn
	res.RefreshToken = authRes.RefreshToken
	res.AuthID = authRes.User.Id

	acc, err := s.FetchUserByAuthID(ctx, authRes.User.Id)
	if err != nil {
		return nil, err
	}

	// for auth context
	res.DisplayName = acc.Account.FirstName + " " + acc.Account.MiddleName + " " + acc.Account.LastName
	res.AccountId = acc.Account.ID

	return &res, nil
}

func (s *AccountService) UpdatePassword(ctx *context.Context, in *types.UpdatePasswordRequest) (*types.AccountResponse, error) {
	old_encrpyptedPassword, err := s.GetUserAuth(*ctx, in.AuthID)
	if err != nil {
		return nil, err
	}
	if !verifyPassword(old_encrpyptedPassword, in.OldPassword) {
		return nil, fmt.Errorf("old password is incorrect")
	}
	new_encryptedPassword, err := hashPassword(in.NewPassword)
	if err != nil {
		return nil, err
	}
	_, err = s.db.Exec(*ctx, "UPDATE auth.users SET encrypted_password = $1 WHERE id = $2", new_encryptedPassword, in.AuthID)
	if err != nil {
		return nil, err
	}
	res, _ := s.FetchUserByAuthID(ctx, in.AuthID)
	return res, nil
}

func (s *AccountService) UpdateUser(ctx *context.Context, in *types.UpdateAccountRequest) (*types.Account, error) {
	query := `
		UPDATE account SET first_name = $1, middle_name = $2, last_name = $3, 
		email = $4, phone_number = $5, address = $6 WHERE id = $7
	`
	_, err := s.db.Exec(*ctx, query, in.FirstName, in.MiddleName, in.LastName, in.Email, in.Phone, in.Address, in.AccountID)
	if err != nil {
		return nil, err
	}
	res, _ := s.FetchUserById(ctx, in.AccountID)

	return &res.Account, nil
}
func (s *AccountService) UpdateUserDetails(ctx *context.Context, in *types.UpdateAccountDetailsRequest) (*types.Account, error) {
	var query string
	var args []any
	const (
		UpdateTypeEmail   = "Email"
		UpdateTypePhone   = "Phone"
		UpdateTypeAddress = "Address"
	)
	switch in.Type {
	case UpdateTypeEmail:
		err := s.UpdateAuthEmail(*ctx, in.AccountID, in.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to update email in auth service: %v", err)
		}
		query = "UPDATE account SET email = $1 WHERE id = $2"
		args = []any{in.Email, in.AccountID}
	case UpdateTypePhone:
		query = "UPDATE account SET phone_number = $1 WHERE id = $2"
		args = []any{in.Phone, in.AccountID}
	case UpdateTypeAddress:
		query = "UPDATE account SET address = $1 WHERE id = $2"
		args = []any{in.Address, in.AccountID}
	default:
		return nil, fmt.Errorf("invalid update type: %s", in.Type)
	}

	_, err := s.db.Exec(*ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %v", err)
	}

	res, err := s.FetchUserById(ctx, in.AccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated account: %v", err)
	}

	return &res.Account, nil
}

func (s *AccountService) UpdateAccountStatus(c *context.Context, in *types.UpdateAccountStatusRequest) (*types.Account, error) {
	var query string
	const (
		DEACTIVATE_ACCOUNT = "DEACTIVATE"
		REACTIVATE_ACCOUNT = "ACTIVATE"
		SUSPEND_ACCOUNT    = "SUSPEND"
	)
	switch in.Type {
	case DEACTIVATE_ACCOUNT:
		query = `UPDATE account SET account_status = 'Closed' WHERE id = $1`
	case REACTIVATE_ACCOUNT:
		query = `UPDATE account set account_status = 'Active' WHERE id = $1`
	case SUSPEND_ACCOUNT:
		query = `UPDATE account set account_status = 'Suspended' WHERE id = $1`
	}
	_, err := s.db.Exec(*c, query, in.AccountID)
	if err != nil {
		return nil, err
	}
	res, err := s.FetchUserById(c, in.AccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated account: %v", err)
	}

	return &res.Account, nil
}

// Helper functions
func (s *AccountService) GetUserAuth(ctx context.Context, authID string) (string, error) {
	var encrypted_password string
	query := `SELECT encrypted_password FROM auth.users WHERE id = $1;`
	err := s.db.QueryRow(ctx, query, authID).Scan(&encrypted_password)
	if err != nil {
		return "", fmt.Errorf("error querying auth user: %v", err)
	}
	return encrypted_password, nil
}
func (s *AccountService) UpdateAuthEmail(ctx context.Context, id, email string) error {
	var auth_id string
	err := s.db.QueryRow(ctx, "SELECT auth_id from account WHERE id = $1", id).Scan(&auth_id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(ctx, "UPDATE auth.users SET email = $1 WHERE id = $2", email, id)
	if err != nil {
		return err
	}
	return nil
}
