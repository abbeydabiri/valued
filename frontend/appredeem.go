package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"

	"regexp"
	"strings"
)

type AppRedeem struct {
	functions.Templates
	pageMap     map[string]interface{}
	mapAppCache map[string]interface{}
	GOSESSID    *http.Cookie
}

func (this *AppRedeem) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	this.GOSESSID, _ = httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(this.GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "view":

		viewHTML := this.view(html.EscapeString(httpReq.FormValue("control")), curdb)
		if viewHTML != "" {
			httpRes.Write([]byte(`{"pageTitle":"The Reward","pageContent":` + viewHTML + `}`))
		} else {
			httpRes.Write([]byte(`{}`))
		}
		return

	case "redeem":
		sMessage, _ := this.ValidateEligibility(html.EscapeString(httpReq.FormValue("reward")), curdb)
		if len(sMessage) > 0 {
			httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
			return
		}

		this.viewMemberPin(httpRes, httpReq, curdb)
		return

	case "checkValuedCode":
		this.checkValuedCode(httpRes, httpReq, curdb)
		return

	case "save":
		switch httpReq.FormValue("step") {
		case "memberpin":
			this.saveMemberPin(httpRes, httpReq, curdb)
			return
		case "merchantpin":
			this.saveMerchantPin(httpRes, httpReq, curdb)
			return
		case "promocode":
			this.savePromoCode(httpRes, httpReq, curdb)
			return
		case "feedback":
			this.saveFeedback(httpRes, httpReq, curdb)
			return
		}
	}
}

func (this *AppRedeem) SetAppCache(mapAppCache map[string]interface{}) {
	this.mapAppCache = mapAppCache
}

func (this *AppRedeem) ValidateEligibility(control string, curdb database.Database) (sMessage, sClass string) {

	//Using Redemption Table
	sClass = ""
	sMessage = ""

	RewardSQL := fmt.Sprintf(`select maxuse, maxperuser, maxpermonth, enddate from reward where control = '%s'`, control)
	RewardRes, _ := curdb.Query(RewardSQL)
	if RewardRes["1"] != nil {

		xDocReward := RewardRes["1"].(map[string]interface{})
		fRewardMaxUse := xDocReward["maxuse"].(int64)
		fRewardMaxPerUser := xDocReward["maxperuser"].(int64)
		fRewardMaxPerMonth := xDocReward["maxpermonth"].(int64)

		////Check if Reward MaxUse is exceeded
		MaxUseSQL := fmt.Sprintf(`select count(*) as maxuse from redemption where rewardcontrol = '%s' `, control)
		MaxUseRes, _ := curdb.Query(MaxUseSQL)
		xDocMaxUse := MaxUseRes["1"].(map[string]interface{})
		fRedemptionMaxUse := xDocMaxUse["maxuse"].(int64)

		if fRewardMaxUse > 0 && fRewardMaxUse <= fRedemptionMaxUse {
			sClass = "expired forever"
			sMessage = "Reward Quota Exceeded"
		}

		if xDocMaxUse["enddate"] != nil {
			if functions.GetDifferenceInSeconds("", xDocMaxUse["enddate"].(string)) > 0 {
				sClass = "expired forever"
				sMessage = "Reward Quota Exceeded"
			}
		}

		if len(sMessage) > 0 {
			//@TODO Update Reward set workflow to inactive
			xDocReward["workflow"] = "inactive"
			new(database.Reward).Update("root", xDocReward, curdb)
			return
		}

		//

		////Check if Reward MaxPerMonth is exceeded

		sCurrentDate := functions.GetSystemTime()
		sCurrentMonth := `%%` + sCurrentDate[2:8] + `%%`

		MaxPerMonthSQL := fmt.Sprintf(`select count(*) as maxpermonth from redemption where rewardcontrol = '%s' and createdate like '%s' `, control, sCurrentMonth)
		MaxPerMonthRes, _ := curdb.Query(MaxPerMonthSQL)
		xDocMaxPerMonth := MaxPerMonthRes["1"].(map[string]interface{})
		fRedemptionMaxPerMonth := xDocMaxPerMonth["maxpermonth"].(int64)

		if fRewardMaxPerMonth > 0 && fRewardMaxPerMonth <= fRedemptionMaxPerMonth {
			sClass = "expired month"
			sMessage = "Reward Monthly Quota Exceeded"
			return
		}

		//

		////Check if Reward MaxPerUser is exceeded
		MaxPerUserSQL := fmt.Sprintf(`select count(*) as maxperuser from redemption where rewardcontrol = '%s' and membercontrol = '%s' `, control, this.mapAppCache["control"])
		MaxPerUserRes, _ := curdb.Query(MaxPerUserSQL)
		xDocMaxPerUser := MaxPerUserRes["1"].(map[string]interface{})
		fRedemptionMaxPerUser := xDocMaxPerUser["maxperuser"].(int64)

		if fRewardMaxPerUser > 0 && fRewardMaxPerUser <= fRedemptionMaxPerUser {
			sClass = "expired user"
			sMessage = "Your Reward Quota Exceeded"
		}

	} else {
		sMessage = "Cannot Find Reward"
	}

	return
}

