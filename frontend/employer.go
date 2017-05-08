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

type Employer struct {
	EmployerControl string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Employer) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/admin", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.EmployerControl = this.mapCache["control"].(string)
	if this.mapCache["company"].(string) != "Yes" {
		this.EmployerControl = this.mapCache["employercontrol"].(string)
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
		this.pageMap["employer-search"] = searchResult

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Employers","mainpanelContentSearch":` + contentHTML + `}`))
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

	/*
		case "search":
			searchResult := this.search(httpReq, curdb)
			contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
			httpRes.Write([]byte(`{"searchresult":` + contentHTML + `}`))
			return


				case "searchMember":
					subSearch := new(Member).search(httpReq, curdb)
					httpRes.Write([]byte(`{"subsearchresult":` + strconv.Quote(string(this.Generate(subSearch, nil))) + `}`))
					return

				case "searchSubscription":
					subSearch := new(Member).search(httpReq, curdb)
					httpRes.Write([]byte(`{"subsearchresult":` + strconv.Quote(string(this.Generate(subSearch, nil))) + `}`))
					return

				case "viewMember":
					this.viewMember(httpRes, httpReq, curdb)
					return

				case "viewSubscription":
					this.viewMember(httpRes, httpReq, curdb)
					return
	*/
	case "new":
		newHtml := this.new(httpReq, curdb)
		httpRes.Write([]byte(`{"mainpanelContent":` + newHtml + `}`))
		return

	case "save":
		this.save(httpRes, httpReq, curdb)
		return

	case "view":
		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), functions.TrimEscape(httpReq.FormValue("subview")), curdb)
		httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
		return

	case "viewUser":
		this.viewUser(httpRes, httpReq, curdb)
		return

	case "searchUser":
		this.searchUser(httpRes, httpReq, curdb)
		return

	case "linkedUserActivate":
		this.linkedUserActivate(httpRes, httpReq, curdb)
		return

	case "linkedUserDeactivate":
		this.linkedUserDeactivate(httpRes, httpReq, curdb)
		return

	case "viewEmployee":
		this.viewEmployee(httpRes, httpReq, curdb)
	case "searchEmployee":
		this.searchEmployee(httpRes, httpReq, curdb)
	case "newEmployee":
		this.newEmployee(httpRes, httpReq, curdb)
	case "editEmployee":
		this.editEmployee(httpRes, httpReq, curdb)
	case "saveEmployee":
		this.saveEmployee(httpRes, httpReq, curdb)
	case "activateEmployee":
		this.activateEmployee(httpRes, httpReq, curdb)
	case "deactivateEmployee":
		this.deactivateEmployee(httpRes, httpReq, curdb)

	case "viewSubscription":
		this.viewSubscription(httpRes, httpReq, curdb)
	case "searchSubscription":
		this.searchSubscription(httpRes, httpReq, curdb)
	case "newSubscription":
		this.newSubscription(httpRes, httpReq, curdb)
	case "editSubscription":
		this.editSubscription(httpRes, httpReq, curdb)
	case "saveSubscription":
		this.saveSubscription(httpRes, httpReq, curdb)
	case "activateSubscription":
		this.activateSubscription(httpRes, httpReq, curdb)
	case "deactivateSubscription":
		this.deactivateSubscription(httpRes, httpReq, curdb)

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

	case "sendWelcomeMail":
		this.sendWelcomeMail(httpRes, httpReq, curdb)
		return

	case "delete":
		this.delete(httpRes, httpReq, curdb)
		return
	}
}

func (this *Employer) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblEmployer := new(database.Profile)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["role"] = "merchant"
	xDocrequest["company"] = "Yes"
	xDocrequest["employercode"] = "main"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocresult := tblEmployer.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["title"] = "No Employers Found"
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Employer) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	tblEmployer := new(database.Profile)
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

	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	// xDocrequest["industry"] = functions.TrimEscape(httpReq.FormValue("industry"))
	xDocrequest["lastname"] = functions.TrimEscape(httpReq.FormValue("lastname"))
	xDocrequest["firstname"] = functions.TrimEscape(httpReq.FormValue("firstname"))
	xDocrequest["status"] = functions.TrimEscape(httpReq.FormValue("status"))
	xDocrequest["employercode"] = "main"
	xDocrequest["company"] = "Yes"
	xDocrequest["role"] = "employer"

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblEmployer.Search(xDocrequest, curdb)
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

	xDocresult := tblEmployer.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#employer-search-result"] = xDoc

		switch xDoc["status"].(string) {
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

	return
}

