package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
)

type User struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *User) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.pageMap["user-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Users","mainpanelContentSearch":` + contentHTML + `}`))
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

	case "delete":
		this.delete(httpRes, httpReq, curdb)
		return
	}
}

func (this *User) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblProfile := new(database.Profile)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocresult := tblProfile.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		xDoc["title"] = "No Users Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))
	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *User) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblUser := new(database.Profile)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))

	xDocrequest["role"] = html.EscapeString(httpReq.FormValue("role"))
	xDocrequest["email"] = html.EscapeString(httpReq.FormValue("email"))
	xDocrequest["mobile"] = html.EscapeString(httpReq.FormValue("mobile"))
	xDocrequest["username"] = html.EscapeString(httpReq.FormValue("username"))

	xDocrequest["workflow"] = html.EscapeString(httpReq.FormValue("status"))

	xDocresult := tblUser.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#user-search-result"] = xDoc

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

func (this *User) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})
	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"

	switch functions.TrimEscape(httpReq.FormValue("role")) {
	default:
		formSelection["role"] = "member"
	case "member", "employer", "merchant":
		formSelection["role"] = functions.TrimEscape(httpReq.FormValue("role"))
	}

	sqlProfile := fmt.Sprintf(`select title as title, firstname as firstname, lastname as lastname,
			control as profilecontrol from profile where control = '%s'`, functions.TrimEscape(httpReq.FormValue("control")))

	defaultMap, _ := curdb.Query(sqlProfile)
	if defaultMap["1"] != nil {
		xDoc := defaultMap["1"].(map[string]interface{})
		for sTitle, iValue := range xDoc {
			formSelection[sTitle] = iValue
		}
	}

	formNew["user-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *User) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("profile")) == "" {
		sMessage += "Profile is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("username")) == "" {
		sMessage += "Username is missing <br>"
	} else {

		//--> Check if Email/Username exists in Login Table
		sqlCheckEmail := fmt.Sprintf(`select control from login where username = '%s'`,
			functions.TrimEscape(httpReq.FormValue("username")))
		mapCheckEmail, _ := curdb.Query(sqlCheckEmail)

		if mapCheckEmail["1"] != nil {
			sMessage += "Email already exists <br>"
		}
		//--> Check if Email/Username exists in Login Table
	}

	if html.EscapeString(httpReq.FormValue("role")) == "" {
		sMessage += "Role is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("pin")) != "" {
		if len(html.EscapeString(httpReq.FormValue("pin"))) != 4 {
			sMessage += "Pin must be only 4 digits <br>"
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblUser := new(database.Profile)
	xDoc := make(map[string]interface{})
	xDoc["role"] = html.EscapeString(httpReq.FormValue("role"))
	xDoc["profilecontrol"] = html.EscapeString(httpReq.FormValue("profile"))

	if html.EscapeString(httpReq.FormValue("pin")) != "" {
		xDoc["pincode"] = html.EscapeString(httpReq.FormValue("pin"))
	}

	xDoc["username"] = html.EscapeString(httpReq.FormValue("username"))
	xDoc["workflow"] = html.EscapeString(httpReq.FormValue("workflow"))
	xDoc["password"] = html.EscapeString(httpReq.FormValue("password"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblUser.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblUser.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"User <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *User) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblUser := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblUser.Read(xDocRequest, curdb)
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

func (this *User) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblUser := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = html.EscapeString(httpReq.FormValue("control"))

	ResMap := tblUser.Read(xDocRequest, curdb)
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

func (this *User) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"User Activated","triggerSearch":true}`))
}

func (this *User) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"User Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *User) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"User De-Activated","triggerSearch":true}`))
}

func (this *User) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"User De-Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *User) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	fmt.Println(controlList)
	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update profile set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"User De-Activated","triggerSearch":true}`))
}

func (this *User) delete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}
	//
	curdb.Query(fmt.Sprintf(`delete from login where control = '%s'`, functions.TrimEscape(httpReq.FormValue("control"))))
	httpRes.Write([]byte(`{"triggerSearch":true}`))
}
