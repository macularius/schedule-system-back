package controllers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"myapp/app"
	"strings"

	"github.com/revel/revel"
)

// Authenticate отвечает за аутентификацию
type Authenticate struct {
	*revel.Controller
	actualNonce map[string]string // [:Session ID:]:nonce:
}

// Login digest auth
func (c *Authenticate) Login() revel.Result {
	if app.IsExistBySID(c.Session.ID()) {
		c.Redirect((*GUI).Index)
	}

	if c.Request.GetHttpHeader("Authorization") != "" {
		username, realmVal, nonceVal, digestURIVal, responseVal := getDigestHeaders(c.Request.GetHttpHeader("Authorization"))
		method := c.Request.Method
		var password string

		if c.actualNonce[c.Session.ID()] == nonceVal {
			db, err := sql.Open("postgres", app.GetConnectionString())
			if err != nil {
				return c.RenderJSON(Failed(err))
			}
			defer db.Close()

			row := db.QueryRow(sqlGetUserString(username))
			if err = row.Scan(&password); err != nil {
				return c.RenderJSON(Failed(err))
			}

			ha1 := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", username, realmVal, password))))
			ha2 := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", method, digestURIVal))))
			serverResp := fmt.Sprintf("%x", md5.Sum([]byte(strings.Join([]string{ha1, nonceVal, ha2}, ":"))))

			if serverResp == responseVal {
				fmt.Print("\nAll right\n", serverResp, "\n", responseVal, "\n\n")

				app.Add(c.Session.ID(), username, responseVal)

				fmt.Println("SID ", c.Session.ID())
				return c.Redirect((*GUI).Index)
			}

		}
	}

	c.Response.Status = 401
	nonce, digestString, err := GetDigestString("users@schedules", c.ClientIP)
	if err != nil {
		return c.RenderJSON(Failed(err))
	}
	c.Response.Out.Header().Add("WWW-Authenticate", digestString)
	if c.actualNonce == nil {
		c.actualNonce = make(map[string]string, 0)
	}
	c.actualNonce[c.Session.ID()] = nonce

	return c.Redirect((*Authenticate).Login)
}

// Logout user's logout
func (c *Authenticate) Logout() revel.Result {
	err := app.DeleteBySID(c.Session.ID())
	if err != nil {
		return c.RenderError(err)
	}
	delete(c.actualNonce, c.Session.ID())

	return c.Redirect((*Authenticate).Login)
}

func sqlGetUserString(username string) string {
	return fmt.Sprintf("SELECT password FROM users WHERE login='%s'", username)
}
