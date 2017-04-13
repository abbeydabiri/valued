package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Category struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Category) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if !strings.Contains(httpReq.FormValue("action"), "download") {
		if httpReq.Method != "POST" {
			http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
		}
	} else {
		println("Action:" + httpReq.FormValue("action"))
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})
		this.pageMap["category-search"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Categories","mainpanelContentSearch":` + contentHTML + `}`))
		return

	case "fetchCategory":
		this.fetchMerchantCategory(httpRes, httpReq, curdb)
		// viewKeywordsHTML := this.fetchKeywords(httpRes, httpReq, curdb)
		// viewMerchantHTML := this.fetchMerchantCategory(httpRes, httpReq, curdb)
		// httpRes.Write([]byte(`{"reward-edit-merchant":` + viewMerchantHTML + `"categorylink":` + viewKeywordsHTML + `}`))
		return

	case "fetchKeywords":
		this.fetchKeywords(httpRes, httpReq, curdb)
		return

	case "quicksearch":
		this.quicksearch(httpRes, httpReq, curdb)
		return

	case "search":
		searchResult := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		httpRes.Write([]byte(`{"searchresult":` + contentHTML + `}`))
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

	case "importcsv":
		this.importcsv(httpRes, httpReq, curdb)
		return

	case "importcsvdownload":
		this.importcsvdownload(httpRes, httpReq, curdb)
		return

	case "importcsvsave":
		this.importcsvsave(httpRes, httpReq, curdb)
		return

	}
}

func (this *Category) fetchMerchantCategory(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sqlMerchant := fmt.Sprintf(`select control, categorycontrol, subcategorycontrol from profile where company = 'Yes' AND control = '%s'`, httpReq.FormValue("merchant"))
	xDocResult, _ := curdb.Query(sqlMerchant)
	xDoc := xDocResult["1"].(map[string]interface{})

	categoryList := this.ListAll(curdb)
	if xDoc["subcategorycontrol"] != nil && xDoc["subcategorycontrol"].(string) != "" {
		if categoryList[xDoc["subcategorycontrol"].(string)] != nil {
			xDocCategory := categoryList[xDoc["categorycontrol"].(string)].(map[string]interface{})
			xDoc["categorytitle"] = xDocCategory["title"]
		}
	}

	if xDoc["subcategorycontrol"] != nil && xDoc["subcategorycontrol"].(string) != "" {
		if categoryList[xDoc["subcategorycontrol"].(string)] != nil {
			xDocCategory := categoryList[xDoc["subcategorycontrol"].(string)].(map[string]interface{})
			xDoc["subcategorytitle"] = xDocCategory["title"]
		}
	}

	mapHTML := make(map[string]interface{})
	mapHTML["reward-edit-category"] = xDoc
	viewHTML := strconv.Quote(string(this.Generate(mapHTML, nil)))
	httpRes.Write([]byte(`{"reward-edit-category":` + viewHTML + `}`))
}

func (this *Category) fetchKeywords(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("category")) == "" {
		sMessage += "Category is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	lReadonly := false
	if functions.TrimEscape(httpReq.FormValue("readonly")) != "" {
		lReadonly = true
	}

	merchantKeywords := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("merchant")) != "" {
		sqlMerchant := fmt.Sprintf(`select categorycontrol from categorylink where merchantcontrol = '%s'`, httpReq.FormValue("merchant"))

		xDocResult, _ := curdb.Query(sqlMerchant)
		for _, xDoc := range xDocResult {
			xDoc := xDoc.(map[string]interface{})
			merchantKeywords[xDoc["categorycontrol"].(string)] = xDoc
		}
	}

	rewardKeywords := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("reward")) != "" {
		sqlReward := fmt.Sprintf(`select categorycontrol from categorylink where rewardcontrol = '%s'`, httpReq.FormValue("reward"))

		xDocResult, _ := curdb.Query(sqlReward)
		for _, xDoc := range xDocResult {
			xDoc := xDoc.(map[string]interface{})
			rewardKeywords[xDoc["categorycontrol"].(string)] = xDoc
		}
	}

	//Check The Review Categories and Link Accordingly
	tblCategory := new(database.Category)
	xDocrequest := make(map[string]interface{})
	xDocrequest["workflow"] = "active"
	xDocrequest["category"] = functions.TrimEscape(httpReq.FormValue("category"))

	checkBoxView := make(map[string]interface{})
	xDocresult := tblCategory.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})

		sState := ""
		if merchantKeywords[xDoc["control"].(string)] != nil {
			sState += " checked "
		}

		if rewardKeywords[xDoc["control"].(string)] != nil {
			sState += " checked "
		}

		if !lReadonly {
			checkBoxView[cNumber+"#categorylink-checkbox"] = xDoc
		} else {
			if len(sState) > 0 {
				checkBoxView[cNumber+"#categorylink-checkbox"] = xDoc
			}
			sState += " disabled "
		}
		xDoc["state"] = sState

	}
	//Check The Review Categories and Link Accordingly

	viewHTML := strconv.Quote(string(this.Generate(checkBoxView, nil)))
	httpRes.Write([]byte(`{"categorylink":` + viewHTML + `}`))
}

