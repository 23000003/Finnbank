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
	"golang.org/x/crypto/bcrypt"
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
	var acc types.Account
	var res types.AddAccountResponse
	createQuery := `
	INSERT INTO account (
		email, first_name, middle_name, last_name, phone_number, address, nationality,
		account_type, account_number, has_card, balance, auth_id, birthdate, national_id
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING 
		id, email, first_name, middle_name, last_name, phone_number, address, nationality,
		account_type, account_number, has_card, balance, date_created, date_updated, auth_id, birthdate, national_id, account_status
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
		false, 0.00,
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
			&acc.Balance,
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
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.NationalID,
		&acc.AccountStatus,
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
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.NationalID,
		&acc.AccountStatus,
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
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.NationalID,
		&acc.AccountStatus,
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
		&acc.Balance,
		&acc.AccountType,
		&acc.DateCreated,
		&acc.DateUpdated,
		&acc.Nationality,
		&acc.NationalID,
		&acc.AccountStatus,
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
	if !VerifyPassword(old_encrpyptedPassword, in.OldPassword) {
		return nil, fmt.Errorf("old password is incorrect")
	}
	new_encryptedPassword, err := HashPassword(in.NewPassword)
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

// TODO: This is not implemented yet, but i think it should be done by tommorow
// func (s* AccountService) UpdateUser(ctx* context.Context, in * types.UpdateAccountRequest)(*types.UpdateAccountResponse, error) {

// }

// This too i guess
func (s *AccountService) GetUserAuth(ctx context.Context, authID string) (string, error) {
	var encrypted_password string
	query := `SELECT encrypted_password FROM auth.users WHERE id = $1;`
	err := s.db.QueryRow(ctx, query, authID).Scan(&encrypted_password)
	if err != nil {
		return "", fmt.Errorf("error querying auth user: %v", err)
	}
	return encrypted_password, nil
}

// These are just helper functions
func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
func HashPassword(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
