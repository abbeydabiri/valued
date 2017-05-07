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

type Member struct {
	functions.Templates
	EmployerControl string
	mapCache        map[string]interface{}
	pageMap         map[string]interface{}
}

func (this *Member) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
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

		sTag := "member-search"
		switch this.mapCache["role"].(string) {
		case "employer", "merchant":
			sTag = "member-employee-search"
		}
		this.pageMap[sTag] = searchResult

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Members","mainpanelContentSearch":` + contentHTML + `}`))
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

	case "viewSubscription":
		this.viewSubscription(httpRes, httpReq, curdb)
	case "searchSubscription":
		this.searchSubscription(httpRes, httpReq, curdb)
	case "saveSubscription":
		this.saveSubscription(httpRes, httpReq, curdb)
	case "newSubscription":
		this.newSubscription(httpRes, httpReq, curdb)
	case "editSubscription":
		this.editSubscription(httpRes, httpReq, curdb)
	case "activateSubscription":
		this.activateSubscription(httpRes, httpReq, curdb)
	case "deactivateSubscription":
		this.deactivateSubscription(httpRes, httpReq, curdb)
	case "deleteSubscription":
		this.deleteSubscription(httpRes, httpReq, curdb)

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

	case "sendSubscriptionMail":
		this.sendSubscriptionMail(httpRes, httpReq, curdb)
		return

	case "delete":
		this.delete(httpRes, httpReq, curdb)
		return
	}
}

func (this *Member) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblMember := new(database.Profile)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["member"] = "Yes"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))

	xDocresult := tblMember.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		xDoc["title"] = "No Members Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Member) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	sTag := "member-search-result"
	switch this.mapCache["role"].(string) {
	case "employer", "merchant":
		sTag = "member-employee-search-result"
	}

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	// tblMember := new(database.Profile)
	tblMember := new(database.Profile)
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

	xDocrequest["firstname"] = functions.TrimEscape(httpReq.FormValue("firstname"))
	xDocrequest["lastname"] = functions.TrimEscape(httpReq.FormValue("lastname"))
	xDocrequest["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDocrequest["username"] = functions.TrimEscape(httpReq.FormValue("username"))
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))
	xDocrequest["status"] = functions.TrimEscape(httpReq.FormValue("status"))
	xDocrequest["role"] = functions.TrimEscape(httpReq.FormValue("role"))

	xDocrequest["employertitle"] = functions.TrimEscape(httpReq.FormValue("employer"))
	xDocrequest["employercontrol"] = functions.TrimEscape(httpReq.FormValue("employercontrol"))

	//Check SignedIn-User Role and Limit Results
	switch this.mapCache["role"].(string) {
	case "merchant", "employer":
		delete(xDocrequest, "role")
		delete(xDocrequest, "employertitle")
		xDocrequest["employercontrol"] = this.EmployerControl
	}
	//Check SignedIn-User Role and Limit Results

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblMember.Search(xDocrequest, curdb)
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

	xDocresult := tblMember.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "pending", "expired", "subscribed-pending":
			mapButton := make(map[string]interface{})
			mapButton["action"] = "paid"
			mapButton["actionColor"] = "success"
			mapButton["actionLabel"] = "Paid"
			xDoc["workflow-button"] = mapButton
		}

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

		formSearch[cNumber+"#"+sTag] = xDoc
	}

	return
}

