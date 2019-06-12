package mappers

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"time"

	_ "github.com/lib/pq"
)

// ScheduleMapper type of ScheduleProvider's mapper
type ScheduleMapper struct {
	Employee entities.Employee
	Days     map[string]string // карта дней, 	  где key-дата,        value-промежуток времени, вида XXXX-XXXX
	Template map[string]string // карта шаблонов, где key-день недели, value-промежуток времени, вида XXXX-XXXX
}

// DaySeconds 24 часа = 86400 секунд
const (
	DaySeconds = 86400
)

// Init initialize employee days and temlates with id
func (m *ScheduleMapper) Init(dayRows *sql.Rows, templateRow *sql.Row) error {
	err := m.daysInit(dayRows)
	if err != nil {
		return fmt.Errorf("days error %s", err.Error())
	}
	err = m.templatesInit(templateRow)
	if err != nil {
		return fmt.Errorf("template error %s", err.Error())
	}

	return nil
}

// GetSchedule возвращает расписание на основе дней и шаблонов
func (m *ScheduleMapper) GetSchedule(dateStart time.Time, dateEnd time.Time) []entities.Day {
	schedule := make([]entities.Day, 0)

	for day := dateStart; day.Unix() <= dateEnd.Unix(); day = time.Unix(day.Unix()+DaySeconds, 64) {
		var newday entities.Day
		if curDay, exist := m.Days[day.Format("01.02.2006")]; exist {
			newday.Date = day.Format("01.02.2006")
			newday.Timerange = curDay
			schedule = append(schedule, newday)
		} else {
			curWeekday := day.Weekday().String()
			newday.Date = day.Format("01.02.2006")
			newday.Timerange = m.Template[curWeekday]

			schedule = append(schedule, newday)
		}
	}

	fmt.Print("\n", schedule, "\n")

	return schedule
}

// GetEmployee возвращает данные сотрудника
func (m *ScheduleMapper) GetEmployee(employeeRow *sql.Row) entities.GroupEmployee {
	var eid int64
	var lastname string
	var firstname string
	var middlename string

	employeeRow.Scan(&eid, &lastname, &firstname, &middlename)

	employee := entities.GroupEmployee{
		EID:        eid,
		Lastname:   lastname,
		Firstname:  firstname,
		Middlename: middlename,
	}

	return employee
}

// #TODO вынос в структуру
func (m *ScheduleMapper) daysInit(rows *sql.Rows) error {
	m.Days = make(map[string]string)
	for rows.Next() {
		var date time.Time
		var timerange string
		err := rows.Scan(&date, &timerange)
		if err != nil {
			return err
		}
		m.Days[date.Format("01.02.2006")] = timerange
	}

	return nil
}

// #TODO вынос в структуру
func (m *ScheduleMapper) templatesInit(row *sql.Row) error {
	m.Template = make(map[string]string)
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	var ranges [7]string
	err := row.Scan(&ranges[0], &ranges[1], &ranges[2], &ranges[3], &ranges[4], &ranges[5], &ranges[6])
	if err != nil {
		return err
	}
	for i, day := range weekdays {
		m.Template[day] = ranges[i]
	}

	return nil
}
