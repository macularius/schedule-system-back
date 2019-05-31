package controllers

import (
	"database/sql"
	"log"
	"myapp/app"
	"myapp/app/models/providers"

	"github.com/revel/revel"
)

// Schedule controller struct
type Schedule struct {
	*revel.Controller
	provider providers.ScheduleProvider
}

// func getScheduleById(id int) {

// }

// GetSchedule get schedule action
func (c Schedule) GetSchedule() revel.Result {
	eid := c.Params.Values["id"][0]
	// params := c.Params.Values

	db, err := sql.Open("postgres", app.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	c.provider.Init(eid, db)
	schedule := c.provider.GetSchedule()

	return c.RenderJSON(schedule)
}
