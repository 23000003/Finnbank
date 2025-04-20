package services

import (
	"context"
	pb "finnbank/common/grpc/auth"
	"finnbank/graphql-api/types"
	"fmt"
	"math/rand"
	"time"

	"finnbank/common/utils"

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

// This just generates a random 16 digit account number
// this works for now until i come up with a better solution
func GenAccNum() string {
	rand.Seed(time.Now().UnixNano())
	var accNum string
	for i := 0; i < 16; i++ {
		accNum += string(rune('0' + rand.Intn(10)))
	}
	return accNum
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
	
	accNum := GenAccNum()

	// TODO: This seems really bad, will have to find a better way for this somehow
	authRes, err := s.authService.SignUpUser(*ctx, req)
	if err != nil {
		return nil, err
	}
	if in.AccountType != "Personal" && in.AccountType != "Business" {
		return nil, fmt.Errorf("account type must be either Personal or Business")
	}
	var res types.AddAccountResponse
	createQuery := `
	INSERT INTO account (
		email, full_name, phone_number, address, nationality,
		account_type, account_number, has_card, balance, auth_id
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING 
		id, email, full_name, phone_number, address, nationality,
		account_type, account_number, has_card, balance, date_created, auth_id
	`

	s.l.Info("Creating account for user: %s", authRes)
	s.l.Info("Account number: %s", accNum)

	err = s.db.QueryRow(*ctx,
		createQuery,
		in.Email,
		in.FullName,
		in.PhoneNumber,
		in.Address,
		in.Nationality,
		in.AccountType,
		accNum,
		false, 0.00,
		authRes.User.Id).
		Scan(
			&res.ID,
			&res.Email,
			&res.FullName,
			&res.PhoneNumber,
			&res.Address,
			&res.Nationality,
			&res.AccountType,
			&res.AccountNumber,
			&res.HasCard,
			&res.Balance,
			&res.DateCreated,
			&res.AuthID,
		)
	if err != nil {
		return nil, err
	}

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
		&acc.FullName,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
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
		&acc.FullName,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
	)
	if err != nil {
		return nil, err
	}
	return &types.AccountResponse{Account: acc}, nil
}

func (s *AccountService) FetchUserByPhone(ctx *context.Context, req string) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE phone_number = $1
	`
	err := s.db.QueryRow(*ctx, query, req).Scan(
		&acc.ID,
		&acc.AuthID,
		&acc.Email,
		&acc.FullName,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
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
		&acc.FullName,
		&acc.PhoneNumber,
		&acc.HasCard,
		&acc.AccountNumber,
		&acc.Address,
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
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

	acc , err := s.FetchUserByAuthID(ctx, authRes.User.Id)
	if err != nil {
		return nil, err
	}

	// for auth context
	res.FullName = acc.Account.FullName
	res.AccountId = acc.Account.ID

	return &res, nil
}

func (s *AccountService) UpdatePassword(ctx *context.Context, in *types.UpdatePasswordRequest) (*types.AccountResponse, error) {
	req := &pb.UpdatePasswordRequest{
		AuthID:      in.AuthID,
		NewPassword: in.NewPassword,
		OldPassword: in.OldPassword,
	}
	authRes, err := s.authService.HashAndEncryptPassowrd(*ctx, req)
	if err != nil {
		return nil, err
	}
	updateQuery := `
		UPDATE auth.users SET encrypted_password = $1 WHERE id = $2
	`
	_, err = s.db.Exec(*ctx, updateQuery, authRes.EncryptedPassword, in.AuthID)
	if err != nil {
		return nil, err
	}
	// Just calling this to fetch the account again, we can remove this later, but for now im using this
	// in testing
	res, _ := s.FetchUserByAuthID(ctx, in.AuthID)
	return res, nil
}
