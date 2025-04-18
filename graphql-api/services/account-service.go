package services

// Use this for resolvers business logic
// Planning on putting helper functions here that can basically do CRUD to the DB

import (
	"finnbank/common/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)


type AccountService struct {
	db     *pgxpool.Pool
	l *utils.Logger
}

func NewAccountService(db *pgxpool.Pool, logger *utils.Logger) *AccountService {
	return &AccountService{
		db:     db,
		l: logger,
	}
}
