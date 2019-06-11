package providers

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"myapp/app/models/mappers"
	"time"

	_ "github.com/lib/pq"
)

// Days30Seconds 30 дней = 2592000 секунд
const (
	readableDateFormat = "Mon Jan 02 2019 00:00:00 GMT+0400 (GMT+04:00)"
	Days30Seconds      = 2592000
)

// ScheduleProvider structere of auth provider
type ScheduleProvider struct {
	DB     *sql.DB
	mapper *mappers.ScheduleMapper
}

// Init initialize mapper by user's id
func (p *ScheduleProvider) Init(eid string, db *sql.DB) error {
	p.mapper = new(mappers.ScheduleMapper)

	templatesRow := db.QueryRow(selectTemplatesQueryString(eid))
	// fmt.Print("\ndays\n", selectDaysQueryString(eid), "\n")
	dayRows, err := db.Query(selectDaysQueryString(eid))
	if err != nil {
		return fmt.Errorf("dayRows error %s", err.Error())
	}
	defer dayRows.Close()
	err = p.mapper.Init(dayRows, templatesRow)
	if err != nil {
		return err
	}

	return nil
}

// GetSchedule return days of schedule initializing employee
func (p *ScheduleProvider) GetSchedule(dateNumberStart time.Time, dateNumberEnd time.Time) []entities.Day {

	fmt.Printf("\nData after\n%s\n%s\n", dateNumberEnd.String(), dateNumberStart.String())

	// Если левая граница временного промежутка отсутствует, то вернуть 30 дней от текущего дня
	// Иначе, если правая отсутствует то вернуть расписание одного дня
	if dateNumberStart.IsZero() {
		dateNumberStart = time.Now()
		dateNumberEnd = time.Unix(time.Now().Unix()+Days30Seconds, 64) // 30-ый день от текущего
	} else {
		if dateNumberEnd.IsZero() {
			dateNumberEnd = dateNumberStart
		}
	}

	fmt.Printf("\nData before\n%s\n%s\n", dateNumberEnd.String(), dateNumberStart.String())
	return p.mapper.GetSchedule(dateNumberStart, dateNumberEnd)
}

func selectTemplatesQueryString(eid string) string {
	return fmt.Sprintf("SELECT mon, tue, wed, thu, fri, sat, sun FROM templates WHERE eid='%s';", eid)
}
func selectDaysQueryString(eid string) string {
	return fmt.Sprintf("SELECT date, range FROM schedules WHERE eid='%s';", eid) // #TODO отрезать дни, которые прошли
}
func selectDaysQueryStringByRange(eid string, start time.Time, end time.Time) string {
	return fmt.Sprintf("SELECT date, range FROM schedules WHERE eid='%s' AND date >= '%s' AND date <= '%s';", eid, start, end)
}
