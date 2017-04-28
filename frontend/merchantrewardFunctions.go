package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"

	"strings"
	"time"
)

func (this *MerchantReward) viewRedeemed(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchRedeemed := make(map[string]interface{})
	subSearchRedeemed["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))
	subSearchRedeemed["merchant"] = this.MerchantControl
	subSearch["merchantrewardview-redeemed-search"] = subSearchRedeemed

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *MerchantReward) searchRedeemed(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	// tblMerchantRewardRedeemed := new(database.Redemption)
	// quickSearch := make(map[string]interface{})
	// xDocrequest := make(map[string]interface{})

	// xDocrequest["limit"] = "10"
	// xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	// xDocresult := tblMerchantRewardRedeemed.Search(xDocrequest, curdb)

	// for cNumber, xDoc := range xDocresult {
	// 	xDoc := xDoc.(map[string]interface{})
	// 	xDoc["number"] = cNumber
	// 	xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
	// 	quickSearch[cNumber+"#quick-search-result"] = xDoc
	// }

	quickSearch := make(map[string]interface{})
	quickSearch["merchantrewardview-redeemed-search"] = make(map[string]interface{})

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewDropdownHtml + `}`))
}

func (this *MerchantReward) viewCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Coupon Code is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	subSearchStore["merchant"] = this.MerchantControl
	subSearch["merchantrewardview-coupon-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *MerchantReward) searchCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Valued Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocrequest := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("code")) != "" {
		xDocrequest["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	}

	xDocrequest["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))

	xDocrequest["limit"] = "10"
	if functions.TrimEscape(httpReq.FormValue("limit")) != "" {
		xDocrequest["limit"] = functions.TrimEscape(httpReq.FormValue("limit"))
	}

	xDocrequest["offset"] = "0"
	if functions.TrimEscape(httpReq.FormValue("offset")) != "" {
		offset, err := strconv.Atoi(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("offset"))))
		if err == nil {
			if offset > 1 {
				xDocrequest["offset"] = fmt.Sprintf("%s", offset-1)
			}
		}
	}

	tblCoupon := new(database.Coupon)
	subSearch := make(map[string]interface{})
	xDocrequest["member"] = functions.TrimEscape(httpReq.FormValue("member"))
	xDocresult := tblCoupon.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "activateCoupon"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "deactivateCoupon"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#rewardview-coupon-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *MerchantReward) newCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	formNew := make(map[string]interface{})
	formNewFields := make(map[string]interface{})
	formNewFields["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))
	formNew["merchantrewardview-coupon-edit"] = formNewFields

	viewHTML := strconv.Quote(string(this.Generate(formNew, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *MerchantReward) importCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	formNew := make(map[string]interface{})
	formNewFields := make(map[string]interface{})
	formNewFields["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))
	formNew["merchantrewardview-coupon-import"] = formNewFields

	viewHTML := strconv.Quote(string(this.Generate(formNew, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *MerchantReward) importcouponcsvdownload(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	httpRes.Header().Set("Content-Type", "text/csv")
	httpRes.Header().Set("Content-Disposition", "attachment;filename=Valued_ImportCouponTemplate.csv")
	httpRes.Write([]byte(strings.Join([]string{"uniquecode", "uniquecode", "uniquecode", "etc..."}, "\n")))
}

func (this *MerchantReward) importcouponcsvsave(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	mapFiletypes := make(map[string]string)
	mapFiletypes["csvfile"] = "text/plain"
	mapFiles, sMessage := functions.UploadFile([]string{"csvfile"}, mapFiletypes, httpReq)

	if mapFiles["csvfile"] != nil {
		mapCoupon := mapFiles["csvfile"].(map[string]interface{})
		couponList := string(mapCoupon["filebytes"].([]byte))
		couponList = strings.Replace(couponList, "\r", "", -1)

		sliceRow := strings.Split(couponList, "\n")
		tblCoupon := new(database.Coupon)

		if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
			sMessage += "MerchantReward is missing <br>"

		}

		if len(sMessage) > 0 {
			httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
			return
		}

		go func() {
			for index, stringCols := range sliceRow {

				stringCols = strings.TrimSpace(stringCols)
				if index == 0 || len(stringCols) == 0 {
					continue
				}

				sliceCols := strings.Split(stringCols, ",")

				xDoc := make(map[string]interface{})
				xDoc["code"] = strings.TrimSpace(sliceCols[0])
				xDoc["title"] = strings.TrimSpace(sliceCols[0])
				xDoc["rewardcontrol"] = functions.TrimEscape(httpReq.FormValue("reward"))
				xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))

				tblCoupon.Create(this.mapCache["username"].(string), xDoc, curdb)
				<-time.Tick(time.Millisecond * 15)
			}
		}()

		sMessage = fmt.Sprintf("Importing <b>%d</b> Coupon Records", len(sliceRow)-1)
		httpRes.Write([]byte(`{"error":"` + sMessage + `","triggerSubSearch":true,"toggleSubForm":"true"}`))

	} else {
		sMessage += "File is empty "
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
	}

	return
}

func (this *MerchantReward) editCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		sMessage += "Coupon is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblCoupon := new(database.Coupon)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblCoupon.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})
	xDocRequest["reward"] = xDocRequest["rewardcontrol"]

	formView := make(map[string]interface{})
	formView["merchantrewardview-coupon-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *MerchantReward) saveCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "MerchantReward is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("code")) == "" {
		sMessage += "Coupon Code is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblCoupon := new(database.Coupon)
	xDoc := make(map[string]interface{})
	xDoc["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDoc["title"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("workflow"))

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblCoupon.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["rewardcontrol"] = functions.TrimEscape(httpReq.FormValue("reward"))
		tblCoupon.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	httpRes.Write([]byte(`{"error":"Record Saved","triggerSubSearch":true,"toggleSubForm":"true"}`))
	return
}

func (this *MerchantReward) activateCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a Coupon"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update coupon set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Coupon Activated","triggerSubSearch":true}`))
}

func (this *MerchantReward) deactivateCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a Coupon"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update coupon set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Coupon Deactivated","triggerSubSearch":true}`))
}

func (this *MerchantReward) viewStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	subSearchStore["merchant"] = this.MerchantControl
	subSearch["merchantrewardview-store-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *MerchantReward) searchStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocrequest := make(map[string]interface{})
	xDocrequest["limit"] = "10"
	if html.EscapeString(httpReq.FormValue("limit")) == "" {
		xDocrequest["limit"] = html.EscapeString(httpReq.FormValue("limit"))
	}

	xDocrequest["offset"] = "0"
	if html.EscapeString(httpReq.FormValue("offset")) == "" {
		offset, err := strconv.Atoi(strings.TrimSpace(html.EscapeString(httpReq.FormValue("offset"))))
		if err == nil {
			if offset > 1 {
				xDocrequest["offset"] = fmt.Sprintf("%s", offset-1)
			}
		}
	}

	tblMerchantRewardStore := new(database.RewardStore)
	subSearch := make(map[string]interface{})
	xDocrequest["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	xDocrequest["merchantcontrol"] = this.MerchantControl
	xDocresult := tblMerchantRewardStore.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "linkedStoreActivate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "linkedStoreDeactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#rewardview-store-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *MerchantReward) linkStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("store")) == "" {
		sMessage += "Store is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblMerchantRewardStore := new(database.RewardStore)
	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "inactive"
	xDoc["rewardcontrol"] = html.EscapeString(httpReq.FormValue("reward"))
	xDoc["storecontrol"] = html.EscapeString(httpReq.FormValue("store"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblMerchantRewardStore.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblMerchantRewardStore.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	this.searchStore(httpRes, httpReq, curdb)
}

func (this *MerchantReward) linkedStoreDelete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`delete from rewardstore where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Store Deleted","triggerSubSearch":true}`))
}

func (this *MerchantReward) linkedStoreActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardstore set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Store Activated","triggerSubSearch":true}`))
}