func (this *Member) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"

	sTag := "member-edit"
	switch this.mapCache["role"].(string) {
	case "employer", "merchant":
		sTag = "member-employee-edit"
	}

	//Get Employer Title
	if functions.TrimEscape(httpReq.FormValue("merchant")) != "" {
		sqlEmployer := fmt.Sprintf(`select title as employertitle, control as employercontrol from profile where control = '%s'`, functions.TrimEscape(httpReq.FormValue("merchant")))
		defaultMap, _ := curdb.Query(sqlEmployer)
		if defaultMap["1"] != nil {
			formSelection["employertitle"] = defaultMap["1"].(map[string]interface{})["employertitle"]
			formSelection["employercontrol"] = defaultMap["1"].(map[string]interface{})["employercontrol"]
		}
	}
	//Get Employer Title

	//Check The Profile and Link Accordingly
	tblRole := new(database.Role)
	xDocRoleReq := make(map[string]interface{})
	xDocRoleReq["workflow"] = "active"
	listMapRole := tblRole.Search(xDocRoleReq, curdb)

	for cNumber, xDoc := range listMapRole {
		xDoc := xDoc.(map[string]interface{})
		if sTag == "member-employee-edit" {
			switch xDoc["code"].(string) {
			case "member":
				xDoc["state"] = "checked"
				formSelection[cNumber+"#role-edit-checkbox"] = xDoc

			case "employer", "merchant":
				if this.mapCache["role"].(string) == xDoc["code"].(string) {
					xDoc["state"] = "checked"
					formSelection[cNumber+"#role-edit-checkbox"] = xDoc
				}
			}
		} else {
			if xDoc["code"].(string) == "member" {
				xDoc["state"] = "checked"
			}
			if xDoc["code"].(string) == functions.TrimEscape(httpReq.FormValue("role")) {
				xDoc["state"] = "checked"
			}

			formSelection[cNumber+"#role-edit-checkbox"] = xDoc
		}
	}
	//Check The Profile and Link Accordingly

	//Check The Member Groups and Link Accordingly
	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = ""
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		formSelection[cNumber+"#member-edit-group"] = xDoc
	}
	//Check The Member Groups and Link Accordingly

	formNew[sTag] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Member) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("status")) == "" {
		sMessage += "Status is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("workflow")) == "" {
		sMessage += "Workflow is missing <br>"
	}

	if len(httpReq.Form["profilerole"]) == 0 {
		sMessage += "Select Profile Role <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		if this.mapCache["role"].(string) == "admin" {
			sMessage += "Employer is missing <br>"
		}
	}

	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += "Title is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("firstname")) == "" {
		sMessage += "First Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("lastname")) == "" {
		sMessage += "Surname is missing <br>"
	}

	// if functions.TrimEscape(httpReq.FormValue("dob")) == "" {
	// 	sMessage += "Date of Birth is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("nationality")) == "" {
	// 	sMessage += "Nationality is missing <br>"
	// }

	// if functions.TrimEscape(httpReq.FormValue("phone")) == "" {
	// 	sMessage += "Contact Number is missing <br>"
	// }

	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "EMail is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	defaultMap, _ := curdb.Query(`select control from category where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Category <br>"
	} else {
		xDoc["categorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
		xDoc["subcategorycontrol"] = xDoc["categorycontrol"]
	}

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

	if functions.TrimEscape(httpReq.FormValue("pincode")) == "" {
		xDoc["pincode"] = "1234"
	} else {
		xDoc["pincode"] = functions.TrimEscape(httpReq.FormValue("pincode"))
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

	tblMember := new(database.Profile)
	xDoc["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	// xDoc["pincode"] = functions.TrimEscape(httpReq.FormValue("pincode"))
	// xDoc["username"] = functions.TrimEscape(httpReq.FormValue("username"))
	// xDoc["password"] = functions.TrimEscape(httpReq.FormValue("password"))

	xDoc["title"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("title")))
	xDoc["status"] = functions.TrimEscape(httpReq.FormValue("status"))

	xDoc["firstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("firstname")))
	xDoc["lastname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("lastname")))

	xDoc["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDoc["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))
	if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
		xDoc["phonecode"] = "+971"
	} else {
		xDoc["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
	}

	xDoc["dob"] = functions.TrimEscape(httpReq.FormValue("dob"))
	xDoc["nationality"] = functions.TrimEscape(httpReq.FormValue("nationality"))
	xDoc["company"] = ""

	if this.mapCache["role"] != "admin" {
		xDoc["employercontrol"] = this.EmployerControl
	} else {
		xDoc["employercontrol"] = functions.TrimEscape(httpReq.FormValue("employer"))
	}

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("member_%s_%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {

		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))
		tblMember.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["workflow"] = "registered"
		xDoc["control"] = tblMember.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	//Handle Profile Role Link
	tblProfileRoleLink := new(database.ProfileRole)
	curdb.Query(fmt.Sprintf(`delete from profilerole where profilecontrol = '%s'`, xDoc["control"]))

	allowedRoles := make(map[string]string)
	for _, control := range httpReq.Form["profilerole"] {
		allowedRoles[control] = control
	}

	if allowedRoles["1.00000000001"] != "" {
		allowedRoles = make(map[string]string)
		allowedRoles["1.00000000001"] = "1.00000000001"
	}

	for _, control := range allowedRoles {
		xDocProfileRoleReqLink := make(map[string]interface{})
		xDocProfileRoleReqLink["profilecontrol"] = xDoc["control"]
		xDocProfileRoleReqLink["rolecontrol"] = control
		xDocProfileRoleReqLink["workflow"] = "active"
		tblProfileRoleLink.Create(this.mapCache["username"].(string), xDocProfileRoleReqLink, curdb)
	}
	//Handle Profile Role Link

	//Handle Group Link
	if this.mapCache["role"].(string) == "admin" {
		tblMemberGroup := new(database.MemberGroup)
		curdb.Query(fmt.Sprintf(`delete from membergroup where membercontrol = '%s'`, xDoc["control"]))
		for _, control := range httpReq.Form["membergroup"] {
			xDocMemberGroup := make(map[string]interface{})
			xDocMemberGroup["membercontrol"] = xDoc["control"]
			xDocMemberGroup["workflow"] = "active"
			xDocMemberGroup["groupcontrol"] = control
			tblMemberGroup.Create(this.mapCache["username"].(string), xDocMemberGroup, curdb)
		}
	}
	//Handle Group Link

	sendEmail := ""
	if functions.TrimEscape(httpReq.FormValue("sendEmail")) == "Yes" {
		sendEmail = fmt.Sprintf(`"getform":"/member?action=sendWelcomeMail&control=%s",`, xDoc["control"])
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{` + sendEmail + `"error":"Record Saved","mainpanelContent":` + viewHTML + `}`))
	return
}

func (this *Member) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	sTag := "member-view"
	switch this.mapCache["role"].(string) {
	case "employer", "merchant":
		sTag = "member-employee-view"
	}

	tblMember := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblMember.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocResult := ResMap["1"].(map[string]interface{})

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

	nCurYear, _ := strconv.Atoi(functions.GetSystemTime()[6:10])
	nDobYear := nCurYear
	if len(xDocResult["dob"].(string)) == 10 {
		nDobYear, _ = strconv.Atoi(xDocResult["dob"].(string)[6:10])
	}
	xDocResult["age"] = nCurYear - nDobYear

	//Get Employer Info via EmployerControl
	/*
		tblEmployer := new(database.Profile)
		employerReq := make(map[string]interface{})
		employerReq["searchvalue"] = xDocResult["employercontrol"]
		employerRes := tblEmployer.Read(employerReq, curdb)
		if employerRes["1"] != nil {
			employerXdoc := employerRes["1"].(map[string]interface{})
			for cFieldname, iFieldvalue := range employerXdoc {
				xDocResult["employer"+cFieldname] = iFieldvalue
			}
		}
	*/
	//Get Employer Info via EmployerControl

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

		if sTag == "member-employee-view" {
			switch xDoc["code"].(string) {
			case "member":
				xDocResult[cNumber+"#role-view-checkbox"] = xDoc

			case "employer", "merchant":
				if this.mapCache["role"].(string) == xDoc["code"].(string) {
					xDocResult[cNumber+"#role-view-checkbox"] = xDoc
				}
			}
		} else {
			xDocResult[cNumber+"#role-view-checkbox"] = xDoc
		}
	}
	//Check Profile Roles and Link Accordingly

	//Check The Member Groups and Link Accordingly
	tblMemberGroup := new(database.MemberGroup)
	xDocMemberGroupReq := make(map[string]interface{})
	linkedMemberGroup := make(map[string]interface{})

	xDocMemberGroupReq["member"] = control
	linkListMap := tblMemberGroup.Search(xDocMemberGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedMemberGroup[xDoc["groupcontrol"].(string)] = xDoc["control"]
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedMemberGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#member-view-group"] = xDoc
	}
	//Check The Member Groups and Link Accordingly

	formView := make(map[string]interface{})
	formView[sTag] = xDocResult

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Member) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	sTag := "member-edit"
	switch this.mapCache["role"].(string) {
	case "employer", "merchant":
		sTag = "member-employee-edit"
	}

	tblMember := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblMember.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocResult := ResMap["1"].(map[string]interface{})

	//Get Employer Info via EmployerControl
	/*
		tblEmployer := new(database.Profile)
		employerReq := make(map[string]interface{})
		employerReq["searchvalue"] = xDocResult["employercontrol"]
		employerRes := tblEmployer.Read(employerReq, curdb)
		if employerRes["1"] != nil {
			employerXdoc := employerRes["1"].(map[string]interface{})
			for cFieldname, iFieldvalue := range employerXdoc {
				xDocResult["employer"+cFieldname] = iFieldvalue
			}
		}
	*/
	//Get Employer Info via EmployerControl

	//Check Profile Roles and Link Accordingly
	tblProfileRole := new(database.ProfileRole)
	xDocProfileRole := make(map[string]interface{})
	linkedRole := make(map[string]interface{})

	xDocProfileRole["profile"] = functions.TrimEscape(httpReq.FormValue("control"))
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

		if sTag == "member-employee-edit" {
			switch xDoc["code"].(string) {
			case "member":
				xDocResult[cNumber+"#role-edit-checkbox"] = xDoc

			case "employer", "merchant":
				if this.mapCache["role"].(string) == xDoc["code"].(string) {
					xDocResult[cNumber+"#role-edit-checkbox"] = xDoc
				}
			}
		} else {
			xDocResult[cNumber+"#role-edit-checkbox"] = xDoc
		}
	}
	//Check Profile Roles and Link Accordingly

	//Check The Member Groups and Link Accordingly
	tblMemberGroup := new(database.MemberGroup)
	xDocMemberGroupReq := make(map[string]interface{})
	linkedMemberGroup := make(map[string]interface{})

	xDocMemberGroupReq["member"] = xDocResult["control"]
	linkListMap := tblMemberGroup.Search(xDocMemberGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedMemberGroup[xDoc["groupcontrol"].(string)] = xDoc["control"]
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedMemberGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#member-edit-group"] = xDoc
	}
	//Check The Member Groups and Link Accordingly

	xDocResult[strings.Replace(xDocResult["phonecode"].(string), "+", "", 1)] = "selected"

	xDocResult["formtitle"] = "Edit"
	formView := make(map[string]interface{})
	formView[sTag] = xDocResult

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Member) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a User"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Member) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a User"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Member Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Member) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a member"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Member) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a member"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Member Deactivated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Member) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a Member"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Member De-Activated","triggerSearch":true}`))
}

