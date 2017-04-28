package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
)

type AppFavorite struct {
	functions.Templates
	mapAppCache    map[string]interface{}
	mapSearchCache map[string]interface{}
	pageMap        map[string]interface{}
}

func (this *AppFavorite) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")
	this.mapSearchCache = curdb.GetSession(GOSESSID.Value, "mapSearchCache")

	if this.mapSearchCache == nil {
		this.mapSearchCache = make(map[string]interface{})
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})
		appFavorite := this.search(httpReq, curdb)

		appNavbar := new(AppNavbar)
		appFavorite["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
		appFavorite["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

		appFavorite["app-searchdiv"] = make(map[string]interface{})
		appFavorite["app-slidebar"] = make(map[string]interface{})

		appFooter := make(map[string]interface{})
		appFavorite["app-footer"] = appFooter
		appFooter["favourite"] = "white"

		this.pageMap["app-favorite"] = appFavorite
		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Valued Rewards","pageContent":` + contentHTML + `}`))
		return

	case "search":
		searchResult := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		httpRes.Write([]byte(`{"searchresult":` + contentHTML + `}`))
		return

	case "add":
		this.add(httpRes, httpReq, curdb)
		return

	case "remove":
		this.remove(httpRes, httpReq, curdb)
		return
	}
}

func (this *AppFavorite) List(mapAppCache map[string]interface{}, curdb database.Database) map[string]interface{} {
	favoriteList := make(map[string]interface{})

	if mapAppCache["control"] == nil {
		return favoriteList
	}

	sqlFavorite := `select control, rewardcontrol, merchantcontrol from favorite where profilecontrol = '%s' order by control desc`
	sqlFavorite = fmt.Sprintf(sqlFavorite, mapAppCache["control"])
	xDocresult, _ := curdb.Query(sqlFavorite)

	control := ""
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		if xDoc["merchantcontrol"] != nil || xDoc["merchantcontrol"].(string) != "" {
			control = xDoc["merchantcontrol"].(string)
			favoriteList["merchant-"+control] = xDoc
		}

		if xDoc["rewardcontrol"] != nil || xDoc["rewardcontrol"].(string) != "" {
			control = xDoc["rewardcontrol"].(string)
			favoriteList["reward-"+control] = xDoc
		}
	}

	return favoriteList
}

