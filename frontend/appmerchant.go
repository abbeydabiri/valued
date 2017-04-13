package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
)

type AppMerchant struct {
	functions.Templates
	mapSearchCache map[string]interface{}
	mapAppCache    map[string]interface{}
	pageMap        map[string]interface{}
}

func (this *AppMerchant) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		this.mapSearchCache = make(map[string]interface{})
		this.pageMap = make(map[string]interface{})
		appMerchant := this.search(httpReq, curdb)

		appNavbar := new(AppNavbar)
		appMerchant["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
		appMerchant["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

		appMerchant["app-searchdiv"] = make(map[string]interface{})
		appMerchant["app-slidebar"] = make(map[string]interface{})

		appFooter := make(map[string]interface{})
		appMerchant["app-footer"] = appFooter
		appFooter["brand"] = "white"

		this.pageMap["app-merchant"] = appMerchant
		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Valued Rewards","pageContent":` + contentHTML + `}`))
		return

	case "search":
		searchResult := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		httpRes.Write([]byte(`{"searchresult":` + contentHTML + `}`))
		return
	}
}

func (this *AppMerchant) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	sLimit := "10"
	if html.EscapeString(httpReq.FormValue("limit")) != "" {
		sLimit = html.EscapeString(httpReq.FormValue("limit"))
	}

	sOffset := "0"
	if html.EscapeString(httpReq.FormValue("offset")) != "" {
		sOffset = html.EscapeString(httpReq.FormValue("offset"))
	}

	sqlSearch := `select m.code as merchantcode,
					m.control as control, m.title as title, m.image as image,
					m.subcategorycontrol as subcategorycontrol,
					m.categorycontrol as categorycontrol, 
					m.createdate as createdate


					from profile as m, category as c %s
					where m.status = 'active' AND lower(m.company) = 'yes'
					AND m.code != 'main' 
					AND m.categorycontrol = c.control 

					%s %s

					%s %s %s AND m.control in (select distinct merchantcontrol from reward where control in (select rewardcontrol from rewardscheme where schemecontrol in (select control from scheme where code in ('lite','lifestyle')   ))) order by m.title limit %s offset %s`

	sqlSearchtext := fmt.Sprintf(`AND (lower(c.title) like lower('%%%s%%') OR lower(m.title) like lower('%%%s%%') OR lower(m.description) like lower('%%%s%%'))`, html.EscapeString(httpReq.FormValue("searchtext")), html.EscapeString(httpReq.FormValue("searchtext")), html.EscapeString(httpReq.FormValue("searchtext")))

	sqlCategory := ""
	sqlSubCategory := ""
	sqlKeywordFilter := ""
	sqlKeywordFilterTable := ""
	sqlKeywordFilterCondition := ""

	if this.mapSearchCache["category"] != nil {
		sqlCategory = fmt.Sprintf(` AND m.categorycontrol = '%s' `, this.mapSearchCache["category"])
	}

	if this.mapSearchCache["subcategory"] != nil {
		sqlSubCategory = fmt.Sprintf(` AND m.subcategorycontrol = '%s' `, this.mapSearchCache["subcategory"])
	}

	if this.mapSearchCache["keyword"] != nil {
		sqlKeywordFilterTable = ", categorylink as cl"
		sqlKeywordFilterCondition = " AND cl.merchantcontrol = m.control "
		sqlKeywordFilter = fmt.Sprintf(` AND cl.categorycontrol = '%s' `, this.mapSearchCache["keyword"])
	}

	sqlSearch = fmt.Sprintf(sqlSearch, sqlKeywordFilterTable, sqlKeywordFilterCondition, sqlSearchtext,
		sqlCategory, sqlSubCategory, sqlKeywordFilter, sLimit, sOffset)
	xDocresult, _ := curdb.Query(sqlSearch)

	favoriteItem := ""
	categoryList := new(Category).ListAll(curdb)
	favoriteList := new(AppFavorite).List(this.mapAppCache, curdb)

	formSearch := make(map[string]interface{})
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		if xDoc["merchantcode"] != nil && xDoc["merchantcode"].(string) == "none" {
			continue
		}

		if this.mapAppCache["control"] != nil {
			xDoc["signedIn"] = "yes"
		}

		favoriteItem = fmt.Sprintf("merchant-%s", xDoc["control"])
		if favoriteList[favoriteItem] == nil {
			xDoc["heart"] = "heart.png"
		} else {
			xDoc["heart"] = "heart_filled.png"
		}

		if xDoc["categorycontrol"] != nil {
			if categoryList[xDoc["categorycontrol"].(string)] != nil {
				xDocCategory := categoryList[xDoc["categorycontrol"].(string)].(map[string]interface{})
				xDoc["categorytitle"] = xDocCategory["title"]
			}
		}

		if xDoc["subcategorycontrol"] != nil {
			if categoryList[xDoc["subcategorycontrol"].(string)] != nil {
				xDocCategory := categoryList[xDoc["subcategorycontrol"].(string)].(map[string]interface{})
				xDoc["subcategorytitle"] = xDocCategory["title"]
			}
		}

		formSearch[cNumber+"#app-merchant-list"] = xDoc
	}

	return formSearch
}
