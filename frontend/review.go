package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"
	"fmt"
	"html"
	"net/http"
	"strconv"

	"regexp"
	"strings"
)

type Review struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Review) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.pageMap["review-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Reviews","mainpanelContentSearch":` + contentHTML + `}`))
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

	case "viewRedeemed":
		this.viewRedeemed(httpRes, httpReq, curdb)
		return

	case "searchRedeemed":
		this.searchRedeemed(httpRes, httpReq, curdb)
		return

	case "viewStore":
		this.viewStore(httpRes, httpReq, curdb)
		return

	case "searchStore":
		this.searchStore(httpRes, httpReq, curdb)
		return

	case "linkStore":
		this.linkStore(httpRes, httpReq, curdb)
		return

	case "linkedStoreDelete":
		this.linkedStoreDelete(httpRes, httpReq, curdb)
		return

	case "linkedStoreActivate":
		this.linkedStoreActivate(httpRes, httpReq, curdb)
		return

	case "linkedStoreDeactivate":
		this.linkedStoreDeactivate(httpRes, httpReq, curdb)
		return

	case "linkedStoreDeactivateAll":
		this.linkedStoreDeactivateAll(httpRes, httpReq, curdb)
		return

	case "viewScheme":
		this.viewScheme(httpRes, httpReq, curdb)
		return

	case "searchScheme":
		this.searchScheme(httpRes, httpReq, curdb)
		return

	case "linkScheme":
		this.linkScheme(httpRes, httpReq, curdb)
		return

	case "linkedSchemeDelete":
		this.linkedSchemeDelete(httpRes, httpReq, curdb)
		return

	case "linkedSchemeActivate":
		this.linkedSchemeActivate(httpRes, httpReq, curdb)
		return

	case "linkedSchemeDeactivate":
		this.linkedSchemeDeactivate(httpRes, httpReq, curdb)
		return

	case "linkedSchemeDeactivateAll":
		this.linkedSchemeDeactivateAll(httpRes, httpReq, curdb)
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

func (this *Review) viewRedeemed(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("review")) == "" {
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
	subSearchRedeemed := make(map[string]interface{})
	subSearchRedeemed["review"] = html.EscapeString(httpReq.FormValue("review"))
	subSearchRedeemed["merchant"] = html.EscapeString(httpReq.FormValue("merchant"))
	subSearch["reviewview-redeemed-search"] = subSearchRedeemed

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Review) searchRedeemed(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	// tblReviewRedeemed := new(database.Redemption)
	// quickSearch := make(map[string]interface{})
	// xDocrequest := make(map[string]interface{})

	// xDocrequest["limit"] = "10"
	// xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	// xDocresult := tblReviewRedeemed.Search(xDocrequest, curdb)

	// for cNumber, xDoc := range xDocresult {
	// 	xDoc := xDoc.(map[string]interface{})
	// 	xDoc["number"] = cNumber
	// 	xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
	// 	quickSearch[cNumber+"#quick-search-result"] = xDoc
	// }

	quickSearch := make(map[string]interface{})
	quickSearch["reviewview-redeemed-search"] = make(map[string]interface{})

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewDropdownHtml + `}`))
}

func (this *Review) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblReview := new(database.Review)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocresult := tblReview.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		xDoc["title"] = "No Reviews Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"reviewDropdown":` + viewDropdownHtml + `}`))
}

func (this *Review) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblReview := new(database.Review)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["merchant"] = html.EscapeString(httpReq.FormValue("merchant"))
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDocrequest["category"] = html.EscapeString(httpReq.FormValue("industry"))
	xDocrequest["email"] = html.EscapeString(httpReq.FormValue("email"))

	xDocrequest["workflow"] = html.EscapeString(httpReq.FormValue("status"))

	xDocresult := tblReview.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#review-search-result"] = xDoc

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

func (this *Review) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["review-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Review) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("title")) == "" {
		sMessage += "Review Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("category")) == "" {
		sMessage += "Category is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("startdate")) == "" {
		sMessage += "Start Date missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("enddate")) == "" {
		sMessage += "End Date missing <br>"
	}

	if strings.TrimSpace(html.EscapeString(httpReq.FormValue("totalavailable"))) == "" {
		sMessage += "Total Available is missing <br>"
	} else {
		regexNumber := regexp.MustCompile("^-*([0-9]+)*\\.*([0-9]+)$")
		if regexNumber.MatchString(strings.TrimSpace(html.EscapeString(httpReq.FormValue("totalavailable")))) == false {
			sMessage += "Total Available must be Numeric!! <br>"
		}
	}

	if strings.TrimSpace(html.EscapeString(httpReq.FormValue("txpermember"))) == "" {
		sMessage += "Transactions per Member is missing <br>"
	} else {
		regexNumber := regexp.MustCompile("^-*([0-9]+)*\\.*([0-9]+)$")
		if regexNumber.MatchString(strings.TrimSpace(html.EscapeString(httpReq.FormValue("txpermember")))) == false {
			sMessage += "Transactions per Member must be Numeric!! <br>"
		}
	}

	if html.EscapeString(httpReq.FormValue("beneficiary")) == "" {
		sMessage += "Beneficiary is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblReview := new(database.Review)
	xDoc := make(map[string]interface{})
	xDoc["title"] = html.EscapeString(httpReq.FormValue("title"))
	xDoc["workflow"] = html.EscapeString(httpReq.FormValue("workflow"))

	xDoc["merchantcontrol"] = html.EscapeString(httpReq.FormValue("merchant"))
	xDoc["categorycontrol"] = html.EscapeString(httpReq.FormValue("category"))

	xDoc["startdate"] = html.EscapeString(httpReq.FormValue("startdate"))
	xDoc["enddate"] = html.EscapeString(httpReq.FormValue("enddate"))

	txpermember, err := strconv.ParseFloat(strings.TrimSpace(html.EscapeString(httpReq.FormValue("txpermember"))), 10)
	if err != nil {
		sMessage += err.Error() + "!!! <br>"
	}

	totalavailable, err := strconv.ParseFloat(strings.TrimSpace(html.EscapeString(httpReq.FormValue("totalavailable"))), 10)
	if err != nil {
		sMessage += err.Error() + "!!! <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc["txpermember"] = txpermember
	xDoc["totalavailable"] = totalavailable
	xDoc["beneficiary"] = html.EscapeString(httpReq.FormValue("beneficiary"))

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("review-%s-%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblReview.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblReview.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"Review <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *Review) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblReview := new(database.Review)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblReview.Read(xDocRequest, curdb)
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
	formView["review-view"] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Review) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblReview := new(database.Review)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = html.EscapeString(httpReq.FormValue("control"))

	ResMap := tblReview.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})

	xDocRequest["formtitle"] = "Edit"
	formView["review-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Review) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a review"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update review set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Review Activated","triggerSearch":true}`))
}

func (this *Review) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a review"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update review set workflow = 'active' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Review Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Review) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a review"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update review set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Review De-Activated","triggerSearch":true}`))
}

func (this *Review) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a review"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update review set workflow = 'inactive' where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Review Deactivated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Review) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Reviews Not De-Activated"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update review set workflow = 'inactive' where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
	// httpRes.Write([]byte(`{"error":"Review De-Activated","triggerSearch":true}`))
}
