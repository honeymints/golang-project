package data

import (
	"context"
	"database/sql"
	"time"
)

type ListModel struct {
	DB *sql.DB
}
type Lists struct {
	ID          int64
	User_ID     int64
	CreatedAt   time.Time
	Title       string
	Description string
}

func (m ListModel) Insert(list *Lists) error {
	query := `
	INSERT INTO lists (title, description, created_at, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`
	args := []interface{}{list.Title, list.Description, list.CreatedAt, list.User_ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "users_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&list.ID, &list.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "list_id_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
