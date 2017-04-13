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

func (this *Reward) viewRedeemed(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchRedeemed := make(map[string]interface{})
	subSearchRedeemed["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))
	subSearchRedeemed["merchant"] = functions.TrimEscape(httpReq.FormValue("merchant"))
	subSearch["rewardview-redeemed-search"] = subSearchRedeemed

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Reward) searchRedeemed(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	// tblRewardRedeemed := new(database.Redemption)
	// quickSearch := make(map[string]interface{})
	// xDocrequest := make(map[string]interface{})

	// xDocrequest["limit"] = "10"
	// xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	// xDocresult := tblRewardRedeemed.Search(xDocrequest, curdb)

	// for cNumber, xDoc := range xDocresult {
	// 	xDoc := xDoc.(map[string]interface{})
	// 	xDoc["number"] = cNumber
	// 	xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
	// 	quickSearch[cNumber+"#quick-search-result"] = xDoc
	// }

	quickSearch := make(map[string]interface{})
	quickSearch["rewardview-redeemed-search"] = make(map[string]interface{})

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewDropdownHtml + `}`))
}

func (this *Reward) replaceCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("code")) == "" {
		sMessage += "Coupon Code is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sqlDisable := fmt.Sprintf(`update coupon set workflow = 'inactive' where workflow = 'active' and rewardcontrol = '%s'`, functions.TrimEscape(httpReq.FormValue("reward")))
	curdb.Query(sqlDisable)

	xDocCoupon := make(map[string]interface{})
	xDocCoupon["workflow"] = "active"
	xDocCoupon["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDocCoupon["title"] = xDocCoupon["code"]
	xDocCoupon["rewardcontrol"] = functions.TrimEscape(httpReq.FormValue("reward"))
	new(database.Coupon).Create(this.mapCache["username"].(string), xDocCoupon, curdb)

	sMessage = "Coupon Code replaced successfully"

	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("reward")), curdb)
	httpRes.Write([]byte(`{"error":"` + sMessage + `","mainpanelContent":` + viewHTML + `}`))
}

func (this *Reward) generateCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	intGenerate := 0
	var sliceRow []string
	tblCoupon := new(database.Coupon)

	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("generate")) == "" {
		sMessage += "Generate Coupons is missing <br>"
	} else {
		intGenerateErr, err := strconv.Atoi(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("generate"))))
		if err != nil {
			sMessage += err.Error() + "!!! <br>"
		} else {
			intGenerate = intGenerateErr
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sqlCoupon := `select coupon.code as code from coupon left join reward on reward.control = coupon.rewardcontrol 
					where reward.merchantcontrol in (select merchantcontrol from reward where control = '%s')`
	sqlCoupon = fmt.Sprintf(sqlCoupon, functions.TrimEscape(httpReq.FormValue("reward")))
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

	if len(sliceRow) > 0 {
		go func() {
			for _, stringCols := range sliceRow {

				stringCols = strings.TrimSpace(stringCols)
				sliceCols := strings.Split(stringCols, ",")

				xDocCoupon := make(map[string]interface{})
				xDocCoupon["code"] = strings.TrimSpace(sliceCols[0])
				xDocCoupon["title"] = strings.TrimSpace(sliceCols[0])
				xDocCoupon["rewardcontrol"] = functions.TrimEscape(httpReq.FormValue("reward"))
				xDocCoupon["workflow"] = "active"

				tblCoupon.Create(this.mapCache["username"].(string), xDocCoupon, curdb)
				<-time.Tick(time.Millisecond * 15)
			}
		}()
	}

	sMessage = fmt.Sprintf(" <b>%d</b> Coupon Codes Generated", len(sliceRow))
	viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("reward")), curdb)
	httpRes.Write([]byte(`{"error":"` + sMessage + `","mainpanelContent":` + viewHTML + `}`))
}

func (this *Reward) importCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
			sMessage += "Reward is missing <br>"

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
				xDoc["workflow"] = "active"

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

