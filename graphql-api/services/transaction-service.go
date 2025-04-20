package services

// Use this for resolvers business logic

// GetTransactionByUserId, (Query)
// GetTransactionByTimeStampByUserId (Query) 
// CreateTransactionByUserId (Mutation)

import (
	"finnbank/common/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)


type TransactionService struct {
	db     *pgxpool.Pool
	l *utils.Logger
}

func NewTransactionService(db *pgxpool.Pool, logger *utils.Logger) *AccountService {
	return &AccountService{
		db:     db,
		l: logger,
	}
}
