package app

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DBConfig получение config'ов подключения к db
type DBConfig struct {
}

const (
	server   = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	database = "schedule_db"
	sslmode  = "disable"
)

// GetConnectionString возвращает string сконфигурированную строку подключения к Postgre
func GetConnectionString() string {
	return fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", user, password, database)
}

// GetNewConnectionString получение новой строки подключения
func GetNewConnectionString(username string) (string, error) {
	var err error
	var rolename string

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", user, password, database))
	if err != nil {
		// log.Fatal("Error creating connection: ", err.Error())
		return "", err
	}
	defer db.Close()

	// Проверка существования роли
	roles, err := db.Query(fmt.Sprintf("SELECT rolname FROM pg_roles WHERE rolname = '%s';", username))
	if err != nil {
		// log.Fatal("Error creating role: ", err.Error())
		return "", err
	}
	defer roles.Close()
	if roles.Next() {
		err = roles.Scan(&rolename)
	}
	if err != nil {
		return "", err
	}

	if rolename == "" {
		_, err = db.Query(fmt.Sprintf("create role \"%s\" login;", username))
		if err != nil {
			// log.Fatal("Error creating role: ", err.Error())
			return "", err
		}
	}

	return fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", username, password, database), nil
}