func (this *Member) sendWelcomeMail(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a Merchant"}`))
		return
	}

	//SEND AN EMAIL USING TEMPLATE
	sqlMember := fmt.Sprintf(`select profile.firstname as firstname, profile.lastname as lastname, employer.title as employertitle, 
		profile.username as username, profile.password as password, profile.email as email from profile left join 
		profile as employer on employer.control = profile.employercontrol where profile.control in ('0'%s)`, controlList)

	emailTo := ""
	emailFrom := "rewards@valued.com"
	emailFromName := "VALUED ADMIN"
	emailTemplate := "member-welcome"
	emailSubject := fmt.Sprintf("Welcome to VALUED - Portal Registration")

	resMember, _ := curdb.Query(sqlMember)
	for _, xDoc := range resMember {
		emailFields := xDoc.(map[string]interface{})

		title := ""
		firstname := ""
		lastname := ""

		if emailFields["title"] != nil {
			title = fmt.Sprintf(`%v`, emailFields["title"])
		}

		if emailFields["title"] != nil {
			firstname = fmt.Sprintf(`%v`, emailFields["firstname"])
		}

		if emailFields["title"] != nil {
			lastname = fmt.Sprintf(`%v`, emailFields["lastname"])
		}

		emailFields["fullname"] = fmt.Sprintf(`%v %v %v`, functions.CamelCase(title, firstname, lastname))

		if emailFields["email"] != nil && emailFields["email"].(string) != "" {
			emailTo = emailFields["email"].(string)
			go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, "", emailFields)
		}
	}
	//SEND AN EMAIL USING TEMPLATE

	sMessage := fmt.Sprintf(`{"error":"%v Welcome Email(s) Sent"}`, len(resMember))
	httpRes.Write([]byte(sMessage))

}

