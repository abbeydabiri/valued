package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
)

type Subscription struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Subscription) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.pageMap["subscription-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Subscriptions","mainpanelContentSearch":` + contentHTML + `}`))
		return

	case "quicksearch":
		this.quicksearch(httpRes, httpReq, curdb)
		return

	case "search":
		searchResult := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		httpRes.Write([]byte(`{"searchresult":` + contentHTML + `}`))
		return

	case "new":
		newHtml := this.new(httpReq, curdb)
		httpRes.Write([]byte(`{"mainpanelContent":` + newHtml + `}`))
		return

	case "save":
		this.save(httpRes, httpReq, curdb)
		return

	case "view":
		viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
		httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
		return

	case "edit":
		this.edit(httpRes, httpReq, curdb)
		return

	case "activate":
		this.activate(httpRes, httpReq, curdb)
		return

	case "activateView":
		this.activateView(httpRes, httpReq, curdb)
		return

	case "deactivate":
		this.deactivate(httpRes, httpReq, curdb)
		return

	case "deactivateView":
		this.deactivateView(httpRes, httpReq, curdb)
		return

	case "deactivateAll":
		this.deactivateAll(httpRes, httpReq, curdb)
		return
	}
}

func (this *Subscription) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblSubscription := new(database.Subscription)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocrequest["merchantcontrol"] = html.EscapeString(httpReq.FormValue("merchantcontrol"))
	xDocresult := tblSubscription.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		xDoc["title"] = "No Subscriptions Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Subscription) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblSubscription := new(database.Subscription)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["merchant"] = html.EscapeString(httpReq.FormValue("merchant"))
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocrequest["code"] = html.EscapeString(httpReq.FormValue("code"))
	xDocrequest["city"] = html.EscapeString(httpReq.FormValue("city"))

	xDocrequest["workflow"] = html.EscapeString(httpReq.FormValue("status"))

	xDocresult := tblSubscription.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#subscription-search-result"] = xDoc

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "activate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "deactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}
	}

	return formSearch
}

func (this *Subscription) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["subscription-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Subscription) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("title")) == "" {
		sMessage += "Subscription Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("contact")) == "" {
		sMessage += "Contact Person is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("phone")) == "" {
		sMessage += "Contact Number is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("street")) == "" {
		sMessage += "Address is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("address")) == "" {
		sMessage += "Address is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("city")) == "" {
		sMessage += "City is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("country")) == "" {
		sMessage += "Country is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("hoursmontofri")) == "" {
		sMessage += "Hours Mon-Fri is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("hourssat")) == "" {
		sMessage += "Hours Saturday is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("hourssun")) == "" {
		sMessage += "Hours Sun is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("hoursholiday")) == "" {
		sMessage += "Hours Holiday is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("gpslat")) == "" {
		sMessage += "GPS Latitude is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("gpslong")) == "" {
		sMessage += "GPS Longitude is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblSubscription := new(database.Subscription)
	xDoc := make(map[string]interface{})
	xDoc["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDoc["workflow"] = html.EscapeString(httpReq.FormValue("workflow"))

	xDoc["merchantcontrol"] = html.EscapeString(httpReq.FormValue("merchant"))

	xDoc["contact"] = html.EscapeString(httpReq.FormValue("contact"))
	xDoc["phone"] = html.EscapeString(httpReq.FormValue("phone"))

	xDoc["address"] = html.EscapeString(httpReq.FormValue("address"))
	xDoc["city"] = html.EscapeString(httpReq.FormValue("city"))
	xDoc["street"] = html.EscapeString(httpReq.FormValue("street"))
	xDoc["country"] = html.EscapeString(httpReq.FormValue("country"))

	xDoc["address"] = html.EscapeString(httpReq.FormValue("address"))
	xDoc["city"] = html.EscapeString(httpReq.FormValue("city"))
	xDoc["street"] = html.EscapeString(httpReq.FormValue("street"))
	xDoc["country"] = html.EscapeString(httpReq.FormValue("country"))

	xDoc["hoursholiday"] = html.EscapeString(httpReq.FormValue("hoursholiday"))
	xDoc["hoursmontofri"] = html.EscapeString(httpReq.FormValue("hoursmontofri"))
	xDoc["hourssat"] = html.EscapeString(httpReq.FormValue("hourssat"))
	xDoc["hourssun"] = html.EscapeString(httpReq.FormValue("hourssun"))

	xDoc["gpslat"] = html.EscapeString(httpReq.FormValue("gpslat"))
	xDoc["gpslong"] = html.EscapeString(httpReq.FormValue("gpslong"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblSubscription.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblSubscription.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"Subscription <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *Subscription) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblSubscription := new(database.Subscription)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblSubscription.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	switch xDocRequest["workflow"].(string) {
	case "inactive":
		xDocRequest["actionView"] = "activateView"
		xDocRequest["actionColor"] = "success"
		xDocRequest["actionLabel"] = "Activate"

	case "active":
		xDocRequest["actionView"] = "deactivateView"
		xDocRequest["actionColor"] = "danger"
		xDocRequest["actionLabel"] = "De-Activate"
	}

	xDocRequest["createdate"] = xDocRequest["createdate"].(string)[0:19]

	formView := make(map[string]interface{})
	formView["subscription-view"] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Subscription) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblSubscription := new(database.Subscription)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = html.EscapeString(httpReq.FormValue("control"))

	ResMap := tblSubscription.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})

	xDocRequest["formtitle"] = "Edit"
	formView["subscription-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Subscription) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscription"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Subscription) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscription"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Scheme Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Subscription) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscription"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Subscription) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscription"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Subscription De-Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Subscription) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscription"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update subscription set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}
