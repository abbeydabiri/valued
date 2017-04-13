package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	// "html"
	"net/http"
	"strconv"
)

type AppHome struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppHome) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		this.pageMap = make(map[string]interface{})
		this.pageMap["app-home"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Valued Home","pageContent":` + contentHTML + `}`))
		return
	}
}

func (this *AppHome) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})
	appFooter := make(map[string]interface{})
	appFooter["home"] = "white"
	formSearch["app-footer"] = appFooter

	appNavbar := new(AppNavbar)
	formSearch["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
	formSearch["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

	// if this.mapAppCache["control"] == nil {
	// 	formSearch["title"] = "REGISTER"
	// 	formSearch["onclick"] = "app-login?p=signup"
	// } else {

	// 	if this.mapAppCache["subscriptioncontrol"] == nil {
	// 		formSearch["title"] = "SUBSCRIBE"
	// 		formSearch["onclick"] = "app-scheme"
	// 	} else {
	// 		formSearch["title"] = "GIFT"
	// 		formSearch["onclick"] = "/app-scheme?action=gift"
	// 	}
	// }

	tblAppHome := new(database.Category)
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "100"
	xDocrequest["offset"] = "0"

	xDocrequest["sub"] = "false"
	xDocrequest["workflow"] = "active"

	xDocresult := tblAppHome.Search(xDocrequest, curdb)
	aSortedMap := this.SortMap(xDocresult)

	iMaxCounter := 3
	for _, cNumber := range aSortedMap {
		xDocInterface := xDocresult[cNumber]

		xDoc := xDocInterface.(map[string]interface{})
		iNumber, _ := strconv.Atoi(cNumber)
		switch iNumber % iMaxCounter {
		case 1:
			xDoc["class"] = "appCatLeft"
		case 2:
			xDoc["class"] = "appCatRight"
		default:
			xDoc["clearfix"] = "<div class='clearfix'></div>"
			xDoc["class"] = "appCatFull"

		}
		xDoc["onclick"] = fmt.Sprintf("getForm('/app-reward?category=%s')", xDoc["control"])
		formSearch[cNumber+"#app-home-category"] = xDoc
	}

	return formSearch
}
