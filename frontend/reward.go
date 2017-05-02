package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"time"

	"regexp"
	"strings"
)

type Reward struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Reward) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	if !strings.Contains(httpReq.FormValue("action"), "download") {
		if httpReq.Method != "POST" {
			sUrl := "/?a=app-reward"
			if this.mapCache["control"] != nil {
				sUrl = "/admin"
			}
			http.Redirect(httpRes, httpReq, sUrl, http.StatusMovedPermanently)
		}
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
		this.pageMap["reward-search"] = searchResult

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Rewards","mainpanelContentSearch":` + contentHTML + `}`))
		return

	case "fetchStore":
		this.fetchRewardStore(httpRes, httpReq, curdb)
		return

	case "fetchScheme":
		this.fetchRewardScheme(httpRes, httpReq, curdb)
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

	case "viewRedeemed":
		this.viewRedeemed(httpRes, httpReq, curdb)
		return

	case "searchRedeemed":
		this.searchRedeemed(httpRes, httpReq, curdb)
		return

	case "viewCoupon":
		this.viewCoupon(httpRes, httpReq, curdb)
		/*
			case "searchCoupon":
				this.searchCoupon(httpRes, httpReq, curdb)

			case "saveCoupon":
				this.saveCoupon(httpRes, httpReq, curdb)

			case "newCoupon":
				this.newCoupon(httpRes, httpReq, curdb)
		*/

	case "replaceCoupon":
		this.replaceCoupon(httpRes, httpReq, curdb)
		return

	case "generateCoupon":
		this.generateCoupon(httpRes, httpReq, curdb)
		return

	case "importCoupon":
		this.importCoupon(httpRes, httpReq, curdb)

	case "importcouponcsvdownload":
		this.importcouponcsvdownload(httpRes, httpReq, curdb)
		return

	case "downloadActiveCoupon":
		this.downloadActiveCoupon(httpRes, httpReq, curdb)
		return

	// case "editCoupon":
	// 	this.editCoupon(httpRes, httpReq, curdb)
	// case "activateCoupon":
	// 	this.activateCoupon(httpRes, httpReq, curdb)
	// case "deactivateCoupon":
	// 	this.deactivateCoupon(httpRes, httpReq, curdb)

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

	case "delete":
		this.delete(httpRes, httpReq, curdb)
		return
	}
}

func (this *Reward) fetchRewardStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("merchant")) == "" {
		sMessage += "Please Select Merchant First<br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	lReadonly := false
	if functions.TrimEscape(httpReq.FormValue("readonly")) != "" {
		lReadonly = true
	}

	rewardKeywords := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("reward")) != "" {
		sqlReward := fmt.Sprintf(`select storecontrol from rewardstore where rewardcontrol = '%s'`, httpReq.FormValue("reward"))

		xDocResult, _ := curdb.Query(sqlReward)
		for _, xDoc := range xDocResult {
			xDoc := xDoc.(map[string]interface{})
			rewardKeywords[xDoc["storecontrol"].(string)] = xDoc
		}
	}

	//Check The Store Table and Link Accordingly
	tblStore := new(database.Store)
	xDocrequest := make(map[string]interface{})
	xDocrequest["workflow"] = "active"
	xDocrequest["merchantcontrol"] = functions.TrimEscape(httpReq.FormValue("merchant"))

	checkBoxView := make(map[string]interface{})
	xDocresult := tblStore.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})

		sState := ""
		if rewardKeywords[xDoc["control"].(string)] != nil {
			sState += " checked "
		}

		if !lReadonly {
			checkBoxView[cNumber+"#rewardstore-checkbox"] = xDoc
		} else {
			if len(sState) > 0 {
				checkBoxView[cNumber+"#rewardstore-checkbox"] = xDoc
			}
			sState += " disabled "
		}
		xDoc["state"] = sState
	}
	//Check The Store Table and Link Accordingly

	viewHTML := strconv.Quote(string(this.Generate(checkBoxView, nil)))
	httpRes.Write([]byte(`{"rewardstorecheckbox":` + viewHTML + `}`))
}

