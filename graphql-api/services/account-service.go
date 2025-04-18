package services

// Use this for resolvers business logic
// Planning on putting helper functions here that can basically do CRUD to the DB

import (
	"context"
	pb "finnbank/common/grpc/auth"
	"finnbank/graphql-api/types"
	"fmt"
	"math/rand"
	"time"

	"finnbank/common/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountService struct {
	db *pgxpool.Pool
	l  *utils.Logger
}

func NewAccountService(db *pgxpool.Pool, logger *utils.Logger) *AccountService {
	return &AccountService{
		db: db,
		l:  logger,
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

func CreateUser(ctx *context.Context, in *types.AddAccountRequest, DB *pgx.Conn, authServer pb.AuthServiceClient) (*types.AddAccountResponse, error) {
	req := &pb.SignUpRequest{
		Email:    in.Email,
		Password: in.Password,
	}
	// TODO: This seems really bad, will have to find a better way for this somehow
	authRes, err := authServer.SignUpUser(*ctx, req)
	if err != nil {
		return nil, err
	}
	if in.AccountType != "Personal" && in.AccountType != "Business" {
		return nil, fmt.Errorf("account type must be either Personal or Business")
	}
	accNum := GenAccNum()
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
	err = DB.QueryRow(*ctx,
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

func FetchUserByAccountNumber(ctx *context.Context, req string, DB *pgx.Conn) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE account_number = $1
	`
	err := DB.QueryRow(*ctx, query, req).Scan(
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

func FetchUserByEmail(ctx *context.Context, req string, DB *pgx.Conn) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE email = $1
	`
	err := DB.QueryRow(*ctx, query, req).Scan(
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

func FetchUserByPhone(ctx *context.Context, req string, DB *pgx.Conn) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE phone_number = $1
	`
	err := DB.QueryRow(*ctx, query, req).Scan(
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
func FetchUserByAuthID(ctx *context.Context, req string, DB *pgx.Conn) (*types.AccountResponse, error) {
	var acc types.Account
	query := `
		SELECT * FROM account WHERE auth_id = $1
	`
	err := DB.QueryRow(*ctx, query, req).Scan(
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
func Login(ctx *context.Context, in *types.LoginRequest, authServer pb.AuthServiceClient) (*types.LoginResponse, error) {
	req := &pb.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	}
	authRes, err := authServer.LoginUser(*ctx, req)
	if err != nil {
		return nil, err
	}
	var res types.LoginResponse
	res.AccessToken = authRes.AccessToken
	res.TokenType = authRes.TokenType
	res.ExpiresIn = authRes.ExpiresIn
	res.RefreshToken = authRes.RefreshToken
	res.AuthID = authRes.User.Id
	res.Email = authRes.User.Email

	return &res, nil
}
