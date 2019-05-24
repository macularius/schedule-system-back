package schedule

import (
	"myapp/app/models/entities"
	"time"
)

// ScheduleMapper type of ScheduleProvider's mapper
type ScheduleMapper struct {
	Employee entities.Employee
}

// Init initialize employee days and temlates with id
func (m *ScheduleMapper) Init(id string) {
	switch id {
	case "0":
		m.Employee.Name = "Коваценко Игорь Николаевич"

		m.Employee.GroupID = "0"

		m.Employee.Days = []entities.Day{
			{Date: time.Parse(time.RFC3339, "2019-05-26T20:00:00.000Z"), RangeStart: "16", RangeEnd: "18:30"},
			{Date: "2019-05-19T20:00:00.000Z", RangeStart: "16", RangeEnd: "18:30"},
		}
		m.Employee.Templates = []entities.Template{
			{Mon: "14 - 18", Tue: "14 - 18", Wed: "9 - 18", Thu: "", Fri: "9 - 18", Sat: "", Sun: ""},
		}
	case "1":
		m.Employee.Name = "Иванов Иван Иванович"

		m.Employee.GroupID = "1"
		m.Employee.GroupName = "Группа 1"

		m.Employee.Days = []entities.Day{
			{Date: "2019-05-23T20:00:00.000Z", RangeStart: "12", RangeEnd: "18"},
			{Date: "2019-05-19T20:00:00.000Z", RangeStart: "8", RangeEnd: "17"},
		}
		m.Employee.Templates = []entities.Template{
			{Mon: "9 - 18", Tue: "9 - 18", Wed: "9 - 18", Thu: "9 - 18", Fri: "9 - 18", Sat: "", Sun: ""},
		}
	}
}
