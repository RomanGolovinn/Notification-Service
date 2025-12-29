package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"Notification-Service/internal/service"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, notification *service.Notification) error {
	query := `
		INSERT INTO notifications (type, address, message, created_at, scheduled_at, is_sended)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query, notification.Type, notification.Address,
		notification.Message, notification.CreatedAt).Scan(&notification.ID)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetPending(ctx context.Context) ([]service.Notification, error) {
	//will be called once every time to check whether it is time to send a notification
	query := `SELECT id, type, address, message FROM notifications 
              WHERE is_sent = false AND scheduled_at <= NOW() LIMIT 100`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []service.Notification
	for rows.Next() {
		var n service.Notification
		rows.Scan(&n.ID, &n.Type, &n.Address, &n.Message)
		list = append(list, n)
	}
	return list, nil
}

func (r *PostgresRepository) MarkAsSent(ctx context.Context, id int64) error {
	now := time.Now()
	query := `
        UPDATE notifications 
        SET is_sended = true, sent_at = $1 
        WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, now, id)
	return err
}
