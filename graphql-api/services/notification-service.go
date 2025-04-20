package services

// Use this for resolvers business logic

// GetAllNotificationByUserId, (Query)
// GetNotificationByUserId, (Query)
// GenerateNotification, (Mutation)
// ReadNotificationByUserId (Mutation)

import (
	"finnbank/common/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)


type NotificationService struct {
	db     *pgxpool.Pool
	l *utils.Logger
}

func NewNotificationService(db *pgxpool.Pool, logger *utils.Logger) *AccountService {
	return &AccountService{
		db:     db,
		l: logger,
	}
}
