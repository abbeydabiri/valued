package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	// "html"
	"net/http"
	"strconv"
)

type AppSearch struct {
	functions.Templates
	mapSearchCache map[string]interface{}
	pageMap        map[string]interface{}
}

func (this *AppSearch) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapSearchCache = curdb.GetSession(GOSESSID.Value, "mapSearchCache")

	if this.mapSearchCache == nil {
		this.mapSearchCache = make(map[string]interface{})
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})
		this.pageMap["app-searchdiv"] = this.search("subcategory", httpReq.FormValue("category"), curdb)

		// searchSubCategory
		searchSubCategory := ``
		// if this.mapSearchCache["category"] != nil && this.mapSearchCache["subcategory"] != nil {
		if this.mapSearchCache["category"] != nil {
			searchSubCategoryMap := this.search("keyword", this.mapSearchCache["category"].(string), curdb)
			searchSubCategory = `,"appSearchSubCategory":` + strconv.Quote(string(this.Generate(searchSubCategoryMap, nil)))
		}
		// searchSubCategory

		// searchKeyword
		searchKeyword := ``
		if this.mapSearchCache["subcategory"] != nil && this.mapSearchCache["keyword"] != nil {
			searchKeywordMap := this.search("end", this.mapSearchCache["subcategory"].(string), curdb)
			searchKeyword = `,"appSearchKeyword":` + strconv.Quote(string(this.Generate(searchKeywordMap, nil)))
		}
		// searchKeyword

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchDiv":` + contentHTML + `,"appSearchTags":` + searchTagsHTML + searchSubCategory + searchKeyword + `}`))
		return

	case "slider":
		switch functions.TrimEscape(httpReq.FormValue("type")) {
		default:
			delete(this.mapSearchCache, "type")
		case "PERKS":
			this.mapSearchCache["type"] = "Perk"
		case "PRIVILEGES":
			this.mapSearchCache["type"] = "Privilege"
		case "LIFESTYLE":
			this.mapSearchCache["type"] = "Lifestyle"
		case "LITE":
			this.mapSearchCache["type"] = "Lite"
		}
		curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)
		httpRes.Write([]byte(`{"triggerAppSearch":"triggerAppSearch"}`))
		return

	case "category":

		searchResult := this.search("subcategory", httpReq.FormValue("category"), curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))

		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchCategory":` + contentHTML + `,"appSearchSubCategory":"","appSearchKeyword":""
			,"ppSearchTags":` + searchTagsHTML + `,"triggerAppSearch":"triggerAppSearch"}`))
		return

	case "subcategory":

		//Populate SearchTags
		delete(this.mapSearchCache, "keyword")
		delete(this.mapSearchCache, "keywordtitle")
		delete(this.mapSearchCache, "subcategory")
		delete(this.mapSearchCache, "subcategorytitle")
		delete(this.mapSearchCache, "category")
		delete(this.mapSearchCache, "categorytitle")

		this.mapSearchCache["category"] = httpReq.FormValue("category")
		sql := fmt.Sprintf(`select title from category where control = '%s'`, this.mapSearchCache["category"])

		defaultMap, _ := curdb.Query(sql)
		if defaultMap["1"] != nil {
			this.mapSearchCache["categorytitle"] = defaultMap["1"].(map[string]interface{})["title"]
		}

		curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)
		//Populate SearchTags

		searchResult := this.search("keyword", httpReq.FormValue("category"), curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))

		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchSubCategory":` + contentHTML + `,"appSearchKeyword":"","appSearchTags":` + searchTagsHTML + `,"triggerAppSearch":"triggerAppSearch"}`))
		//httpRes.Write([]byte(`{"appSearchSubCategory":` + contentHTML + `,"appSearchKeyword":"","appSearchTags":` + searchTagsHTML + `}`))

	case "keyword":

		//Populate SearchTags
		delete(this.mapSearchCache, "keyword")
		delete(this.mapSearchCache, "keywordtitle")
		delete(this.mapSearchCache, "subcategory")
		delete(this.mapSearchCache, "subcategorytitle")

		this.mapSearchCache["subcategory"] = httpReq.FormValue("category")
		sql := fmt.Sprintf(`select title from category where control = '%s'`, this.mapSearchCache["subcategory"])

		defaultMap, _ := curdb.Query(sql)
		if defaultMap["1"] != nil {
			this.mapSearchCache["subcategorytitle"] = defaultMap["1"].(map[string]interface{})["title"]
		}

		curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)
		//Populate SearchTags

		searchResult := this.search("end", httpReq.FormValue("category"), curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))

		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchKeyword":` + contentHTML + `,"appSearchTags":` + searchTagsHTML + `,"triggerAppSearch":"triggerAppSearch"}`))
		//httpRes.Write([]byte(`{"appSearchKeyword":` + contentHTML + `,"appSearchTags":` + searchTagsHTML + `}`))
		return

	case "end":

		//Populate SearchTags
		delete(this.mapSearchCache, "keyword")
		delete(this.mapSearchCache, "keywordtitle")

		this.mapSearchCache["keyword"] = httpReq.FormValue("category")
		sql := fmt.Sprintf(`select title from category where control = '%s'`, this.mapSearchCache["keyword"])

		defaultMap, _ := curdb.Query(sql)
		if defaultMap["1"] != nil {
			this.mapSearchCache["keywordtitle"] = defaultMap["1"].(map[string]interface{})["title"]
		}

		curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)
		//Populate SearchTags

		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `,"triggerAppSearch":"triggerAppSearch"}`))
		//httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `}`))
		return

	case "clearTagCategory":
		delete(this.mapSearchCache, "keyword")
		delete(this.mapSearchCache, "keywordtitle")
		delete(this.mapSearchCache, "subcategory")
		delete(this.mapSearchCache, "subcategorytitle")
		delete(this.mapSearchCache, "category")
		delete(this.mapSearchCache, "categorytitle")
		curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)

		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `,"appSearchSubCategory":"","appSearchKeyword":"", "appSearchClearTag":"appSearchCategory","triggerAppSearch":"triggerAppSearch"}`))
		//httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `,"appSearchSubCategory":"","appSearchKeyword":"", "appSearchClearTag":"appSearchCategory"}`))
		return

	case "clearTagSubCategory":
		delete(this.mapSearchCache, "keyword")
		delete(this.mapSearchCache, "keywordtitle")
		delete(this.mapSearchCache, "subcategory")
		delete(this.mapSearchCache, "subcategorytitle")
		curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)

		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `,"appSearchKeyword":"", "appSearchClearTag":"appSearchSubCategory","triggerAppSearch":"triggerAppSearch"}`))
		//httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `,"appSearchKeyword":"", "appSearchClearTag":"appSearchSubCategory"}`))
		return

	case "clearTagKeyword":
		delete(this.mapSearchCache, "keyword")
		delete(this.mapSearchCache, "keywordtitle")
		curdb.SetSession(GOSESSID.Value, "mapSearchCache", this.mapSearchCache, false)

		searchTagsHTML := strconv.Quote(string(this.Generate(this.searchTags(), nil)))
		httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `, "appSearchClearTag":"appSearchKeyword","triggerAppSearch":"triggerAppSearch"}`))
		//httpRes.Write([]byte(`{"appSearchTags":` + searchTagsHTML + `, "appSearchClearTag":"appSearchKeyword"}`))
		return
	}
}

