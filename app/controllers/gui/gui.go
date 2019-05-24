package controllers

import (
	"github.com/revel/revel"
)

// GUI controller struct
type GUI struct {
	*revel.Controller
}

// Index action name
func (c *GUI) Index() revel.Result {
	return c.Render()
}