func (this *Category) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblCategory := new(database.Category)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"

	if functions.TrimEscape(httpReq.FormValue("sub")) != "" {
		xDocrequest["sub"] = functions.TrimEscape(httpReq.FormValue("sub"))
	}

	if functions.TrimEscape(httpReq.FormValue("category")) != "" {
		xDocrequest["category"] = functions.TrimEscape(httpReq.FormValue("category"))
	}

	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))

	xDocresult := tblCategory.Search(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})

		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["title"] = "No Records Found"
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Category) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	tblCategory := new(database.Category)
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = functions.TrimEscape(httpReq.FormValue("limit"))
	xDocrequest["offset"] = functions.TrimEscape(httpReq.FormValue("offset"))
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))

	xDocresult := tblCategory.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#category-search-result"] = xDoc

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

	return formSearch
}

func (this *Category) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["category-edit"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Category) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("title")) == "" {
		sMessage += " Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("workflow")) == "" {
		sMessage += "Status is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblCategory := new(database.Category)
	xDoc := make(map[string]interface{})
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))
	xDoc["placement"] = functions.TrimEscape(httpReq.FormValue("placement"))
	xDoc["description"] = functions.TrimEscape(httpReq.FormValue("description"))
	xDoc["categorycontrol"] = functions.TrimEscape(httpReq.FormValue("category"))

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("category_%s_%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblCategory.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblCategory.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"Category <b>` + xDoc["title"].(string) + `</b> Saved","mainpanelContent":` + viewHTML + `}`))
}

func (this *Category) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblCategory := new(database.Category)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblCategory.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	switch xDocRequest["workflow"].(string) {
	case "inactive":
		xDocRequest["actionView"] = "activateView"
		xDocRequest["actionColor"] = "success"
		xDocRequest["actionLabel"] = "Activate"

	case "active":
		xDocRequest["actionView"] = "deactivateView"
		xDocRequest["actionColor"] = "danger"
		xDocRequest["actionLabel"] = "De-Activate"
	}

	xDocRequest["createdate"] = xDocRequest["createdate"].(string)[0:19]

	formView := make(map[string]interface{})
	formView["category-view"] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Category) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblCategory := new(database.Category)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblCategory.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	formView := make(map[string]interface{})

	xDocRequest["formtitle"] = "Edit"
	formView["category-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Category) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update category set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Category) activateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update category set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Category Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Category) deactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update category set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Category) deactivateView(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update category set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
	httpRes.Write([]byte(`{"error":"Category De-Activated","mainpanelContent":` + viewHTML + `}`))
}

func (this *Category) deactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a category"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update category set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"triggerSearch":true}`))
}

