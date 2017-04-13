package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"

	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
)

type Profile struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Profile) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		viewHTML := this.view(this.mapCache["control"].(string), curdb)
		httpRes.Write([]byte(`{"pageTitle":"My Profile","mainpanelContent":` + viewHTML + `}`))
		return

	case "quicksearch":
		this.quicksearch(httpRes, httpReq, curdb)
		return

	case "save":
		this.save(httpRes, httpReq, curdb)
		return

	case "edit":
		this.edit(httpRes, httpReq, curdb)
		return
	}
}

func (this *Profile) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblProfile := new(database.Profile)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "20"
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))

	xDocresult := tblProfile.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))

		xDoc["title"] = fmt.Sprintf("%s %s %s", xDoc["title"], xDoc["firstname"], xDoc["lastname"])

		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = html.EscapeString(httpReq.FormValue("tag"))
		xDoc["title"] = "No Profiles Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Profile) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	tblProfile := new(database.Profile)
	xDoc := make(map[string]interface{})

	switch this.mapCache["role"].(string) {
	case "merchant":
		// if functions.TrimEscape(httpReq.FormValue("industry")) == "" {
		// 	sMessage += "Industry is missing <br>"
		// }

		// if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		// 	sMessage += "Company Name is missing <br>"
		// }
		break

	case "employer":
		xDoc["keywords"] = functions.TrimEscape(httpReq.FormValue("keywords"))
		break

	case "member":
		if functions.TrimEscape(httpReq.FormValue("firstname")) == "" {
			sMessage += "First Name is missing <br>"
		}

		if functions.TrimEscape(httpReq.FormValue("lastname")) == "" {
			sMessage += "Last Name is missing <br>"
		}
	}

	if functions.TrimEscape(httpReq.FormValue("phone")) == "" {
		sMessage += "Contact Number is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Primary Notification Email is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("username")) == "" {
		sMessage += "Username is missing <br>"
	}

	switch this.mapCache["role"].(string) {
	case "employer", "merchant":
		if functions.TrimEscape(httpReq.FormValue("website")) == "" {
			sMessage += "Website is missing <br>"
		}
		break
	}

	switch this.mapCache["role"].(string) {
	case "merchant":
		if functions.TrimEscape(httpReq.FormValue("description")) == "" {
			sMessage += "Description is missing <br>"
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//--> Check if Email/Username exists in Login Table
	sqlCheckEmail := fmt.Sprintf(`select control from profile where email = '%s' and control != '%s' `,
		functions.TrimEscape(httpReq.FormValue("email")), this.mapCache["control"])
	mapCheckEmail, _ := curdb.Query(sqlCheckEmail)

	if mapCheckEmail["1"] != nil {
		sMessage += "Email already exists <br>"
	}

	sqlCheckUsername := fmt.Sprintf(`select control from profile where username = '%s' and control != '%s' `,
		functions.TrimEscape(httpReq.FormValue("username")), this.mapCache["control"])
	mapCheckUsername, _ := curdb.Query(sqlCheckUsername)

	if mapCheckUsername["1"] != nil {
		sMessage += "Username already exists <br>"
	}
	//--> Check if Email/Username exists in Login Table

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc["control"] = this.mapCache["control"]
	// xDoc["title"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("title")))
	// xDoc["firstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("firstname")))
	// xDoc["lastname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("lastname")))

	xDoc["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDoc["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))

	if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
		xDoc["phonecode"] = "+971"
	} else {
		xDoc["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
	}

	xDoc["description"] = functions.TrimEscape(httpReq.FormValue("description"))
	xDoc["username"] = functions.TrimEscape(httpReq.FormValue("username"))

	switch this.mapCache["role"].(string) {
	case "employer", "merchant":
		xDoc["emailsecondary"] = functions.TrimEscape(httpReq.FormValue("emailsecondary"))
		xDoc["emailalternate"] = functions.TrimEscape(httpReq.FormValue("emailalternate"))
		xDoc["website"] = functions.TrimEscape(httpReq.FormValue("website"))
		// xDoc["industrycontrol"] = functions.TrimEscape(httpReq.FormValue("industry"))

		switch this.mapCache["role"].(string) {
		case "employer":
			break
		case "merchant":
			// xDoc["subindustrycontrol"] = functions.TrimEscape(httpReq.FormValue("subindustry"))
			// xDoc["keywords"] = functions.TrimEscape(httpReq.FormValue("keywords"))

			//Handle Reward Category Link
			// tblRewardCategoryLink := new(database.CategoryLink)
			// curdb.Query(fmt.Sprintf(`delete from categorylink where merchantcontrol = '%s'`, xDoc["control"]))
			// for _, control := range httpReq.Form["reviewcategory"] {
			// 	xDocRewardCategoryReqLink := make(map[string]interface{})
			// 	xDocRewardCategoryReqLink["merchantcontrol"] = xDoc["control"]
			// 	xDocRewardCategoryReqLink["workflow"] = "active"
			// 	xDocRewardCategoryReqLink["categorycontrol"] = control
			// 	tblRewardCategoryLink.Create(this.mapCache["username"].(string), xDocRewardCategoryReqLink, curdb)
			// }
			//Handle Reward Category Link

			//Handle Review Category Link
			// tblReviewCategoryLink := new(database.ReviewCategoryLink)
			// curdb.Query(fmt.Sprintf(`delete from reviewcategorylink where merchantcontrol = '%s'`, xDoc["control"]))
			// for _, control := range httpReq.Form["reviewcategory"] {
			// 	xDocReviewCategoryReqLink := make(map[string]interface{})
			// 	xDocReviewCategoryReqLink["merchantcontrol"] = xDoc["control"]
			// 	xDocReviewCategoryReqLink["workflow"] = "active"
			// 	xDocReviewCategoryReqLink["reviewcategorycontrol"] = control
			// 	tblReviewCategoryLink.Create(this.mapCache["username"].(string), xDocReviewCategoryReqLink, curdb)
			// }
			//Handle Review Category Link
			break
		}
	}

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("profile_%s_%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	tblProfile.Update(this.mapCache["username"].(string), xDoc, curdb)

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"Record Saved","mainpanelContent":` + viewHTML + `}`))
	return
}

func (this *Profile) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblProfile := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblProfile.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	xDocRequest["createdate"] = xDocRequest["createdate"].(string)[0:10]

	nCurYear, _ := strconv.Atoi(functions.GetSystemTime()[6:10])
	nDobYear := nCurYear
	if len(xDocRequest["dob"].(string)) == 10 {
		nDobYear, _ = strconv.Atoi(xDocRequest["dob"].(string)[6:10])
	}
	xDocRequest["age"] = nCurYear - nDobYear

	switch this.mapCache["role"].(string) {
	case "employer":

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
				xDocRequest[cNumber+"#profile-view-employer-group"] = xDoc
			}
		}
		//Check The Employer Groups and Link Accordingly

	case "merchant":

		//Check The Reward Categories and Link Accordingly
		tblRewardCategoryLink := new(database.CategoryLink)
		xDocRewardCategoryReqLink := make(map[string]interface{})
		linkedRewardCategoryegory := make(map[string]interface{})

		xDocRewardCategoryReqLink["merchant"] = control
		linkListMap := tblRewardCategoryLink.Search(xDocRewardCategoryReqLink, curdb)
		for _, xDoc := range linkListMap {
			xDoc := xDoc.(map[string]interface{})
			linkedRewardCategoryegory[xDoc["categorycontrol"].(string)] = xDoc["control"]
		}

		tblRewardCategory := new(database.Category)
		xDocRewardCategoryReq := make(map[string]interface{})
		xDocRewardCategoryReq["workflow"] = "active"
		listMap := tblRewardCategory.Search(xDocRewardCategoryReq, curdb)

		sliceKeyword := []string{}
		for _, xDoc := range listMap {
			xDoc := xDoc.(map[string]interface{})
			if linkedRewardCategoryegory[xDoc["control"].(string)] != nil {
				sliceKeyword = append(sliceKeyword, xDoc["title"].(string))
			}
		}
		xDocRequest["keywords"] = strings.Join(sliceKeyword, ",")
		//Check The Reward Categories and Link Accordingly

		//Check The Review Categories and Link Accordingly
		tblReviewCategoryLink := new(database.ReviewCategoryLink)
		xDocReviewCategoryReqLink := make(map[string]interface{})
		linkedReviewCategory := make(map[string]interface{})

		xDocReviewCategoryReqLink["merchant"] = control
		linkListMap = tblReviewCategoryLink.Search(xDocReviewCategoryReqLink, curdb)
		for _, xDoc := range linkListMap {
			xDoc := xDoc.(map[string]interface{})
			linkedReviewCategory[xDoc["reviewcategorycontrol"].(string)] = xDoc["control"]
		}

		tblReviewCategory := new(database.ReviewCategory)
		xDocReviewCategoryReq := make(map[string]interface{})
		xDocReviewCategoryReq["workflow"] = "active"
		listMap = tblReviewCategory.Search(xDocReviewCategoryReq, curdb)

		for cNumber, xDoc := range listMap {
			xDoc := xDoc.(map[string]interface{})
			if linkedReviewCategory[xDoc["control"].(string)] != nil {
				xDoc["state"] = "checked"
			}

			xDocRequest[cNumber+"#merchant-view-checkbox"] = xDoc
		}
		//Check The Review Categories and Link Accordingly
	}

	formView := make(map[string]interface{})
	cTemplateFile := fmt.Sprintf("profile-view-%s", this.mapCache["role"])
	formView[cTemplateFile] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Profile) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblProfile := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = this.mapCache["control"]

	ResMap := tblProfile.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})
	xDocRequest["formtitle"] = "Edit Profile"

	switch this.mapCache["role"].(string) {
	case "merchant":
		//Check The Reward Categories and Link Accordingly
		tblRewardCategoryLink := new(database.CategoryLink)
		xDocRewardCategoryReqLink := make(map[string]interface{})
		linkedRewardCategoryegory := make(map[string]interface{})

		xDocRewardCategoryReqLink["merchant"] = this.mapCache["control"].(string)
		linkListMap := tblRewardCategoryLink.Search(xDocRewardCategoryReqLink, curdb)
		for _, xDoc := range linkListMap {
			xDoc := xDoc.(map[string]interface{})
			linkedRewardCategoryegory[xDoc["categorycontrol"].(string)] = xDoc["control"]
		}

		tblRewardCategory := new(database.Category)
		xDocRewardCategoryReq := make(map[string]interface{})
		xDocRewardCategoryReq["workflow"] = "active"
		listMap := tblRewardCategory.Search(xDocRewardCategoryReq, curdb)

		sliceKeyword := []string{}
		// sliceKeywordControl := []string{}
		for _, xDoc := range listMap {
			xDoc := xDoc.(map[string]interface{})
			if linkedRewardCategoryegory[xDoc["control"].(string)] != nil {
				sliceKeyword = append(sliceKeyword, xDoc["control"].(string))
				// sliceKeywordControl = append(sliceKeyword, xDoc["title"].(string))
			}
		}
		xDocRequest["keywords"] = strings.Join(sliceKeyword, ",")
		//Check The Reward Categories and Link Accordingly

		//Check The Review Categories and Link Accordingly
		tblReviewCategoryLink := new(database.ReviewCategoryLink)
		xDocReviewCategoryReqLink := make(map[string]interface{})
		linkedReviewCategory := make(map[string]interface{})

		xDocReviewCategoryReqLink["merchant"] = this.mapCache["control"].(string)
		linkListMap = tblReviewCategoryLink.Search(make(map[string]interface{}), curdb)
		for _, xDoc := range linkListMap {
			xDoc := xDoc.(map[string]interface{})
			linkedReviewCategory[xDoc["reviewcategorycontrol"].(string)] = xDoc["control"]
		}

		tblReviewCategory := new(database.ReviewCategory)
		xDocReviewCategoryReq := make(map[string]interface{})
		xDocReviewCategoryReq["workflow"] = "active"
		listMap = tblReviewCategory.Search(xDocReviewCategoryReq, curdb)

		for cNumber, xDoc := range listMap {
			xDoc := xDoc.(map[string]interface{})
			if linkedReviewCategory[xDoc["control"].(string)] != nil {
				xDoc["state"] = "checked"
			}
			xDocRequest[cNumber+"#merchant-edit-checkbox"] = xDoc
		}
		//Check The Review Categories and Link Accordingly
		break
	}

	xDocRequest[strings.Replace(xDocRequest["phonecode"].(string), "+", "", 1)] = "selected"

	cTemplateFile := fmt.Sprintf("profile-edit-%s", this.mapCache["role"])
	formView[cTemplateFile] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}
