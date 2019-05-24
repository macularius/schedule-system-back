package mappers

import "myapp/app/models/entities"

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
			{Date: "Mon May 27 2019 00:00:00 GMT+0400 (GMT+04:00)", RangeStart: "16", RangeEnd: "18:30"},
			{Date: "Mon May 20 2019 00:00:00 GMT+0400 (GMT+04:00)", RangeStart: "16", RangeEnd: "18:30"},
		}
		m.Employee.Templates = []entities.Template{
			{Mon: "14 - 18", Tue: "14 - 18", Wed: "9 - 18", Thu: "", Fri: "9 - 18", Sat: "", Sun: ""},
			{Mon: "14 - 18", Tue: "14 - 18", Wed: "9 - 18", Thu: "", Fri: "9 - 18", Sat: "", Sun: ""},
		}
	case "1":
		m.Employee.Name = "Иванов Иван Иванович"
		m.Employee.GroupID = "1"
		m.Employee.GroupName = "Группа 1"
		m.Employee.Days = []entities.Day{
			{Date: "Mon May 24 2019 00:00:00 GMT+0400 (GMT+04:00)", RangeStart: "12", RangeEnd: "18"},
			{Date: "Mon May 20 2019 00:00:00 GMT+0400 (GMT+04:00)", RangeStart: "8", RangeEnd: "17"},
		}
		m.Employee.Templates = []entities.Template{
			{Mon: "9 - 18", Tue: "9 - 18", Wed: "9 - 18", Thu: "9 - 18", Fri: "9 - 18", Sat: "", Sun: ""},
			{Mon: "9 - 18", Tue: "9 - 18", Wed: "9 - 18", Thu: "9 - 18", Fri: "9 - 18", Sat: "", Sun: ""},
		}

	}
}
