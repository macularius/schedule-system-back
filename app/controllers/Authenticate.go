package controllers

import (
	"fmt"

	"github.com/revel/revel"
)

// Authenticate отвечает за аутентификацию
type Authenticate struct {
	*revel.Controller
}

// Login digest auth
func (c *Authenticate) Login() revel.Result {
	fmt.Print("\nAuthorization\n" + c.Request.GetHttpHeader("Authorization") + "\n\n")

	if c.Request.GetHttpHeader("Authorization") == "" {
		c.Response.Status = 401
		digestString, err := GetDigestString("users@schedules", c.ClientIP)
		if err != nil {
			return c.Render(Failed(err))
		}
		c.Response.Out.Header().Add("WWW-Authenticate", digestString)

		return c.Render()
	}

	return c.RenderJSON(Succes(c.Request.Header.Get("Authorization")))
}

/*
// Auth controller struct
type Auth struct {
	*revel.Controller
	Provider *providers.AuthProvider
}

// Login action name
func (c *Auth) Login(login string, password string) revel.Result {
	var connection *sql.DB
	var exist bool

	_, exist, err := c.Provider.Authentication(login, password)
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
*/