func (this *Reward) fetchRewardScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	lReadonly := false
	if functions.TrimEscape(httpReq.FormValue("readonly")) != "" {
		lReadonly = true
	}

	rewardKeywords := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("reward")) != "" {
		sqlReward := fmt.Sprintf(`select schemecontrol from rewardscheme where rewardcontrol = '%s'`, httpReq.FormValue("reward"))

		xDocResult, _ := curdb.Query(sqlReward)
		for _, xDoc := range xDocResult {
			xDoc := xDoc.(map[string]interface{})
			rewardKeywords[xDoc["schemecontrol"].(string)] = xDoc
		}
	}

	//Check The Scheme Table and Link Accordingly

	tblScheme := new(database.Scheme)
	xDocrequest := make(map[string]interface{})
	xDocrequest["workflow"] = "active"

	checkBoxView := make(map[string]interface{})
	xDocresult := tblScheme.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})

		sState := ""
		if rewardKeywords[xDoc["control"].(string)] != nil {
			sState += " checked "
		}

		if !lReadonly {
			checkBoxView[cNumber+"#rewardscheme-checkbox"] = xDoc
		} else {
			if len(sState) > 0 {
				checkBoxView[cNumber+"#rewardscheme-checkbox"] = xDoc
			}
			sState += " disabled "
		}
		xDoc["state"] = sState
	}

	//Check The Scheme Table and Link Accordingly

	viewHTML := strconv.Quote(string(this.Generate(checkBoxView, nil)))
	httpRes.Write([]byte(`{"rewardschemecheckbox":` + viewHTML + `}`))
}

func (this *Reward) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblReward := new(database.Reward)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["merchantcontrol"] = functions.TrimEscape(httpReq.FormValue("merchantcontrol"))
	xDocresult := tblReward.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		xDoc["title"] = "No Rewards Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"rewardDropdown":` + viewDropdownHtml + `}`))
}

func (this *Reward) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	tblReward := new(database.Reward)
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
	xDocrequest["category"] = functions.TrimEscape(httpReq.FormValue("category"))
	xDocrequest["subcategory"] = functions.TrimEscape(httpReq.FormValue("subcategory"))
	xDocrequest["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDocrequest["type"] = functions.TrimEscape(httpReq.FormValue("type"))
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblReward.Search(xDocrequest, curdb)
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

	xDocresult := tblReward.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#reward-search-result"] = xDoc

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

	return
}

