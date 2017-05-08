package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	// "html"
	"net/http"
	"strconv"
	"strings"

	"regexp"
)

type AppTell struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppTell) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		this.pageMap = make(map[string]interface{})
		AppTell := make(map[string]interface{})

		appNavbar := new(AppNavbar)
		AppTell["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
		AppTell["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

		if this.mapAppCache["control"] != nil {
			// AppTell["nameState"] = "disabled"
			AppTell["yourname"] = fmt.Sprintf("%s %s %s", this.mapAppCache["title"], this.mapAppCache["firstname"], this.mapAppCache["lastname"])
		}

		AppTell["app-footer"] = make(map[string]interface{})
		this.pageMap["app-tell"] = AppTell

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Tell A Friend","pageContent":` + contentHTML + `}`))
		return

	case "tell":
		this.tell(httpRes, httpReq, curdb)
		return

	}
}

func (this *AppTell) tell(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("friendfirstname")) == "" {
		sMessage += "Friend Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("friendlastname")) == "" {
		sMessage += "Friend Surname is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("friendemail")) == "" {
		sMessage += "Friend Email is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	regexValidEmail := "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"
	lValidEmail, _ := regexp.MatchString(regexValidEmail, functions.TrimEscape(httpReq.FormValue("friendemail")))
	if !lValidEmail {
		sMessage += "Friend Email is not a valid email address <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//--> Check if Email/Username exists in Login Table
	sqlCheckEmail := fmt.Sprintf(`select control, username, title from profile where username = '%s'`,
		functions.TrimEscape(httpReq.FormValue("friendemail")))
	mapCheckEmail, _ := curdb.Query(sqlCheckEmail)

	if mapCheckEmail["1"] != nil {
		sMessage += fmt.Sprintf(`Your Friend with email <b>%s</b> is already a Valued Member<br>`,
			functions.TrimEscape(httpReq.FormValue("friendemail")))
	}
	//--> Check if Email/Username exists in Login Table

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sCouponCode := ""
	sCouponMessage := ""
	//Generate A Refferal Coupon Code

	//--//Search for Existing Refferal Coupon Code
	sqlFindLastCoupon := `select distinct(code) as code from referral where lower(friendemail) = lower('%s') AND lower(email) = lower('%s')`
	sqlFindLastCoupon = fmt.Sprintf(sqlFindLastCoupon, functions.TrimEscape(httpReq.FormValue("friendemail")), this.mapAppCache["username"])
	sqlFindLastCouponRes, _ := curdb.Query(sqlFindLastCoupon)
	if sqlFindLastCouponRes["1"] != nil {
		xDocFindLastCoupon := sqlFindLastCouponRes["1"].(map[string]interface{})
		if xDocFindLastCoupon["code"] != nil && xDocFindLastCoupon["code"].(string) != "" {
			sCouponCode = xDocFindLastCoupon["code"].(string)
		}
	}
	//--//Search for Existing Refferal Coupon Code

	//--//Search for tellafriend Reward if Active generate a CouponCode
	if sCouponCode == "" {
		sqlFindReward := `select reward.method as rewardmethod, reward.control as rewardcontrol from reward 
						left join profile as merchant on merchant.control = reward.merchantcontrol
						where lower(reward.code) = 'tellafriend' AND lower(merchant.code) = 'main' AND reward.workflow = 'active'`

		sqlFindRewardRes, _ := curdb.Query(sqlFindReward)
		if sqlFindRewardRes["1"] != nil {
			xDocFindRewardRes := sqlFindRewardRes["1"].(map[string]interface{})

			switch xDocFindRewardRes["rewardmethod"].(string) {
			case "Pin", "Valued Code":

				xDoc := make(map[string]interface{})
				sCouponCode = functions.RandomString(6)

				xDoc["code"] = sCouponCode
				xDoc["workflow"] = "active"
				xDoc["rewardcontrol"] = xDocFindRewardRes["rewardcontrol"]
				new(database.Coupon).Create(this.mapAppCache["username"].(string), xDoc, curdb)

			case "Client Code":
				sqlFindCoupon := `select coupon.code as code from coupon where rewardcontrol = '%s' AND workflow = 'active' order by control limit 1`
				sqlFindCoupon = fmt.Sprintf(sqlFindCoupon, xDocFindRewardRes["rewardcontrol"])

				sqlFindCouponRes, _ := curdb.Query(sqlFindCoupon)
				if sqlFindCouponRes["1"] != nil {
					xDocFindCoupon := sqlFindCouponRes["1"].(map[string]interface{})
					if xDocFindCoupon["code"] != nil {
						sCouponCode = xDocFindCoupon["code"].(string)
					}
				}
			}

		}
	}
	//--//Search for tellafriend Reward if Active generate a CouponCode

	if sCouponCode != "" {
		sCouponMessage = fmt.Sprintf(`Register and Subscribe with Coupon Code <b>%s</b> to receive a discount. <br><br> `, sCouponCode)
	}

	//Generate A Refferal Coupon Code

	sFriendMessage := functions.TrimEscape(httpReq.FormValue("friendmessage"))
	sFriendMessage = strings.Replace(sFriendMessage, "\r", "", -1)
	sFriendMessage = strings.Replace(sFriendMessage, "\n", "<br>", -1)

	yourname := fmt.Sprintf("%s %s %s", this.mapAppCache["title"], this.mapAppCache["firstname"], this.mapAppCache["lastname"])

	//generateTellAFriendEmail
	emailFrom := "rewards@valued.com"
	emailFromName := "VALUED.COM Rewards"
	emailTo := functions.TrimEscape(httpReq.FormValue("friendemail"))
	emailSubject := fmt.Sprintf("%s thinks you will appreciate this", functions.CamelCase(yourname))
	emailTemplate := "app-tell"

	emailFields := make(map[string]interface{})
	emailFields["email"] = functions.TrimEscape(httpReq.FormValue("friendemail"))
	emailFields["sendername"] = functions.CamelCase(yourname)
	emailFields["receiverfirstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("friendfirstname")))
	emailFields["message"] = sFriendMessage
	emailFields["couponmessage"] = sCouponMessage

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, nil, emailFields)
	//generateTellAFriendEmail

	//saveReferral
	xDocReferral := make(map[string]interface{})

	xDocReferral["code"] = sCouponCode

	if this.mapAppCache["control"] != nil {
		xDocReferral["title"] = this.mapAppCache["title"]
		xDocReferral["firstname"] = this.mapAppCache["firstname"]
		xDocReferral["lastname"] = this.mapAppCache["lastname"]
		xDocReferral["email"] = this.mapAppCache["username"]
		xDocReferral["profilecontrol"] = this.mapAppCache["control"]
	}

	xDocReferral["friendemail"] = functions.TrimEscape(httpReq.FormValue("friendemail"))
	xDocReferral["friendlastname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("friendlastname")))
	xDocReferral["friendfirstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("friendfirstname")))
	xDocReferral["description"] = sFriendMessage
	xDocReferral["workflow"] = "referred"
	new(database.Referral).Create(this.mapAppCache["username"].(string), xDocReferral, curdb)
	//saveReferral

	sMessage = `THANK YOU - We have forwarded your message to %s %s`
	sMessage = fmt.Sprintf(sMessage, strings.ToUpper(functions.TrimEscape(httpReq.FormValue("friendfirstname"))), functions.TrimEscape(httpReq.FormValue("friendemail")))
	httpRes.Write([]byte(`{"alertSuccess":"` + sMessage + `","getform":"/app-tell"}`))
}