func (this *MerchantReward) linkedStoreDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardstore set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Store Deactivated","triggerSubSearch":true}`))
}

func (this *MerchantReward) linkedStoreDeactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a linked store"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update rewardstore set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Linked Stores Deactivated","triggerSubSearch":true}`))
}

func (this *MerchantReward) viewScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "MerchantReward is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchScheme := make(map[string]interface{})
	subSearchScheme["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	subSearchScheme["merchant"] = html.EscapeString(httpReq.FormValue("merchant"))
	subSearch["merchantrewardview-scheme-search"] = subSearchScheme

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *MerchantReward) searchScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "MerchantReward is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocrequest := make(map[string]interface{})
	xDocrequest["limit"] = "10"
	if html.EscapeString(httpReq.FormValue("limit")) == "" {
		xDocrequest["limit"] = html.EscapeString(httpReq.FormValue("limit"))
	}

	xDocrequest["offset"] = "0"
	if html.EscapeString(httpReq.FormValue("offset")) == "" {
		offset, err := strconv.Atoi(strings.TrimSpace(html.EscapeString(httpReq.FormValue("offset"))))
		if err == nil {
			if offset > 1 {
				xDocrequest["offset"] = fmt.Sprintf("%s", offset-1)
			}
		}
	}

	tblMerchantRewardScheme := new(database.RewardScheme)
	subSearch := make(map[string]interface{})
	xDocrequest["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	xDocresult := tblMerchantRewardScheme.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "linkedSchemeActivate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "linkedSchemeDeactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}
		subSearch[cNumber+"#rewardview-scheme-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *MerchantReward) linkScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Merchant Reward is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("scheme")) == "" {
		sMessage += "Scheme is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblMerchantRewardScheme := new(database.RewardScheme)
	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "inactive"
	xDoc["rewardcontrol"] = html.EscapeString(httpReq.FormValue("reward"))
	xDoc["schemecontrol"] = html.EscapeString(httpReq.FormValue("scheme"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblMerchantRewardScheme.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblMerchantRewardScheme.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	this.searchScheme(httpRes, httpReq, curdb)
}

func (this *MerchantReward) linkedSchemeDelete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`delete from rewardscheme where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Deleted","triggerSubSearch":true}`))
}

func (this *MerchantReward) linkedSchemeActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Linked Scheme Not Activated"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardscheme set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Please select a linked scheme","triggerSubSearch":true}`))
}

func (this *MerchantReward) linkedSchemeDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a linked scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardscheme set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Deactivated","triggerSubSearch":true}`))
}

func (this *MerchantReward) linkedSchemeDeactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Linked Schemes Not De-Activated"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update rewardscheme set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Linked Schemes Deactivated","triggerSubSearch":true}`))
}