func (this *AppRedeem) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblRedeem := new(database.Reward)
	xDocRewardReq := make(map[string]interface{})
	xDocRewardReq["searchvalue"] = control

	ResMap := tblRedeem.Read(xDocRewardReq, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocRewardRes := ResMap["1"].(map[string]interface{})
	xDocRewardRes["rewardtitle"] = xDocRewardRes["title"]
	xDocRewardRes["createdate"] = xDocRewardRes["createdate"].(string)[0:19]

	xDocRewardRes["title"] = this.mapAppCache["title"]
	xDocRewardRes["firstname"] = this.mapAppCache["firstname"]
	xDocRewardRes["lastname"] = this.mapAppCache["lastname"]

	// xDocRewardRes["app-redeem-favorite"] =  make(map[string]interface{})
	//Search Favorite Table to see if this reward is favorited
	sqlFavorite := `select control, rewardcontrol from favorite where rewardcontrol = '%s' AND profilecontrol = '%s' order by control desc`
	sqlFavorite = fmt.Sprintf(sqlFavorite, xDocRewardRes["control"], this.mapAppCache["control"])
	xDocFavoriteRes, _ := curdb.Query(sqlFavorite)

	if xDocFavoriteRes["1"] != nil {
		xDocRewardRes["appredeemfavoriteAction"] = "remove"
		xDocRewardRes["appredeemfavoriteTitle"] = "REMOVE FROM"
	} else {
		xDocRewardRes["appredeemfavoriteAction"] = "add"
		xDocRewardRes["appredeemfavoriteTitle"] = "ADD TO"
	}

	if this.mapAppCache["control"] != nil {
		xDocRewardRes["signedIn"] = "yes"
	}
	//Search Favorite Table to see if this reward is favorited

	formView := make(map[string]interface{})
	appFooter := make(map[string]interface{})
	xDocRewardRes["app-footer"] = appFooter
	appFooter["reward"] = "white"

	formView["app-redeem"] = xDocRewardRes

	//xDocRewardRes -> add to mapAppReward
	curdb.SetSession(this.GOSESSID.Value, "mapAppReward", xDocRewardRes, false)

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *AppRedeem) viewMemberPin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if this.mapAppCache["subscription"] == nil {
		sMessage += "Please purchase a membership to be able to redeem rewards"
	}

	if this.mapAppCache["workflow"] != nil {
		switch this.mapAppCache["workflow"].(string) {
		case "subscribed", "subscribed-pending":
		default:
			sMessage += "Please purchase a membership to be able to redeem rewards"
		}

	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	this.pageMap = make(map[string]interface{})
	appRedeemStep := make(map[string]interface{})
	appRedeemStep["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))
	this.pageMap["app-redeem-memberpin"] = appRedeemStep

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Member Pin","redemptionStep":` + contentHTML + `}`))
}

func (this *AppRedeem) saveMemberPin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sPin := functions.TrimEscape(httpReq.FormValue("pin1")) + functions.TrimEscape(httpReq.FormValue("pin2")) +
		functions.TrimEscape(httpReq.FormValue("pin3")) + functions.TrimEscape(httpReq.FormValue("pin4"))

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if sPin == "" || len(sPin) != 4 {
		sMessage += "Your PIN must be 4-Digits <br>"
	}

	if sPin != this.mapAppCache["pincode"].(string) {
		sMessage += "Incorrect PIN entered <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")

	switch mapAppReward["method"].(string) {
	case "Pin", "Valued Code", "Client Code Bulk", "Client Code Single":
	default:

		httpRes.Write([]byte(`{"error":"Reward Method [` + mapAppReward["method"].(string) + `] is Invalid"}`))
		return
	}

	//Check/Generate for Next Valid Coupon
	switch mapAppReward["method"].(string) {
	case "Pin", "Valued Code":
		xDoc := make(map[string]interface{})
		xDoc["code"] = functions.RandomString(6)
		xDoc["workflow"] = "active"
		xDoc["rewardcontrol"] = functions.TrimEscape(httpReq.FormValue("reward"))
		new(database.Coupon).Create(this.mapAppCache["username"].(string), xDoc, curdb)
	}

	sqlPromoCode := `select code as couponcode, control as couponcontrol, rewardcontrol from coupon where rewardcontrol = '%s' AND workflow = 'active' order by couponcontrol limit 1`
	sqlPromoCode = fmt.Sprintf(sqlPromoCode, functions.TrimEscape(httpReq.FormValue("reward")))
	xDocresult, _ := curdb.Query(sqlPromoCode)
	if xDocresult["1"] != nil {

		promoXdoc := xDocresult["1"].(map[string]interface{})
		for key, value := range promoXdoc {
			mapAppReward[key] = value
		}

		curdb.SetSession(this.GOSESSID.Value, "mapAppReward", mapAppReward, false)

	} else {
		sMessage += "Promo Codes Exhausted - Merchant has been alerted!!"

		//Send an Email to Merchant
		//SEND AN EMAIL USING TEMPLATE
		emailTo := ""
		emailFrom := "redemptions@valued.com"
		emailFromName := "VALUED PROMO CODES"
		emailTemplate := "app-merchant-request-promocode"
		emailSubject := fmt.Sprintf("RELEASE OF NEW PROMO CODES")
		emailFields := make(map[string]interface{})

		sqlMerchant := fmt.Sprintf(`select employer.email as email, employer.title as employertitle from profile as employer
			left join reward on reward.merchantcontrol = employer.control where reward.control = '%s' `,
			functions.TrimEscape(httpReq.FormValue("reward")))

		resMember, _ := curdb.Query(sqlMerchant)
		xDocEmail := resMember["1"].(map[string]interface{})

		if xDocEmail["email"] != nil && xDocEmail["email"].(string) != "" {
			emailTo = xDocEmail["email"].(string)
			emailFields["title"] = xDocEmail["employertitle"]

			go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, "", emailFields)
		}
		//SEND AN EMAIL USING TEMPLATE
		//Send an Email to Merchant

	}
	//Check/Generate for Next Valid Coupon

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	switch mapAppReward["method"].(string) {
	case "Pin":
		this.viewMerchantPin(httpRes, httpReq, curdb)
	case "Valued Code", "Client Code Bulk", "Client Code Single":
		this.viewPromoCode(httpRes, httpReq, curdb)
	}
}

