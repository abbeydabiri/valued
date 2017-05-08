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

type Merchant struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Merchant) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/admin", http.StatusMovedPermanently)
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
		this.pageMap["merchant-search"] = searchResult

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Merchants","mainpanelContentSearch":` + contentHTML + `}`))
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
		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), functions.TrimEscape(httpReq.FormValue("subview")), curdb)
		httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
		return

	case "viewStore":
		this.viewStore(httpRes, httpReq, curdb)
		return

	case "searchStore":
		this.searchStore(httpRes, httpReq, curdb)
		return

	case "viewReward":
		this.viewReward(httpRes, httpReq, curdb)
		return

	case "searchReward":
		this.searchReward(httpRes, httpReq, curdb)
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

	case "activateAll":
		this.activateAll(httpRes, httpReq, curdb)
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

func (this *Merchant) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblMerchant := new(database.Profile)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["role"] = "merchant"
	xDocrequest["company"] = "Yes"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocresult := tblMerchant.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		xDoc["title"] = "No Merchants Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Merchant) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	tblMerchant := new(database.Profile)
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
	xDocrequest["category"] = functions.TrimEscape(httpReq.FormValue("category"))
	xDocrequest["lastname"] = functions.TrimEscape(httpReq.FormValue("lastname"))
	xDocrequest["firstname"] = functions.TrimEscape(httpReq.FormValue("firstname"))
	xDocrequest["status"] = functions.TrimEscape(httpReq.FormValue("status"))
	xDocrequest["employercode"] = "main"
	xDocrequest["company"] = "Yes"
	xDocrequest["role"] = "merchant"

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblMerchant.Search(xDocrequest, curdb)
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

	xDocresult := tblMerchant.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#merchant-search-result"] = xDoc

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

