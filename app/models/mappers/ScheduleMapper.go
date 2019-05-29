package mappers

import (
	"database/sql"
	"log"
	"myapp/app"
	"myapp/app/models/entities"
)

// ScheduleMapper type of ScheduleProvider's mapper
type ScheduleMapper struct {
	Employee entities.Employee
}

// Init initialize employee days and temlates with id
func (m *ScheduleMapper) Init(id string) {
	connstr := app.GetConnectionString()

	_, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

}
