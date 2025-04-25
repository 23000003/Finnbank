package types

import "time"


type Notification struct {
	NotifID       string     `json:"notif_id"`
	NotifType     string     `json:"notif_type"`
	NotifToID     string     `json:"notif_to_id"`
	NotifFromName string     `json:"notif_from_name"`
	Content       string     `json:"content"`
	IsRead        bool       `json:"is_read"`
	RedirectURL   string     `json:"redirect_url"`
	DateNotified  time.Time  `json:"date_notified"`
	DateRead      *time.Time `json:"date_read"` // pointer for nullable
}