//viewMerchantPin or viewMerchantPromoCode
func (this *AppRedeem) viewMerchantPin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	this.pageMap = make(map[string]interface{})
	appRedeemStep := make(map[string]interface{})
	appRedeemStep["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))

	//Find Stores linked to this reward
	storeSQL := fmt.Sprintf(`select s.title as title, s.city as city, s.control as control from 
		rewardstore as rs, store as s where rs.storecontrol = s.control AND rs.rewardcontrol = '%s'`, appRedeemStep["reward"])
	storeMapResult, _ := curdb.Query(storeSQL)
	appRedeemStore := make(map[string]interface{})
	for cNumber, storeXdoc := range storeMapResult {
		xDoc := storeXdoc.(map[string]interface{})
		appRedeemStore[cNumber+"#app-redeem-store-option"] = xDoc
	}
	if len(appRedeemStore) > 0 {
		if len(appRedeemStore) == 1 {
			appRedeemStore["state"] = "hide"
		}
		appRedeemStep["app-redeem-store"] = appRedeemStore
	}
	//Find Stores linked to this reward

	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")
	appRedeemStep["couponcode"] = mapAppReward["couponcode"]

	this.pageMap["app-redeem-merchantpin"] = appRedeemStep

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Merchant Pin","redemptionStep":` + contentHTML + `}`))
}

