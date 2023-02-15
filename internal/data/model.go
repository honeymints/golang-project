package data

import (
	"database/sql"
	"errors"
)

type Models struct {
	Users UserModel // Add a new Users field.
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db}, // Initialize a new UserModel instance.
	}
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)
