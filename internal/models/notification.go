package models

type CustomerNotification struct {
	ID              int64  `db:"id"`
	CustomerID      string `db:"customer_id"`
	EmailNotifyDays bool   `db:"email_notify_days"`
	CreatedAt       string `db:"created_at"`
	UpdatedAt       string `db:"updated_at"`
}
