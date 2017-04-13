package frontend

import (
	"valued/database"
	"valued/functions"

	"net/http"
	"strconv"
)

type Password struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Password) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})
		this.pageMap["password"] = make(map[string]interface{})
		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))

		httpRes.Write([]byte(`{"pageTitle":"Account Password","mainpanelContent":` + contentHTML + `}`))
		return

	case "save":
		this.save(httpRes, httpReq, curdb)
		return
	}

}

func (this *Password) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("current")) == "" {
		sMessage += "Current Password is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("password")) == "" {
		sMessage += "New Password is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("confirm")) == "" {
		sMessage += "Confirm New Password is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("password")) != functions.TrimEscape(httpReq.FormValue("confirm")) {
		sMessage += "New Password does not match <br>"
	}

	tblUser := new(database.Profile)
	xDoc := make(map[string]interface{})
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = this.mapCache["control"]

	ResMap := tblUser.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		sMessage += "Cannot find user to update <br>"
	} else {
		xDocUser := ResMap["1"].(map[string]interface{})
		if xDocUser["password"].(string) != functions.TrimEscape(httpReq.FormValue("current")) {
			sMessage += "Current Password is incorrect <br>"
		}

	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}
	xDoc["control"] = this.mapCache["control"]
	xDoc["password"] = functions.TrimEscape(httpReq.FormValue("password"))
	tblUser.Update(this.mapCache["username"].(string), xDoc, curdb)

	httpRes.Write([]byte(`{"error":"Password Changed Successfully"}`))

}
