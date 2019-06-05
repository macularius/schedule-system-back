package controllers

import (
	"database/sql"
	"errors"
	"myapp/app"
	"myapp/app/models/providers"
	"time"

	"github.com/revel/revel"
)

// Schedule controller struct
type Schedule struct {
	*revel.Controller
	provider *providers.ScheduleProvider
}

// GetSchedule get schedule action
func (c Schedule) GetSchedule() revel.Result {

	eid, err := c.getEIDByParams()
	if err != nil {
		return c.RenderJSON(Failed(err))
	}
	// #TODO db брать из сессии
	db, err := sql.Open("postgres", app.GetConnectionString()) // #TODO пользовательское подключение к базе
	if err != nil {
		return c.RenderJSON(Failed(err))
	}
	defer db.Close()

	if c.provider == nil {
		c.provider = new(providers.ScheduleProvider)
	}
	c.provider.Init(eid, db)

	// Инициализация границ временного промежутка
	start, end, err := c.getRangeByParams()
	if err != nil {
		return c.RenderJSON(Failed(err))
	}

	schedule := c.provider.GetSchedule(start, end)
	return c.RenderJSON(Succes(schedule))
}

// getRangeByParams возвращает левую и правую границу временного промежутка из get параметров
func (c *Schedule) getRangeByParams() (time.Time, time.Time, error) {
	var start time.Time
	var end time.Time

	if c.Params.Values["start"] != nil && c.Params.Values["start"][0] != "" {
		start, err := time.Parse("02.01.2006", c.Params.Values["start"][0])
		if err != nil {
			return start, end, err
		}

		// проверка существования правого ограничения daterange'а
		if !(c.Params.Values["end"] != nil && c.Params.Values["end"][0] != "") {
			end = start
		} else {
			end, err = time.Parse("02.01.2006", c.Params.Values["end"][0])
			if err != nil {
				return start, end, err
			}
		}
	}

	return start, end, nil
}

// getEIDByParams возвращает id сотрудника из get параметров
func (c *Schedule) getEIDByParams() (string, error) {
	if c.Params.Values["id"] != nil && c.Params.Values["id"][0] != "" {
		return c.Params.Values["id"][0], nil
	}

	return "", errors.New("GET параметра 'eid' не существует")
}
