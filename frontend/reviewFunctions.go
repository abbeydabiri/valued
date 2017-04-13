package frontend

import (
	"valued/database"

	"fmt"
	"html"
	"net/http"
	"strconv"

	"strings"
)

func (this *Review) viewStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Review is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	subSearchStore["merchant"] = html.EscapeString(httpReq.FormValue("merchant"))
	subSearch["rewardview-store-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Review) searchStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Review is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocrequest := make(map[string]interface{})
	xDocrequest["limit"] = "10"
	if html.EscapeString(httpReq.FormValue("limit")) == "" {
		xDocrequest["limit"] = html.EscapeString(httpReq.FormValue("limit"))
	}

	xDocrequest["offset"] = "0"
	if html.EscapeString(httpReq.FormValue("offset")) == "" {
		offset, err := strconv.Atoi(strings.TrimSpace(html.EscapeString(httpReq.FormValue("offset"))))
		if err == nil {
			if offset > 1 {
				xDocrequest["offset"] = fmt.Sprintf("%s", offset-1)
			}
		}
	}

	tblReviewStore := new(database.Store)
	subSearch := make(map[string]interface{})
	xDocrequest["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	xDocresult := tblReviewStore.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "linkedStoreActivate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "linkedStoreDeactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#rewardview-store-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *Review) linkStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Review is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("store")) == "" {
		sMessage += "Store is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblReviewStore := new(database.Store)
	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "inactive"
	xDoc["rewardcontrol"] = html.EscapeString(httpReq.FormValue("reward"))
	xDoc["storecontrol"] = html.EscapeString(httpReq.FormValue("store"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblReviewStore.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblReviewStore.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	this.searchStore(httpRes, httpReq, curdb)
}

func (this *Review) linkedStoreDelete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`delete from rewardstore where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Stored Deleted","triggerSubSearch":true}`))
}

func (this *Review) linkedStoreActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardstore set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Stored Activated","triggerSubSearch":true}`))
}

func (this *Review) linkedStoreDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardstore set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Stored Deactivated","triggerSubSearch":true}`))
}

func (this *Review) linkedStoreDeactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update rewardstore set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Linked Stores Deactivated","triggerSubSearch":true}`))
}

func (this *Review) viewScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Review is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchScheme := make(map[string]interface{})
	subSearchScheme["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	subSearchScheme["merchant"] = html.EscapeString(httpReq.FormValue("merchant"))
	subSearch["rewardview-scheme-search"] = subSearchScheme

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Review) searchScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Review is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocrequest := make(map[string]interface{})
	xDocrequest["limit"] = "10"
	if html.EscapeString(httpReq.FormValue("limit")) == "" {
		xDocrequest["limit"] = html.EscapeString(httpReq.FormValue("limit"))
	}

	xDocrequest["offset"] = "0"
	if html.EscapeString(httpReq.FormValue("offset")) == "" {
		offset, err := strconv.Atoi(strings.TrimSpace(html.EscapeString(httpReq.FormValue("offset"))))
		if err == nil {
			if offset > 1 {
				xDocrequest["offset"] = fmt.Sprintf("%s", offset-1)
			}
		}
	}

	tblReviewScheme := new(database.Scheme)
	subSearch := make(map[string]interface{})
	xDocrequest["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	xDocresult := tblReviewScheme.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "linkedSchemeActivate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "linkedSchemeDeactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}
		subSearch[cNumber+"#rewardview-scheme-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *Review) linkScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Review is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("scheme")) == "" {
		sMessage += "Scheme is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblReviewScheme := new(database.Scheme)
	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "inactive"
	xDoc["rewardcontrol"] = html.EscapeString(httpReq.FormValue("reward"))
	xDoc["schemecontrol"] = html.EscapeString(httpReq.FormValue("scheme"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblReviewScheme.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblReviewScheme.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	this.searchScheme(httpRes, httpReq, curdb)
}

func (this *Review) linkedSchemeDelete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`delete from rewardscheme where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Deleted","triggerSubSearch":true}`))
}

func (this *Review) linkedSchemeActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardscheme set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Activated","triggerSubSearch":true}`))
}

func (this *Review) linkedSchemeDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardscheme set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Deactivated","triggerSubSearch":true}`))
}

func (this *Review) linkedSchemeDeactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update rewardscheme set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Linked Schemes Deactivated","triggerSubSearch":true}`))
}
