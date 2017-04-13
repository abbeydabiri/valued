package frontend

import (
	"valued/database"
	"valued/functions"

	"net/http"
	"strconv"
	"strings"
)

type AppPin struct {
	functions.Templates
}

func (this *AppPin) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if httpReq.Method == "GET" {

		pageMap := make(map[string]interface{})
		htmlContent := make(map[string]interface{})

		GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		mapCache := curdb.GetSession(GOSESSID.Value, "mapCache")

		sPayload := "getForm('/apppin')"
		if mapCache["control"] != nil {
			sPayload = "getForm('/apphome')"
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

		pageMap := make(map[string]interface{})
		pageMap["apppin"] = make(map[string]interface{})
		pageHTML := strconv.Quote(string(this.Generate(pageMap, nil)))

		httpRes.Write([]byte(`{"pageTitle":"Valued Login","pageContent":` + pageHTML + `}`))
		return

	case "apppin":
		cUsername := strings.TrimSpace(httpReq.FormValue("user"))
		cPassword := strings.TrimSpace(httpReq.FormValue("pass"))
		if cUsername != "" && cPassword != "" {

			tblProfile := new(database.Profile)
			mapRes := tblProfile.VerifyLogin("member", cUsername, cPassword, curdb)

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
