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

type Store struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Store) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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

		searchResult, searchPagination := this.search(httpReq, curdb)
		for sKey, iPagination := range searchPagination {
			searchResult[sKey] = iPagination
		}
		this.pageMap["store-search"] = searchResult

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Stores","mainpanelContentSearch":` + contentHTML + `}`))
		return

	case "quicksearch":
		this.quicksearch(httpRes, httpReq, curdb)
		return

	case "search":
		searchResult, searchPagination := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		paginationHTML := strconv.Quote(string(this.Generate(searchPagination, nil)))
		httpRes.Write([]byte(`{"searchresult":` + contentHTML + `,"searchPage":` + paginationHTML + `}`))
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

func (this *Store) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblStore := new(database.Store)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["merchantcontrol"] = functions.TrimEscape(httpReq.FormValue("merchantcontrol"))
	xDocresult := tblStore.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		xDoc["title"] = "No Stores Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Store) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	tblStore := new(database.Store)
	xDocrequest := make(map[string]interface{})

	//Get Pagination Limit & Offset
	sLimit := "10"
	if functions.TrimEscape(httpReq.FormValue("limit")) != "" {
		sLimit = functions.TrimEscape(httpReq.FormValue("limit"))
	}

	sOffset := "0"
	if functions.TrimEscape(httpReq.FormValue("offset")) != "" {
		sOffset = functions.TrimEscape(httpReq.FormValue("offset"))
	}

	intLimit, _ := strconv.Atoi(sLimit)
	intOffset, _ := strconv.Atoi(sOffset)

	if intLimit > 0 && intOffset > 0 {
		sOffset = fmt.Sprintf("%v", (intOffset-1)*intLimit)
	}

	xDocrequest["limit"] = sLimit
	xDocrequest["offset"] = sOffset
	//Get Pagination Limit & Offset

	xDocrequest["merchant"] = functions.TrimEscape(httpReq.FormValue("merchant"))
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDocrequest["city"] = functions.TrimEscape(httpReq.FormValue("city"))
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))

	xDocrequest["merchantcontrol"] = functions.TrimEscape(httpReq.FormValue("merchantcontrol"))

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblStore.Search(xDocrequest, curdb)
	if xDocPagination["1"] != nil {
		xDocPagination := xDocPagination["1"].(map[string]interface{})
		if xDocPagination["paginationtotal"] != nil {
			nTotal = xDocPagination["paginationtotal"].(int64)
		}
	}
	delete(xDocrequest, "pagination")

	if nTotal > int64(intLimit) {
		nPage := int64(1)
		nPageMax := int64(nTotal/int64(intLimit)) + 1

		for nPage <= nPageMax {
			sPage := fmt.Sprintf("%v", nPage)
			mapPage := make(map[string]interface{})

			if intOffset > 0 && int64(intOffset) == nPage {
				mapPage["state"] = "selected"
			}
			mapPage["page"] = sPage
			searchPagination[sPage+"#select-page"] = mapPage
			nPage++
		}
	} else {
		searchPagination["1#select-page"] = "1"
	}
	//Set Pagination Limit & Offset

	xDocresult := tblStore.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive", "pending":
			xDoc["action"] = "activate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "deactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		formSearch[cNumber+"#store-search-result"] = xDoc
	}

	return
}

