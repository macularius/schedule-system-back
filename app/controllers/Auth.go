package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"myapp/app"
	"myapp/app/models/mappers"

	"github.com/revel/revel"
)

// Auth controller struct
type Auth struct {
	*revel.Controller
	Mapper *mappers.AuthMapper
}

// Login action name
func (c *Auth) Login(login string, password string) revel.Result {
	var connection *sql.DB
	var exist bool

	_, exist, err := c.Mapper.Authentication(login, password)
	if err != nil {
		errStr := err.Error()
		return c.Render(errStr)
	}
	if !exist {
		errStr := "Failed authentication"
		return c.Render(login, errStr)
	}

	sessions, errStr := app.Add(login)
	if errStr != "" {
		return c.Render(login, errStr)
	}
	for _, s := range sessions {
		if s.Login == login {
			connection = s.Connection
			break
		}
	}
	paramst, err := json.Marshal(c.Params.Values)
	if err != nil {
		errStr = err.Error()
		c.Render(login, errStr)
	}
	sessionsJSON, err := json.MarshalIndent(sessions, "", " ")
	if err != nil {
		errStr = err.Error()
		c.Render(login, errStr)
	}

	params := fmt.Sprintf("%s\n", paramst)
	sessionsStr := fmt.Sprintf("%s\n ", sessionsJSON)

	return c.Render(login, connection, params, sessionsStr, errStr)
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
		c.Render(login, params, errStr)
	}

	sessionsStr := fmt.Sprintf("%s\n ", sessionsJSON)

	return c.Render(login, params, sessionsStr, errStr)
}
