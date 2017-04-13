package frontend

import (
	"fmt"
	"net/http"
	"strconv"

	"valued/database"
	"valued/functions"
)

type Dashboard struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Dashboard) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if httpReq.Method == "GET" {

		pageMap := make(map[string]interface{})
		htmlContent := make(map[string]interface{})

		GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		mapCache := curdb.GetSession(GOSESSID.Value, "mapCache")

		sPayload := "getForm('/dashboard')"
		if mapCache["control"] == nil {
			sPayload = "getForm('/app-home')"
		}

		htmlContent["payload"] = sPayload
		pageMap["html"] = htmlContent

		httpRes.Write(this.Generate(pageMap, nil))
		return
	}

	httpRes.Header().Set("content-type", "application/json")

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.pageMap = make(map[string]interface{})
	this.pageMap["dashboard"] = this.GetMenu(curdb)

	// dashboard := fmt.Sprintf("dashboard-%s", this.mapCache["role"])
	// mainpanelMap := make(map[string]interface{})
	// this.pageMap = make(map[string]interface{})
	// this.pageMap[dashboard] = mainpanelMap

	pageHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageContent":` + pageHTML + `}`))

	// welcomeMessage := fmt.Sprintf("Welcome, %s", this.mapCache["title"])
	// httpRes.Write([]byte(`{"error":"` + welcomeMessage + `" , "pageContent":` + pageHTML + `}`))

	return
}

func (this *Dashboard) GetMenu(curdb database.Database) map[string]interface{} {

	xDocMenu := make(map[string]interface{})
	xDocMenuRequest := make(map[string]interface{})
	xDocMenuRequest["role"] = this.mapCache["role"]

	if this.mapCache["username"] != "root@localhost.com" {
		//Do Permissions Checking Here
	}

	xDocresult := new(database.Menu).Search(xDocMenuRequest, curdb)

	aSortedMap := this.SortMap(xDocresult)
	for _, cNumber := range aSortedMap {
		xDoc := xDocresult[cNumber].(map[string]interface{})
		if xDoc["workflow"].(string) == "active" {
			if xDoc["parent"].(string) == "" {
				if xDoc["code"].(string) == "member" &&
					this.mapCache["company"].(string) != "Yes" &&
					this.mapCache["employercode"].(string) != "main" {
					continue
				}

				sTag := fmt.Sprintf(`%s#dashboard-menuitem`, xDoc["placement"])
				xDocMenu[sTag] = xDoc
			}
		}
	}

	return xDocMenu

}
