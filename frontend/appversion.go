package frontend

import (
	"valued/database"
	"valued/functions"

	// "fmt"
	// "html"
	"net/http"
	"strconv"
)

type AppVersion struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppVersion) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		AppVersion := make(map[string]interface{})
		this.pageMap = make(map[string]interface{})

		appNavbar := new(AppNavbar)
		AppVersion["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
		AppVersion["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

		AppVersion["app-footer"] = make(map[string]interface{})
		this.pageMap["app-version"] = AppVersion

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"APP VERSION","pageContent":` + contentHTML + `}`))
		return

	}
}
