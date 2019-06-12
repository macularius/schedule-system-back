package controllers

import (
	"myapp/app"
	"myapp/app/models/providers"
	"strconv"

	"github.com/revel/revel"
)

// Metadata controller struct
type Metadata struct {
	*revel.Controller
	provider *providers.MetadataProvider
}

// GetMenuMeta action name
func (c *Metadata) GetMenuMeta() revel.Result {
	// Проверка авторизованности
	if !app.IsExistBySID(c.Session.ID()) {
		return c.Redirect((*Authenticate).Login)
	}
	if c.provider == nil {
		c.provider = new(providers.MetadataProvider)
		c.provider.Init()
	}
	session, err := app.GetSessionBySID(c.Session.ID())
	if err != nil {
		return c.RenderJSON(Failed(err))
	}
	groups, err := c.provider.GetMenuMeta(strconv.FormatInt(session.EmployeeID, 10), session.Connection)
	if err != nil {
		return c.RenderJSON(Failed(err))
	}

	return c.RenderJSON(Succes(groups))
}
