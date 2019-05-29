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
func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", user, password, database)
}

// GetNewConnectionString получение новой строки подключения
func GetNewConnectionString(username string) (string, error) {
	var err error
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", user, password, database))
	if err != nil {
		// log.Fatal("Error creating connection: ", err.Error())
		return "Error creating connection: ", err
	}
	defer db.Close()

	// Проверка существования роли
	roles, err := db.Query(fmt.Sprintf("SELECT rolname FROM pg_roles WHERE rolname = '%s';", username))
	if err != nil {
		// log.Fatal("Error creating role: ", err.Error())
		return "Error creating role: ", err
	}
	defer roles.Close()
	var role string
	if roles.Next() {
		err = roles.Scan(&role)
	}
	if err != nil {
		return "Error scanning role: ", err
	}

	if role == "" {
		_, err = db.Query(fmt.Sprintf("CREATE ROLE \"%s\" LOGIN;", username))
		if err != nil {
			// log.Fatal("Error creating role: ", err.Error())
			return "Error creating role: ", err
		}
	}

	return fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", username, password, database), nil
}
