package mappers

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"strconv"
)

// MetadataMapper type of ScheduleProvider's mapper
type MetadataMapper struct {
}

// GetMenuMeta возвращает массив групп для меню
func (m *MetadataMapper) GetMenuMeta(empRows *sql.Rows, groupRows *sql.Rows) ([]*entities.Group, error) {
	groupsMap := make(map[string]*entities.Group, 0)
	groups := make([]*entities.Group, 0)

	userEmployeeGroup := &entities.Group{
		GID:       -1,
		Employees: make([]entities.GroupEmployee, 0),
		Name:      "Мое расписание",
	}
	groups = append(groups, userEmployeeGroup)

	for groupRows.Next() {
		var gid int64
		var name string

		groupRows.Scan(&gid, &name)
		group := entities.Group{
			GID:       gid,
			Employees: make([]entities.GroupEmployee, 0),
			Name:      name,
		}
		groupsMap[strconv.FormatInt(gid, 10)] = &group
	}

	for empRows.Next() {
		var gid int64
		var eid int64
		var lname string
		var fname string
		var mname string

		empRows.Scan(&gid, &eid, &lname, &fname, &mname)
		employee := entities.GroupEmployee{
			EID:        eid,
			Lastname:   lname,
			Firstname:  fname,
			Middlename: mname,
		}

		if groupsMap[strconv.FormatInt(gid, 10)] != nil {
			groupsMap[strconv.FormatInt(gid, 10)].Employees = append(groupsMap[strconv.FormatInt(gid, 10)].Employees, employee)
		}
	}

	for _, group := range groupsMap {
		groups = append(groups, group)
	}

	return groups, nil
}

// GetTitleMeta возвращает карту сотрудников вида, [eid]ФИО
func (m *MetadataMapper) GetTitleMeta(relatedEmployeesRows *sql.Rows) (map[string]string, error) {
	employees := make(map[string]string)

	fmt.Println("Related employees:")

	for relatedEmployeesRows.Next() {
		eid := ""
		lastname := ""
		firstname := ""
		middlename := ""

		relatedEmployeesRows.Scan(&eid, &lastname, &firstname, &middlename)

		fmt.Printf("[%s]: %s %s %s", eid, lastname, firstname, middlename)

		employees[eid] = fmt.Sprintf("%s %s %s", lastname, firstname, middlename)
	}

	return employees, nil
}
