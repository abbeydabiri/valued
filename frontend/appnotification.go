package frontend

import (
	"valued/database"
	"valued/functions"

	// "fmt"
	// "html"
	"net/http"
	"strconv"
)

type AppNotification struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppNotification) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		this.pageMap = make(map[string]interface{})
		AppNotification := make(map[string]interface{})

		appNavbar := new(AppNavbar)
		AppNotification["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
		AppNotification["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

		AppNotification["app-footer"] = make(map[string]interface{})
		this.pageMap["app-profile-notifications"] = AppNotification

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"My Notifications","pageContent":` + contentHTML + `}`))
		return

	}
}
