package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
)

type PendingApproval struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *PendingApproval) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.pageMap["pendingapproval-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search PendingApprovals","mainpanelContent":` + contentHTML + `}`))
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
		this.view(httpRes, httpReq, curdb)
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

func (this *PendingApproval) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblPendingApproval := new(database.PendingApproval)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocrequest["email"] = html.EscapeString(httpReq.FormValue("email"))
	xDocrequest["industry"] = html.EscapeString(httpReq.FormValue("industry"))
	xDocrequest["lastname"] = html.EscapeString(httpReq.FormValue("lastname"))
	xDocrequest["firstname"] = html.EscapeString(httpReq.FormValue("firstname"))
	xDocrequest["expirydate"] = html.EscapeString(httpReq.FormValue("expirydate"))
	xDocrequest["workflow"] = html.EscapeString(httpReq.FormValue("status"))

	xDocresult := tblPendingApproval.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#pendingapproval-search-result"] = xDoc

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

func (this *PendingApproval) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["pendingapproval-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *PendingApproval) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("title")) == "" {
		sMessage += "Company Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("industry")) == "" {
		sMessage += "Industry is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("firstname")) == "" {
		sMessage += "First Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("lastname")) == "" {
		sMessage += "Last Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("expirydate")) == "" {
		sMessage += "Expiry Date is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblPendingApproval := new(database.PendingApproval)
	xDoc := make(map[string]interface{})

	xDoc["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDoc["workflow"] = html.EscapeString(httpReq.FormValue("workflow"))

	xDoc["industry"] = html.EscapeString(httpReq.FormValue("industry"))
	xDoc["firstname"] = html.EscapeString(httpReq.FormValue("firstname"))
	xDoc["lastname"] = html.EscapeString(httpReq.FormValue("lastname"))
	xDoc["expirydate"] = html.EscapeString(httpReq.FormValue("expirydate"))

	xDoc["email"] = html.EscapeString(httpReq.FormValue("email"))
	xDoc["phone"] = html.EscapeString(httpReq.FormValue("phone"))
	xDoc["commercialized"] = html.EscapeString(httpReq.FormValue("commercialized"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblPendingApproval.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblPendingApproval.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = xDoc["control"]

	xDoc = tblPendingApproval.Read(xDocRequest, curdb)["1"].(map[string]interface{})
	formView := make(map[string]interface{})
	formView["pendingapproval-view"] = xDoc
	viewHtml := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"error":"PendingApproval <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHtml + `}`))
}

func (this *PendingApproval) view(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblPendingApproval := new(database.PendingApproval)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = html.EscapeString(httpReq.FormValue("control"))

	ResMap := tblPendingApproval.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
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
	formView["pendingapproval-view"] = xDocRequest

	viewHtml := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHtml + `}`))
}

func (this *PendingApproval) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblPendingApproval := new(database.PendingApproval)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = html.EscapeString(httpReq.FormValue("control"))

	ResMap := tblPendingApproval.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})

	xDocRequest["formtitle"] = "Edit"
	formView["pendingapproval-edit"] = xDocRequest

	viewHtml := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHtml + `}`))
}

func (this *PendingApproval) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a pending-approval"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update pendingapproval set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"PendingApproval Activated","triggerSearch":true}`))
}

func (this *PendingApproval) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a pending-approval"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update pendingapproval set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	this.view(httpRes, httpReq, curdb)
}

func (this *PendingApproval) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a pending-approval"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update pendingapproval set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"PendingApproval De-Activated","triggerSearch":true}`))
}

func (this *PendingApproval) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a pending-approval"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update pendingapproval set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	this.view(httpRes, httpReq, curdb)
}

func (this *PendingApproval) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	fmt.Println(controlList)
	if controlList == "" {
		httpRes.Write([]byte(`{"error":"PendingApprovals Not De-Activated"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update pendingapproval set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"PendingApproval De-Activated","triggerSearch":true}`))
}