func (this *Merchant) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"

	//Check The Profile and Link Accordingly
	tblRole := new(database.Role)
	xDocRoleReq := make(map[string]interface{})
	xDocRoleReq["workflow"] = "active"
	listMapRole := tblRole.Search(xDocRoleReq, curdb)

	for cNumber, xDoc := range listMapRole {
		xDoc := xDoc.(map[string]interface{})
		if xDoc["code"].(string) == "merchant" {
			xDoc["state"] = "checked"
		}
		formSelection[cNumber+"#role-edit-checkbox"] = xDoc
	}
	//Check The Profile and Link Accordingly

	//Check The Review Categories and Link Accordingly
	tblReviewCategory := new(database.ReviewCategory)
	xDocReviewCategoryReq := make(map[string]interface{})
	xDocReviewCategoryReq["workflow"] = "active"
	listMap := tblReviewCategory.Search(xDocReviewCategoryReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		formSelection[cNumber+"#merchant-edit-checkbox"] = xDoc
	}
	//Check The Review Categories and Link Accordingly

	formNew["merchant-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Merchant) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("status")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(httpReq.Form["profilerole"]) == 0 {
		sMessage += "Select Profile Role <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("category")) == "" {
		sMessage += "Category is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("subcategory")) == "" {
		sMessage += "Sub-Category is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += "Company Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("phone")) == "" {
		sMessage += "Contact Number is missing <br>"
	}

	// if functions.TrimEscape(httpReq.FormValue("firstname")) == "" {
	// 	sMessage += "First Name is missing <br>"
	// }

	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Primary Email is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("website")) == "" {
		sMessage += "Website is missing <br>"
	}

	if len(httpReq.Form["rewardcategory"]) == 0 {
		sMessage += "Select Extra Filters <br>"
	}

	// if functions.TrimEscape(httpReq.FormValue("keywords")) == "" {
	// 	sMessage += "Keywords is missing <br>"
	// }

	if functions.TrimEscape(httpReq.FormValue("description")) == "" {
		sMessage += "Description is missing <br>"
	}

	if len(httpReq.Form["reviewcategory"]) != 4 {
		sMessage += "Select only 4 Review Categories <br>"
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
	// 	if functions.TrimEscape(httpReq.FormValue("industry")) == "" {
	// 		xDoc["industrycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	// 	}

	// 	if functions.TrimEscape(httpReq.FormValue("subindustry")) == "" {
	// 		xDoc["subindustrycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	// 	}
	// }

	defaultMap, _ = curdb.Query(`select control from category where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Category <br>"
	} else {
		if functions.TrimEscape(httpReq.FormValue("category")) == "" {
			xDoc["categorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
		}

		if functions.TrimEscape(httpReq.FormValue("subcategory")) == "" {
			xDoc["subcategorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
		}
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

	tblMerchant := new(database.Profile)
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDoc["status"] = functions.TrimEscape(httpReq.FormValue("status"))

	// xDoc["firstname"] = functions.TrimEscape(httpReq.FormValue("firstname"))
	// xDoc["lastname"] = functions.TrimEscape(httpReq.FormValue("lastname"))

	xDoc["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDoc["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))

	if functions.TrimEscape(httpReq.FormValue("pincode")) != "" {
		xDoc["pincode"] = functions.TrimEscape(httpReq.FormValue("pincode"))
	} else {
		xDoc["pincode"] = "4321"
	}

	if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
		xDoc["phonecode"] = "+971"
	} else {
		xDoc["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
	}

	xDoc["emailsecondary"] = functions.TrimEscape(httpReq.FormValue("emailsecondary"))
	xDoc["emailalternate"] = functions.TrimEscape(httpReq.FormValue("emailalternate"))

	xDoc["website"] = functions.TrimEscape(httpReq.FormValue("website"))
	xDoc["keywords"] = functions.TrimEscape(httpReq.FormValue("keywords"))

	xDoc["categorycontrol"] = functions.TrimEscape(httpReq.FormValue("category"))
	xDoc["subcategorycontrol"] = functions.TrimEscape(httpReq.FormValue("subcategory"))

	sDescription := functions.TrimEscape(httpReq.FormValue("description"))
	sDescription = strings.Replace(sDescription, "\r", "", -1)
	sDescription = strings.Replace(sDescription, "\n", "<br>", -1)
	xDoc["description"] = sDescription
	xDoc["company"] = "Yes"

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("merchant_%s_%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblMerchant.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblMerchant.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	//Handle Reward Category Link
	tblRewardCategoryLink := new(database.CategoryLink)
	curdb.Query(fmt.Sprintf(`delete from categorylink where merchantcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardcategory"] {
		xDocRewardCategoryReqLink := make(map[string]interface{})
		xDocRewardCategoryReqLink["merchantcontrol"] = xDoc["control"]
		xDocRewardCategoryReqLink["workflow"] = "active"
		xDocRewardCategoryReqLink["categorycontrol"] = control
		tblRewardCategoryLink.Create(this.mapCache["username"].(string), xDocRewardCategoryReqLink, curdb)
	}
	//Handle Reward Category Link

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

	//Handle Review Category Link
	tblReviewCategoryLink := new(database.ReviewCategoryLink)
	curdb.Query(fmt.Sprintf(`delete from reviewcategorylink where merchantcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["reviewcategory"] {
		xDocReviewCategoryReqLink := make(map[string]interface{})
		xDocReviewCategoryReqLink["merchantcontrol"] = xDoc["control"]
		xDocReviewCategoryReqLink["workflow"] = "active"
		xDocReviewCategoryReqLink["reviewcategorycontrol"] = control
		tblReviewCategoryLink.Create(this.mapCache["username"].(string), xDocReviewCategoryReqLink, curdb)
	}
	//Handle Review Category Link

	sendEmail := ""
	if functions.TrimEscape(httpReq.FormValue("sendEmail")) == "Yes" {
		sendEmail = fmt.Sprintf(`"getform":"/merchant?action=sendWelcomeMail&control=%s"`, xDoc["control"])
	}

	viewHTML := this.view(xDoc["control"].(string), "", curdb)
	httpRes.Write([]byte(`{` + sendEmail + `"error":"Record Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *Merchant) view(control, subview string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblMerchant := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblMerchant.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocResult := ResMap["1"].(map[string]interface{})

	switch subview {
	case "reward":
		xDocResult["subview"] = "viewReward"

	case "store":
		xDocResult["subview"] = "viewStore"

	default:
		//case "user":
		xDocResult["subview"] = "viewUser"
	}

	// employerTag := "Employer" + xDocResult["employer"].(string)
	// xDocResult[employerTag] = "checked"

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

	//Search for Category & Sub Category Title
	//new(Category).QuickSearchTitle(xDocResult, curdb)
	//Search for Category & Sub Category Title

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

	//Check The Review Categories and Link Accordingly
	tblReviewCategoryLink := new(database.ReviewCategoryLink)
	xDocReviewCategoryReqLink := make(map[string]interface{})
	linkedReviewCategory := make(map[string]interface{})

	xDocReviewCategoryReqLink["merchant"] = control
	linkListMap := tblReviewCategoryLink.Search(xDocReviewCategoryReqLink, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedReviewCategory[xDoc["reviewcategorycontrol"].(string)] = xDoc["control"]
	}

	tblReviewCategory := new(database.ReviewCategory)
	xDocReviewCategoryReq := make(map[string]interface{})
	xDocReviewCategoryReq["workflow"] = "active"
	listMap := tblReviewCategory.Search(xDocReviewCategoryReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedReviewCategory[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#merchant-view-checkbox"] = xDoc
	}
	//Check The Review Categories and Link Accordingly

	formView := make(map[string]interface{})
	formView["merchant-view"] = xDocResult

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Merchant) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblMerchant := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblMerchant.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocResult := ResMap["1"].(map[string]interface{})
	formView := make(map[string]interface{})
	xDocResult["formtitle"] = "Edit"

	// employerTag := "Employer" + xDocResult["employer"].(string)
	// xDocResult[employerTag] = "checked"

	sDescription := xDocResult["description"].(string)
	xDocResult["description"] = strings.Replace(sDescription, "<br>", "\n", -1)

	//Search for Category & Sub Category Title
	//new(Category).QuickSearchTitle(xDocResult, curdb)
	//Search for Category & Sub Category Title

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

	//Check The Review Categories and Link Accordingly
	linkedReviewCategory := make(map[string]interface{})
	sqlReviewCategory := fmt.Sprintf(`select reviewcategorycontrol from reviewcategorylink where merchantcontrol = '%s'`, httpReq.FormValue("control"))
	linkListMap, _ := curdb.Query(sqlReviewCategory)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedReviewCategory[xDoc["reviewcategorycontrol"].(string)] = xDoc
	}

	tblReviewCategory := new(database.ReviewCategory)
	xDocReviewCategoryReq := make(map[string]interface{})
	xDocReviewCategoryReq["workflow"] = "active"
	listMap := tblReviewCategory.Search(xDocReviewCategoryReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedReviewCategory[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#merchant-edit-checkbox"] = xDoc
	}
	//Check The Review Categories and Link Accordingly

	xDocResult[strings.Replace(xDocResult["phonecode"].(string), "+", "", 1)] = "selected"

	formView["merchant-edit"] = xDocResult
	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Merchant) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a Merchant"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Merchant) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a Merchant"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), "", curdb)
	httpRes.Write([]byte(`{"error":"Merchant Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Merchant) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a Merchant"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Merchant) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a Merchant"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), "", curdb)
	httpRes.Write([]byte(`{"error":"Merchant Deactivated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Merchant) activateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a Merchant"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update profile set status = 'active', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Merchants Activated","triggerSearch":true}`))
}

func (this *Merchant) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a Merchant"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update profile set status = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Merchants De-Activated","triggerSearch":true}`))
}

func (this *Merchant) sendWelcomeMail(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a User"}`))
		return
	}

	emailTemplate := "merchant-company-welcome"

	//SEND AN EMAIL USING TEMPLATE

	sqlMember := fmt.Sprintf(`select profile.email as email, profile.title as title, profile.username as username, 
		profile.password as password from profile where profile.control in ('0'%s)`, controlList)

	// sqlMember := fmt.Sprintf(`select employer.email as email, employer.title as employertitle, profile.username as username,
	// 	profile.password as password from profile left join profile as employer on employer.control = profile.employercontrol
	// 	where profile.control in ('0'%s)`, controlList)

	emailTo := ""
	emailFrom := "partnership@valued.com"
	emailFromName := "VALUED PARTNERSHIP"
	emailSubject := fmt.Sprintf("VALUED GOES LIVE")
	// emailSubject := fmt.Sprintf("WELCOME TO VALUED - PARTNERSHIP")

	//
	//
	//Add Temporary Free Merchant Subscription for Merchants
	//--//Search for tellafriend Reward if Active generate a CouponCode
	sqlFindReward := `select reward.method as rewardmethod, reward.control as rewardcontrol from reward 
						left join profile as merchant on merchant.control = reward.merchantcontrol
						where lower(reward.code) = 'free.merchant' AND lower(merchant.code) = 'main' AND reward.workflow = 'active'`
	sqlFindRewardRes, _ := curdb.Query(sqlFindReward)
	//--//Search for tellafriend Reward if Active generate a CouponCode
	//Add Temporary Free Merchant Subscription for Merchants
	//
	//

	sUserDetail := `USER: <br> Username: %v <br> Password: %v <br><br>`
	resMember, _ := curdb.Query(sqlMember)
	for _, xDoc := range resMember {
		xDocEmail := xDoc.(map[string]interface{})
		emailFields := make(map[string]interface{})
		if xDocEmail["email"] != nil && xDocEmail["email"].(string) != "" {
			emailFields["email"] = xDocEmail["email"]

			//
			//
			sCoupon := ""
			//Add Temporary Free Merchant Subscription for Merchants
			//--//Search for tellafriend Reward if Active generate a CouponCode
			if sqlFindRewardRes["1"] != nil {
				xDocFindRewardRes := sqlFindRewardRes["1"].(map[string]interface{})

				switch xDocFindRewardRes["rewardmethod"].(string) {
				case "Pin", "Valued Code":

					xDocCoupon := make(map[string]interface{})
					sCoupon = functions.RandomString(6)

					xDocCoupon["code"] = sCoupon
					xDocCoupon["workflow"] = "active"
					xDocCoupon["rewardcontrol"] = xDocFindRewardRes["rewardcontrol"]
					new(database.Coupon).Create(this.mapCache["username"].(string), xDocCoupon, curdb)

				case "Client Code":
					sqlFindCoupon := `select coupon.code as code from coupon where rewardcontrol = '%s' AND workflow = 'active' order by control limit 1`
					sqlFindCoupon = fmt.Sprintf(sqlFindCoupon, xDocFindRewardRes["rewardcontrol"])

					sqlFindCouponRes, _ := curdb.Query(sqlFindCoupon)
					if sqlFindCouponRes["1"] != nil {
						xDocFindCoupon := sqlFindCouponRes["1"].(map[string]interface{})
						if xDocFindCoupon["code"] != nil {
							sCoupon = xDocFindCoupon["code"].(string)
						}
					}
				}
			}
			emailFields["coupon"] = sCoupon
			//--//Search for tellafriend Reward if Active generate a CouponCode
			//Add Temporary Free Merchant Subscription for Merchants
			//
			//

			emailTo = emailFields["email"].(string)
			emailFields["title"] = xDocEmail["title"]
			emailFields["userdata"] = fmt.Sprintf(sUserDetail, xDocEmail["username"], xDocEmail["password"])
			go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, []string, emailFields)

		}
	}
	//SEND AN EMAIL USING TEMPLATE

	sMessage := fmt.Sprintf(`{"error":"%v Welcome Email(s) Sent"}`, len(resMember))
	httpRes.Write([]byte(sMessage))

}

func (this *Merchant) delete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Merchant Not Deleted}`))
		return
	}
	//
	curdb.Query(fmt.Sprintf(`delete from profile where control = '%s'`, functions.TrimEscape(httpReq.FormValue("control"))))
	httpRes.Write([]byte(`{"triggerSearch":true}`))
}