func (this *AppFavorite) add(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if (functions.TrimEscape(httpReq.FormValue("rewardcontrol")) == "" &&
		functions.TrimEscape(httpReq.FormValue("merchantcontrol")) == "") ||
		functions.TrimEscape(httpReq.FormValue("title")) == "" ||
		this.mapAppCache["control"] == nil {

		sMessage += fmt.Sprintf("Error Adding <b>%s</b> to Favourite!<br>", functions.TrimEscape(httpReq.FormValue("title")))
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblFavorite := new(database.Favorite)
	xDocFavorite := make(map[string]interface{})
	xDocFavorite["workflow"] = "active"
	if functions.TrimEscape(httpReq.FormValue("merchantcontrol")) != "" {
		xDocFavorite["profilecontrol"] = this.mapAppCache["control"]
		xDocFavorite["merchantcontrol"] = functions.TrimEscape(httpReq.FormValue("merchantcontrol"))
		tblFavorite.Create(this.mapAppCache["username"].(string), xDocFavorite, curdb)
	}

	if functions.TrimEscape(httpReq.FormValue("rewardcontrol")) != "" {
		xDocFavorite["profilecontrol"] = this.mapAppCache["control"]
		xDocFavorite["rewardcontrol"] = functions.TrimEscape(httpReq.FormValue("rewardcontrol"))
		tblFavorite.Create(this.mapAppCache["username"].(string), xDocFavorite, curdb)
	}

	sMessage += fmt.Sprintf("<b>%s</b> <br> Added to Favourite!<br>", functions.TrimEscape(httpReq.FormValue("title")))
	httpRes.Write([]byte(`{"warning":"` + sMessage + `"}`))
	return
}

func (this *AppFavorite) remove(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""

	if (functions.TrimEscape(httpReq.FormValue("rewardcontrol")) == "" &&
		functions.TrimEscape(httpReq.FormValue("merchantcontrol")) == "") ||
		functions.TrimEscape(httpReq.FormValue("title")) == "" ||
		this.mapAppCache["control"] == nil {

		sMessage += fmt.Sprintf("Error Removing <b>%s</b> from Favourite!<br>", functions.TrimEscape(httpReq.FormValue("title")))
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	if functions.TrimEscape(httpReq.FormValue("merchantcontrol")) != "" {
		sql := `delete from favorite where profilecontrol = '%s' and merchantcontrol = '%s' `
		sql = fmt.Sprintf(sql, this.mapAppCache["control"], functions.TrimEscape(httpReq.FormValue("merchantcontrol")))
		curdb.Query(sql)
	}

	if functions.TrimEscape(httpReq.FormValue("rewardcontrol")) != "" {
		sql := `delete from favorite where profilecontrol = '%s' and rewardcontrol = '%s' `
		sql = fmt.Sprintf(sql, this.mapAppCache["control"], functions.TrimEscape(httpReq.FormValue("rewardcontrol")))
		curdb.Query(sql)
	}

	sMessage += fmt.Sprintf("<b>%s</b> <br> Removed from Favourite!<br>", functions.TrimEscape(httpReq.FormValue("title")))
	httpRes.Write([]byte(`{"warning":"` + sMessage + `"}`))
	return
}

func (this *AppFavorite) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	sLimit := "300"
	if html.EscapeString(httpReq.FormValue("limit")) != "" {
		sLimit = html.EscapeString(httpReq.FormValue("limit"))
	}

	sOffset := "0"
	if html.EscapeString(httpReq.FormValue("offset")) != "" {
		sOffset = html.EscapeString(httpReq.FormValue("offset"))
	}

	sqlSearchMerchant := `select 
					m.control as control, m.title as title, m.image as image, 
					m.categorycontrol as categorycontrol,
					f.createdate as createdate

					from favorite as f, profile as m %s

					where f.profilecontrol = '%s' AND m.status = 'active' 
					AND m.control = f.merchantcontrol %s %s

					%s %s %s order by createdate desc limit %s offset %s`

	sqlSearchtextMerchant := fmt.Sprintf(` AND lower(m.title) like lower('%%%s%%') `, html.EscapeString(httpReq.FormValue("searchtext")))

	//

	sqlSearchReward := `select 
					r.title as title, r.control as control, 
					r.discount as discount, r.discounttype as discounttype, 
					r.categorycontrol as categorycontrol, 
					m.title as merchanttitle, m.image as merchantimage, 
					f.createdate as createdate

					from favorite as f, profile as m, reward as r %s

					where f.profilecontrol = '%s' AND f.rewardcontrol = r.control AND r.workflow = 'active' 
					AND m.control = r.merchantcontrol %s %s

					%s %s %s order by createdate desc limit %s offset %s`

	sqlSearchtextReward := fmt.Sprintf(` AND (lower(r.title) like lower('%%%s%%') OR lower(m.title) like lower('%%%s%%')) `,
		html.EscapeString(httpReq.FormValue("searchtext")), html.EscapeString(httpReq.FormValue("searchtext")))

	//

	var sqlRewardCategory, sqlRewardSubCategory, sqlRewardKeywordFilter, sqlRewardKeywordFilterTable, sqlRewardKeywordFilterCondition string
	var sqlMerchantCategory, sqlMerchantSubCategory, sqlMerchantKeywordFilter, sqlMerchantKeywordFilterTable, sqlMerchantKeywordFilterCondition string

	//

	if this.mapSearchCache["category"] != nil {
		sqlRewardCategory = fmt.Sprintf(` AND r.categorycontrol = '%s' `, this.mapSearchCache["category"])
		sqlMerchantCategory = fmt.Sprintf(` AND m.categorycontrol = '%s' `, this.mapSearchCache["category"])
	}

	if this.mapSearchCache["subcategory"] != nil {
		sqlRewardSubCategory = fmt.Sprintf(` AND r.subcategorycontrol = '%s' `, this.mapSearchCache["subcategory"])
		sqlMerchantSubCategory = fmt.Sprintf(` AND m.subcategorycontrol = '%s' `, this.mapSearchCache["subcategory"])
	}

	if this.mapSearchCache["keyword"] != nil {
		sqlRewardKeywordFilterTable = ", categorylink as cl"
		sqlRewardKeywordFilterCondition = " AND r.control == cl.rewardcontrol"
		sqlRewardKeywordFilter = fmt.Sprintf(` AND cl.categorycontrol = '%s' `, this.mapSearchCache["keyword"])

		sqlMerchantKeywordFilterTable = ", categorylink as cl"
		sqlMerchantKeywordFilterCondition = " AND m.control == cl.merchantcontrol"
		sqlMerchantKeywordFilter = fmt.Sprintf(` AND cl.categorycontrol = '%s' `, this.mapSearchCache["keyword"])
	}

	formSearch := make(map[string]interface{})
	categoryList := new(Category).ListAll(curdb)

	sqlSearchMerchant = fmt.Sprintf(sqlSearchMerchant, sqlMerchantKeywordFilterTable, this.mapAppCache["control"],
		sqlMerchantKeywordFilterCondition, sqlSearchtextMerchant, sqlMerchantCategory, sqlMerchantSubCategory, sqlMerchantKeywordFilter, sLimit, sOffset)
	xDocresult, _ := curdb.Query(sqlSearchMerchant)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["signedIn"] = "yes"
		xDoc["heart"] = "heart_filled.png"
		xDoc["favorite"] = `getForm('/app-favorite')`

		if xDoc["categorycontrol"] != nil {
			if categoryList[xDoc["categorycontrol"].(string)] != nil {
				xDocCategory := categoryList[xDoc["categorycontrol"].(string)].(map[string]interface{})
				xDoc["categorytitle"] = xDocCategory["title"]
			}
		}
		formSearch[cNumber+"#app-merchant-list"] = xDoc
	}

	//-->

	myAppRedeem := new(AppRedeem)
	myAppRedeem.SetAppCache(this.mapAppCache)

	sqlSearchReward = fmt.Sprintf(sqlSearchReward, sqlRewardKeywordFilterTable, this.mapAppCache["control"],
		sqlRewardKeywordFilterCondition, sqlSearchtextReward, sqlRewardCategory, sqlRewardSubCategory, sqlRewardKeywordFilter, sLimit, sOffset)
	xDocresult, _ = curdb.Query(sqlSearchReward)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["signedIn"] = "yes"
		xDoc["heart"] = "heart_filled.png"
		xDoc["favorite"] = `getForm('/app-favorite')`

		_, xDoc["state"] = myAppRedeem.ValidateEligibility(xDoc["control"].(string), curdb)

		if xDoc["categorycontrol"] != nil {
			if categoryList[xDoc["categorycontrol"].(string)] != nil {
				xDocCategory := categoryList[xDoc["categorycontrol"].(string)].(map[string]interface{})
				xDoc["categorytitle"] = xDocCategory["title"]
			}
		}
		formSearch[cNumber+"#app-reward-list"] = xDoc
	}

	return formSearch
}

