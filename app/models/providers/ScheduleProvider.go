package providers

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"myapp/app/models/mappers"
	"time"
)

const (
	readableDateFormat = "Mon Jan 02 2019 00:00:00 GMT+0400 (GMT+04:00)"
)

// ScheduleProvider structere of auth provider
type ScheduleProvider struct {
	EID    string
	DB     *sql.DB
	mapper mappers.ScheduleMapper
}

// Init initialize mapper by user's id
func (p *ScheduleProvider) Init(eid string, db *sql.DB) error {
	fmt.Println("EID = " + eid)
	p.EID = eid

	templatesRow := db.QueryRow(selectTemplatesQueryString(eid))
	daysRows, err := db.Query(selectDaysQueryString(eid))
	if err != nil {
		return err
	}
	defer daysRows.Close()
	p.mapper.Init(daysRows, templatesRow)

	return nil
}

// GetSchedule return days of schedule initializing employee
func (p *ScheduleProvider) GetSchedule() []entities.Day {
	return p.mapper.GetSchedule()
}

// GetScheduleByRange return days of schedule initializing employee
func (p *ScheduleProvider) GetScheduleByRange(dateNumberStart time.Time, dateNumberEnd time.Time) []entities.Day {
	return p.mapper.GetScheduleByRange(dateNumberStart, dateNumberEnd)
}

func selectTemplatesQueryString(eid string) string {
	return fmt.Sprintf("SELECT mon, tue, wed, thu, fri, sat, sun FROM templates WHERE eid='%s';", eid)
}
func selectDaysQueryString(eid string) string {
	return fmt.Sprintf("SELECT date, range FROM schedules WHERE eid='%s'", eid) // #TODO отрезать дни, которые прошли
}
func selectDaysQueryStringByRange(eid string, start time.Time, end time.Time) string {
	return fmt.Sprintf("SELECT date, range FROM schedules WHERE eid='%s' AND date >= '%s' AND date <= '%s'", eid, start, end)
}
