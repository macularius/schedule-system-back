package controllers

import (
	"encoding/json"
	"fmt"
	"myapp/app"

	"github.com/revel/revel"
)

// Auth controller struct
type Auth struct {
	*revel.Controller
}

// Login action name
func (c *Auth) Login(login string) revel.Result {
	_, sessions := app.Add(login)
	paramst, err := json.Marshal(c.Params.Values)
	errStr := ""
	sessionsJSON, err := json.MarshalIndent(sessions, "", " ")
	if err != nil {
		errStr = err.Error()
		c.Render(login, "", nil, errStr)
	}

	params := fmt.Sprintf("%s\n", paramst)
	sessionsStr := fmt.Sprintf("%s\n ", sessionsJSON)

	return c.Render(login, params, sessionsStr, errStr)
}

// Logout action name
func (c *Auth) Logout(login string) revel.Result {
	_, sessions := app.Add(login)
	app.DeleteByLogin(login)
	params := c.Params.Values
	errStr := ""
	sessionsJSON, err := json.MarshalIndent(sessions, "", " ")
	if err != nil {
		errStr = err.Error()
		c.Render(login, params, nil, errStr)
	}

	sessionsStr := fmt.Sprintf("%s\n ", sessionsJSON)

	return c.Render(login, params, sessionsStr, errStr)
}
