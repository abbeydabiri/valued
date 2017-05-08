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

type MerchantStore struct {
	MerchantControl string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *MerchantStore) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.MerchantControl = this.mapCache["control"].(string)
	if this.mapCache["company"].(string) != "Yes" {
		this.MerchantControl = this.mapCache["employercontrol"].(string)
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})

		searchResult, searchPagination := this.search(httpReq, curdb)
		for sKey, iPagination := range searchPagination {
			searchResult[sKey] = iPagination
		}
		this.pageMap["merchantstore-search"] = searchResult

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

	/*
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
	*/

	case "deactivateAll":
		this.deactivateAll(httpRes, httpReq, curdb)
		return
	}
}

func (this *MerchantStore) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblMerchantStore := new(database.Store)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["merchantcontrol"] = this.MerchantControl
	xDocresult := tblMerchantStore.Search(xDocrequest, curdb)

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

func (this *MerchantStore) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	tblMerchantStore := new(database.Store)
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

	xDocrequest["offset"] = functions.TrimEscape(httpReq.FormValue("offset"))
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDocrequest["city"] = functions.TrimEscape(httpReq.FormValue("city"))
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))

	xDocrequest["merchantcontrol"] = this.MerchantControl

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblMerchantStore.Search(xDocrequest, curdb)
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

	xDocresult := tblMerchantStore.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#merchantstore-search-result"] = xDoc

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
	}

	return
}

