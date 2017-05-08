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

type MerchantReward struct {
	MerchantControl string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *MerchantReward) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if !strings.HasSuffix(httpReq.FormValue("action"), "download") {
		if httpReq.Method != "POST" {
			http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
		}
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
		// if this.mapCache["rewardSearch"] != nil {
		// 	httpReq = this.mapCache["rewardSearch"].(*http.Request)
		// }

		this.pageMap = make(map[string]interface{})

		searchResult, searchPagination := this.search(httpReq, curdb)
		for sKey, iPagination := range searchPagination {
			searchResult[sKey] = iPagination
		}
		this.pageMap["merchantreward-search"] = searchResult

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Rewards","mainpanelContentSearch":` + contentHTML + `}`))
		return

	case "fetchMerchantRewardStore":
		this.fetchMerchantRewardStore(httpRes, httpReq, curdb)
		return

	case "fetchMerchantRewardScheme":
		this.fetchMerchantRewardScheme(httpRes, httpReq, curdb)
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
	case "searchCoupon":
		this.searchCoupon(httpRes, httpReq, curdb)
	case "saveCoupon":
		this.saveCoupon(httpRes, httpReq, curdb)
	case "newCoupon":
		this.newCoupon(httpRes, httpReq, curdb)
	case "importCoupon":
		this.importCoupon(httpRes, httpReq, curdb)

	case "importcouponcsvdownload":
		this.importcouponcsvdownload(httpRes, httpReq, curdb)
		return

	case "importcouponcsvsave":
		this.importcouponcsvsave(httpRes, httpReq, curdb)
		return

	case "editCoupon":
		this.editCoupon(httpRes, httpReq, curdb)
	case "activateCoupon":
		this.activateCoupon(httpRes, httpReq, curdb)
	case "deactivateCoupon":
		this.deactivateCoupon(httpRes, httpReq, curdb)

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

	case "requestChanges":
		this.requestChanges(httpRes, httpReq, curdb)
		return

	case "deactivateAll":
		this.deactivateAll(httpRes, httpReq, curdb)
		return
	}
}

func (this *MerchantReward) fetchMerchantRewardStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	lReadonly := false
	if functions.TrimEscape(httpReq.FormValue("readonly")) != "" {
		lReadonly = true
	}

	rewardKeywords := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("reward")) != "" {
		sqlMerchantReward := fmt.Sprintf(`select storecontrol from rewardstore where rewardcontrol = '%s'`, httpReq.FormValue("reward"))

		xDocResult, _ := curdb.Query(sqlMerchantReward)
		for _, xDoc := range xDocResult {
			xDoc := xDoc.(map[string]interface{})
			rewardKeywords[xDoc["storecontrol"].(string)] = xDoc
		}
	}

	//Check The Store Table and Link Accordingly
	tblStore := new(database.Store)
	xDocrequest := make(map[string]interface{})
	xDocrequest["workflow"] = "active"
	xDocrequest["merchantcontrol"] = this.MerchantControl

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

