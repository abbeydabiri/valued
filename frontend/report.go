package frontend

import (
	"valued/database"
	"valued/functions"

	"net/http"
)

type Report struct {
	AdminCreatedate string
	AdminControl    string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Report) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.AdminControl = this.mapCache["control"].(string)
	this.AdminCreatedate = this.mapCache["createdate"].(string)[:10]
	if this.mapCache["company"].(string) != "Yes" {
		this.AdminControl = this.mapCache["employercontrol"].(string)
		this.AdminCreatedate = this.mapCache["employercreatedate"].(string)[:10]
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "", "summary":
		this.summary(httpRes, httpReq, curdb)
		return

	case "subscription":
		this.subscription(httpRes, httpReq, curdb)
		return

	case "merchant":
		this.merchant(httpRes, httpReq, curdb)
		return

	}
}