func (this *AppRedeem) saveMerchantPin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""

	sPin := functions.TrimEscape(httpReq.FormValue("pin1")) + functions.TrimEscape(httpReq.FormValue("pin2")) +
		functions.TrimEscape(httpReq.FormValue("pin3")) + functions.TrimEscape(httpReq.FormValue("pin4"))

	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if sPin == "" || len(sPin) != 4 {
		sMessage += "Your PIN must be 4-Digits <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("transactionvalue")) == "" {
		sMessage += "Transaction Value is missing <br>"
	} else {
		regexNumber := regexp.MustCompile("^-*([0-9]+)*\\.*([0-9]+)$")
		if regexNumber.MatchString(strings.TrimSpace(html.EscapeString(httpReq.FormValue("transactionvalue")))) == false {
			sMessage += "Transaction Value must be Numeric!! <br>"
		}
	}

	//Search for MerchantPin via Reward
	sqlMerchantPin := "select m.pincode as merchantpin from reward as r left join profile as m on m.control = r.merchantcontrol where r.control = '%s'"
	sqlMerchantPin = fmt.Sprintf(sqlMerchantPin, functions.TrimEscape(httpReq.FormValue("reward")))
	mapMerchantPin, _ := curdb.Query(sqlMerchantPin)

	if mapMerchantPin["1"] != nil {
		merchantPin := mapMerchantPin["1"].(map[string]interface{})["merchantpin"].(string)

		if merchantPin != sPin {
			sMessage += "Incorrect PIN entered <br>"
		}

	} else {
		sMessage += "Incorrect PIN entered <br>"
	}

	//Search for MerchantPin via Reward

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")
	xDoc["control"] = mapAppReward["couponcontrol"]
	xDoc["workflow"] = "approved"
	new(database.Coupon).Update(this.mapAppCache["username"].(string), xDoc, curdb)

	this.viewFeedback(httpRes, httpReq, curdb)
}

