package services

// Use this for resolvers business logic

// GetBankCardOfUserById, (Query)
// CreateBankCardForUser,  (Mutation)
// UpdateBankcardExpiryDateByUserId  (Mutation)

import (
	"context"
	"finnbank/common/utils"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AllBankCardRequests struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CardType  string `json:"card_type"`
}

type Requester struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CardType  string `json:"card_type"`
	BirthDate string `json:"birth_date"`
}

type BankcardService struct {
	db *pgxpool.Pool
	l  *utils.Logger
}

func NewBankcardService(db *pgxpool.Pool, logger *utils.Logger) *BankcardService {
	return &BankcardService{
		db: db,
		l:  logger,
	}
}

func (b *BankcardService) CreateCardRequest(ctx context.Context, req Requester) error {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
        INSERT INTO card_request (first_name, last_name, card_type)
        VALUES ($1, $2, $3)
    `, req.FirstName, req.LastName, req.CardType)

	if err != nil {
		return fmt.Errorf("failed to insert card request: %w", err)
	}

	return nil
}

func (b *BankcardService) GetAllBankCardRequestsById(ctx context.Context, id int) ([]AllBankCardRequests, error) {

	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`select first_name, last_name, card_type from card_request where request_id = $1`,
		id,
	)
	defer rows.Close()

	var results []AllBankCardRequests
	for rows.Next() {
		var acc AllBankCardRequests
		if err := rows.Scan(
			&acc.FirstName,
			&acc.LastName,
			&acc.CardType,
		); err != nil {
			return nil, err
		}
		results = append(results, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
