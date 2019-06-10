package controllers

import (
	"myapp/app"

	"github.com/revel/revel"
)

// GUI controller struct
type GUI struct {
	*revel.Controller
}

// Index action name
func (c *GUI) Index() revel.Result {
	// Проверка авторизованности
	if !app.IsExistBySID(c.Session.ID()) {
		return c.Redirect((*Authenticate).Login)
	}
	return c.Render()
}
