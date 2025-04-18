package services

// Use this for resolvers business logic

// GenerateStatementForUserByTimestamp // Query (calls to transaction-service)

import (
	"finnbank/common/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)


type StatementService struct {
	db     *pgxpool.Pool
	l *utils.Logger
}

func NewStatementService(db *pgxpool.Pool, logger *utils.Logger) *AccountService {
	return &AccountService{
		db:     db,
		l: logger,
	}
}
