package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"

	"net/http"
	"strconv"
	"strings"
)

type Admin struct {
	functions.Templates
}

func (this *Admin) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if httpReq.Method == "GET" {

		pageMap := make(map[string]interface{})
		htmlContent := make(map[string]interface{})

		GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		mapCache := curdb.GetSession(GOSESSID.Value, "mapCache")

		frontend := strings.Split(httpReq.URL.Path[1:], "/")
		sPayload := fmt.Sprintf("getForm('/%s')", frontend[0])
		if mapCache["control"] != nil {
			sPayload = "getForm('/dashboard')"
		}

		htmlContent["payload"] = sPayload
		pageMap["html"] = htmlContent

		httpRes.Write(this.Generate(pageMap, nil))
		return
	}

	httpRes.Header().Set("content-type", "application/json")

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":

		frontend := strings.Split(httpReq.URL.Path[1:], "/")
		adminMap := make(map[string]interface{})
		adminMap["role"] = frontend[0]

		GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		mapCache := curdb.GetSession(GOSESSID.Value, "mapCache")

		if mapCache["role"] != nil {
			adminMap["role"] = mapCache["role"]
		}

		pageMap := make(map[string]interface{})
		pageMap["admin"] = adminMap
		curdb.SetSession(GOSESSID.Value, "mapCache", make(map[string]interface{}), false)

		pageHTML := strconv.Quote(string(this.Generate(pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Admin Login","pageContent":` + pageHTML + `}`))
		return

	case "admin":
		cUsername := strings.TrimSpace(httpReq.FormValue("user"))
		cPassword := strings.TrimSpace(httpReq.FormValue("pass"))
		cRole := strings.TrimSpace(httpReq.FormValue("role"))
		if cUsername != "" && cPassword != "" {

			tblProfile := new(database.Profile)
			mapRes := tblProfile.VerifyLogin(cRole, cUsername, cPassword, curdb)

			if mapRes["1"] != nil {
				GOSESSID, _ := httpReq.Cookie(_COOKIE_)
				xDoc := mapRes["1"].(map[string]interface{})
				curdb.SetSession(GOSESSID.Value, "mapCache", xDoc, false)

				new(Dashboard).Process(httpRes, httpReq, curdb)
				return
			}
		}

		httpRes.Write([]byte(`{"error":"Invalid Login Details"}`))
		return
	}
}