func (this *MerchantReward) fetchMerchantRewardScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	lReadonly := false
	if functions.TrimEscape(httpReq.FormValue("readonly")) != "" {
		lReadonly = true
	}

	rewardKeywords := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("reward")) != "" {
		sqlMerchantReward := fmt.Sprintf(`select schemecontrol from rewardscheme where rewardcontrol = '%s'`, httpReq.FormValue("reward"))

		xDocResult, _ := curdb.Query(sqlMerchantReward)
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

func (this *MerchantReward) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblMerchantReward := new(database.Reward)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["merchantcontrol"] = this.MerchantControl
	xDocresult := tblMerchantReward.Search(xDocrequest, curdb)

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

func (this *MerchantReward) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	tblMerchantReward := new(database.Reward)
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
	xDocrequest["category"] = functions.TrimEscape(httpReq.FormValue("category"))
	xDocrequest["subcategory"] = functions.TrimEscape(httpReq.FormValue("subcategory"))
	xDocrequest["email"] = functions.TrimEscape(httpReq.FormValue("email"))

	xDocrequest["merchantcontrol"] = this.MerchantControl
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblMerchantReward.Search(xDocrequest, curdb)
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

	xDocresult := tblMerchantReward.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#merchantreward-search-result"] = xDoc

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

func (this *MerchantReward) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formSelection["merchantcontrol"] = this.MerchantControl

	//Check The MerchantReward Groups and Link Accordingly
	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		formSelection[cNumber+"#reward-edit-group"] = xDoc
	}
	//Check The MerchantReward Groups and Link Accordingly

	formSelection["Public"] = "selected"

	formNew["merchantreward-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *MerchantReward) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""

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

	// if functions.TrimEscape(httpReq.FormValue("rewardscheme")) == "" && functions.TrimEscape(httpReq.FormValue("rewardgroup")) == "" {
	// 	sMessage += "Please Select a Scheme or Group <br>"
	// } else {
	// 	if functions.TrimEscape(httpReq.FormValue("rewardscheme")) != "" && functions.TrimEscape(httpReq.FormValue("rewardgroup")) != "" {
	// 		sMessage += "This Reward cannot belong to Groups and Schemes <br>"
	// 	}
	// }

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

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		if functions.TrimEscape(httpReq.FormValue("orderby")) == "" {
			xDoc["orderby"] = 100
		} else {
			orderby, err := strconv.ParseFloat(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("orderby"))), 10)
			if err != nil {
				sMessage += err.Error() + "!!! <br>"
			}
			xDoc["orderby"] = orderby
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//
	// If Reward Method changes Disable All Coupons
	// if functions.TrimEscape(httpReq.FormValue("control")) != "" {
	// 	sqlMethod := fmt.Sprintf(`select control from reward where control = '%s' and method = '%s'`,
	// 		functions.TrimEscape(httpReq.FormValue("control")), functions.TrimEscape(httpReq.FormValue("method")))
	// 	resMethod, _ := curdb.Query(sqlMethod)
	// 	if resMethod["1"] == nil {
	// 		sqlDisable := fmt.Sprintf(`update coupon set workflow = 'inactive' where rewardcontrol = %s'`, functions.TrimEscape(httpReq.FormValue("control")))
	// 		curdb.Query(sqlDisable)
	// 	}
	// }
	// If Reward Method changes Disable All Coupons
	//

	tblMerchantReward := new(database.Reward)
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDoc["workflow"] = "inactive"

	xDoc["categorycontrol"] = functions.TrimEscape(httpReq.FormValue("category"))
	xDoc["subcategorycontrol"] = functions.TrimEscape(httpReq.FormValue("subcategory"))
	xDoc["merchantcontrol"] = this.MerchantControl

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
			fileName := fmt.Sprintf("merchantreward-%s-%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	sAction := ""
	sMessage = `Your reward addition request has been successfully submitted for approval`

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		sMessage = fmt.Sprintf("Change Requested")
		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
		httpRes.Write([]byte(`{"error":"` + sMessage + `","mainpanelContent":` + viewHTML + `}`))
		return

		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		// tblMerchantReward.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		sAction = "Approval for"
		xDoc["code"] = fmt.Sprintf("RWD-%s-%s", functions.RandomString(4), functions.RandomString(4))
		xDoc["control"] = tblMerchantReward.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	if xDoc["control"] == nil {
		sMessage = fmt.Sprintf("Error Saving Reward Record")
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sMessage = fmt.Sprintf("Reward %s Saved", xDoc["title"])

	//Handle MerchantReward Category Link
	tblMerchantRewardCategoryLink := new(database.CategoryLink)
	curdb.Query(fmt.Sprintf(`delete from categorylink where rewardcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardcategory"] {
		xDocMerchantRewardCategoryReqLink := make(map[string]interface{})
		xDocMerchantRewardCategoryReqLink["rewardcontrol"] = xDoc["control"]
		xDocMerchantRewardCategoryReqLink["workflow"] = "active"
		xDocMerchantRewardCategoryReqLink["categorycontrol"] = control
		tblMerchantRewardCategoryLink.Create(this.mapCache["username"].(string), xDocMerchantRewardCategoryReqLink, curdb)
	}
	//Handle MerchantReward Category Link

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
			//Check if coupon code is unique or matches existing code
			sqlCoupon := fmt.Sprintf(`select control from coupon where workflow = 'active' and rewardcontrol = '%s'`, functions.TrimEscape(httpReq.FormValue("control")))
			resCoupon, _ := curdb.Query(sqlCoupon)
			if resCoupon["1"] == nil {
				sliceRow = append(sliceRow, functions.TrimEscape(httpReq.FormValue("coupon")))
			}
			//Check if coupon code is unique or matches existing code
		}
	}

	if len(sliceRow) > 0 {
		go func() {
			for _, stringCols := range sliceRow {

				stringCols = strings.TrimSpace(stringCols)
				sliceCols := strings.Split(stringCols, ",")

				if sliceCols[0] == "" {
					continue
				}

				xDocCoupon := make(map[string]interface{})
				xDocCoupon["code"] = strings.TrimSpace(sliceCols[0])
				xDocCoupon["title"] = strings.TrimSpace(sliceCols[0])
				xDocCoupon["rewardcontrol"] = xDoc["control"]
				xDocCoupon["workflow"] = "active"

				tblCoupon.Create(this.mapCache["username"].(string), xDocCoupon, curdb)
				<-time.Tick(time.Millisecond * 15)
			}
		}()

		sMessage += fmt.Sprintf("Coupon Codes Imported")
	}
	//Handle Coupon Code Import

	//Handle Scheme Link
	tblMerchantRewardScheme := new(database.RewardScheme)
	curdb.Query(fmt.Sprintf(`delete from rewardscheme where rewardcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardscheme"] {
		xDocMerchantRewardScheme := make(map[string]interface{})
		xDocMerchantRewardScheme["rewardcontrol"] = xDoc["control"]
		xDocMerchantRewardScheme["workflow"] = "active"
		xDocMerchantRewardScheme["schemecontrol"] = control
		tblMerchantRewardScheme.Create(this.mapCache["username"].(string), xDocMerchantRewardScheme, curdb)
	}
	//Handle Scheme Link

	//Handle Group Link - RewardGroup Saving disabled for merchants
	// tblMerchantRewardGroup := new(database.RewardGroup)
	// curdb.Query(fmt.Sprintf(`delete from rewardgroup where rewardcontrol = '%s'`, xDoc["control"]))
	// for _, control := range httpReq.Form["rewardgroup"] {
	// 	xDocMerchantRewardGroup := make(map[string]interface{})
	// 	xDocMerchantRewardGroup["rewardcontrol"] = xDoc["control"]
	// 	xDocMerchantRewardGroup["workflow"] = "active"
	// 	xDocMerchantRewardGroup["groupcontrol"] = control
	// 	tblMerchantRewardGroup.Create(this.mapCache["username"].(string), xDocMerchantRewardGroup, curdb)
	// }
	//Handle Group Link - RewardGroup Saving disabled for merchants

	//Handle Store Link
	tblMerchantRewardStore := new(database.RewardStore)
	curdb.Query(fmt.Sprintf(`delete from rewardstore where rewardcontrol = '%s'`, xDoc["control"]))
	for _, control := range httpReq.Form["rewardstore"] {
		xDocMerchantRewardStore := make(map[string]interface{})
		xDocMerchantRewardStore["rewardcontrol"] = xDoc["control"]
		xDocMerchantRewardStore["workflow"] = "active"
		xDocMerchantRewardStore["storecontrol"] = control
		tblMerchantRewardStore.Create(this.mapCache["username"].(string), xDocMerchantRewardStore, curdb)
	}
	//Handle Store Link

	//SEND AN EMAIL USING TEMPLATE
	emailFields := make(map[string]interface{})

	sqlStoreOLD := fmt.Sprintf(`select workflow, title as reward, method, (select title from profile where control = merchantcontrol) 
		as merchant from reward where control = '%s'`, xDoc["control"])
	resStoreOLD, _ := curdb.Query(sqlStoreOLD)
	xDocOLD := make(map[string]interface{})
	if resStoreOLD["1"] != nil {
		xDocOLD = resStoreOLD["1"].(map[string]interface{})
		emailFields["reward"] = xDocOLD["reward"]
		emailFields["merchant"] = xDocOLD["merchant"]
		emailFields["method"] = xDocOLD["method"]
	}

	emailFields["action"] = sAction
	emailFields["username"] = this.mapCache["username"]

	emailTo := "rewards@valued.com"
	emailFrom := "notifications@valued.com"
	emailFromName := "VALUED ADMIN NOTIFICATIONS"
	emailSubject := fmt.Sprintf("Merchant %s Reward Change Request", emailFields["merchant"])
	if sAction != "" {
		emailSubject = fmt.Sprintf("Merchant %s Added New Reward", emailFields["merchant"])
	}
	emailTemplate := "merchantstore-save"
	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, "", emailFields)
	//SEND AN EMAIL USING TEMPLATE

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"` + sMessage + `","mainpanelContent":` + viewHTML + `}`))
}

func (this *MerchantReward) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblMerchantReward := new(database.Reward)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblMerchantReward.Read(xDocRequest, curdb)
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

	//Check The MerchantReward Groups and Link Accordingly
	tblMerchantRewardGroup := new(database.RewardGroup)
	xDocMerchantRewardGroupReq := make(map[string]interface{})
	linkedMerchantRewardGroup := make(map[string]interface{})

	xDocMerchantRewardGroupReq["reward"] = control
	linkListMap := tblMerchantRewardGroup.Search(xDocMerchantRewardGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedMerchantRewardGroup[xDoc["groupcontrol"].(string)] = xDoc["control"]
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedMerchantRewardGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#reward-view-group"] = xDoc
	}
	//Check The MerchantReward Groups and Link Accordingly

	formView := make(map[string]interface{})
	formView["merchantreward-view"] = xDocResult
	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *MerchantReward) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblMerchantReward := new(database.Reward)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblMerchantReward.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocResult := ResMap["1"].(map[string]interface{})

	// sMessage := ""
	// switch xDocResult["workflow"].(string) {
	// case "active":
	// 	sMessage = "Cannot Edit <b>Active</b> Rewards"
	// }

	// if sMessage != "" {
	// 	httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
	// 	return
	// }

	rewardType := xDocResult["type"].(string)
	xDocResult[rewardType] = "checked"

	rewardMethod := strings.Replace(xDocResult["method"].(string), " ", "", -1)
	xDocResult[rewardMethod] = "checked"

	rewardVisibleTo := xDocResult["visibleto"].(string)
	xDocResult[rewardVisibleTo] = "checked"

	rewardDiscontType := xDocResult["discounttype"].(string)
	xDocResult[rewardDiscontType] = "selected"

	xDocResult["formtitle"] = "Edit"

	sDiscount := xDocResult["discount"].(string)
	sDiscount = strings.Replace(strings.ToLower(sDiscount), "%", "", -1)
	sDiscount = strings.Replace(strings.ToLower(sDiscount), "Off", "", -1)
	sDiscount = strings.Replace(strings.ToLower(sDiscount), "off", "", -1)
	xDocResult["discount"] = functions.TrimEscape(sDiscount)

	sDescription := xDocResult["description"].(string)
	xDocResult["description"] = strings.Replace(sDescription, "<br>", "\n", -1)

	sRestriction := xDocResult["restriction"].(string)
	xDocResult["restriction"] = strings.Replace(sRestriction, "<br>", "\n", -1)

	//Search for Category & Sub Category Title
	new(Category).QuickSearchTitle(xDocResult, curdb)
	//Search for Category & Sub Category Title

	//Check The MerchantReward Groups and Link Accordingly
	tblMerchantRewardGroup := new(database.RewardGroup)
	xDocMerchantRewardGroupReq := make(map[string]interface{})
	linkedMerchantRewardGroup := make(map[string]interface{})

	xDocMerchantRewardGroupReq["reward"] = xDocResult["control"]
	linkListMap := tblMerchantRewardGroup.Search(xDocMerchantRewardGroupReq, curdb)
	for _, xDoc := range linkListMap {
		xDoc := xDoc.(map[string]interface{})
		linkedMerchantRewardGroup[xDoc["groupcontrol"].(string)] = xDoc["control"]
	}

	tblGroup := new(database.Groups)
	xDocGroupReq := make(map[string]interface{})
	xDocGroupReq["workflow"] = "active"
	listMap := tblGroup.Search(xDocGroupReq, curdb)

	for cNumber, xDoc := range listMap {
		xDoc := xDoc.(map[string]interface{})
		if linkedMerchantRewardGroup[xDoc["control"].(string)] != nil {
			xDoc["state"] = "checked"
		}
		xDocResult[cNumber+"#reward-edit-group"] = xDoc
	}
	//Check The MerchantReward Groups and Link Accordingly

	formView := make(map[string]interface{})
	formView["merchantreward-edit"] = xDocResult

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

/*

	func (this *MerchantReward) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a Reward"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update reward set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		httpRes.Write([]byte(`{"triggerSearch":true}`))
	}

	func (this *MerchantReward) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a Reward"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update reward set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
		httpRes.Write([]byte(`{"error":"Reward Activated","mainpanelContent":` + viewHTML + `}`))
	}

	func (this *MerchantReward) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a Reward"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update reward set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		httpRes.Write([]byte(`{"triggerSearch":true}`))
	}

	func (this *MerchantReward) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if functions.TrimEscape(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a Reward"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update reward set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
		httpRes.Write([]byte(`{"error":"MerchantReward Deactivated","mainpanelContent":` + viewHTML + `}`))
	}
*/

func (this *MerchantReward) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
	sqlDeactivate := fmt.Sprintf(`select reward.workflow as workflow, reward.title as reward, merchant.title as merchant from reward 
		left join profile as merchant on merchant.control = reward.merchantcontrol where reward.control in ('0'%s)`, controlList)
	resDeactivate, _ := curdb.Query(sqlDeactivate)

	sTableRows := ""
	sRow := `<tr> <td>%s</td> <td>%s</td> <td>%s</td> <td>%s</td></tr>`

	sMerchant := ""
	for _, xDoc := range resDeactivate {
		xDoc := xDoc.(map[string]interface{})
		sTableRows += fmt.Sprintf(sRow, "inactive", xDoc["workflow"], xDoc["reward"], xDoc["merchant"])
		sMerchant = xDoc["merchant"].(string)
	}

	emailFields["rows"] = sTableRows
	emailFields["username"] = this.mapCache["username"]
	emailFields["message"] = functions.TrimEscape(httpReq.FormValue("message"))

	emailTo := "rewards@valued.com"
	emailCC := this.mapCache["email"].(string)
	emailFrom := "notifications@valued.com"
	emailFromName := "VALUED ADMIN NOTIFICATIONS"
	emailSubject := "Merchant " + sMerchant + " Reward Deactivate Request"
	emailTemplate := "merchantreward-deactivate"
	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, emailCC, emailFields)
	//SEND AN EMAIL USING TEMPLATE

	httpRes.Write([]byte(`{"error":"Your request has been successfully submitted"}`))
}

func (this *MerchantReward) requestChanges(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		sMessage += "Please a store"
	}

	if functions.TrimEscape(httpReq.FormValue("message")) == "" {
		sMessage += "Please enter your changes"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//SEND AN EMAIL USING TEMPLATE
	emailFields := make(map[string]interface{})
	sqlReward := fmt.Sprintf(`select reward.workflow as workflow, reward.title as reward, merchant.title as merchant from reward 
		left join profile as merchant on merchant.control = reward.merchantcontrol where reward.control = '%s'`, functions.TrimEscape(httpReq.FormValue("control")))
	resReward, _ := curdb.Query(sqlReward)

	for _, xDoc := range resReward {
		xDoc := xDoc.(map[string]interface{})
		emailFields["reward"] = xDoc["reward"]
		emailFields["merchant"] = xDoc["merchant"]
	}

	emailFields["username"] = this.mapCache["username"]
	emailFields["message"] = functions.TrimEscape(httpReq.FormValue("message"))

	emailTo := "rewards@valued.com"
	emailCC := this.mapCache["email"].(string)
	emailFrom := "notifications@valued.com"
	emailFromName := "VALUED ADMIN NOTIFICATIONS"
	emailSubject := "Merchant " + emailFields["merchant"].(string) + " Reward Change Request"
	emailTemplate := "merchantreward-change"
	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, emailCC, emailFields)
	//SEND AN EMAIL USING TEMPLATE

	httpRes.Write([]byte(`{"error":"Your request has been successfully submitted"}`))
}