func (this *Store) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"

	//Get Employer Title
	if functions.TrimEscape(httpReq.FormValue("merchant")) != "" {
		sqlEmployer := fmt.Sprintf(`select title as merchanttitle, control as merchantcontrol from profile where control = '%s'`, functions.TrimEscape(httpReq.FormValue("merchant")))
		defaultMap, _ := curdb.Query(sqlEmployer)
		if defaultMap["1"] != nil {
			formSelection["merchanttitle"] = defaultMap["1"].(map[string]interface{})["merchanttitle"]
			formSelection["merchantcontrol"] = defaultMap["1"].(map[string]interface{})["merchantcontrol"]
		}
	}
	//Get Employer Title

	formNew["store-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Store) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += "Store Name is missing <br>"
	}

	// if functions.TrimEscape(httpReq.FormValue("contact")) == "" {
	// 	sMessage += "Contact Person is missing <br>"
	// }

	if functions.TrimEscape(httpReq.FormValue("phone")) == "" {
		sMessage += "Contact Number is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Email is missing <br>"
	}

	// if functions.TrimEscape(httpReq.FormValue("address")) == "" {
	// 	sMessage += "Address is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("city")) == "" {
	// 	sMessage += "City is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("country")) == "" {
	// 	sMessage += "Country is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("hoursmontofri")) == "" {
	// 	sMessage += "Opening Hours is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("hourssat")) == "" {
	// 	sMessage += "Hours Saturday is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("hourssun")) == "" {
	// 	sMessage += "Hours Sunday is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("hoursholiday")) == "" {
	// 	sMessage += "Hours Holiday is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("gpslat")) == "" {
	// 	sMessage += "GPS Latitude is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("gpslong")) == "" {
	// 	sMessage += "GPS Longitude is missing <br>"
	// }

	if functions.TrimEscape(httpReq.FormValue("description")) == "" {
		sMessage += "Description is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("flagship")) == "Yes" {
		sqlCheckFlagship := fmt.Sprintf(`select control from store where flagship = 'Yes' and merchantcontrol = '%s' and control != '%s'`,
			functions.TrimEscape(httpReq.FormValue("merchant")), functions.TrimEscape(httpReq.FormValue("control")))
		resMap, _ := curdb.Query(sqlCheckFlagship)

		if resMap["1"] != nil {
			sMessage += "Multiple Flagship Stores Not Allowed<br>"
		}

	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblStore := new(database.Store)
	xDoc := make(map[string]interface{})
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))

	xDoc["merchantcontrol"] = functions.TrimEscape(httpReq.FormValue("merchant"))
	xDoc["flagship"] = functions.TrimEscape(httpReq.FormValue("flagship"))

	xDoc["contact"] = functions.TrimEscape(httpReq.FormValue("contact"))
	xDoc["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))
	if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
		xDoc["phonecode"] = "+971"
	} else {
		xDoc["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
	}

	xDoc["email"] = functions.TrimEscape(httpReq.FormValue("email"))

	xDoc["address"] = functions.TrimEscape(httpReq.FormValue("address"))
	xDoc["city"] = functions.TrimEscape(httpReq.FormValue("city"))
	xDoc["country"] = functions.TrimEscape(httpReq.FormValue("country"))

	// xDoc["hoursholiday"] = functions.TrimEscape(httpReq.FormValue("hoursholiday"))
	// xDoc["hourssat"] = functions.TrimEscape(httpReq.FormValue("hourssat"))
	// xDoc["hourssun"] = functions.TrimEscape(httpReq.FormValue("hourssun"))
	// xDoc["gpslat"] = functions.TrimEscape(httpReq.FormValue("gpslat"))
	// xDoc["gpslong"] = functions.TrimEscape(httpReq.FormValue("gpslong"))

	sOpeningHours := functions.TrimEscape(httpReq.FormValue("hoursmontofri"))
	sOpeningHours = strings.Replace(sOpeningHours, "\r", "", -1)
	sOpeningHours = strings.Replace(sOpeningHours, "\n", "<br>", -1)
	xDoc["hoursmontofri"] = sOpeningHours

	sDescription := functions.TrimEscape(httpReq.FormValue("description"))
	sDescription = strings.Replace(sDescription, "\r", "", -1)
	sDescription = strings.Replace(sDescription, "\n", "<br>", -1)
	xDoc["description"] = sDescription

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("store-%s-%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblStore.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblStore.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"Store <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *Store) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblStore := new(database.Store)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblStore.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocResult := ResMap["1"].(map[string]interface{})

	switch xDocResult["workflow"].(string) {
	case "inactive", "pending":
		xDocResult["actionView"] = "activateView"
		xDocResult["actionColor"] = "success"
		xDocResult["actionLabel"] = "Activate"

	case "active":
		xDocResult["actionView"] = "deactivateView"
		xDocResult["actionColor"] = "danger"
		xDocResult["actionLabel"] = "De-Activate"
	}

	flagshipTag := "Flagship" + xDocResult["flagship"].(string)
	xDocResult[flagshipTag] = "checked"

	xDocResult["createdate"] = xDocResult["createdate"].(string)[0:19]

	formView := make(map[string]interface{})
	formView["store-view"] = xDocResult

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Store) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblStore := new(database.Store)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblStore.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocResult := ResMap["1"].(map[string]interface{})
	formView := make(map[string]interface{})
	xDocResult["formtitle"] = "Edit"

	flagshipTag := "Flagship" + xDocResult["flagship"].(string)
	xDocResult[flagshipTag] = "checked"

	sDescription := xDocResult["description"].(string)
	xDocResult["description"] = strings.Replace(sDescription, "<br>", "\n", -1)

	xDocResult[strings.Replace(xDocResult["phonecode"].(string), "+", "", 1)] = "selected"

	formView["store-edit"] = xDocResult

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Store) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update store set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Store) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update store set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Scheme Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Store) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update store set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Store) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update store set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Store De-Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Store) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update store set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}
