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
	ID        int64
	User_ID   int64
	CreatedAt time.Time
	Title     string
}

func (m ListModel) Insert(list *Lists) error {
	query := `
  INSERT INTO lists (title, created_at, user_id)
  VALUES ($1, $2, $3)
  RETURNING id, created_at`
	args := []interface{}{list.Title, list.CreatedAt, list.User_ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

/* func (m ListModel) GetByUser(ID int64) (*Lists, error) {
	query := `SELECT id, user_id, created_at, title
	from lists
	WHERE user_id=$1`
	var list Lists
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, ID).Scan(
		&list.ID,
		&list.CreatedAt,
		&list.Title,
		&list.User_ID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.

	return &list, nil
} */

func (m ListModel) GetByUser(userID int64) ([]*Lists, error) {
	query := `
		SELECT id, user_id, created_at, title
		FROM lists
		WHERE user_id=$1
		ORDER BY created_at DESC
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query and return an error if there is one.
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	// Close the rows object when the function returns.
	defer rows.Close()

	// Iterate through the rows and scan each row into a Lists struct.
	var lists []*Lists
	for rows.Next() {
		list := &Lists{}
		err := rows.Scan(
			&list.ID,
			&list.User_ID,
			&list.CreatedAt,
			&list.Title,
		)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	// If there was an error while iterating through the rows, return it.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the Lists slice.
	return lists, nil
}