func (this *Employer) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	//Check The Profile and Link Accordingly
	tblRole := new(database.Role)
	xDocRoleReq := make(map[string]interface{})
	xDocRoleReq["workflow"] = "active"
	listMapRole := tblRole.Search(xDocRoleReq, curdb)

	formSelection := make(map[string]interface{})
	for cNumber, xDoc := range listMapRole {
		xDoc := xDoc.(map[string]interface{})
		if xDoc["code"].(string) == "employer" {
			xDoc["state"] = "checked"
		}
		formSelection[cNumber+"#role-edit-checkbox"] = xDoc
	}
	//Check The Profile and Link Accordingly

	formSelection["formtitle"] = "Add"
	formNew["employer-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Employer) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("status")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(httpReq.Form["profilerole"]) == 0 {
		sMessage += "Select Profile Role <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += "Company Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("phone")) == "" {
		sMessage += "Contact Number is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Primary Email is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("website")) == "" {
		sMessage += "Website is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("description")) == "" {
		sMessage += "Description is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	defaultMap, _ := curdb.Query(`select control from profile where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Employer <br>"
	} else {
		xDoc["employercontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	}

	// defaultMap, _ = curdb.Query(`select control from industry where code = 'main'`)
	// if defaultMap["1"] == nil {
	// 	sMessage += "Missing Default Industry <br>"
	// } else {
	// 	if functions.TrimEscape(httpReq.FormValue("subindustry")) == "" {
	// 		xDoc["subindustrycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	// 	}
	// }

	// defaultMap, _ = curdb.Query(`select control from category where code = 'main'`)
	// if defaultMap["1"] == nil {
	// 	sMessage += "Missing Default Category <br>"
	// } else {

	// 	if functions.TrimEscape(httpReq.FormValue("category")) == "" {
	// 		xDoc["categorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	// 	}

	// 	if functions.TrimEscape(httpReq.FormValue("subcategory")) == "" {
	// 		xDoc["subcategorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	// 	}
	// }

	if functions.TrimEscape(httpReq.FormValue("username")) == "" {
		xDoc["username"] = functions.TrimEscape(httpReq.FormValue("email"))
	} else {
		xDoc["username"] = functions.TrimEscape(httpReq.FormValue("username"))
	}

	if functions.TrimEscape(httpReq.FormValue("password")) == "" {
		xDoc["password"] = xDoc["username"]
	} else {
		xDoc["password"] = functions.TrimEscape(httpReq.FormValue("password"))
	}

	//--> Check if Email/Username exists in Login Table
	sqlCheckEmail := fmt.Sprintf(`select control from profile where email = '%s' and control != '%s' `,
		functions.TrimEscape(httpReq.FormValue("email")), functions.TrimEscape(httpReq.FormValue("control")))
	mapCheckEmail, _ := curdb.Query(sqlCheckEmail)

	if mapCheckEmail["1"] != nil {
		sMessage += "Email already exists <br>"
	}

	sqlCheckUsername := fmt.Sprintf(`select control from profile where username = '%s' and control != '%s' `,
		xDoc["username"], functions.TrimEscape(httpReq.FormValue("control")))
	mapCheckUsername, _ := curdb.Query(sqlCheckUsername)

	if mapCheckUsername["1"] != nil {
		sMessage += "Username already exists <br>"
	}
	//--> Check if Email/Username exists in Login Table

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblEmployer := new(database.Profile)
	xDoc["title"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("title")))
	xDoc["status"] = functions.TrimEscape(httpReq.FormValue("status"))

	// xDoc["firstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("firstname")))
	// xDoc["lastname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("lastname")))

	xDoc["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDoc["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))

	if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
		xDoc["phonecode"] = "+971"
	} else {
		xDoc["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
	}

	xDoc["emailsecondary"] = functions.TrimEscape(httpReq.FormValue("emailsecondary"))
	xDoc["emailalternate"] = functions.TrimEscape(httpReq.FormValue("emailalternate"))

	xDoc["website"] = functions.TrimEscape(httpReq.FormValue("website"))
	xDoc["description"] = functions.TrimEscape(httpReq.FormValue("description"))
	xDoc["company"] = "Yes"

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("employer_%s_%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblEmployer.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblEmployer.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	if xDoc["control"] == nil {
		sMessage = fmt.Sprintf("Error Saving Employer Record")
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//Handle Profile Role Link
	tblProfileRoleLink := new(database.ProfileRole)
	curdb.Query(fmt.Sprintf(`delete from profilerole where profilecontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["profilerole"] {
		xDocProfileRoleReqLink := make(map[string]interface{})
		xDocProfileRoleReqLink["profilecontrol"] = xDoc["control"]
		xDocProfileRoleReqLink["rolecontrol"] = control
		xDocProfileRoleReqLink["workflow"] = "active"
		tblProfileRoleLink.Create(this.mapCache["username"].(string), xDocProfileRoleReqLink, curdb)
	}
	//Handle Profile Role Link

	//Handle Group Link
	tblEmployerGroup := new(database.EmployerGroup)
	sPrimaryGroup := functions.TrimEscape(httpReq.FormValue("primarygroup"))
	curdb.Query(fmt.Sprintf(`delete from employergroup where employercontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["employergroup"] {
		xDocEmployerGroup := make(map[string]interface{})
		xDocEmployerGroup["employercontrol"] = xDoc["control"]
		xDocEmployerGroup["workflow"] = "active"
		xDocEmployerGroup["groupcontrol"] = control

		if sPrimaryGroup == "" {
			sPrimaryGroup = control
		}

		if sPrimaryGroup == control {
			xDocEmployerGroup["code"] = "primary"
		}

		tblEmployerGroup.Create(this.mapCache["username"].(string), xDocEmployerGroup, curdb)
	}
	//Handle Group Link

	sendEmail := ""
	if functions.TrimEscape(httpReq.FormValue("sendEmail")) == "Yes" {
		sendEmail = fmt.Sprintf(`"getform":"/employer?action=sendWelcomeMail&control=%s",`, xDoc["control"])
	}

	viewHTML := this.view(xDoc["control"].(string), "", curdb)
	httpRes.Write([]byte(`{` + sendEmail + `"error":"Record Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *Employer) view(control string, subview string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblEmployer := new(database.Profile)
	xDocResult := make(map[string]interface{})
	xDocResult["searchvalue"] = control

	ResMap := tblEmployer.Read(xDocResult, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocResult = ResMap["1"].(map[string]interface{})

	switch subview {
	case "employee":
		xDocResult["subview"] = "viewEmployee"

	case "subscription":
		xDocResult["subview"] = "viewSubscription"

	default:
		//case "user":
		xDocResult["subview"] = "viewUser"
	}

	switch xDocResult["status"].(string) {
	case "inactive":
		xDocResult["actionView"] = "activateView"
		xDocResult["actionColor"] = "success"
		xDocResult["actionLabel"] = "Activate"

	case "active":
		xDocResult["actionView"] = "deactivateView"
		xDocResult["actionColor"] = "danger"
		xDocResult["actionLabel"] = "De-Activate"
	}

	xDocResult["createdate"] = xDocResult["createdate"].(string)[0:19]

	//Check Profile Roles and Link Accordingly
	tblProfileRole := new(database.ProfileRole)
	xDocProfileRole := make(map[string]interface{})
	linkedRole := make(map[string]interface{})

	xDocProfileRole["profile"] = control
	linkRoleMap := tblProfileRole.Search(xDocProfileRole, curdb)
	for _, xDoc := range linkRoleMap {
		xDoc := xDoc.(map[string]interface{})
		linkedRole[xDoc["rolecontrol"].(string)] = xDoc["control"]
	}

	tblRole := new(database.Role)
	xDocRoleReq := make(map[string]interface{})
	xDocRoleReq["workflow"] = "active"
	listMapRole := tblRole.Search(xDocRoleReq, curdb)

	for cNumber, xDoc := range listMapRole {
		xDoc := xDoc.(map[string]interface{})
		if linkedRole[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#role-view-checkbox"] = xDoc
	}
	//Check Profile Roles and Link Accordingly

	//Check The Employer Groups and Link Accordingly
	tblEmployerGroup := new(database.EmployerGroup)
	xDocEmployerGroupReq := make(map[string]interface{})
	linkedEmployerGroup := make(map[string]interface{})

	xDocEmployerGroupReq["employer"] = control
	linkListMap := tblEmployerGroup.Search(xDocEmployerGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedEmployerGroup[xDoc["groupcontrol"].(string)] = xDoc
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedEmployerGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
			xDoc["isprimary"] = linkedEmployerGroup[xDoc["control"].(string)].(map[string]interface{})["code"]
			xDocResult[cNumber+"#employer-view-group"] = xDoc
		}
	}
	//Check The Employer Groups and Link Accordingly

	formView := make(map[string]interface{})
	formView["employer-view"] = xDocResult
	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Employer) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblEmployer := new(database.Profile)
	xDocResult := make(map[string]interface{})
	xDocResult["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblEmployer.Read(xDocResult, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocResult = ResMap["1"].(map[string]interface{})

	xDocResult[strings.Replace(xDocResult["phonecode"].(string), "+", "", 1)] = "selected"

	//Check The Profile and Link Accordingly
	linkedRole := make(map[string]interface{})
	sqlRole := fmt.Sprintf(`select rolecontrol from profilerole where profilecontrol = '%s'`, httpReq.FormValue("control"))
	linkRoleMap, _ := curdb.Query(sqlRole)
	for _, xDoc := range linkRoleMap {
		xDoc := xDoc.(map[string]interface{})
		linkedRole[xDoc["rolecontrol"].(string)] = xDoc
	}

	tblRole := new(database.Role)
	xDocRoleReq := make(map[string]interface{})
	xDocRoleReq["workflow"] = "active"
	listMapRole := tblRole.Search(xDocRoleReq, curdb)

	for cNumber, xDoc := range listMapRole {
		xDoc := xDoc.(map[string]interface{})
		if linkedRole[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#role-edit-checkbox"] = xDoc
	}
	//Check The Profile and Link Accordingly

	//Check The Employer Groups and Link Accordingly
	tblEmployerGroup := new(database.EmployerGroup)
	xDocEmployerGroupReq := make(map[string]interface{})
	linkedEmployerGroup := make(map[string]interface{})

	xDocEmployerGroupReq["employer"] = xDocResult["control"]
	linkListMap := tblEmployerGroup.Search(xDocEmployerGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedEmployerGroup[xDoc["groupcontrol"].(string)] = xDoc
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedEmployerGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
			xDoc["selection"] = "selected"
			xDoc["isprimary"] = linkedEmployerGroup[xDoc["control"].(string)].(map[string]interface{})["code"]
		}
		xDocResult[cNumber+"#employer-edit-group"] = xDoc
		xDocResult[cNumber+"#employer-edit-select"] = xDoc
	}
	//Check The Employer Groups and Link Accordingly

	formView := make(map[string]interface{})

	xDocResult["formtitle"] = "Edit"
	formView["employer-edit"] = xDocResult

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Employer) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select an employer"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Employer) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select an employer"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), "", curdb)
	httpRes.Write([]byte(`{"error":"Employer Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Employer) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select an Employer"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Employer) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select an Employer"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), "", curdb)
	httpRes.Write([]byte(`{"error":"Employer Deactivated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Employer) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Employers Not De-Activated"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Employer De-Activated","triggerSearch":true}`))
}

func (this *Employer) sendWelcomeMail(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a User"}`))
		return
	}

	//SEND AN EMAIL USING TEMPLATE
	sqlMember := fmt.Sprintf(`select employer.email as email, employer.title as employertitle, profile.username as username, 
		profile.password as password from profile left join profile as employer on employer.control = profile.employercontrol 
		where profile.control in ('0'%s)`, controlList)

	emailTo := ""
	emailFrom := "partnership@valued.com"
	emailFromName := "VALUED EMPLOYERS"
	emailTemplate := "employer-company-welcome"
	emailSubject := fmt.Sprintf("WELCOME TO VALUED - EMPLOYER")

	sUserDetail := `USER: <br> Username: %v <br> Password: %v < br><br>`
	resMember, _ := curdb.Query(sqlMember)
	for _, xDoc := range resMember {
		xDocEmail := xDoc.(map[string]interface{})
		emailFields := make(map[string]interface{})
		if xDocEmail["email"] != nil && xDocEmail["email"].(string) != "" {
			emailFields["email"] = xDocEmail["email"]

			emailTo = emailFields["email"].(string)
			emailFields["title"] = xDocEmail["employertitle"]
			emailFields["userdata"] = fmt.Sprintf(sUserDetail, xDocEmail["username"], xDocEmail["password"])
			go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, []string, emailFields)

		}
	}
	//SEND AN EMAIL USING TEMPLATE

	sMessage := fmt.Sprintf(`{"error":"%v Welcome Email(s) Sent"}`, len(resMember))
	httpRes.Write([]byte(sMessage))

}

func (this *Employer) delete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Employer Not Deleted}`))
		return
	}
	//
	curdb.Query(fmt.Sprintf(`delete from profile where control = '%s'`, functions.TrimEscape(httpReq.FormValue("control"))))
	httpRes.Write([]byte(`{"triggerSearch":true}`))
}
