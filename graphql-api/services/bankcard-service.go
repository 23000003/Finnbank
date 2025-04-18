package services

// Use this for resolvers business logic

// GetBankCardOfUserById, (Query)
// CreateBankCardForUser,  (Mutation)
// UpdateBankcardExpiryDateByUserId  (Mutation)

import (
	"finnbank/common/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)


type BankcardService struct {
	db     *pgxpool.Pool
	l *utils.Logger
}

func NewBankcardService(db *pgxpool.Pool, logger *utils.Logger) *AccountService {
	return &AccountService{
		db:     db,
		l: logger,
	}
}
