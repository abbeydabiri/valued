package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
)

type AppReward struct {
	functions.Templates
	mapSearchCache map[string]interface{}
	mapAppCache    map[string]interface{}
	pageMap        map[string]interface{}
}

func (this *AppReward) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		appReward := this.search(httpReq, curdb)

		appNavbar := new(AppNavbar)
		appReward["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
		appReward["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

		appReward["app-searchdiv"] = make(map[string]interface{})

		appReward["app-slidebar"] = this.GetSliderBar()

		appFooter := make(map[string]interface{})
		appReward["app-footer"] = appFooter
		appFooter["reward"] = "white"

		this.pageMap["app-reward"] = appReward
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

func (this *AppReward) GetSliderBar() (mapSlider map[string]interface{}) {
	mapSlider = make(map[string]interface{})
	mapSliderButton := make(map[string]interface{})

	if this.mapAppCache["control"] != nil {
		mapSliderButton["title"] = "PERKS"
		mapSlider["1#app-sliderbar-button"] = mapSliderButton

		mapSliderButton = make(map[string]interface{})

		mapSliderButton["title"] = "PRIVILEGES"
		mapSlider["2#app-sliderbar-button"] = mapSliderButton
	} else {
		mapSliderButton["title"] = "LITE"
		mapSlider["1#app-sliderbar-button"] = mapSliderButton

		mapSliderButton = make(map[string]interface{})

		mapSliderButton["title"] = "LIFESTYLE"
		mapSlider["2#app-sliderbar-button"] = mapSliderButton
	}

	return
}

func (this *AppReward) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	//Update SearchTags
	if html.EscapeString(httpReq.FormValue("category")) != "" {
		this.mapSearchCache["category"] = httpReq.FormValue("category")
		sql := fmt.Sprintf(`select title from category where control = '%s'`, this.mapSearchCache["category"])

		defaultMap, _ := curdb.Query(sql)
		if defaultMap["1"] != nil {
			this.mapSearchCache["categorytitle"] = defaultMap["1"].(map[string]interface{})["title"]
		}
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)
	//Update SearchTags

	sLimit := "30"
	if html.EscapeString(httpReq.FormValue("limit")) != "" {
		sLimit = html.EscapeString(httpReq.FormValue("limit"))
	}

	sOffset := "0"
	if html.EscapeString(httpReq.FormValue("offset")) != "" {
		sOffset = html.EscapeString(httpReq.FormValue("offset"))
	}

	sqlSearch := `select 
					r.title as title, r.control as control, r.orderby as orderby,
					r.discount as discount, r.discounttype as discounttype, 
					r.categorycontrol as categorycontrol, 
					r.subcategorycontrol as subcategorycontrol,
					m.image as merchantimage, m.title as merchanttitle,  m.code as merchantcode,
					r.createdate as createdate

					from reward as r, profile as m %s %s
					where r.workflow = 'active' AND r.merchantcontrol = m.control AND m.status = 'active' AND m.code != 'main'
					%s %s

					%s %s %s %s %s AND r.control in (select rewardcontrol from rewardscheme where schemecontrol in (select control from scheme where code in ('lite','lifestyle') )) order by orderby, control desc limit %s offset %s`

	sqlSearchtextSub := ""
	sqlSearchtext := ` AND r.control not in ( select distinct(rewardgroup.rewardcontrol) as rewardcontrol from rewardgroup 
							left join groups on rewardgroup.groupcontrol = groups.control where groups.workflow = 'active' %s ) `

	if this.mapAppCache["control"] != nil {
		sqlSearchtextSub = fmt.Sprintf(` and groups.control not in ( select groupcontrol from membergroup where membergroup.membercontrol = '%s')`, this.mapAppCache["control"])
	}
	sqlSearchtext = fmt.Sprintf(sqlSearchtext, sqlSearchtextSub)

	sqlSearchtext += fmt.Sprintf(` AND (lower(r.keywords) like lower('%%%s%%') OR lower(r.title) like lower('%%%s%%') OR lower(m.title) like lower('%%%s%%')) `,
		html.EscapeString(httpReq.FormValue("searchtext")), html.EscapeString(httpReq.FormValue("searchtext")),
		html.EscapeString(httpReq.FormValue("searchtext")))

	sqlScheme := ""
	sqlSchemeTable := ""

	sqlType := ""
	sqlMerchant := ""
	sqlCategory := ""
	sqlSubCategory := ""

	sqlExtraFilter := ""
	sqlExtraFilterTable := ""

	if html.EscapeString(httpReq.FormValue("merchant")) != "" {
		sqlMerchant = fmt.Sprintf(` AND r.merchantcontrol = '%s' `, html.EscapeString(httpReq.FormValue("merchant")))
	} else if html.EscapeString(httpReq.FormValue("scheme")) != "" {

		sqlSchemeTable = ", rewardscheme as rs"
		sqlScheme = fmt.Sprintf(` AND r.control = rs.rewardcontrol AND rs.workflow = 'active' AND rs.schemecontrol = '%s'  `, html.EscapeString(httpReq.FormValue("scheme")))

	} else {

		if this.mapSearchCache["type"] != nil {
			switch this.mapSearchCache["type"].(string) {
			case "Perk", "Privilege":
				sqlType = fmt.Sprintf(` AND r.type = '%s' `, this.mapSearchCache["type"])
			case "Lite", "Lifestyle":
				sqlSchemeControl := fmt.Sprintf("select control from scheme where lower(title) like lower('%%%s%%') limit 1", this.mapSearchCache["type"])
				defaultMap, _ := curdb.Query(sqlSchemeControl)
				if defaultMap["1"] != nil {
					sqlSchemeTable = ", rewardscheme as rs"
					sqlScheme = fmt.Sprintf(` AND r.control = rs.rewardcontrol AND rs.workflow = 'active' AND rs.schemecontrol = '%s'  `, defaultMap["1"].(map[string]interface{})["control"])
				}

			}
		}

		if this.mapSearchCache["category"] != nil {
			sqlCategory = fmt.Sprintf(` AND r.categorycontrol = '%s' `, this.mapSearchCache["category"])
		}

		if this.mapSearchCache["subcategory"] != nil {
			sqlSubCategory = fmt.Sprintf(` AND r.subcategorycontrol = '%s' `, this.mapSearchCache["subcategory"])
		}

		if this.mapSearchCache["keyword"] != nil {
			sqlExtraFilterTable = ", categorylink as cl"
			//sqlExtraFilterCondition = " AND cl.rewardcontrol == r.control "
			sqlExtraFilter = fmt.Sprintf(` AND cl.rewardcontrol = r.control AND cl.categorycontrol = '%s' `, this.mapSearchCache["keyword"])
		}
	}

	//Search Rewards

	sGroup := "33"
	formSearch := make(map[string]interface{})
	sqlSearch = fmt.Sprintf(sqlSearch, sqlSchemeTable, sqlExtraFilterTable, sqlSearchtext,
		sqlScheme, sqlMerchant, sqlType, sqlCategory, sqlSubCategory, sqlExtraFilter, sLimit, sOffset)

	favoriteItem := ""
	categoryList := new(Category).ListAll(curdb)
	favoriteList := new(AppFavorite).List(this.mapAppCache, curdb)

	myAppRedeem := new(AppRedeem)
	myAppRedeem.SetAppCache(this.mapAppCache)

	if this.mapAppCache["control"] != nil {
		sqlSearchGroup := strings.Replace(sqlSearch, "AND r.control not in ( select distinct(", "AND r.control in ( select distinct(", 1)
		sqlSearchGroup = strings.Replace(sqlSearchGroup, "and groups.control not in ( select groupcontrol", "and groups.control in ( select groupcontrol", 1)
		sqlSearchGroup = strings.Replace(sqlSearchGroup, "AND r.control in (select rewardcontrol from rewardscheme", "AND r.control not in (select rewardcontrol from rewardscheme", 1)

		xDocresultGroup, _ := curdb.Query(sqlSearchGroup)
		for cNumber, xDoc := range xDocresultGroup {
			xDoc := xDoc.(map[string]interface{})
			xDoc["number"] = cNumber

			if xDoc["merchantcode"] != nil && xDoc["merchantcode"].(string) == "none" {
				continue
			}

			if this.mapAppCache["control"] != nil {
				xDoc["signedIn"] = "yes"
				_, xDoc["state"] = myAppRedeem.ValidateEligibility(xDoc["control"].(string), curdb)
			}

			favoriteItem = fmt.Sprintf("reward-%s", xDoc["control"])
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

			formSearch[cNumber+"#app-reward-list"] = xDoc
		}

	}

	//Logic to Manage Grouping
	xDocresult, _ := curdb.Query(sqlSearch)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		if xDoc["merchantcode"] != nil && xDoc["merchantcode"].(string) == "none" {
			continue
		}

		if this.mapAppCache["control"] != nil {
			xDoc["signedIn"] = "yes"
			_, xDoc["state"] = myAppRedeem.ValidateEligibility(xDoc["control"].(string), curdb)
		}

		favoriteItem = fmt.Sprintf("reward-%s", xDoc["control"])
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

		formSearch[sGroup+cNumber+"#app-reward-list"] = xDoc
	}

	return formSearch
}
