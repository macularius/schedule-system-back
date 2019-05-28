package app

import (
	"database/sql"
	"fmt"
	"log"

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
func GetNewConnectionString(username string) string {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", user, password, database))
	if err != nil {
		log.Fatal("Error creating connection: ", err.Error())
	}
	defer db.Close()

	rows, err := db.Query("CREATE ROLE " + username)
	if err != nil {
		log.Fatal("Error creating role: ", err.Error())
	}
	defer rows.Close()

	return fmt.Sprintf("user=%s password=%s database=%s sslmode=disable", username, password, database)
}
