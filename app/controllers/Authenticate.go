package controllers

import (
	"crypto/md5"
	"fmt"
	"io"
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

	if c.Request.GetHttpHeader("Authorization") == "" {
		c.Response.Status = 401
		str, _ := GetDigestString("users@schedules", c.ClientIP)
		fmt.Print("\nGetDigestString\n" + str + "\n\n")
		digestString, err := GetDigestString("users@schedules", c.ClientIP)
		if err != nil {
			return c.Render(Failed(err))
		}
		c.Response.Out.Header().Add("WWW-Authenticate", digestString)

		return c.Render()
	}

	username, realmVal, nonceVal, digestURIVal, responseVal := getDigestHeaders(c.Request.GetHttpHeader("Authorization"))
	method := c.Request.Method
	password := "ikov"
	nonceVal = fmt.Sprintf("'%s'", nonceVal)

	// fmt.Print("\n\n", "username=", username, "\nrealm=", realmVal, "\nnonce=", nonceVal, "\nuriVal=", digestURIVal, "\nuri=", c.Request.GetPath(), "\nmethod=", method, "\n\n")

	ha1 := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:'%s':%s", username, realmVal, password))))
	ha2 := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s", method, digestURIVal))))

	serverResp := fmt.Sprintf("%x", md5.Sum([]byte(strings.Join([]string{ha1, nonceVal, ha2}, ":"))))

	// fmt.Print("\n\n", "ha1=", ha1, "\nha2=", ha2, "\nresp=", serverResp, "\n\n")

	if serverResp == responseVal {

		fmt.Print("\nAll right\n", serverResp, "\n", responseVal, "\n\n")
	} else {

		fmt.Print("\nAll bad\n", serverResp, "\n", responseVal, "\n\n")
	}

	return c.RenderJSON(Succes(c.Request.Header.Get("Authorization")))
}

func h(data string) []byte {
	h := md5.New()
	io.WriteString(h, data)

	return h.Sum(nil)
}
func kd(secret string, data string) []byte {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s:%s", secret, data))

	return h.Sum(nil)
}

func getDigestHeaders(response string) (string, string, string, string, string) {
	resp := strings.Replace(response, "Digest ", "", 1)
	respStrs := strings.Split(resp, ", ")

	respKeysVals := make(map[string]string, 0)
	for _, str := range respStrs {
		keyVal := strings.Split(str, "=")
		respKeysVals[keyVal[0]] = strings.Trim(keyVal[1], "'\"")
	}

	fmt.Print("\n\n", respKeysVals, "\n\n")

	username := respKeysVals["username"]
	realm := respKeysVals["realm"]
	nonce := respKeysVals["nonce"]
	uri := respKeysVals["uri"]
	responseVal := respKeysVals["response"]

	return username, realm, nonce, uri, responseVal
}