func (this *AppSearch) searchTags() map[string]interface{} {
	searchTags := make(map[string]interface{})

	if this.mapSearchCache["category"] != nil {
		tag := make(map[string]interface{})
		tag["action"] = "clearTagCategory"
		tag["title"] = this.mapSearchCache["categorytitle"]
		searchTags["1#app-searchdiv-tags"] = tag
	}

	if this.mapSearchCache["subcategory"] != nil {
		tag := make(map[string]interface{})
		tag["action"] = "clearTagSubCategory"
		tag["title"] = this.mapSearchCache["subcategorytitle"]
		searchTags["2#app-searchdiv-tags"] = tag
	}

	if this.mapSearchCache["keyword"] != nil {
		tag := make(map[string]interface{})
		tag["action"] = "clearTagKeyword"
		tag["title"] = this.mapSearchCache["keywordtitle"]
		searchTags["3#app-searchdiv-tags"] = tag
	}

	return searchTags
}

func (this *AppSearch) search(nextaction string, category string, curdb database.Database) map[string]interface{} {

	tblAppCategory := new(database.Category)

	formSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})
	xDocrequest["workflow"] = "active"

	xDocrequest["category"] = category //httpReq.FormValue("category")
	if nextaction == "subcategory" {
		xDocrequest["sub"] = "false"
	}

	xDocresult := tblAppCategory.Search(xDocrequest, curdb)

	prevaction := ""
	switch nextaction {
	case "keyword":
		prevaction = "subcategory"
	case "end": //i.e keyword
		prevaction = "keyword"
	}

	if prevaction != "" {
		xDoc := make(map[string]interface{})
		xDoc["title"] = "SUB-CATEGORY"
		if prevaction == "keyword" {
			xDoc["title"] = "ADDITIONAL FILTER"
		}

		cNumber := "0"
		xDoc["number"] = cNumber
		xDoc["action"] = prevaction
		xDoc["isselected"] = "category"
		xDoc["control"] = category //httpReq.FormValue("category")
		formSearch[cNumber+"#app-searchdiv-result"] = xDoc
	}

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["action"] = nextaction

		if this.mapSearchCache["category"] != nil {
			if xDoc["control"].(string) == this.mapSearchCache["category"].(string) {
				xDoc["isselected"] = "active"
			}
		}

		if this.mapSearchCache["subcategory"] != nil {
			if xDoc["control"].(string) == this.mapSearchCache["subcategory"].(string) {
				xDoc["isselected"] = "active"
			}
		}

		if this.mapSearchCache["keyword"] != nil {
			if xDoc["control"].(string) == this.mapSearchCache["keyword"].(string) {
				xDoc["isselected"] = "active"
			}
		}

		formSearch[cNumber+"#app-searchdiv-result"] = xDoc
	}

	return formSearch
}
