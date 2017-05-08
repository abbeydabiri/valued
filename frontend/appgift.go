package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AppGift struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppGift) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		mapAppGift := make(map[string]interface{})
		curdb.SetSession(GOSESSID.Value, "mapAppGift", mapAppGift, false)

		this.pageMap = make(map[string]interface{})
		AppGift := make(map[string]interface{})

		//Search Scheme and List with Price
		schemeMapResult, _ := curdb.Query(`select control, title, price from scheme  where code in ('lite','lifestyle') order by price desc`)
		for cNumber, schemeXdoc := range schemeMapResult {
			xDoc := schemeXdoc.(map[string]interface{})

			title := xDoc["title"].(string)
			if strings.Contains(title, "LIFESTYLE") {
				xDoc["state"] = "checked"
			}
			xDoc["schemeprice"] = functions.Round(xDoc["price"].(float64))
			AppGift[cNumber+"#app-subscribe-option"] = xDoc
		}
		//Search Scheme and List with Price

		if this.mapAppCache["control"] != nil {
			AppGift["sendersname"] = fmt.Sprintf("%s %s %s", this.mapAppCache["title"],
				this.mapAppCache["firstname"], this.mapAppCache["lastname"])
			AppGift["sendersemail"] = this.mapAppCache["username"]
		}

		appFooter := make(map[string]interface{})
		AppGift["app-footer"] = appFooter
		appFooter["home"] = "white"

		this.pageMap["app-gift"] = AppGift

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Gift Membership","pageContent":` + contentHTML + `}`))
		return

	case "gift":
		this.gift(httpRes, httpReq, curdb)
		return

	case "paynow":
		this.paynow(httpRes, httpReq, curdb)
		return

	case "telr":
		this.verify(httpRes, httpReq, curdb)
		return
	}
}