func (this *MerchantStore) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["merchantstore-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *MerchantStore) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += "Merchant Store Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("contact")) == "" {
		sMessage += "Contact Person is missing <br>"
	}

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

	// if functions.TrimEscape(httpReq.FormValue("description")) == "" {
	// 	sMessage += "Description is missing <br>"
	// }

	if functions.TrimEscape(httpReq.FormValue("flagship")) == "Yes" {
		sqlCheckFlagship := fmt.Sprintf(`select control from store where flagship = 'Yes' and merchantcontrol = '%s' and control != '%s'`,
			this.MerchantControl, functions.TrimEscape(httpReq.FormValue("control")))
		resMap, _ := curdb.Query(sqlCheckFlagship)
		if resMap["1"] != nil {
			sMessage += "Multiple Flagship Stores Not Allowed<br>"
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblMerchantStore := new(database.Store)
	xDoc := make(map[string]interface{})
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))

	xDoc["merchantcontrol"] = this.MerchantControl
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
			fileName := fmt.Sprintf("merchantstore-%s-%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	sAction := ""
	sMerchant := ""
	sMessage = `Your store addition request has been successfully submitted for approval`
	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		// tblMerchantStore.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		sAction = "Approval for"
		xDoc["control"] = tblMerchantStore.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	//SEND AN EMAIL USING TEMPLATE
	emailFields := make(map[string]interface{})

	sqlStoreOLD := fmt.Sprintf(`select workflow, title, flagship, contact, phonecode, phone, email, address, city, 
		country, gpslat, gpslong, hoursmontofri, (select title from profile where control = merchantcontrol) as merchant from store where control = '%s'`, xDoc["control"])
	resStoreOLD, _ := curdb.Query(sqlStoreOLD)
	xDocOLD := make(map[string]interface{})
	if resStoreOLD["1"] != nil {
		xDocOLD = resStoreOLD["1"].(map[string]interface{})

		emailFields["workflowOLD"] = xDocOLD["workflow"]
		emailFields["titleOLD"] = xDocOLD["title"]
		emailFields["flagshipOLD"] = xDocOLD["flagship"]
		emailFields["contactOLD"] = xDocOLD["contact"]
		emailFields["phonecodeOLD"] = xDocOLD["phonecode"]
		emailFields["phoneOLD"] = xDocOLD["phone"]
		emailFields["emailOLD"] = xDocOLD["email"]
		emailFields["addressOLD"] = xDocOLD["address"]
		emailFields["cityOLD"] = xDocOLD["city"]
		emailFields["countryOLD"] = xDocOLD["country"]
		emailFields["gpslatOLD"] = xDocOLD["gpslat"]
		emailFields["gpslongOLD"] = xDocOLD["gpslong"]
		emailFields["openingOLD"] = xDocOLD["opening"]

		sMerchant = xDocOLD["merchant"].(string)
	}

	emailFields["action"] = sAction
	emailFields["username"] = this.mapCache["username"]

	emailFields["workflowNEW"] = xDoc["workflow"]
	emailFields["titleNEW"] = xDoc["title"]
	emailFields["flagshipNEW"] = xDoc["flagship"]
	emailFields["contactNEW"] = xDoc["contact"]
	emailFields["phonecodeNEW"] = xDoc["phonecode"]
	emailFields["phoneNEW"] = xDoc["phone"]
	emailFields["emailNEW"] = xDoc["email"]
	emailFields["addressNEW"] = xDoc["address"]
	emailFields["cityNEW"] = xDoc["city"]
	emailFields["countryNEW"] = xDoc["country"]
	emailFields["gpslatNEW"] = xDoc["gpslat"]
	emailFields["gpslongNEW"] = xDoc["gpslong"]
	emailFields["openingNEW"] = xDoc["opening"]

	emailTo := "rewards@valued.com"
	emailFrom := "notifications@valued.com"
	emailFromName := "VALUED ADMIN NOTIFICATIONS"
	emailSubject := fmt.Sprintf("Merchant %s Store Change Request", sMerchant)
	if sAction != "" {
		emailSubject = fmt.Sprintf("Merchant %s Added New Store", sMerchant)
	}
	emailTemplate := "merchantstore-save"
	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, []string, emailFields)
	//SEND AN EMAIL USING TEMPLATE

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"` + sMessage + `","mainpanelContent":` + viewHTML + `}`))
}

func (this *MerchantStore) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblMerchantStore := new(database.Store)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblMerchantStore.Read(xDocRequest, curdb)
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
	formView["merchantstore-view"] = xDocResult

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *MerchantStore) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblMerchantStore := new(database.Store)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblMerchantStore.Read(xDocRequest, curdb)
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

	formView["merchantstore-edit"] = xDocResult

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

/*
	func (this *MerchantStore) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a store"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update store set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		httpRes.Write([]byte(`{"triggerSearch":true}`))
	}

	func (this *MerchantStore) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a store"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update store set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
		httpRes.Write([]byte(`{"error":"Scheme Activated","mainpanelContent":` + viewHTML + `}`))
	}

	func (this *MerchantStore) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a store"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update store set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		httpRes.Write([]byte(`{"triggerSearch":true}`))
	}

	func (this *MerchantStore) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a store"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update store set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
		httpRes.Write([]byte(`{"error":"MerchantStore De-Activated","mainpanelContent":` + viewHTML + `}`))
	}
*/

func (this *MerchantStore) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	sMessage := ""
	if controlList == "" {
		sMessage += "Please select a store"
	}

	if functions.TrimEscape(httpReq.FormValue("message")) == "" {
		sMessage += "Please enter your reason"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//SEND AN EMAIL USING TEMPLATE
	emailFields := make(map[string]interface{})
	sqlDeactivate := fmt.Sprintf(`select store.workflow as workflow, store.title as store, merchant.title as merchant from store 
		left join profile as merchant on merchant.control = store.merchantcontrol where store.control in ('0'%s)`, controlList)
	resDeactivate, _ := curdb.Query(sqlDeactivate)

	sTableRows := ""
	sRow := `<tr> <td>%s</td> <td>%s</td> <td>%s</td> <td>%s</td></tr>`

	sMerchant := ""
	for _, xDoc := range resDeactivate {
		xDoc := xDoc.(map[string]interface{})
		sTableRows += fmt.Sprintf(sRow, "inactive", xDoc["workflow"], xDoc["store"], xDoc["merchant"])
		sMerchant = xDoc["merchant"].(string)
	}

	emailFields["rows"] = sTableRows
	emailFields["username"] = this.mapCache["username"]

	emailTo := "rewards@valued.com"
	emailCC := []string{this.mapCache["email"].(string)}
	emailFrom := "notifications@valued.com"
	emailFromName := "VALUED ADMIN NOTIFICATIONS"
	emailSubject := "Merchant " + sMerchant + "Store Deactivate Request"
	emailTemplate := "merchantstore-deactivate"
	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, emailCC, emailFields)
	//SEND AN EMAIL USING TEMPLATE

	httpRes.Write([]byte(`{"error":"Your request has been successfully submitted"}`))
}
