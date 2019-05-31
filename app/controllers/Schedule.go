package controllers

import (
	"database/sql"
	"log"
	"myapp/app"
	"myapp/app/models/entities"
	"myapp/app/models/providers"
	"time"

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
	var schedule []entities.Day
	if c.Params.Values["start"] != nil && c.Params.Values["start"][0] != "" && c.Params.Values["end"] != nil && c.Params.Values["end"][0] != "" {
		var start time.Time
		var end time.Time

		start, err = time.Parse("02.01.2006", c.Params.Values["start"][0])
		if err != nil {
			log.Fatal(err)
		}
		end, err = time.Parse("02.01.2006", c.Params.Values["end"][0])
		if err != nil {
			log.Fatal(err)
		}

		schedule = c.provider.GetScheduleByRange(start, end)
	} else {
		schedule = c.provider.GetSchedule()
	}

	return c.RenderJSON(schedule)
}
