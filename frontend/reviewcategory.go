package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"

	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ReviewCategory struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *ReviewCategory) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.pageMap["reviewcategory-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Categories","mainpanelContentSearch":` + contentHTML + `}`))
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
		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
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

func (this *ReviewCategory) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblReviewCategory := new(database.ReviewCategory)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocresult := tblReviewCategory.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	if viewDropdownHtml == "" {
		viewDropdownHtml = "No Categories Found"
	}

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *ReviewCategory) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblReviewCategory := new(database.ReviewCategory)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = functions.TrimEscape(httpReq.FormValue("offset"))
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))

	xDocresult := tblReviewCategory.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#reviewcategory-search-result"] = xDoc

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

func (this *ReviewCategory) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["reviewcategory-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *ReviewCategory) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += " Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblReviewCategory := new(database.ReviewCategory)
	xDoc := make(map[string]interface{})
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))
	xDoc["placement"] = functions.TrimEscape(httpReq.FormValue("placement"))
	xDoc["description"] = functions.TrimEscape(httpReq.FormValue("description"))

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("revcat-%s-%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblReviewCategory.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblReviewCategory.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"ReviewCategory <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *ReviewCategory) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblReviewCategory := new(database.ReviewCategory)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblReviewCategory.Read(xDocRequest, curdb)
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
	formView["reviewcategory-view"] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *ReviewCategory) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblReviewCategory := new(database.ReviewCategory)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblReviewCategory.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})

	xDocRequest["formtitle"] = "Edit"
	formView["reviewcategory-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *ReviewCategory) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update reviewcategory set workflow = 'active' where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Category Activated","triggerSearch":true}`))
}

func (this *ReviewCategory) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update reviewcategory set workflow = 'active' where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Category Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *ReviewCategory) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a review category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update reviewcategory set workflow = 'inactive' where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Category De-Activated","triggerSearch":true}`))
}

func (this *ReviewCategory) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a review category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update reviewcategory set workflow = 'inactive' where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Category De-Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *ReviewCategory) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update reviewcategory set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Category De-Activated","triggerSearch":true}`))
}
