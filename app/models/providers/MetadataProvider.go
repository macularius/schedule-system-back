package providers

import (
	"database/sql"
	"fmt"
	"myapp/app/models/entities"
	"myapp/app/models/mappers"
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
	empRows, err := db.Query(getGroupEmployeesConnectionString(eid))
	if err != nil {
		return nil, err
	}
	groupRows, err := db.Query(getGroupsConnectionString(eid))
	if err != nil {
		return nil, err
	}

	groups, err := p.mapper.GetMenuMeta(empRows, groupRows)
	return groups, nil
}

func getGroupEmployeesConnectionString(eid string) string {
	return fmt.Sprintf("select g.gid, e.eid, e.lastname, e.firstname, e.middlename from employees as e, groups as g, grouplist as gl where gl.gid = g.gid and gl.eid = e.eid and g.leadid = %s;", eid)
}
func getGroupsConnectionString(eid string) string {
	return fmt.Sprintf("select gid, name from groups where leadid = %s;", eid)
}
