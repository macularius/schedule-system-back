package providers

import (
	"myapp/app/models/mappers"
	"time"
)

const (
	ReadableDateFormat = "Mon Jan 02 2019 00:00:00 GMT+0400 (GMT+04:00)"
)

type ScheduleProvider struct {
	ID     string
	Mapper mappers.ScheduleMapper
}

// Init initialize mapper by employee's id
func (p *ScheduleProvider) Init(ID string) {
	p.Init(ID)
}

// GetSchedule return days of schedule initializing employee
func (p *ScheduleProvider) GetSchedule(dateNumberStart, dateNumberEnd time.Time) {
	// days := p.Mapper.Employee.Days
	// templates := p.Mapper.Employee.Templates[0]

	// var result []entities.Day
}
