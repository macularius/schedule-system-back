package mappers

import (
	"database/sql"
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