func (this *AppGift) gift(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	mapAppGift := make(map[string]interface{})

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	mapAppGift = curdb.GetSession(GOSESSID.Value, "mapAppGift")

	if mapAppGift["schemecontrol"] == nil {
		if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
			sMessage += "Please Select a Scheme<br>"
		}

		if len(sMessage) > 0 {
			httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
			return
		}

		sqlScheme := fmt.Sprintf(`select control as schemecontrol, title as schemetitle, price as schemeprice, price as totalprice from scheme where control = '%s' and workflow = 'active'`,
			functions.TrimEscape(httpReq.FormValue("scheme")))
		schemeMapResult, _ := curdb.Query(sqlScheme)

		if schemeMapResult["1"] == nil {
			sMessage += "Please Select a Valid Scheme<br>"
		}

		if functions.TrimEscape(httpReq.FormValue("sendersname")) == "" {
			sMessage += "Senders Name is Missing<br>"
		}

		if functions.TrimEscape(httpReq.FormValue("sendersemail")) == "" {
			sMessage += "Senders Email is Missing<br>"
		}

		if functions.TrimEscape(httpReq.FormValue("friendfirstname")) == "" {
			sMessage += "Friends Name is Missing<br>"
		}

		if functions.TrimEscape(httpReq.FormValue("friendlastname")) == "" {
			sMessage += "Friends Name is Missing<br>"
		}

		if functions.TrimEscape(httpReq.FormValue("friendemail")) == "" {
			sMessage += "Friends Email is Missing<br>"
		}

		sFriendMessage := "You're VALUED - enjoy the rewards this brings!"
		if functions.TrimEscape(httpReq.FormValue("message")) != "" {
			sFriendMessage = functions.TrimEscape(httpReq.FormValue("message"))
		}

		if functions.TrimEscape(httpReq.FormValue("above18")) == "" {
			sMessage += "Please Confirm You are above 18 Years<br>"
		}

		if functions.TrimEscape(httpReq.FormValue("terms")) == "" {
			sMessage += "Please Accept Terms and Conditions <br>"
		}

		if len(sMessage) > 0 {
			httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
			return
		}

		mapAppGift = schemeMapResult["1"].(map[string]interface{})
		mapAppGift["sendersname"] = functions.TrimEscape(httpReq.FormValue("sendersname"))
		mapAppGift["sendersemail"] = functions.TrimEscape(httpReq.FormValue("sendersemail"))
		mapAppGift["friendfirstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("friendfirstname")))
		mapAppGift["friendlastname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("friendlastname")))
		mapAppGift["friendemail"] = functions.TrimEscape(httpReq.FormValue("friendemail"))

		if this.mapAppCache["control"] != nil {
			mapAppGift["logincontrol"] = this.mapAppCache["control"]
			mapAppGift["profilecontrol"] = this.mapAppCache["control"]
		} else {
			mapAppGift["logincontrol"] = ""
			mapAppGift["profilecontrol"] = ""
		}

		sFriendMessage = strings.Replace(sFriendMessage, "\r", "", -1)
		sFriendMessage = strings.Replace(sFriendMessage, "\n", "<br>", -1)
		mapAppGift["message"] = sFriendMessage

		curdb.SetSession(GOSESSID.Value, "mapAppGift", mapAppGift, false)
	}

	this.pageMap = make(map[string]interface{})
	AppGift := make(map[string]interface{})

	if mapAppGift["couponcontrol"] == nil {
		AppGift["app-subscribe-coupon"] = make(map[string]interface{})
	}

	for key, value := range mapAppGift {
		AppGift[key] = value
	}

	//

	//Calculate Total Price and Apply Discount if Needed
	AppGift["totalpriceString"] = functions.Round(AppGift["totalprice"].(float64))
	if mapAppGift["coupondiscountvalue"] != nil {
		floatSchemeprice := mapAppGift["schemeprice"].(float64)
		floatDiscountValue := mapAppGift["coupondiscountvalue"].(float64)

		switch mapAppGift["coupondiscounttype"].(string) {
		case "%":
			AppGift["totalprice"] = functions.RoundUp(((100.00-floatDiscountValue)/100.00)*floatSchemeprice, 2)
			AppGift["totalpriceString"] = functions.Round(AppGift["totalprice"].(float64))

		case "Off":
			if floatDiscountValue >= floatSchemeprice {
				AppGift["totalprice"] = functions.RoundUp(((100.00-floatDiscountValue)/100.00)*floatSchemeprice, 2)
				AppGift["totalpriceString"] = functions.Round(AppGift["totalprice"].(float64))
			} else {
				AppGift["totalprice"] = functions.RoundUp(floatSchemeprice-floatDiscountValue, 2)
				AppGift["totalpriceString"] = functions.Round(AppGift["totalprice"].(float64))
			}
		}

		mapAppGift["totalprice"] = AppGift["totalprice"]
		curdb.SetSession(GOSESSID.Value, "mapAppGift", mapAppGift, false)
	}
	//Calculate Total Price and Apply Discount if Needed

	AppGift["app-footer"] = make(map[string]interface{})
	this.pageMap["app-gift-paynow"] = AppGift

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Pay Now","pageContent":` + contentHTML + `}`))
}

func (this *AppGift) paynow(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	mapAppGift := curdb.GetSession(GOSESSID.Value, "mapAppGift")

	//if coupon code info is provided = validate coupon code and add to mapAppGift
	couponXdoc := make(map[string]interface{})
	if functions.TrimEscape(httpReq.FormValue("coupon")) != "" {

		sqlCoupon := `select coupon.control as couponcontrol, reward.control as rewardcontrol, coupon.workflow as couponworkflow, 
						reward.workflow as rewardworkflow, reward.method as rewardmethod, merchant.code as merchantcode, reward.discount as coupondiscount, 
						reward.discountvalue as coupondiscountvalue, reward.discounttype as coupondiscounttype
						from coupon, reward, profile as merchant where 
						coupon.rewardcontrol = reward.control AND reward.merchantcontrol = merchant.control AND coupon.code = '%s'`

		sqlCoupon = fmt.Sprintf(sqlCoupon, functions.TrimEscape(httpReq.FormValue("coupon")))
		couponMapResult, _ := curdb.Query(sqlCoupon)

		if couponMapResult["1"] == nil {
			sMessage += fmt.Sprintf("Coupon <b>%s</b> is Invalid<br>", functions.TrimEscape(httpReq.FormValue("coupon")))
		} else {
			couponXdoc = couponMapResult["1"].(map[string]interface{})
			if couponXdoc["merchantcode"].(string) == "none" {
				switch couponXdoc["rewardworkflow"].(string) {
				case "active":
					switch couponXdoc["couponworkflow"].(string) {
					case "active":
						switch couponXdoc["coupondiscounttype"].(string) {
						case "%", "Off":
							sMessage += "Coupon Code Applied Sucessfully<br>"

							for key, value := range couponXdoc {
								mapAppGift[key] = value
							}
							curdb.SetSession(GOSESSID.Value, "mapAppGift", mapAppGift, false)

						default:
							sMessage += "Coupon Code Discount is Invalid<br>"
						}
						break

					default:
						sMessage += "Coupon Code is Invalid<br>"
						break
					}
					break
				default:
					sMessage += "Coupon Code is Invalid (Reward no longer active)<br>"
					break
				}
			} else {
				sMessage += "Coupon Code is Invalid<br>"
			}

		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `", "getform":"/app-gift?action=gift"}`))
		return
	}

	//
	//Mark Coupon as Used Here ->
	if mapAppGift["couponcontrol"] != nil && mapAppGift["couponcontrol"].(string) != "" &&
		mapAppGift["rewardmethod"].(string) != "Client Code Single" {
		xDocCoupon := make(map[string]interface{})
		xDocCoupon["workflow"] = "approved"
		xDocCoupon["control"] = mapAppGift["couponcontrol"]
		new(database.Coupon).Update(this.mapAppCache["username"].(string), xDocCoupon, curdb)
	}
	//Mark Coupon as Used Here ->
	//

	//If Coupon Price == 0
	if mapAppGift["totalprice"].(float64) < 1.0 {
		sMessage = this.subscribeFriend(mapAppGift, httpReq, curdb)
		if sMessage == "" {
			sMessage = "Hello Valued member, Your payment was successful. Start saving today. <br>"
		}
		httpRes.Write([]byte(`{"sticky":"` + sMessage + `",` + new(AppProfile).View(httpRes, httpReq, curdb) + `}`))
		return
	}
	//If Coupon Price == 0

	sUrl := new(PaymentGatewayTELR).CreateOrder("app-gift", httpReq, curdb)
	if sUrl != "" {
		sMessage = "Please wait while we prepare your order!"
		httpRes.Write([]byte(`{"sticky":"` + sMessage + `","redirect":"` + sUrl + `", "startLoading":"startLoading" }`))
		return

	}

	sMessage = "Payment Connectivity Error<br> Please Try Again or Email Support"
	httpRes.Write([]byte(`{"error":"` + sMessage + sUrl + `"}`))
	return
}