func (this *Reward) importcouponcsvdownload(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	httpRes.Header().Set("Content-Type", "text/csv")
	httpRes.Header().Set("Content-Disposition", "attachment;filename=Valued_ImportCouponTemplate.csv")
	httpRes.Write([]byte(strings.Join([]string{"code"}, ",")))
}

func (this *Reward) downloadActiveCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	var sMerchant, sReward string

	sqlFilename := fmt.Sprintf(`select reward.title as reward, profile.title as merchant from
		reward left join profile on profile.control = reward.merchantcontrol
		where reward.control = '%s'
		`, functions.TrimEscape(httpReq.FormValue("reward")))
	resFilename, _ := curdb.Query(sqlFilename)
	if resFilename["1"] != nil {
		xdocFilename := resFilename["1"].(map[string]interface{})
		sReward = xdocFilename["reward"].(string)
		sMerchant = xdocFilename["merchant"].(string)
	}

	sReward = "Download-Active-Coupons"

	sqlCoupon := fmt.Sprintf(`select code from coupon where workflow = 'active' and rewardcontrol = '%s'
					`, functions.TrimEscape(httpReq.FormValue("reward")))
	resCoupon, _ := curdb.Query(sqlCoupon)

	var activeCSV []string
	for _, xDoc := range resCoupon {
		xDocCoupon := xDoc.(map[string]interface{})
		activeCSV = append(activeCSV, xDocCoupon["code"].(string))
	}

	sFilename := fmt.Sprintf("%s-%s.csv", sMerchant, sReward)
	httpRes.Header().Set("Content-Type", "text/csv")
	httpRes.Header().Set("Content-Disposition", "attachment;filename="+sFilename)
	httpRes.Write([]byte(strings.Join(activeCSV, "\n")))
}

func (this *Reward) viewCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	mapCoupon := make(map[string]interface{})
	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})

	sqlUsed := fmt.Sprintf(`select count(control) from coupon where workflow != 'active' and rewardcontrol = '%s'`,
		functions.TrimEscape(httpReq.FormValue("reward")))

	sqlActive := fmt.Sprintf(`select count(control) from coupon where workflow = 'active' and rewardcontrol = '%s'`,
		functions.TrimEscape(httpReq.FormValue("reward")))

	resUsed, _ := curdb.Query(sqlUsed)
	resActive, _ := curdb.Query(sqlActive)

	mapCoupon["used"] = 0
	mapCoupon["active"] = 0

	if resUsed["1"] != nil {
		mapCoupon["used"] = resUsed["1"].(map[string]interface{})["count"]
	}

	if resActive["1"] != nil {
		mapCoupon["active"] = resActive["1"].(map[string]interface{})["count"]
	}

	sqlRewardMethod := fmt.Sprintf(`select method from reward where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("reward")))
	resRewardMethod, _ := curdb.Query(sqlRewardMethod)
	if resRewardMethod["1"] != nil {

		sMethod := resRewardMethod["1"].(map[string]interface{})["method"].(string)

		switch sMethod {
		case "Pin", "Valued":
			subSearchStore["rewardview-coupon-search-valued"] = mapCoupon

		case "Client Code Bulk":
			subSearchStore["rewardview-coupon-search-bulk"] = mapCoupon

		case "Client Code Single":
			sqlCoupon := fmt.Sprintf(`select code from coupon where workflow = 'active' and rewardcontrol = '%s'`,
				functions.TrimEscape(httpReq.FormValue("reward")))

			resCoupon, _ := curdb.Query(sqlCoupon)
			if resCoupon["1"] != nil {
				mapCoupon["coupon"] = resCoupon["1"].(map[string]interface{})["code"]
			}

			subSearchStore["rewardview-coupon-search-single"] = mapCoupon
		}
	}

	subSearchStore["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	subSearchStore["merchant"] = functions.TrimEscape(httpReq.FormValue("merchant"))
	subSearch["rewardview-coupon-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

/*
	func (this *Reward) searchCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
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

	func (this *Reward) newCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		formNew["rewardview-coupon-edit"] = formNewFields

		viewHTML := strconv.Quote(string(this.Generate(formNew, nil)))
		httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
	}

	func (this *Reward) editCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
		formView["rewardview-coupon-edit"] = xDocRequest

		viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
		httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
	}

	func (this *Reward) saveCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		sMessage := ""
		if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
			sMessage += "Reward is missing <br>"
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

	func (this *Reward) activateCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if html.EscapeString(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a coupon"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update coupon set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		httpRes.Write([]byte(`{"error":"Coupon Activated","triggerSubSearch":true}`))
	}

	func (this *Reward) deactivateCoupon(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

		if html.EscapeString(httpReq.FormValue("control")) == "" {
			httpRes.Write([]byte(`{"error":"Please select a coupon"}`))
			return
		}

		curdb.Query(fmt.Sprintf(`update coupon set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
			this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

		httpRes.Write([]byte(`{"error":"Coupon Deactivated","triggerSubSearch":true}`))
	}
