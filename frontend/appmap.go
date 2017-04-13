package frontend

import (
	"valued/database"
	"valued/functions"

	// "fmt"
	"html"
	"net/http"
	"strconv"
)

type AppMap struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *AppMap) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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

		appMap := this.search(httpReq, curdb)
		appMap["app-navbar"] = make(map[string]interface{})
		appMap["app-slidebar"] = appMap["app-navbar"]
		appMap["app-searchdiv"] = appMap["app-navbar"]
		appMap["app-footer"] = appMap["app-navbar"]
		this.pageMap["app-map"] = appMap

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Valued Rewards","pageContent":` + contentHTML + `}`))
		return

	case "search":
		searchResult := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		httpRes.Write([]byte(`{"searchresult":` + contentHTML + `}`))
		return
	}
}

func (this *AppMap) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblAppStore := new(database.Store)
	locatorMarker := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["category"] = html.EscapeString(httpReq.FormValue("category"))
	xDocrequest["limit"] = html.EscapeString(httpReq.FormValue("limit"))
	xDocrequest["workflow"] = "active"

	xDocresult := tblAppStore.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		locatorMarker[cNumber+"#app-map-locator-marker"] = xDoc
	}

	switch len(xDocresult) {
	default:
		locatorMarker["zoom"] = 9
	case 1:
		locatorMarker["zoom"] = 13
	}

	formSearch["app-map-locator"] = locatorMarker

	return formSearch
}
