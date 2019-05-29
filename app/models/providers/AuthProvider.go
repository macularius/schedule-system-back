package providers

import (
	"database/sql"
	"fmt"
	"log"
	"myapp/app"
	"myapp/app/models/mappers"
)

// AuthProvider structure of auth provider
type AuthProvider struct {
	Mapper *mappers.AuthMapper
}

// Authentication authentication
func (m *AuthProvider) Authentication(login string, password string) (int, bool, error) {
	db, err := sql.Open("postgres", app.GetConnectionString())
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT uid FROM users WHERE login='%s' AND password='%s'", login, password))
	if err != nil {
		return 0, false, err
	}
	defer rows.Close()

	if rows.Next() {
		var uid int
		err = rows.Scan(&uid)
		if err != nil {
			return 0, false, err
		}
		return uid, true, nil
	}
	return 0, false, nil
}
