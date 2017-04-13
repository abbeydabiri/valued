package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
)

type LandingPage struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *LandingPage) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.pageMap["landingpage-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Landing Page","pageContent":` + contentHTML + `}`))
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

func (this *LandingPage) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblRewards := new(database.Reward)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocresult := tblRewards.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		xDoc["title"] = "No LandingPages Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))
	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *LandingPage) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblLandingPage := new(database.Reward)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))

	xDocrequest["role"] = html.EscapeString(httpReq.FormValue("role"))
	xDocrequest["email"] = html.EscapeString(httpReq.FormValue("email"))
	xDocrequest["mobile"] = html.EscapeString(httpReq.FormValue("mobile"))
	xDocrequest["username"] = html.EscapeString(httpReq.FormValue("username"))

	xDocrequest["workflow"] = html.EscapeString(httpReq.FormValue("status"))

	xDocresult := tblLandingPage.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#landingpage-search-result"] = xDoc

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

func (this *LandingPage) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["user-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *LandingPage) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("title")) == "" {
		sMessage += "LandingPage Full-Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("role")) == "" {
		sMessage += "Role is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("email")) == "" {
		sMessage += "E-Mail Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("mobile")) == "" {
		sMessage += "Mobile is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("username")) == "" {
		sMessage += "LandingPagename is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblLandingPage := new(database.Reward)
	xDoc := make(map[string]interface{})
	xDoc["code"] = html.EscapeString(httpReq.FormValue("code"))
	xDoc["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDoc["workflow"] = html.EscapeString(httpReq.FormValue("workflow"))

	xDoc["role"] = html.EscapeString(httpReq.FormValue("role"))
	xDoc["email"] = html.EscapeString(httpReq.FormValue("email"))
	xDoc["mobile"] = html.EscapeString(httpReq.FormValue("mobile"))
	xDoc["username"] = html.EscapeString(httpReq.FormValue("username"))
	xDoc["password"] = html.EscapeString(httpReq.FormValue("password"))
	xDoc["description"] = html.EscapeString(httpReq.FormValue("description"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblLandingPage.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblLandingPage.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"LandingPage <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *LandingPage) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblLandingPage := new(database.Reward)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblLandingPage.Read(xDocRequest, curdb)
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
	formView["user-view"] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *LandingPage) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblLandingPage := new(database.Reward)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = html.EscapeString(httpReq.FormValue("control"))

	ResMap := tblLandingPage.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})

	xDocRequest["formtitle"] = "Edit"
	formView["user-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *LandingPage) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a landing page"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update user set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"LandingPage Activated","triggerSearch":true}`))
}

func (this *LandingPage) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a landing page"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update user set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"LandingPage Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *LandingPage) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a LandingPage"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update user set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"LandingPage De-Activated","triggerSearch":true}`))
}

func (this *LandingPage) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a LandingPage"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update user set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"LandingPage De-Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *LandingPage) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	fmt.Println(controlList)
	if controlList == "" {
		httpRes.Write([]byte(`{"error":"LandingPages Not De-Activated"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update user set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"LandingPage De-Activated","triggerSearch":true}`))
}