func (this *AppFavorite) searchOLD(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	tblFavorite := new(database.Favorite)
	formSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["limit"] = html.EscapeString(httpReq.FormValue("limit"))
	xDocrequest["workflow"] = "active"

	xDocrequest["category"] = html.EscapeString(httpReq.FormValue("category"))
	xDocrequest["subcategory"] = html.EscapeString(httpReq.FormValue("subcategory"))
	// xDocrequest["subcategoryfilter"] = html.EscapeString(httpReq.FormValue("filter"))

	//
	//
	//

	//Search Favorite Merchant and Get Details from Profile
	xDocrequest["reftable"] = "merchant"
	xDocresult := tblFavorite.Search(xDocrequest, curdb)

	var merchantList []string
	merchantListOrder := make(map[string]interface{})
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		merchantListOrder[xDoc["merchantcontrol"].(string)] = cNumber
		merchantList = append(merchantList, xDoc["merchantcontrol"].(string))
	}

	if len(merchantList) > 0 {
		xDocrequest = make(map[string]interface{})
		xDocrequest["workflow"] = "active"
		xDocrequest["company"] = "Yes"
		xDocrequest["profilecontrol"] = merchantList
		xDocrequest["limit"] = "2000"

		xDocresult = new(database.Profile).Search(xDocrequest, curdb)
		for cNumber, xDoc := range xDocresult {
			xDoc := xDoc.(map[string]interface{})
			cNumber = merchantListOrder[xDoc["control"].(string)].(string)
			xDoc["number"] = cNumber
			xDoc["signedIn"] = "yes"
			xDoc["heart"] = "heart_filled.png"
			xDoc["favorite"] = `getForm('/app-favorite')`
			formSearch[cNumber+"#app-merchant-list"] = xDoc
		}
	}
	//Search Favorite Merchant and Get Details from Profile

	//
	//
	//

	//Search Favorite Reward and Get Details from Reward
	delete(xDocrequest, "merchant")
	xDocrequest["reftable"] = "reward"
	xDocresult = tblFavorite.Search(xDocrequest, curdb)

	var rewardList []string
	rewardListOrder := make(map[string]interface{})
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		rewardListOrder[xDoc["rewardcontrol"].(string)] = cNumber
		rewardList = append(rewardList, xDoc["rewardcontrol"].(string))
	}

	if len(rewardList) > 0 {
		xDocrequest = make(map[string]interface{})
		xDocrequest["workflow"] = "active"
		xDocrequest["rewardcontrol"] = rewardList
		xDocrequest["limit"] = "2000"

		xDocresult = new(database.Reward).Search(xDocrequest, curdb)
		for cNumber, xDoc := range xDocresult {
			xDoc := xDoc.(map[string]interface{})
			cNumber = rewardListOrder[xDoc["control"].(string)].(string)
			xDoc["number"] = cNumber
			xDoc["signedIn"] = "yes"
			xDoc["heart"] = "heart_filled.png"
			xDoc["favorite"] = `getForm('/app-favorite')`
			formSearch[cNumber+"#app-reward-list"] = xDoc
		}
	}
	//Search Favorite Reward and Get Details from Reward

	return formSearch
}