func (this *AppGift) verify(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	xDocTelrOrder := new(PaymentGatewayTELR).VerifyPayment(httpReq.FormValue("telr"), httpReq, curdb)

	if xDocTelrOrder["error"] != nil {
		sMessage = xDocTelrOrder["error"].(string)

		httpRes.Write([]byte(`{"sticky":"` + sMessage + `","getform":"/app-home"}`))
		//httpRes.Write([]byte(`{"sticky":"` + sMessage + `",` + new(AppProfile).View(httpRes, httpReq, curdb) + `}`))
		return
	}

	if xDocTelrOrder["telrtext"] != nil {
		sMessage = fmt.Sprintf("Your Subscription has been <b>%s</b> <br>", xDocTelrOrder["telrtext"])

		switch xDocTelrOrder["telrtext"].(string) {
		case "Paid", "Authorised":
			sMessage += this.subscribeFriend(xDocTelrOrder, httpReq, curdb)

		default:
			sMessage = "Your payment was not successful"
		}
	} else {
		sMessage = "Your Subscription <b>Failed</b>"
	}

	httpRes.Write([]byte(`{"sticky":"` + sMessage + `","getform":"/app-home"}`))
	// httpRes.Write([]byte(`{"sticky":"` + sMessage + `",` + new(AppProfile).View(httpRes, httpReq, curdb) + `}`))
}

