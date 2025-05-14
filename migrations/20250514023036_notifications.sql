-- +goose Up
CREATE TABLE notifications (
  notif_id UUID PRIMARY KEY,
  notif_type TEXT NOT NULL,
  notif_to_id VARCHAR NOT NULL,
  notif_from_name VARCHAR NOT NULL,
  content TEXT NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  redirect_url TEXT,
  date_notified TIMESTAMPTZ DEFAULT now(),
  date_read TIMESTAMPTZ
);

-- +goose Down
DROP TABLE notifications;