func (this *Reward) new(httpReq *http.Request, curdb database.Database) string {
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

	//Check The Reward Groups and Link Accordingly
	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		formSelection[cNumber+"#reward-edit-group"] = xDoc
	}
	//Check The Reward Groups and Link Accordingly

	formSelection["Public"] = "selected"

	formNew["reward-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Reward) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += "Title is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("type")) == "" {
		sMessage += "Is it a Perk or Privilege is Missing <br>"
	}

	//Handle Coupon Code Import
	mapFiletypes := make(map[string]string)
	mapFiletypes["csvfile"] = "text/plain"
	mapFiles, sMessageImport := functions.UploadFile([]string{"csvfile"}, mapFiletypes, httpReq)
	//Handle Coupon Code Import

	intGenerate := 0
	if functions.TrimEscape(httpReq.FormValue("method")) == "" {
		sMessage += "Redemption Method is Missing <br>"
	} else {

		switch functions.TrimEscape(httpReq.FormValue("method")) {
		case "Client Code Bulk":
			switch functions.TrimEscape(httpReq.FormValue("methodbulk")) {
			case "Import":
				if sMessageImport != "" {
					sMessage += sMessageImport
				}
				if mapFiles["csvfile"] == nil {
					sMessage += "CSV File is Missing"
				}
			case "Generate":
				if functions.TrimEscape(httpReq.FormValue("generate")) == "" {
					sMessage += "Generate Coupons is Missing <br>"
				} else {
					intGenerateErr, err := strconv.Atoi(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("generate"))))
					if err != nil {
						sMessage += err.Error() + "!!! <br>"
					} else {
						intGenerate = intGenerateErr
					}
				}
			}

			// case "Client Code Single":
			// 	if functions.TrimEscape(httpReq.FormValue("coupon")) == "" {
			// 		sMessage += "Redemption Method is Missing <br>"
			// 	}
		}
	}

	regexNumber := regexp.MustCompile("^-*([0-9]+)*\\.*([0-9]+)$")
	if strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxuse"))) == "" {
		sMessage += "Total Uses is missing <br>"
	} else {
		if regexNumber.MatchString(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxuse")))) == false {
			sMessage += "Total Uses must be Numeric!! <br>"
		}
	}

	if strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxperuser"))) == "" {
		sMessage += "Transactions Per Member is missing <br>"
	} else {
		if regexNumber.MatchString(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxperuser")))) == false {
			sMessage += "Transactions Per Member must be Numeric!! <br>"
		}
	}

	if strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxpermonth"))) == "" {
		sMessage += "Transactions per Month is missing <br>"
	} else {
		if regexNumber.MatchString(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxpermonth")))) == false {
			sMessage += "Transactions per Month must be Numeric!! <br>"
		}
	}

	if functions.TrimEscape(httpReq.FormValue("visibleto")) == "" {
		sMessage += "Visible To missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("startdate")) == "" {
		sMessage += "Start Date missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("enddate")) == "" {
		sMessage += "End Date missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("discount")) == "" {
		sMessage += "Discount is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("discounttype")) == "" {
		sMessage += "Discount Type is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("discountvalue")) == "" {
		sMessage += "Discount Value is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("category")) == "" {
		sMessage += "Reward Category is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("subcategory")) == "" {
		sMessage += "Reward Sub-Category is missing <br>"
	}

	if len(httpReq.Form["rewardcategory"]) == 0 || len(httpReq.Form["rewardcategory"]) > 4 {
		sMessage += "Select maximum of 4 Keywords Filter <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("keywords")) == "" {
		sMessage += "Keywords is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("description")) == "" {
		sMessage += "Description is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("rewardstore")) == "" {
		sMessage += "Please Select a Store <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("rewardscheme")) == "" && functions.TrimEscape(httpReq.FormValue("rewardgroup")) == "" {
		sMessage += "Please Select a Scheme or Group <br>"
	} else {
		if functions.TrimEscape(httpReq.FormValue("rewardscheme")) != "" && functions.TrimEscape(httpReq.FormValue("rewardgroup")) != "" {
			sMessage += "This Reward cannot belong to Groups and Schemes <br>"
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	maxuse, err := strconv.ParseFloat(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxuse"))), 10)
	if err != nil {
		sMessage += err.Error() + "!!! <br>"
	}
	xDoc["maxuse"] = maxuse

	maxperuser, err := strconv.ParseFloat(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxperuser"))), 10)
	if err != nil {
		sMessage += err.Error() + "!!! <br>"
	}
	xDoc["maxperuser"] = maxperuser

	maxpermonth, err := strconv.ParseFloat(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("maxpermonth"))), 10)
	if err != nil {
		sMessage += err.Error() + "!!! <br>"
	}
	xDoc["maxpermonth"] = maxpermonth

	if functions.TrimEscape(httpReq.FormValue("discountvalue")) != "" {
		discountvalue, err := strconv.ParseFloat(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("discountvalue"))), 10)
		if err != nil {
			sMessage += err.Error() + "!!! <br>"
		}
		xDoc["discountvalue"] = discountvalue
	}

	if functions.TrimEscape(httpReq.FormValue("orderby")) == "" {
		xDoc["orderby"] = 100
	} else {
		orderby, err := strconv.ParseFloat(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("orderby"))), 10)
		if err != nil {
			sMessage += err.Error() + "!!! <br>"
		}
		xDoc["orderby"] = orderby
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//
	// If Reward Method changes Disable All Coupons
	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		sqlMethod := fmt.Sprintf(`select control from reward where control = '%s' and method = '%s'`,
			functions.TrimEscape(httpReq.FormValue("control")), functions.TrimEscape(httpReq.FormValue("method")))
		resMethod, _ := curdb.Query(sqlMethod)
		if resMethod["1"] == nil {
			// sqlDisable := fmt.Sprintf(`update coupon set workflow = 'inactive' where workflow = 'active' and rewardcontrol = %s'`, functions.TrimEscape(httpReq.FormValue("control")))
			sqlDelete := fmt.Sprintf(`delete from coupon where rewardcontrol = %s'`, functions.TrimEscape(httpReq.FormValue("control")))
			curdb.Query(sqlDelete)
		}
	}
	// If Reward Method changes Disable All Coupons
	//

	tblReward := new(database.Reward)
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))

	xDoc["categorycontrol"] = functions.TrimEscape(httpReq.FormValue("category"))
	xDoc["subcategorycontrol"] = functions.TrimEscape(httpReq.FormValue("subcategory"))
	xDoc["merchantcontrol"] = functions.TrimEscape(httpReq.FormValue("merchant"))

	xDoc["startdate"] = functions.TrimEscape(httpReq.FormValue("startdate"))
	xDoc["enddate"] = functions.TrimEscape(httpReq.FormValue("enddate"))

	xDoc["type"] = functions.TrimEscape(httpReq.FormValue("type"))
	xDoc["method"] = functions.TrimEscape(httpReq.FormValue("method"))
	xDoc["keywords"] = functions.TrimEscape(httpReq.FormValue("keywords"))
	xDoc["beneficiary"] = functions.TrimEscape(httpReq.FormValue("beneficiary"))

	sDescription := functions.TrimEscape(httpReq.FormValue("description"))
	sDescription = strings.Replace(sDescription, "\r", "", -1)
	sDescription = strings.Replace(sDescription, "\n", "<br>", -1)
	xDoc["description"] = sDescription

	sRestriction := functions.TrimEscape(httpReq.FormValue("restriction"))
	sRestriction = strings.Replace(sRestriction, "\r", "", -1)
	sRestriction = strings.Replace(sRestriction, "\n", "<br>", -1)
	xDoc["restriction"] = sRestriction

	xDoc["visibleto"] = functions.TrimEscape(httpReq.FormValue("visibleto"))
	xDoc["discount"] = functions.TrimEscape(httpReq.FormValue("discount"))
	xDoc["discounttype"] = functions.TrimEscape(httpReq.FormValue("discounttype"))

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("reward-%s-%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	xDoc["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	if xDoc["code"].(string) == "" {
		xDoc["code"] = fmt.Sprintf("%s-%s", functions.RandomString(4), functions.RandomString(4))
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblReward.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblReward.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	if xDoc["control"] == nil {
		sMessage = fmt.Sprintf("Error Saving Reward Record")
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sMessage = fmt.Sprintf("Reward %s Saved", xDoc["title"])

	//Handle Reward Category Link
	tblRewardCategoryLink := new(database.CategoryLink)
	curdb.Query(fmt.Sprintf(`delete from categorylink where rewardcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardcategory"] {
		xDocRewardCategoryReqLink := make(map[string]interface{})
		xDocRewardCategoryReqLink["rewardcontrol"] = xDoc["control"]
		xDocRewardCategoryReqLink["workflow"] = "active"
		xDocRewardCategoryReqLink["categorycontrol"] = control
		tblRewardCategoryLink.Create(this.mapCache["username"].(string), xDocRewardCategoryReqLink, curdb)
	}
	//Handle Reward Category Link

	//Handle Coupon Code Import
	var sliceRow []string
	tblCoupon := new(database.Coupon)

	if functions.TrimEscape(httpReq.FormValue("method")) != "" {
		switch functions.TrimEscape(httpReq.FormValue("method")) {
		case "Client Code Bulk":
			switch functions.TrimEscape(httpReq.FormValue("methodbulk")) {
			case "Import":
				if mapFiles["csvfile"] != nil {
					mapCoupon := mapFiles["csvfile"].(map[string]interface{})
					couponList := string(mapCoupon["filebytes"].([]byte))
					couponList = strings.Replace(couponList, "\r", "", -1)
					sliceRow = strings.Split(couponList, "\n")
				}
			case "Generate":
				//Find All Coupons Belonging to Merchant

				sqlCoupon := "select coupon.code as code from coupon left join reward on reward.control = coupon.rewardcontrol where reward.merchantcontrol = '%s'"
				sqlCoupon = fmt.Sprintf(sqlCoupon, xDoc["merchantcontrol"])
				mapResult, _ := curdb.Query(sqlCoupon)
				mapExistingCoupons := make(map[string]bool)

				for _, xDocCoupon := range mapResult {
					xDocCoupon := xDocCoupon.(map[string]interface{})
					mapExistingCoupons[xDocCoupon["code"].(string)] = true
				}

				if intGenerate > 0 {
					sCoupon := functions.RandomString(6)
					intGenerateCounter := 0
					for intGenerateCounter < intGenerate {
						for mapExistingCoupons[sCoupon] {
							sCoupon = functions.RandomString(6)
						}
						mapExistingCoupons[sCoupon] = true
						sliceRow = append(sliceRow, sCoupon)
						intGenerateCounter++
					}
				}
				//Generate and Compare with Existing List
			}
		case "Client Code Single":
			if functions.TrimEscape(httpReq.FormValue("coupon")) != "" {
				//Check if coupon code is unique or matches existing code
				sqlCoupon := fmt.Sprintf(`select control from coupon where workflow = 'active' and rewardcontrol = '%s'`, functions.TrimEscape(httpReq.FormValue("control")))
				resCoupon, _ := curdb.Query(sqlCoupon)
				if resCoupon["1"] == nil {
					sliceRow = append(sliceRow, functions.TrimEscape(httpReq.FormValue("coupon")))
				}
				//Check if coupon code is unique or matches existing code
			}
		}
	}

	if len(sliceRow) > 0 {
		go func() {
			for _, stringCols := range sliceRow {

				stringCols = strings.TrimSpace(stringCols)
				sliceCols := strings.Split(stringCols, ",")

				xDocCoupon := make(map[string]interface{})
				xDocCoupon["code"] = strings.TrimSpace(sliceCols[0])
				xDocCoupon["title"] = strings.TrimSpace(sliceCols[0])
				xDocCoupon["rewardcontrol"] = xDoc["control"]
				xDocCoupon["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))

				tblCoupon.Create(this.mapCache["username"].(string), xDocCoupon, curdb)
				<-time.Tick(time.Millisecond * 15)
			}
		}()

		sMessage += fmt.Sprintf(" <b>%d</b> Coupon Codes Imported", len(sliceRow))
	}
	//Handle Coupon Code Import

	//Handle Scheme Link

	tblRewardScheme := new(database.RewardScheme)
	curdb.Query(fmt.Sprintf(`delete from rewardscheme where rewardcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardscheme"] {
		xDocRewardScheme := make(map[string]interface{})
		xDocRewardScheme["rewardcontrol"] = xDoc["control"]
		xDocRewardScheme["workflow"] = "active"
		xDocRewardScheme["schemecontrol"] = control
		tblRewardScheme.Create(this.mapCache["username"].(string), xDocRewardScheme, curdb)
	}

	//Handle Scheme Link

	//Handle Group Link
	tblRewardGroup := new(database.RewardGroup)
	curdb.Query(fmt.Sprintf(`delete from rewardgroup where rewardcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardgroup"] {
		xDocRewardGroup := make(map[string]interface{})
		xDocRewardGroup["rewardcontrol"] = xDoc["control"]
		xDocRewardGroup["workflow"] = "active"
		xDocRewardGroup["groupcontrol"] = control
		tblRewardGroup.Create(this.mapCache["username"].(string), xDocRewardGroup, curdb)
	}
	//Handle Group Link

	//Handle Store Link
	tblRewardStore := new(database.RewardStore)
	curdb.Query(fmt.Sprintf(`delete from rewardstore where rewardcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardstore"] {
		xDocRewardStore := make(map[string]interface{})
		xDocRewardStore["rewardcontrol"] = xDoc["control"]
		xDocRewardStore["workflow"] = "active"
		xDocRewardStore["storecontrol"] = control
		tblRewardStore.Create(this.mapCache["username"].(string), xDocRewardStore, curdb)
	}
	//Handle Store Link

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"` + sMessage + `","mainpanelContent":` + viewHTML + `}`))
}

func (this *Reward) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblReward := new(database.Reward)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblReward.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocResult := ResMap["1"].(map[string]interface{})

	switch xDocResult["workflow"].(string) {
	case "inactive":
		xDocResult["actionView"] = "activateView"
		xDocResult["actionColor"] = "success"
		xDocResult["actionLabel"] = "Activate"

	case "active":
		xDocResult["actionView"] = "deactivateView"
		xDocResult["actionColor"] = "danger"
		xDocResult["actionLabel"] = "De-Activate"
	}

	rewardType := xDocResult["type"].(string)
	xDocResult[rewardType] = "checked"

	rewardMethod := strings.Replace(xDocResult["method"].(string), " ", "", -1)
	xDocResult[rewardMethod] = "checked"

	rewardVisibleTo := xDocResult["visibleto"].(string)
	xDocResult[rewardVisibleTo] = "checked"

	xDocResult["createdate"] = xDocResult["createdate"].(string)[0:19]

	//Search for Category & Sub Category Title
	new(Category).QuickSearchTitle(xDocResult, curdb)
	//Search for Category & Sub Category Title

	//Check The Reward Groups and Link Accordingly
	tblRewardGroup := new(database.RewardGroup)
	xDocRewardGroupReq := make(map[string]interface{})
	linkedRewardGroup := make(map[string]interface{})

	xDocRewardGroupReq["reward"] = control
	linkListMap := tblRewardGroup.Search(xDocRewardGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedRewardGroup[xDoc["groupcontrol"].(string)] = xDoc["control"]
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedRewardGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#reward-view-group"] = xDoc
	}
	//Check The Reward Groups and Link Accordingly

	formView := make(map[string]interface{})
	formView["reward-view"] = xDocResult
	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Reward) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblReward := new(database.Reward)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblReward.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocResult := ResMap["1"].(map[string]interface{})

	rewardType := xDocResult["type"].(string)
	xDocResult[rewardType] = "checked"

	rewardMethod := strings.Replace(xDocResult["method"].(string), " ", "", -1)
	xDocResult[rewardMethod] = "checked"

	rewardVisibleTo := xDocResult["visibleto"].(string)
	xDocResult[rewardVisibleTo] = "checked"

	rewardDiscontType := xDocResult["discounttype"].(string)
	xDocResult[rewardDiscontType] = "selected"

	xDocResult["formtitle"] = "Edit"

	sDescription := xDocResult["description"].(string)
	xDocResult["description"] = strings.Replace(sDescription, "<br>", "\n", -1)

	sRestriction := xDocResult["restriction"].(string)
	xDocResult["restriction"] = strings.Replace(sRestriction, "<br>", "\n", -1)

	//Search for Category & Sub Category Title
	new(Category).QuickSearchTitle(xDocResult, curdb)
	//Search for Category & Sub Category Title

	//Check The Reward Groups and Link Accordingly
	tblRewardGroup := new(database.RewardGroup)
	xDocRewardGroupReq := make(map[string]interface{})
	linkedRewardGroup := make(map[string]interface{})

	xDocRewardGroupReq["reward"] = xDocResult["control"]
	linkListMap := tblRewardGroup.Search(xDocRewardGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedRewardGroup[xDoc["groupcontrol"].(string)] = xDoc["control"]
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedRewardGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#reward-edit-group"] = xDoc
	}
	//Check The Reward Groups and Link Accordingly

	formView := make(map[string]interface{})
	formView["reward-edit"] = xDocResult

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Reward) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a reward"}`))
		return
	}

	//Send a Reward Activation Email to Merchant
	//Send a Reward Activation Email to Merchant

	curdb.Query(fmt.Sprintf(`update reward set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Reward) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a reward"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update reward set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Reward Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Reward) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a reward"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update reward set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Reward) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a reward"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update reward set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Reward Deactivated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Reward) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Rewards Not De-Activated"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update reward set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Reward) delete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Reward Not Deleted}`))
		return
	}
	//

	curdb.Query(fmt.Sprintf(`delete from reward where control = '%s'; delete from rewardscheme where rewardcontrol = '%s';
		delete from rewardstore where rewardcontrol = '%s'; delete from coupon where rewardcontrol = '%s';`,
		functions.TrimEscape(httpReq.FormValue("control")),
		functions.TrimEscape(httpReq.FormValue("control")),
		functions.TrimEscape(httpReq.FormValue("control")),
		functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}