func (this *Category) importcsv(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	formImport := make(map[string]interface{})
	formImport["category-importcsv"] = make(map[string]interface{})

	importcsvHtml := strconv.Quote(string(this.Generate(formImport, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + importcsvHtml + `}`))
}

func (this *Category) importcsvdownload(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	httpRes.Header().Set("Content-Type", "text/csv")
	httpRes.Header().Set("Content-Disposition", "attachment;filename=Valued_ImportCategoryTemplate.csv")
	httpRes.Write([]byte(strings.Join([]string{"code", "title", "image", "placement", "category"}, ",")))
}

func (this *Category) importcsvsave(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	mapFiletypes := make(map[string]string)
	mapFiletypes["csvfile"] = "text/plain"
	mapFiles, sMessage := functions.UploadFile([]string{"csvfile"}, mapFiletypes, httpReq)

	if mapFiles["csvfile"] != nil {
		maUploaded := mapFiles["csvfile"].(map[string]interface{})
		categoryList := string(maUploaded["filebytes"].([]byte))
		categoryList = strings.Replace(categoryList, "\r", "", -1)

		sliceRow := strings.Split(categoryList, "\n")
		tblCategory := new(database.Category)

		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// List Category by Code
		mapCategory := make(map[string]string)
		if len(sliceRow) > 1 {
			xDocResCategory, _ := curdb.Query("select code, control from category")
			if xDocResCategory != nil {
				for _, xDoc := range xDocResCategory {
					xDoc := xDoc.(map[string]interface{})
					mapCategory[xDoc["code"].(string)] = xDoc["control"].(string)
				}
			}
		}
		// List Category by Code
		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		workflow := functions.TrimEscape(httpReq.FormValue("status"))

		go func() {
			for index, stringCols := range sliceRow {

				stringCols = strings.TrimSpace(stringCols)
				if index == 0 || len(stringCols) == 0 {
					continue
				}
				sliceCols := strings.Split(stringCols, ",")

				xDoc := make(map[string]interface{})
				xDoc["code"] = strings.TrimSpace(sliceCols[0])

				if len(sliceCols) > 1 {
					xDoc["categorycontrol"] = ""
					xDoc["title"] = strings.TrimSpace(sliceCols[1])
				}
				if len(sliceCols) > 2 {
					xDoc["image"] = strings.TrimSpace(sliceCols[2])
				}
				if len(sliceCols) > 3 {
					xDoc["placement"] = strings.TrimSpace(sliceCols[3])
				}
				if len(sliceCols) > 4 {
					if strings.TrimSpace(sliceCols[4]) != "" {
						xDoc["categorycontrol"] = mapCategory[strings.TrimSpace(sliceCols[4])]
					}
				}
				xDoc["workflow"] = workflow

				sCode := strings.TrimSpace(sliceCols[0])
				sControl := tblCategory.Create(this.mapCache["username"].(string), xDoc, curdb)

				if sCode != "" {
					mapCategory[sCode] = sControl
				}

				<-time.Tick(time.Millisecond * 15)
			}
		}()

		this.pageMap = make(map[string]interface{})
		this.pageMap["category-search"] = this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))

		sMessage = fmt.Sprintf("Importing <b>%d</b> Category Records", len(sliceRow)-1)
		httpRes.Write([]byte(`{"error":"` + sMessage + `","pageTitle":"Search Categories","mainpanelContentSearch":` + contentHTML + `}`))

	} else {
		sMessage += "File is empty "
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
	}

	return
}

func (this *Category) ListAll(curdb database.Database) map[string]interface{} {
	xDocResult := make(map[string]interface{})
	sqlQuery := `select code, control, title, placement, categorycontrol, image from category order by categorycontrol, placement`

	xDocResultCategory, _ := curdb.Query(sqlQuery)
	for _, xDoc := range xDocResultCategory {
		xDoc := xDoc.(map[string]interface{})
		xDocResult[xDoc["control"].(string)] = xDoc
	}
	return xDocResult
}

func (this *Category) QuickSearchTitle(xDocResult map[string]interface{}, curdb database.Database) {

	var sqlCategorySlice []string
	if xDocResult["categorycontrol"] != nil {
		sqlCategorySlice = append(sqlCategorySlice, xDocResult["categorycontrol"].(string))
	}

	if xDocResult["subcategorycontrol"] != nil {
		sqlCategorySlice = append(sqlCategorySlice, xDocResult["subcategorycontrol"].(string))
	}

	sqlQuery := fmt.Sprintf(`select code, control, title from category where control in ('%s')`, strings.Join(sqlCategorySlice, `','`))

	xDocResultCategory, _ := curdb.Query(sqlQuery)
	for _, xDoc := range xDocResultCategory {
		xDoc := xDoc.(map[string]interface{})

		if xDocResult["categorycontrol"] != nil {
			if xDoc["control"].(string) == xDocResult["categorycontrol"].(string) {
				xDocResult["categorytitle"] = xDoc["title"]
			}
		}

		if xDocResult["subcategorycontrol"] != nil {
			if xDoc["control"].(string) == xDocResult["subcategorycontrol"].(string) {
				xDocResult["subcategorytitle"] = xDoc["title"]
			}
		}
	}
}