func (this *AppGift) subscribeFriend(xDocTelrOrder map[string]interface{}, httpReq *http.Request, curdb database.Database) string {

	sMessage := ""

	profilecontrol := ""
	employercontrol := ""

	//Get Default Employer Control
	defaultMap, _ := curdb.Query(`select control from profile where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Employer <br>"
	} else {
		employercontrol = defaultMap["1"].(map[string]interface{})["control"].(string)
	}

	if len(sMessage) > 0 {
		return sMessage
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	mapAppGift := curdb.GetSession(GOSESSID.Value, "mapAppGift")
	// mapAppGift["schemecontrol"]

	//--> Check if Email/Username exists in Login Table
	sqlCheckEmail := fmt.Sprintf(`select control as profilecontrol, employercontrol, username from profile where email = '%s'`, mapAppGift["friendemail"])
	mapCheckEmail, _ := curdb.Query(sqlCheckEmail)
	//--> Check if Email/Username exists in Login Table

	sUsername := mapAppGift["sendersemail"].(string)
	if this.mapAppCache["username"] != nil {
		sUsername = this.mapAppCache["username"].(string)
	}

	todayDate := time.Now()
	oneYear := todayDate.Add(time.Hour * 24 * 365)

	//Get Current Subscription Expiry
	sqlSubscription := `select sub.schemecontrol as schemecontrol, sub.code as code, sub.control as control, sub.expirydate as expirydate, sch.title as schemetitle
						from subscription as sub join scheme as sch on sub.schemecontrol = sch.control
						AND sub.workflow = 'active' AND sch.workflow = 'active' AND sch.code in ('lite','lifestyle')  
						AND '%s'::timestamp between sub.startdate::timestamp and sub.expirydate::timestamp AND sub.schemecontrol = '%s' AND sub.membercontrol = '%s' order by control desc`
	sqlSubscription = fmt.Sprintf(sqlSubscription, todayDate.Format("02/01/2006"), xDocTelrOrder["schemecontrol"], this.mapAppCache["control"])

	curdb.Query("set datestyle = dmy")
	xDocCurSubRes, _ := curdb.Query(sqlSubscription)
	if xDocCurSubRes["1"] != nil {
		xDocCurSub := xDocCurSubRes["1"].(map[string]interface{})
		if xDocCurSub["expirydate"] != nil && xDocCurSub["expirydate"].(string) != "" {
			todayDate, _ = time.Parse("02/01/2006", xDocCurSub["expirydate"].(string))
			oneYear = todayDate.Add(time.Hour * 24 * 365)
		}
	}
	//Get Current Subscription Expiry

	xDocSubscribe := make(map[string]interface{})
	xDocSubscribe["startdate"] = todayDate.Format("02/01/2006")
	xDocSubscribe["expirydate"] = oneYear.Format("02/01/2006")

	sReceiverUsername := mapAppGift["friendemail"].(string)
	if mapCheckEmail["1"] != nil {
		xDocSubscribe["workflow"] = "pending"
		profilecontrol = mapCheckEmail["1"].(map[string]interface{})["profilecontrol"].(string)
		employercontrol = mapCheckEmail["1"].(map[string]interface{})["employercontrol"].(string)
		sReceiverUsername = mapCheckEmail["1"].(map[string]interface{})["username"].(string)
	} else {
		xDocSubscribe["workflow"] = "active"

		//Create Friend Profile
		///////////////////////////////////////////////////////////
		xDocProfile := make(map[string]interface{})
		xDocProfile["status"] = "inactive"
		xDocProfile["workflow"] = "subscribed"

		xDocProfile["firstname"] = mapAppGift["friendfirstname"]
		xDocProfile["lastname"] = mapAppGift["friendlastname"]
		xDocProfile["email"] = mapAppGift["friendemail"]
		xDocProfile["username"] = mapAppGift["friendemail"]
		xDocProfile["password"] = functions.RandomString(8)
		xDocProfile["employercontrol"] = employercontrol
		xDocProfile["company"] = ""

		profilecontrol = new(database.Profile).Create(sUsername, xDocProfile, curdb)
		//Create Friend Profile

		//Create a Member Role Link
		sqlRoleClear := fmt.Sprintf(`delete from profilerole where profilecontrol = '%s'`, xDocProfile["control"])
		curdb.Query(sqlRoleClear)
		mapRole, _ := curdb.Query(`select control from role where code = 'member'`)
		if mapRole["1"] != nil {
			xDocMapRole := mapRole["1"].(map[string]interface{})
			xDocRole := make(map[string]interface{})
			xDocRole["rolecontrol"] = xDocMapRole["control"]
			xDocRole["profilecontrol"] = xDocProfile["control"]
			new(database.ProfileRole).Create("signup", xDocRole, curdb)
		}
		//Create a Member Role Link

		//saveReferral
		xDocReferral := make(map[string]interface{})
		if this.mapAppCache["control"] != nil {
			xDocReferral["title"] = functions.CamelCase(this.mapAppCache["title"].(string))
			xDocReferral["firstname"] = functions.CamelCase(this.mapAppCache["firstname"].(string))
			xDocReferral["lastname"] = functions.CamelCase(this.mapAppCache["lastname"].(string))
			xDocReferral["email"] = this.mapAppCache["username"]
			xDocReferral["profilecontrol"] = this.mapAppCache["control"]
		} else {
			xDocReferral["email"] = mapAppGift["sendersemail"]
			xDocReferral["firstname"] = functions.CamelCase(mapAppGift["sendersname"].(string))
		}

		xDocReferral["friendemail"] = mapAppGift["friendemail"]
		xDocReferral["friendlastname"] = functions.CamelCase(mapAppGift["friendlastname"].(string))
		xDocReferral["friendfirstname"] = functions.CamelCase(mapAppGift["friendfirstname"].(string))
		xDocReferral["description"] = mapAppGift["message"]
		xDocReferral["workflow"] = "subscribed"
		new(database.Referral).Create(sUsername, xDocReferral, curdb)
		//saveReferral
		//

		//generateActivationLink
		tblActivationLink := new(database.ActivationLink)
		xDocLink := make(map[string]interface{})
		xDocLink["code"] = functions.RandomString(32)
		xDocLink["title"] = fmt.Sprintf("User [%s] Account Activation", xDocProfile["username"])
		xDocLink["workflow"] = "draft"
		xDocLink["logincontrol"] = profilecontrol
		tblActivationLink.Create(sUsername, xDocLink, curdb)

		activationLink := fmt.Sprintf("%sapp-activate/?code=%s", httpReq.Referer(), xDocLink["code"])
		//generateActivationLink

		//Send User Registration Here
		//generateRegistrationMail
		emailFrom := "membership@valued.com"
		emailFromName := "VALUED membership"
		emailTo := xDocProfile["username"].(string)
		emailSubject := "Welcome to VALUED - Member Registration"
		emailTemplate := "app-registration"

		emailFields := make(map[string]interface{})
		emailFields["fullname"] = fmt.Sprintf(`%v %v`, functions.CamelCase(xDocProfile["firstname"].(string)),
			functions.CamelCase(xDocProfile["lastname"].(string)))
		emailFields["email"] = xDocProfile["email"]
		emailFields["username"] = xDocProfile["username"]
		emailFields["link"] = activationLink

		emailFields["details"] = fmt.Sprintf(` Login Details: <br> ============== <br> Username: %v <br> Password: %v <br><br>`,
			xDocProfile["username"], xDocProfile["password"])

		go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, []string, emailFields)
		//generateRegistrationMail
		//Send User Registration Here

	}

	schemetitle := ""
	sqlScheme := fmt.Sprintf(`select price, title from scheme where control = '%v'`, mapAppGift["schemecontrol"])
	mapScheme, _ := curdb.Query(sqlScheme)
	if mapScheme["1"] == nil {
		sMessage += "But Scheme is Missing! <br>"
	} else {
		xDocSubscribe["price"] = mapScheme["1"].(map[string]interface{})["price"]
		schemetitle = mapScheme["1"].(map[string]interface{})["title"].(string)
		xDocSubscribe["schemecontrol"] = mapAppGift["schemecontrol"]
	}

	sTitle := "GFT" + functions.RandomString(6)
	xDocSubscribe["membercontrol"] = profilecontrol
	xDocSubscribe["code"] = sTitle
	xDocSubscribe["title"] = sTitle

	new(database.Subscription).Create(sUsername, xDocSubscribe, curdb)

	//generateSubscriptionMail
	emailFrom := "membership@valued.com"
	emailFromName := "VALUED Membership"
	emailTo := mapAppGift["friendemail"].(string)
	emailSubject := "WELCOME TO VALUED"
	emailTemplate := "app-subscribe"

	emailFields := make(map[string]interface{})
	emailFields["firstname"] = mapAppGift["friendfirstname"]
	emailFields["lastname"] = mapAppGift["friendlastname"]
	emailFields["scheme"] = schemetitle

	emailFields["expirydate"] = xDocSubscribe["expirydate"]
	emailFields["username"] = sReceiverUsername

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, []string, emailFields)
	//generateSubscriptionMail

	//generateReceiverMail
	emailFrom = "membership@valued.com"
	emailFromName = "VALUED Membership"
	emailTo = functions.TrimEscape(httpReq.FormValue("friendemail"))
	emailSubject = "YOUR GIFT - A VALUED MEMBERSHIP FROM " + mapAppGift["sendersname"].(string)
	emailTemplate = "app-gift-receiver"

	emailFields = make(map[string]interface{})
	emailFields["scheme"] = schemetitle
	emailFields["sendersname"] = mapAppGift["sendersname"]
	emailFields["friendlastname"] = mapAppGift["friendlastname"]
	emailFields["friendfirstname"] = mapAppGift["friendfirstname"]
	emailFields["message"] = mapAppGift["message"]

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, []string, emailFields)
	//generateReceiverMail

	//generateGifterMail
	emailFrom = "membership@valued.com"
	emailFromName = "VALUED Membership"
	emailTo = functions.TrimEscape(httpReq.FormValue("friendemail"))
	emailSubject = `A Gift For You – VALUED Membership`
	emailTemplate = "app-gift-sender"

	emailFields = make(map[string]interface{})
	emailFields["fullname"] = mapAppGift["sendersname"]
	emailFields["scheme"] = schemetitle
	emailFields["schemeprice"] = mapAppGift["totalpriceString"]
	emailFields["firstname"] = mapAppGift["friendfirstname"]
	emailFields["lastname"] = mapAppGift["friendlastname"]
	emailFields["message"] = mapAppGift["message"]

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, []string, emailFields)
	//generateGifterMail

	sMessage += "Gift Membership has been completed"

	return sMessage

}
