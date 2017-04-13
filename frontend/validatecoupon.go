package frontend

import (
	"fmt"
	"valued/database"
	"valued/functions"

	"net/http"
	"strconv"
	"strings"
)

type ValidateCoupon struct {
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *ValidateCoupon) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})
		this.pageMap["merchant-view-validatecoupon"] = make(map[string]interface{})

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Validate Coupon","mainpanelContent":` + contentHTML + `}`))
		return

	case "check":
		this.check(httpRes, httpReq, curdb)
		return

	case "validate":
		this.validate(httpRes, httpReq, curdb)
		return
	}
}

func (this *ValidateCoupon) check(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("code")) == "" {
		sMessage += "Promotional Code is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sqlCoupon := `select coupon.control as coupon, coupon.code as couponcode, coupon.workflow as workflow,
						reward.control as reward, reward.title as title, 
						reward.enddate as enddate, reward.description as description,
						reward.beneficiary as beneficiary, reward.restriction as restriction

					from coupon left join  reward on reward.control = coupon.rewardcontrol
					where coupon.code = '%s' and reward.merchantcontrol = '%s'
					and reward.method = 'Valued Code'
				`

	sqlCoupon = fmt.Sprintf(sqlCoupon, httpReq.FormValue("code"), this.mapCache["control"])
	xDocresult, _ := curdb.Query(sqlCoupon)

	if xDocresult["1"] == nil {
		sMessage += "This promotional code is INVALID <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocReward := xDocresult["1"].(map[string]interface{})
	xDocReward["couponcode"] = httpReq.FormValue("code")

	if xDocReward["workflow"].(string) != "active" {
		sMessage += "This promotional code has already been used <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	this.pageMap = make(map[string]interface{})
	this.pageMap["merchant-view-validatecoupon-result"] = xDocReward

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	sMessage = "This promotional code is VALID - Please apply preferential rate"
	httpRes.Write([]byte(`{"pageTitle":"Validate Coupon","error":"` + sMessage + `", "mainpanelContent":` + contentHTML + `}`))
}

func (this *ValidateCoupon) validate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("status")) == "" {
		sMessage += "Please Approve or Reject Coupon <br>"
	} else {
		switch functions.TrimEscape(httpReq.FormValue("status")) {
		case "approved", "rejected":
		default:
			sMessage += "Please Approve or Reject Coupon <br>"
		}
	}

	if functions.TrimEscape(httpReq.FormValue("coupon")) == "" {
		sMessage += "Promotional Code is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sqlCoupon := `select coupon.control from coupon left join reward on coupon.rewardcontrol = reward.control 
					where coupon.control = '%s' and reward.merchantcontrol = '%s'
					and coupon.workflow ='active' and reward.method = 'Valued Code'
				`

	sqlCoupon = fmt.Sprintf(sqlCoupon, httpReq.FormValue("coupon"), this.mapCache["control"])
	xDocresult, _ := curdb.Query(sqlCoupon)

	if xDocresult["1"] == nil {
		sMessage += "This promotional code is INVALID <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	xDoc["control"] = functions.TrimEscape(httpReq.FormValue("coupon"))
	xDoc["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))
	new(database.Coupon).Update(this.mapCache["username"].(string), xDoc, curdb)

	this.pageMap = make(map[string]interface{})
	this.pageMap["merchant-view-validatecoupon"] = make(map[string]interface{})

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	sMessage = "This promotional code has been " + strings.ToUpper(functions.TrimEscape(httpReq.FormValue("status")))
	httpRes.Write([]byte(`{"pageTitle":"Validate Coupon","error":"` + sMessage + `", "mainpanelContent":` + contentHTML + `}`))

}
