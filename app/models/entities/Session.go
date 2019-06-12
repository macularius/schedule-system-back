package entities

import (
	"database/sql"
	"time"
)

// Session структура сессии пользователя
type Session struct {
	Token      string
	UserID     int64
	EmployeeID int64
	Connection *sql.DB
	Login      string
	Created    time.Time
	Expiration int64
}
