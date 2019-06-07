package controllers

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"myapp/app"
	"strings"

	"github.com/revel/revel"
)

// Authenticate отвечает за аутентификацию
type Authenticate struct {
	*revel.Controller
}

// Login digest auth
func (c *Authenticate) Login() revel.Result {
	fmt.Print("\nAuthorization\n" + c.Request.GetHttpHeader("Authorization") + "\n\n")

	if c.Request.GetHttpHeader("Authorization") != "" {
		username, realmVal, nonceVal, digestURIVal, responseVal := getDigestHeaders(c.Request.GetHttpHeader("Authorization"))
		method := c.Request.Method
		password := ""

		if exResponse, exist := app.IsExistByLogin(username); exist {
			if exResponse == responseVal {
				return c.Redirect("/")
			}
			return c.RenderJSON(Failed(errors.New("Неверный токен")))
		}

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

			app.Add(username, responseVal)

			return c.RenderJSON(Succes("All right " + "server: " + serverResp + "   val: " + responseVal))
		}
	}

	c.Response.Status = 401
	digestString, err := GetDigestString("users@schedules", c.ClientIP)
	if err != nil {
		return c.RenderJSON(Failed(err))
	}
	c.Response.Out.Header().Add("WWW-Authenticate", digestString)

	return c.Render()
}

func sqlGetUserString(username string) string {
	return fmt.Sprintf("SELECT password FROM users WHERE login='%s'", username)
}
