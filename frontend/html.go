package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"net/http"
)

type Html struct {
	functions.Templates

	pageMap  map[string]interface{}
	mapCache map[string]interface{}
}

//

func (this *Html) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	this.pageMap = make(map[string]interface{})
	htmlContent := make(map[string]interface{})

	switch httpReq.FormValue("init") {
	case "menu":
		new(database.Menu).Initialize(curdb)
		new(database.Permission).Initialize(curdb)

	case "role":
		new(database.Role).Initialize(curdb)
		new(database.ProfileRole).Initialize(curdb)
	}

	//

	httpReq.ParseForm()
	sPayload := "app-home"
	if functions.TrimEscape(httpReq.FormValue("a")) != "" {
		sPayload = functions.TrimEscape(httpReq.FormValue("a")) + "?"
		for sFormkey, sFormvalue := range httpReq.Form {
			switch sFormkey {
			case "a":
			default:
				sPayload += fmt.Sprintf("&%s=%s", sFormkey, functions.TrimEscape(sFormvalue[0]))
			}
		}
	}

	htmlContent["payload"] = fmt.Sprintf("getForm('/%s')", sPayload) //sPayload
	this.pageMap["html"] = htmlContent

	pageBytes := this.Generate(this.pageMap, nil)
	httpRes.Write(pageBytes)

}
