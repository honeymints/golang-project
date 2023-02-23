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

func (m ListModel) DeleteByID(ID int64) error {
	query := `
		DELETE FROM lists
		WHERE id=$1
	`
	result, err := m.DB.Exec(query, ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m ListModel) UpdateByID(ID int64, title string) error {
	query := `
UPDATE lists
SET title = $1
WHERE id = $2
RETURNING version`
	args := []interface{}{title, ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan()
	if err != nil {
		switch {
		default:
			return err
		}
	}
	return nil
}
