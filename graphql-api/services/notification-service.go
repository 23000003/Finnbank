package services

// Use this for resolvers business logic

// GetAllNotificationByUserId, (Query)
// GetNotificationByUserId, (Query)
// GenerateNotification, (Mutation)
// ReadNotificationByUserId (Mutation)

import (
	"context"
	"finnbank/common/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationService struct {
	db *pgxpool.Pool
	l  *utils.Logger
}

func NewNotificationService(db *pgxpool.Pool, logger *utils.Logger) *NotificationService {
	return &NotificationService{
		db: db,
		l:  logger,
	}
}

type Notification struct {
	NotifID       string     `json:"notif_id"`
	NotifType     string     `json:"notif_type"`
	AuthID        string     `json:"auth_id"`
	NotifToID     string     `json:"notif_to_id"`
	NotifFromName string     `json:"notif_from_name"`
	Content       string     `json:"content"`
	IsRead        bool       `json:"is_read"`
	RedirectURL   string     `json:"redirect_url"`
	DateNotified  time.Time  `json:"date_notified"`
	DateRead      *time.Time `json:"date_read"` // pointer for nullable
}

func (s *NotificationService) GetAllNotificationByUserId(notifToID string) ([]Notification, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT notif_id, notif_type, auth_id, notif_to_id, notif_from_name,
		       content, is_read, redirect_url, date_notified, date_read
		FROM notifications
		WHERE notif_to_id = $1
		ORDER BY date_notified DESC
	`, notifToID)
	if err != nil {
		s.l.Error("DB query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notif Notification
		err := rows.Scan(
			&notif.NotifID, &notif.NotifType, &notif.AuthID, &notif.NotifToID,
			&notif.NotifFromName, &notif.Content, &notif.IsRead,
			&notif.RedirectURL, &notif.DateNotified, &notif.DateRead,
		)
		if err != nil {
			s.l.Error("Scan failed: %v", err)
			continue
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}
