package link

import (
	"database/sql"
)

type Link struct {
	ID        int          `json:"id"`
	UserID    uint32       `json:"userId"`
	URL       string       `json:"url"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type ModelDeletedError struct {
	msg string
}

func (e ModelDeletedError) Error() string { return e.msg }