*/

func (this *Reward) viewStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	subSearchStore["merchant"] = html.EscapeString(httpReq.FormValue("merchant"))
	subSearch["rewardview-store-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Reward) searchStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
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

	tblRewardStore := new(database.RewardStore)
	subSearch := make(map[string]interface{})
	xDocrequest["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	xDocresult := tblRewardStore.Search(xDocrequest, curdb)

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

func (this *Reward) linkStore(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("store")) == "" {
		sMessage += "Store is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblRewardStore := new(database.RewardStore)
	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "inactive"
	xDoc["rewardcontrol"] = html.EscapeString(httpReq.FormValue("reward"))
	xDoc["storecontrol"] = html.EscapeString(httpReq.FormValue("store"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblRewardStore.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblRewardStore.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	this.searchStore(httpRes, httpReq, curdb)
}

func (this *Reward) linkedStoreDelete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`delete from rewardstore where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Stored Deleted","triggerSubSearch":true}`))
}

func (this *Reward) linkedStoreActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardstore set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Stored Activated","triggerSubSearch":true}`))
}

func (this *Reward) linkedStoreDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardstore set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Stored Deactivated","triggerSubSearch":true}`))
}

func (this *Reward) linkedStoreDeactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a store"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update rewardstore set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Linked Stores Deactivated","triggerSubSearch":true}`))
}

func (this *Reward) viewScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
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
	subSearch["rewardview-scheme-search"] = subSearchScheme

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Reward) searchScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
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

	tblRewardScheme := new(database.RewardScheme)
	subSearch := make(map[string]interface{})
	xDocrequest["reward"] = html.EscapeString(httpReq.FormValue("reward"))
	xDocresult := tblRewardScheme.Search(xDocrequest, curdb)

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

func (this *Reward) linkScheme(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
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

	tblRewardScheme := new(database.RewardScheme)
	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "inactive"
	xDoc["rewardcontrol"] = html.EscapeString(httpReq.FormValue("reward"))
	xDoc["schemecontrol"] = html.EscapeString(httpReq.FormValue("scheme"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblRewardScheme.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblRewardScheme.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	this.searchScheme(httpRes, httpReq, curdb)
}

func (this *Reward) linkedSchemeDelete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`delete from rewardscheme where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Deleted","triggerSubSearch":true}`))
}

func (this *Reward) linkedSchemeActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardscheme set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Activated","triggerSubSearch":true}`))
}

func (this *Reward) linkedSchemeDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update rewardscheme set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Schemed Deactivated","triggerSubSearch":true}`))
}

func (this *Reward) linkedSchemeDeactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a scheme"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update rewardscheme set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Linked Schemes Deactivated","triggerSubSearch":true}`))
}