func (this *Member) sendSubscriptionMail(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a User"}`))
		return
	}

	//generateSubscriptionMail
	sqlMember := fmt.Sprintf(`select profile.email as email, profile.title as title, profile.firstname as firstname, profile.lastname as lastname, profile.username as username, scheme.title as scheme, subscription.expirydate as expirydate from subscription  left join profile on profile.control = subscription.membercontrol left join scheme on scheme.control = subscription.schemecontrol where subscription.control in ('0'%s)`, controlList)

	emailTo := ""
	emailFrom := "membership@valued.com"
	emailFromName := "VALUED Membership"
	emailSubject := "WELCOME TO VALUED"
	emailTemplate := "app-subscribe"

	resMember, _ := curdb.Query(sqlMember)
	for _, xDoc := range resMember {
		emailFields := xDoc.(map[string]interface{})
		if emailFields["email"] != nil && emailFields["email"].(string) != "" {
			emailTo = emailFields["email"].(string)
			go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, "", emailFields)
		}
	}
	//generateSubscriptionMail

	sMessage := fmt.Sprintf(`{"error":"%v Subscription Email(s) Sent"}`, len(resMember))
	httpRes.Write([]byte(sMessage))

}

func (this *Member) delete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a Member"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`delete from profile where control in ('0'%s)`, controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Member Deleted","triggerSearch":true}`))
}