func (this *AppRedeem) viewPromoCode(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	this.pageMap = make(map[string]interface{})
	appRedeemStep := make(map[string]interface{})
	appRedeemStep["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))

	//Find Stores linked to this reward
	storeSQL := fmt.Sprintf(`select s.title as title, s.city as city, s.control as control from 
		rewardstore as rs, store as s where rs.storecontrol = s.control AND rs.rewardcontrol = '%s'`, appRedeemStep["reward"])
	storeMapResult, _ := curdb.Query(storeSQL)
	appRedeemStore := make(map[string]interface{})
	for cNumber, storeXdoc := range storeMapResult {
		xDoc := storeXdoc.(map[string]interface{})
		appRedeemStore[cNumber+"#app-redeem-store-option"] = xDoc
	}

	if len(appRedeemStore) > 0 {
		if len(appRedeemStore) == 1 {
			appRedeemStore["state"] = "hide"
		}
		appRedeemStep["app-redeem-store"] = appRedeemStore
	}
	//Find Stores linked to this reward

	//

	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")
	appRedeemStep["couponcode"] = mapAppReward["couponcode"]

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	if mapAppReward["method"].(string) == "Valued Code" {
		appRedeemStep["app-redeem-promocode-valuedcode-check"] = make(map[string]interface{})
	}

	this.pageMap["app-redeem-promocode"] = appRedeemStep

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Merchant Promo Code","redemptionStep":` + contentHTML + `}`))
}

func (this *AppRedeem) checkValuedCode(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")

	sqlPromoCode := `select workflow from coupon where control = '%s'`
	sqlPromoCode = fmt.Sprintf(sqlPromoCode, mapAppReward["couponcontrol"])
	xDocresult, _ := curdb.Query(sqlPromoCode)
	if xDocresult["1"] == nil {
		this.viewRejected(httpRes, httpReq, curdb)
		return
	}

	sWorkflow := xDocresult["1"].(map[string]interface{})["workflow"].(string)
	if sWorkflow == "rejected" {
		this.viewRejected(httpRes, httpReq, curdb)
		return
	}

	if sWorkflow == "approved" {
		sMessage := "This promotional code is VALID - <br>Please Enter Transaction Value To Complete Redemption"
		httpRes.Write([]byte(`{"sticky":"` + sMessage + `","checkValuedCode":"<script>clearInterval(checkValuedCodeTimer);</script>"}`))
		return
	}

	httpRes.Write([]byte(`{}`))
}

func (this *AppRedeem) savePromoCode(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("transactionvalue")) == "" {
		sMessage += "Transaction Value is missing <br>"
	} else {
		regexNumber := regexp.MustCompile("^-*([0-9]+)*\\.*([0-9]+)$")
		if regexNumber.MatchString(strings.TrimSpace(html.EscapeString(httpReq.FormValue("transactionvalue")))) == false {
			sMessage += "Transaction Value must be Numeric!! <br>"
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//CHECK IF VALUED CODE HAS BEEN VALIDATED BY MERCHANT
	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")
	switch mapAppReward["method"].(string) {
	case "Valued Code":

		sqlPromoCode := `select workflow from coupon where control = '%s'`
		sqlPromoCode = fmt.Sprintf(sqlPromoCode, mapAppReward["couponcontrol"])
		xDocresult, _ := curdb.Query(sqlPromoCode)
		if xDocresult["1"] == nil {
			this.viewRejected(httpRes, httpReq, curdb)
			return
		}

		sWorkflow := xDocresult["1"].(map[string]interface{})["workflow"].(string)
		if sWorkflow == "rejected" {
			this.viewRejected(httpRes, httpReq, curdb)
			return
		}

		if sWorkflow != "approved" {
			sMessage = "Coupon has not been Approved"
		}
	}
	//CHECK IF VALUED CODE HAS BEEN VALIDATED BY MERCHANT

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	this.viewFeedback(httpRes, httpReq, curdb)
}

//viewMerchantPin or viewMerchantPromoCode
func (this *AppRedeem) viewFeedback(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("reward")) == "" {
		sMessage += "Reward is missing <br>"
	}

	transactionvalue, err := strconv.ParseFloat(functions.TrimEscape(httpReq.FormValue("transactionvalue")), 10)
	if err != nil {
		sMessage += "Transaction Value must be Numeric!! <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("store")) == "" {
		sMessage += "Store is missing! <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	this.pageMap = make(map[string]interface{})
	// appRedeemStep := make(map[string]interface{})
	// appRedeemStep["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))

	tblRedeem := new(database.Reward)
	xDocRewardReq := make(map[string]interface{})
	xDocRewardReq["searchvalue"] = functions.TrimEscape(httpReq.FormValue("reward"))

	ResMap := tblRedeem.Read(xDocRewardReq, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"Reward Not Found"}`))
		return
	}

	xDocRewardRes := ResMap["1"].(map[string]interface{})
	xDocRewardRes["createdate"] = xDocRewardRes["createdate"].(string)[0:19]

	//
	//Get Reward Savings

	if xDocRewardRes["discountvalue"] != nil {
		fDiscount := xDocRewardRes["discountvalue"].(float64)
		sDiscountType := xDocRewardRes["discounttype"].(string)

		switch strings.TrimSpace(sDiscountType) {
		default:
			sMessage = "Reward has no Valid Discount Setup!"

		case "%":
			xDocRewardRes["savingsvalue"] = functions.RoundUp((fDiscount/100.00)*transactionvalue, 2)

		case "Off":
			xDocRewardRes["savingsvalue"] = functions.RoundUp(fDiscount, 2)
			// xDocRewardRes["savingsvalue"] = functions.RoundUp(transactionvalue-fDiscount, 2)

		}
	}

	/*
		switch {
		default:
			sMessage = "Reward has no Valid Discount Setup!"

		case strings.HasSuffix(sDiscount, `%`):
			discountvalue, err := strconv.ParseFloat(strings.Replace(strings.ToLower(sDiscount), "%", "", -1), 10)
			if err != nil {
				sMessage += "Discount must be Numeric with %% sign !! <br>"
			} else {
				xDocRewardRes["savingsvalue"] = functions.RoundUp((discountvalue/100.00)*transactionvalue, 2)
			}

		case strings.HasSuffix(strings.ToLower(sDiscount), `off`):

			discountvalue, err := strconv.ParseFloat(strings.Replace(strings.ToLower(sDiscount), "off", "", -1), 10)
			if err != nil {
				sMessage += "Discount must be Numeric e.g 100 OFF !! <br>"
			} else {
				xDocRewardRes["savingsvalue"] = functions.RoundUp(transactionvalue-discountvalue, 2)
			}

		case strings.ToLower(sDiscount) == "free":
			xDocRewardRes["savingsvalue"] = transactionvalue
		}
	*/
	xDocRewardRes["transactionvalue"] = transactionvalue

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	if xDocRewardRes["savingsvalue"] != nil {
		xDocRewardRes["savings"] = functions.Round(xDocRewardRes["savingsvalue"].(float64))
	}

	xDocRewardRes["storecontrol"] = functions.TrimEscape(httpReq.FormValue("store"))

	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")
	xDocRewardRes["couponcode"] = mapAppReward["couponcode"]
	xDocRewardRes["couponcontrol"] = mapAppReward["couponcontrol"]

	// if functions.TrimEscape(httpReq.FormValue("couponcontrol")) != "" {
	// 	xDocRewardRes["couponcontrol"] = functions.TrimEscape(httpReq.FormValue("couponcontrol"))
	// }

	//savePromoCode

	//xDocRewardRes -> add to mapAppReward
	curdb.SetSession(this.GOSESSID.Value, "mapAppReward", xDocRewardRes, false)

	//Get Reward Savings
	//

	//Get Review Categories for Feedback Purposes
	sRvNumber := ""
	tblReview := new(database.ReviewCategoryLink)
	xDocReviewReq := make(map[string]interface{})
	xDocReviewReq["merchant"] = xDocRewardRes["merchantcontrol"]
	xDocReviewRes := tblReview.Search(xDocReviewReq, curdb)
	for cNumber, xDoc := range xDocReviewRes {
		xDoc := xDoc.(map[string]interface{})
		sRvNumber = ""
		switch cNumber {
		case "1":
			sRvNumber = "One"
		case "2":
			sRvNumber = "Two"
		case "3":
			sRvNumber = "Three"
		case "4":
			sRvNumber = "Four"
		}

		if sRvNumber != "" {
			// appRedeemStep["rv"+sRvNumber] = xDoc["reviewcategorytitle"]
			// appRedeemStep["rvc"+sRvNumber] = xDoc["reviewcategorycontrol"]
			xDocRewardRes["rv"+sRvNumber] = xDoc["reviewcategorytitle"]
			xDocRewardRes["rvc"+sRvNumber] = xDoc["reviewcategorycontrol"]
		}
	}
	//Get Review Categories for Feedback Purposes
	// this.pageMap["app-redeem-feedback"] = appRedeemStep

	xDocRewardRes["employertitle"] = this.mapAppCache["employertitle"]

	this.pageMap["app-redeem-feedback"] = xDocRewardRes
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Give Feedback","redemptionStep":` + contentHTML + `}`))
}

