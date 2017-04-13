package frontend

import (
	"valued/database"
	"valued/functions"

	// "fmt"
	// "html"
	"net/http"
	"strconv"
)

type AppTerms struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppTerms) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		this.pageMap = make(map[string]interface{})
		AppTerms := make(map[string]interface{})

		appNavbar := new(AppNavbar)
		AppTerms["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
		AppTerms["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

		AppTerms["app-footer"] = make(map[string]interface{})
		this.pageMap["app-terms"] = AppTerms

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Valued Terms and Conditions","pageContent":` + contentHTML + `}`))
		return

	}
}
