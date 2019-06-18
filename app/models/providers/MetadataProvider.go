package providers

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"myapp/app/models/mappers"
	"strconv"
)

// MetadataProvider metadata provider struct
type MetadataProvider struct {
	mapper *mappers.MetadataMapper
}

// Init инициализирует provider
func (p *MetadataProvider) Init() {
	p.mapper = new(mappers.MetadataMapper)
}

// GetMenuMeta возвращает массив групп для меню
func (p *MetadataProvider) GetMenuMeta(eid string, db *sql.DB) ([]*entities.Group, error) {
	empRows, err := db.Query(selectGroupEmployeesConnectionString(eid))
	if err != nil {
		return nil, err
	}
	groupRows, err := db.Query(selectGroupsConnectionString(eid))
	if err != nil {
		return nil, err
	}

	groups, err := p.mapper.GetMenuMeta(empRows, groupRows)

	eidInt, err := strconv.ParseInt(eid, 10, 64)
	if err != nil {
		return nil, err
	}
	fmt.Print("\nGroups\n")
	fmt.Print(groups)
	groups[0].Employees = append(groups[0].Employees, entities.GroupEmployee{
		EID:        eidInt,
		Lastname:   "",
		Firstname:  "",
		Middlename: "",
	})

	return groups, nil
}

// GetTitleMeta возвращает карту [eid] = фио,
func (p *MetadataProvider) GetTitleMeta(eid string, db *sql.DB) (map[string]string, error) {
	userRelatedEmployeesRows, err := db.Query(selectRelatedEmployees(eid))
	if err != nil {
		return nil, err
	}

	relatedEmployees, err := p.mapper.GetTitleMeta(userRelatedEmployeesRows)
	if err != nil {
		return nil, err
	}

	return relatedEmployees, nil
}

func selectGroupEmployeesConnectionString(eid string) string {
	return fmt.Sprintf("select g.gid, e.eid, e.lastname, e.firstname, e.middlename from employees as e, groups as g, grouplist as gl where gl.gid = g.gid and gl.eid = e.eid and g.leadid = %s;", eid)
}
func selectGroupsConnectionString(eid string) string {
	return fmt.Sprintf("select gid, name from groups where leadid = %s;", eid)
}
func selectRelatedEmployees(eid string) string {
	return fmt.Sprintf("select e.eid, e.lastname, e.firstname, e.middlename from employees as e, groups as g, grouplist as gl where gl.gid = g.gid and gl.eid = e.eid and g.leadid = %s union select e.eid, e.lastname, e.firstname, e.middlename from employees as e where e.eid = %s;", eid, eid)
}