func (this *AppRedeem) saveFeedback(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("feedbackRate")) == "" {
		sMessage = "Please Rate Merchant to Complete Redemption Process! <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("reviewcategory")) == "" {
		sMessage = "Please Rate Merchant to Complete Redemption Process! <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")

	xDocRedemption := make(map[string]interface{})
	xDocRedemption["employercontrol"] = this.mapAppCache["employercontrol"]
	xDocRedemption["membercontrol"] = this.mapAppCache["control"]
	xDocRedemption["merchantcontrol"] = mapAppReward["merchantcontrol"]

	xDocRedemption["schemecontrol"] = ""
	xDocRedemption["couponcontrol"] = mapAppReward["couponcontrol"]

	xDocRedemption["rewardcontrol"] = mapAppReward["control"]
	xDocRedemption["storecontrol"] = mapAppReward["storecontrol"]
	xDocRedemption["savingsvalue"] = mapAppReward["savingsvalue"]
	xDocRedemption["transactionvalue"] = mapAppReward["transactionvalue"]

	//New Fields to allow User Deletion
	xDocRedemption["reward"] = mapAppReward["title"]
	xDocRedemption["location"] = mapAppReward["address"]
	xDocRedemption["discount"] = mapAppReward["discount"]
	xDocRedemption["dob"] = this.mapAppCache["dob"]
	xDocRedemption["gender"] = this.mapAppCache["title"]
	xDocRedemption["nationality"] = this.mapAppCache["nationality"]
	//New Fields to allow User Deletion

	xDocRedemption["workflow"] = "approved"

	xDocRedemption["title"] = fmt.Sprintf(`%s %s`, mapAppReward["merchanttitle"], mapAppReward["title"])
	sRedemptionControl := new(database.Redemption).Create(this.mapAppCache["username"].(string), xDocRedemption, curdb)

	xDocFeedback := make(map[string]interface{})
	xDocFeedback["title"] = fmt.Sprintf(`HOW LIKELY ARE YOU TO RECOMMEND %s TO YOUR FRIENDS`, mapAppReward["merchanttitle"])
	xDocFeedback["answer"] = functions.TrimEscape(httpReq.FormValue("feedbackRate"))
	xDocFeedback["redemptioncontrol"] = sRedemptionControl
	xDocFeedback["workflow"] = "active"
	new(database.Feedback).Create(this.mapAppCache["username"].(string), xDocFeedback, curdb)

	xDocFeedback = make(map[string]interface{})
	xDocFeedback["title"] = fmt.Sprintf(`WHERE WOULD YOU LIKE TO SEE THE MOST IMPROVEMENT WHEN YOU NEXT PURCHASE WITH %s`, mapAppReward["merchanttitle"])
	xDocFeedback["answer"] = functions.TrimEscape(httpReq.FormValue("reviewcategory"))
	xDocFeedback["redemptioncontrol"] = sRedemptionControl
	xDocFeedback["workflow"] = "active"
	new(database.Feedback).Create(this.mapAppCache["username"].(string), xDocFeedback, curdb)

	//Send an Email to Merchant
	//SEND AN EMAIL USING TEMPLATE
	emailTo := ""
	emailFrom := "redemptions@valued.com"
	emailFromName := "VALUED REDEMPTION"
	emailFields := make(map[string]interface{})

	//Query to Find Merchant, Reward and Coupon Details
	sqlRedemption := fmt.Sprintf(`select merchant.email as merchantemail, merchant.title as merchanttitle, reward.title as rewardtitle, 
		reward.method as rewardmethod, coupon.code as couponcode , redemption.createdate as transactiondate, 
		redemption.control as transactionnumber, redemption.transactionvalue as transactionvalue 
		from redemption left join profile as merchant on merchant.control = redemption.merchantcontrol 
		left join reward on reward.control = redemption.rewardcontrol left join coupon on coupon.control = redemption.couponcontrol  
		where redemption.control = '%s'`, sRedemptionControl)

	//Query to Find Merchant, Reward and Coupon Details
	resRedemption, _ := curdb.Query(sqlRedemption)
	xDocEmail := resRedemption["1"].(map[string]interface{})
	if xDocEmail["merchantemail"] != nil && xDocEmail["merchantemail"].(string) != "" {
		emailTo = xDocEmail["merchantemail"].(string)
		emailFields["merchanttitle"] = xDocEmail["merchanttitle"]
		emailFields["transactiondate"] = xDocEmail["transactiondate"]
		emailFields["transactionnumber"] = xDocEmail["transactionnumber"]
		emailFields["transactionvalue"] = fmt.Sprintf(`%v`, xDocEmail["transactionvalue"])
		emailFields["rewardtitle"] = xDocEmail["rewardtitle"]

		emailSubject := fmt.Sprintf("VALUED REDEMPTION - %s", xDocEmail["couponcode"])
		emailTemplate := "app-merchant-redemption-promocode"

		switch xDocEmail["rewardmethod"].(string) {
		case "Pin":
			emailTemplate = "app-merchant-redemption-pincode"
		}

		go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, "", emailFields)
	}
	//SEND AN EMAIL USING TEMPLATE
	//Send an Email to Merchant

	this.viewFinal(httpRes, httpReq, curdb)
}

func (this *AppRedeem) viewFinal(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	this.pageMap = make(map[string]interface{})

	appRedeemStep := make(map[string]interface{})
	appRedeemStep["reward"] = functions.TrimEscape(httpReq.FormValue("reward"))
	this.pageMap["app-redeem-final"] = appRedeemStep

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Redemtion Complete","redemptionStep":` + contentHTML + `}`))
}

func (this *AppRedeem) viewRejected(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	mapAppReward := curdb.GetSession(this.GOSESSID.Value, "mapAppReward")

	this.pageMap = make(map[string]interface{})
	appRedeemStep := make(map[string]interface{})
	appRedeemStep["reward"] = mapAppReward["rewardcontrol"]
	this.pageMap["app-redeem-final-rejected"] = appRedeemStep

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Redemtion Rejected","redemptionStep":` + contentHTML + `}`))
}

//

//
