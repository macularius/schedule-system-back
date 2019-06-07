package entities

import (
	"database/sql"
	"time"
)

// Session структура сессии пользователя
type Session struct {
	UserID     int
	EmployeeID int
	Connection *sql.DB
	Login      string
	Created    time.Time
	Expiration int64
}
