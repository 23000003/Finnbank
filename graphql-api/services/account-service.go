package services

import (
	"context"
	pb "finnbank/common/grpc/auth"
	"finnbank/graphql-api/types"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

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
