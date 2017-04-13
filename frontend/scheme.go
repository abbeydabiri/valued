package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Scheme struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Scheme) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.pageMap["scheme-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Schemes","mainpanelContentSearch":` + contentHTML + `}`))
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

func (this *Scheme) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblScheme := new(database.Scheme)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocresult := tblScheme.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		xDoc["title"] = "No Schemes Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Scheme) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblScheme := new(database.Scheme)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocrequest["workflow"] = html.EscapeString(httpReq.FormValue("status"))

	xDocresult := tblScheme.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#scheme-search-result"] = xDoc

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

func (this *Scheme) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["scheme-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Scheme) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("title")) == "" {
		sMessage += " Name is missing <br>"
	}

	if strings.TrimSpace(html.EscapeString(httpReq.FormValue("price"))) == "" {
		sMessage += "Price is missing <br>"
	} else {
		regexNumber := regexp.MustCompile("^-*([0-9]+)*\\.*([0-9]+)$")
		if regexNumber.MatchString(strings.TrimSpace(html.EscapeString(httpReq.FormValue("price")))) == false {
			sMessage += "Price must be Numeric!! <br>"
		}
	}

	if html.EscapeString(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblScheme := new(database.Scheme)
	xDoc := make(map[string]interface{})
	xDoc["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDoc["workflow"] = html.EscapeString(httpReq.FormValue("workflow"))
	price, err := strconv.ParseFloat(strings.TrimSpace(html.EscapeString(httpReq.FormValue("price"))), 10)
	if err != nil {
		sMessage += "Price must be Numeric!! <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}
	xDoc["price"] = price

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblScheme.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblScheme.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"Scheme <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *Scheme) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblScheme := new(database.Scheme)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblScheme.Read(xDocRequest, curdb)
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
	formView["scheme-view"] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Scheme) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblScheme := new(database.Scheme)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = html.EscapeString(httpReq.FormValue("control"))

	ResMap := tblScheme.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})

	xDocRequest["formtitle"] = "Edit"
	formView["scheme-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Scheme) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update scheme set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Scheme Activated","triggerSearch":true}`))
}

func (this *Scheme) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update scheme set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Scheme Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Scheme) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update scheme set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Scheme De-Activated","triggerSearch":true}`))
}

func (this *Scheme) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update scheme set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Scheme De-Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Scheme) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update scheme set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Scheme De-Activated","triggerSearch":true}`))
}
