package controllers

import (
	"github.com/revel/revel"
)

// Schedule controller struct
type Schedule struct {
	*revel.Controller
}

func getScheduleById(id int) {

}

// GetSchedule get schedule action
func (c Schedule) GetSchedule() revel.Result {
	schedule := c.Params.Values
	return c.RenderJSON(schedule)
}
