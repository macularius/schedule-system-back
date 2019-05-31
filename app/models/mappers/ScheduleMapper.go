package mappers

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"time"
)

// ScheduleMapper type of ScheduleProvider's mapper
type ScheduleMapper struct {
	Employee entities.Employee
	Days     map[string]string // карта дней, 	  где key-дата,        value-промежуток времени, вида XXXX-XXXX
	Template map[string]string // карта шаблонов, где key-день недели, value-промежуток времени, вида XXXX-XXXX
}

// Init initialize employee days and temlates with id
func (m *ScheduleMapper) Init(dayRows *sql.Rows, templateRow *sql.Row) error {
	err := m.daysInit(dayRows)
	if err != nil {
		return err
	}
	err = m.templatesInit(templateRow)
	if err != nil {
		return err
	}

	return nil
}

// GetSchedule возвращает расписание на основе дней и шаблонов
func (m *ScheduleMapper) GetSchedule() map[string]string {
	dateStart := time.Now()
	dateEnd := time.Unix(time.Now().Unix()+2592000, 64) // 30-ый день от текущего

	return m.GetScheduleByRange(dateStart, dateEnd)
}

// GetScheduleByRange возвращает расписание на основе дней и шаблонов
func (m *ScheduleMapper) GetScheduleByRange(dateStart time.Time, dateEnd time.Time) map[string]string {
	schedule := make(map[string]string)
	for day := dateStart; day.Unix() <= dateEnd.Unix(); day = time.Unix(day.Unix()+86400, 64) { // 24 часа = 86400 секунд
		fmt.Println("\n*" + day.Format("2006-01-02") + "*\n")
		if curDay, exist := m.Days[day.Format("2006-01-02")]; exist {
			schedule[day.Format("2006-01-02")] = curDay
		} else {
			curWeekday := day.Weekday().String()
			schedule[day.Format("2006-01-02")] = m.Template[curWeekday]
		}
	}

	return schedule
}

func (m *ScheduleMapper) daysInit(rows *sql.Rows) error {
	m.Days = make(map[string]string)
	fmt.Println("\n*** Days start")
	for rows.Next() {
		var date time.Time
		var timerange string
		err := rows.Scan(&date, &timerange)
		if err != nil {
			return err
		}
		m.Days[date.Format("2006-01-02")] = timerange
		fmt.Println(date.Format("2006-01-02") + " - " + m.Days[date.Format("2006-01-02")])
	}
	fmt.Println("\n*** Days end")

	return nil
}

func (m *ScheduleMapper) templatesInit(row *sql.Row) error {
	m.Template = make(map[string]string)
	weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	fmt.Println("\n*** Templates start")
	var ranges [7]string
	err := row.Scan(&ranges[0], &ranges[1], &ranges[2], &ranges[3], &ranges[4], &ranges[5], &ranges[6])
	if err != nil {
		fmt.Println("ERROR: \n" + err.Error() + "\n")
		return err
	}
	for i, day := range weekdays {
		m.Template[day] = ranges[i]
		fmt.Println(m.Template[day])
	}
	fmt.Println("\n*** Templates end")

	return nil
}
