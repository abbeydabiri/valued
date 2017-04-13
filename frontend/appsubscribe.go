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

type AppSubscribe struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppSubscribe) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		mapAppSubscribe := make(map[string]interface{})
		curdb.SetSession(GOSESSID.Value, "mapAppSubscribe", mapAppSubscribe, false)

		this.pageMap = make(map[string]interface{})
		AppSubscribe := make(map[string]interface{})

		//Search Scheme and List with Price
		schemeMapResult, _ := curdb.Query(`select control, title, price from scheme order by price desc`)
		for cNumber, schemeXdoc := range schemeMapResult {
			xDoc := schemeXdoc.(map[string]interface{})

			title := xDoc["title"].(string)
			if strings.Contains(title, "LIFESTYLE") {
				xDoc["state"] = "checked"
			}
			xDoc["schemeprice"] = functions.Round(xDoc["price"].(float64))
			AppSubscribe[cNumber+"#app-subscribe-option"] = xDoc
		}
		//Search Scheme and List with Price

		AppSubscribe["app-footer"] = make(map[string]interface{})
		this.pageMap["app-subscribe"] = AppSubscribe

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Sign Up For Scheme","pageContent":` + contentHTML + `}`))
		return

	case "subscribe":
		this.Subscribe(httpRes, httpReq, curdb)
		return

	case "paynow":
		this.Paynow(httpRes, httpReq, curdb)
		return

	case "telr":
		this.Verify(httpRes, httpReq, curdb)
		return
	}
}

func (this *AppSubscribe) Subscribe(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	mapAppSubscribe := make(map[string]interface{})

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	mapAppSubscribe = curdb.GetSession(GOSESSID.Value, "mapAppSubscribe")

	if mapAppSubscribe["schemecontrol"] == nil {

		if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
			sMessage += "Please Select a Scheme<br>"
		}

		if len(sMessage) > 0 {
			httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
			return
		}

		sqlScheme := fmt.Sprintf(`select control as schemecontrol, title as schemetitle, price as schemeprice, price as totalprice from scheme where control = '%s'`,
			functions.TrimEscape(httpReq.FormValue("scheme")))
		schemeMapResult, _ := curdb.Query(sqlScheme)

		if schemeMapResult["1"] == nil {
			sMessage += "Please Select a Valid Scheme<br>"
		}

		if len(sMessage) > 0 {
			httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
			return
		}

		mapAppSubscribe = schemeMapResult["1"].(map[string]interface{})
		curdb.SetSession(GOSESSID.Value, "mapAppSubscribe", mapAppSubscribe, false)
	}

	// Check if app-subscribe-complete is needed
	AppSubscribeComplete := make(map[string]interface{})
	sqlProfile := `select image, title, firstname, lastname, phone, email, dob, nationality from profile where control = '%s'`
	sqlProfile = fmt.Sprintf(sqlProfile, this.mapAppCache["control"])
	xDocresult, _ := curdb.Query(sqlProfile)
	if xDocresult["1"] != nil {

		formProfile := xDocresult["1"].(map[string]interface{})
		if formProfile["title"] == nil || formProfile["title"].(string) == "" {
			AppSubscribeCompleteTitle := make(map[string]interface{})
			AppSubscribeCompleteTitle["firstname"] = formProfile["firstname"]
			AppSubscribeCompleteTitle["lastname"] = formProfile["lastname"]
			AppSubscribeComplete["app-subscribe-complete-title"] = AppSubscribeCompleteTitle
		}

		if (formProfile["phonecode"] == nil || formProfile["phonecode"].(string) == "") || (formProfile["phone"] == nil || formProfile["phone"].(string) == "") {
			AppSubscribeComplete["app-subscribe-complete-mobile"] = make(map[string]interface{})
		}

		if formProfile["dob"] == nil || formProfile["dob"].(string) == "" {
			AppSubscribeComplete["app-subscribe-complete-dob"] = make(map[string]interface{})
		}

		if formProfile["nationality"] == nil || formProfile["nationality"].(string) == "" {
			AppSubscribeComplete["app-subscribe-complete-nationality"] = make(map[string]interface{})
		}

	}

	this.pageMap = make(map[string]interface{})
	AppSubscribe := make(map[string]interface{})

	if len(AppSubscribeComplete) > 0 {
		AppSubscribe["app-subscribe-complete"] = AppSubscribeComplete
	}
	// Check if app-subscribe-complete is needed

	if mapAppSubscribe["couponcontrol"] == nil {
		AppSubscribe["app-subscribe-coupon"] = make(map[string]interface{})
	}

	for key, value := range mapAppSubscribe {
		AppSubscribe[key] = value
	}

	//

	//Calculate Total Price and Apply Discount if Needed
	AppSubscribe["totalpriceString"] = functions.Round(AppSubscribe["totalprice"].(float64))
	if mapAppSubscribe["coupondiscountvalue"] != nil {
		floatSchemeprice := mapAppSubscribe["schemeprice"].(float64)
		floatDiscountValue := mapAppSubscribe["coupondiscountvalue"].(float64)

		switch mapAppSubscribe["coupondiscounttype"].(string) {
		case "%":
			AppSubscribe["totalprice"] = functions.RoundUp(((100.00-floatDiscountValue)/100.00)*floatSchemeprice, 2)
			AppSubscribe["totalpriceString"] = functions.Round(AppSubscribe["totalprice"].(float64))

		case "Off":
			if floatDiscountValue >= floatSchemeprice {
				AppSubscribe["totalprice"] = functions.RoundUp(((100.00-floatDiscountValue)/100.00)*floatSchemeprice, 2)
				AppSubscribe["totalpriceString"] = functions.Round(AppSubscribe["totalprice"].(float64))
			} else {
				AppSubscribe["totalprice"] = functions.RoundUp(floatSchemeprice-floatDiscountValue, 2)
				AppSubscribe["totalpriceString"] = functions.Round(AppSubscribe["totalprice"].(float64))
			}
		}

		mapAppSubscribe["totalprice"] = AppSubscribe["totalprice"]
		curdb.SetSession(GOSESSID.Value, "mapAppSubscribe", mapAppSubscribe, false)
	}
	//Calculate Total Price and Apply Discount if Needed

	AppSubscribe["app-footer"] = make(map[string]interface{})
	this.pageMap["app-subscribe-paynow"] = AppSubscribe

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Pay Now","pageContent":` + contentHTML + `}`))
}

func (this *AppSubscribe) Paynow(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	lUpdateProfile := false
	mapProfile := make(map[string]interface{})

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	mapAppSubscribe := curdb.GetSession(GOSESSID.Value, "mapAppSubscribe")

	/////
	sqlProfile := `select title, firstname, lastname, phone, email, dob, nationality, control from profile where control = '%s'`
	sqlProfile = fmt.Sprintf(sqlProfile, this.mapAppCache["control"])
	xDocresult, _ := curdb.Query(sqlProfile)
	if xDocresult["1"] != nil {

		mapProfile = xDocresult["1"].(map[string]interface{})
		if mapProfile["title"] == nil || functions.TrimEscape(mapProfile["title"].(string)) == "" {
			if functions.TrimEscape(httpReq.FormValue("title")) == "" {
				sMessage += "Please select Title<br>"
			} else {
				lUpdateProfile = true
				mapProfile["title"] = functions.TrimEscape(httpReq.FormValue("title"))
			}
		}

		if mapProfile["phone"] == nil || functions.TrimEscape(mapProfile["phone"].(string)) == "" {
			if functions.TrimEscape(httpReq.FormValue("phone")) == "" {
				sMessage += "Please enter Mobile Number<br>"
			} else {
				lUpdateProfile = true
				mapProfile["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))

				if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
					mapProfile["phonecode"] = "+971"
				} else {
					mapProfile["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
				}
			}
		}

		if mapProfile["dob"] == nil || functions.TrimEscape(mapProfile["dob"].(string)) == "" {
			if functions.TrimEscape(httpReq.FormValue("dob")) == "" {
				sMessage += "Please select Date of Birth<br>"
			} else {
				yearsDifference := functions.GetDifferenceInYears("", functions.TrimEscape(httpReq.FormValue("dob")))
				if yearsDifference < 18 {
					sMessage += "You Must be at least 18yrs old to Subscribe<br>"
				} else {
					lUpdateProfile = true
					mapProfile["dob"] = functions.TrimEscape(httpReq.FormValue("dob"))
				}
			}
		}

		if mapProfile["nationality"] == nil || functions.TrimEscape(mapProfile["nationality"].(string)) == "" {
			if functions.TrimEscape(httpReq.FormValue("nationality")) == "" {
				sMessage += "Please select Nationality<br>"
			} else {
				lUpdateProfile = true
				mapProfile["nationality"] = functions.TrimEscape(httpReq.FormValue("nationality"))
			}
		}
	} else {
		sMessage += "Your User Profile is Invalid, Kindly Log-out and Login to continue<br>"
	}
	/////

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//if profile completion info is provided = complete profile
	if lUpdateProfile {
		new(database.Profile).Update(this.mapAppCache["username"].(string), mapProfile, curdb)
		sMessage += "Profile Updated Sucessfully<br>"
	}
	//if profile completion info is provided = complete profile

	//if coupon code info is provided = validate coupon code and add to mapAppSubscribe
	if functions.TrimEscape(httpReq.FormValue("coupon")) != "" {

		sqlCoupon := `select coupon.control as couponcontrol, reward.control as rewardcontrol, coupon.workflow as couponworkflow, 
						reward.workflow as rewardworkflow, merchant.code as merchantcode, reward.discount as coupondiscount, 
						reward.discountvalue as coupondiscountvalue, reward.discounttype as coupondiscounttype
						from coupon, reward, profile as merchant 
						where coupon.rewardcontrol = reward.control AND reward.merchantcontrol  = merchant.control AND coupon.code = '%s'
					`
		sqlCoupon = fmt.Sprintf(sqlCoupon, functions.TrimEscape(httpReq.FormValue("coupon")))
		couponMapResult, _ := curdb.Query(sqlCoupon)

		if couponMapResult["1"] == nil {
			sMessage += fmt.Sprintf("Coupon <b>%s</b> is Invalid<br>", functions.TrimEscape(httpReq.FormValue("coupon")))
		} else {
			couponXdoc := couponMapResult["1"].(map[string]interface{})
			if couponXdoc["merchantcode"].(string) == "main" {
				switch couponXdoc["rewardworkflow"].(string) {
				case "active":
					switch couponXdoc["couponworkflow"].(string) {
					case "active":
						switch couponXdoc["coupondiscounttype"].(string) {
						case "%", "Off":
							sMessage += "Coupon Code Applied Sucessfully<br>"

							for key, value := range couponXdoc {
								mapAppSubscribe[key] = value
							}
							curdb.SetSession(GOSESSID.Value, "mapAppSubscribe", mapAppSubscribe, false)

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
		httpRes.Write([]byte(`{"error":"` + sMessage + `", "getform":"/app-subscribe?action=subscribe"}`))
		return
	}

	mapAppSubscribe["title"] = mapProfile["title"]
	mapAppSubscribe["firstname"] = mapProfile["firstname"]
	mapAppSubscribe["lastname"] = mapProfile["lastname"]
	mapAppSubscribe["profileemail"] = mapProfile["email"]
	mapAppSubscribe["countrycode"] = ""
	curdb.SetSession(GOSESSID.Value, "mapAppSubscribe", mapAppSubscribe, false)

	sUrl := new(PaymentGatewayTELR).CreateOrder("app-subscribe", httpReq, curdb)
	if sUrl != "" {
		sMessage = "Please wait while we prepare your order!"
		httpRes.Write([]byte(`{"sticky":"` + sMessage + `","redirect":"` + sUrl + `", "startLoading":"startLoading" }`))
		return

	}

	sMessage = "Payment Connectivity Error<br> Please Try Again or Email Support"
	httpRes.Write([]byte(`{"error":"` + sMessage + sUrl + `"}`))
	return
}

func (this *AppSubscribe) Verify(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	xDocTelrOrder := new(PaymentGatewayTELR).VerifyPayment(httpReq.FormValue("telr"), httpReq, curdb)

	if xDocTelrOrder["error"] != nil {
		sMessage = xDocTelrOrder["error"].(string)
		httpRes.Write([]byte(`{"sticky":"` + sMessage + `",` + new(AppProfile).View(httpRes, httpReq, curdb) + `}`))
		return
	}

	if xDocTelrOrder["telrtext"] != nil {
		// sMessage = fmt.Sprintf("Your Subscription has been <b>%s</b>", xDocTelrOrder["telrtext"])
		sMessage = "Welcome Valued member, your payment was successful. Start saving today. <br>"

		switch xDocTelrOrder["telrtext"].(string) {
		case "Paid", "Authorised":
			sMessage += this.subscribeMe(xDocTelrOrder, curdb)
		default:
			sMessage = "Your payment was not successful"
		}
	} else {
		sMessage = "Your Subscription <b>Failed</b>"
	}

	httpRes.Write([]byte(`{"sticky":"` + sMessage + `",` + new(AppProfile).View(httpRes, httpReq, curdb) + `}`))
}

func (this *AppSubscribe) subscribeMe(xDocTelrOrder map[string]interface{}, curdb database.Database) string {

	sMessage := ""
	xDoc := make(map[string]interface{})

	sqlEmployer := fmt.Sprintf(`select employercontrol from profile where control = '%s'`, this.mapAppCache["control"])
	mapEmployer, _ := curdb.Query(sqlEmployer)
	if mapEmployer["1"] == nil {
		sMessage += "But Employer is Missing! <br>"
	} else {
		xDoc["employercontrol"] = mapEmployer["1"].(map[string]interface{})["employercontrol"]
	}

	schemetitle := ""
	sqlScheme := fmt.Sprintf(`select price, title from scheme where control = '%s'`, xDocTelrOrder["schemecontrol"])
	mapScheme, _ := curdb.Query(sqlScheme)
	if mapScheme["1"] == nil {
		sMessage += "But Scheme is Missing! <br>"
	} else {
		xDoc["price"] = mapScheme["1"].(map[string]interface{})["price"]
		schemetitle = mapScheme["1"].(map[string]interface{})["title"].(string)
		xDoc["schemecontrol"] = xDocTelrOrder["schemecontrol"]
	}

	if sMessage == "" {

		todayDate := time.Now()
		oneYear := todayDate.Add(time.Hour * 24 * 365)

		xDoc["workflow"] = "active"
		xDoc["membercontrol"] = this.mapAppCache["control"]

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
				xDoc["workflow"] = "paid"
			}
		}
		//Get Current Subscription Expiry

		tblSubscription := new(database.Subscription)
		xDoc["startdate"] = todayDate.Format("02/01/2006")
		xDoc["expirydate"] = oneYear.Format("02/01/2006")

		if xDocTelrOrder["subscriptioncontrol"] != nil &&
			xDocTelrOrder["subscriptioncontrol"].(string) != "" {
			xDoc["control"] = xDocTelrOrder["subscriptioncontrol"]
			tblSubscription.Update(this.mapAppCache["username"].(string), xDoc, curdb)
		} else {
			tblSubscription.Create(this.mapAppCache["username"].(string), xDoc, curdb)
		}

		xDocProfile := make(map[string]interface{})
		xDocProfile["control"] = this.mapAppCache["control"]
		xDocProfile["workflow"] = "subscribed"
		new(database.Profile).Update(this.mapAppCache["username"].(string), xDocProfile, curdb)

		//generateSubscriptionMail
		emailFrom := "membership@valued.com"
		emailFromName := "VALUED Membership"
		emailTo := this.mapAppCache["email"].(string)
		emailSubject := "WELCOME TO VALUED"
		emailTemplate := "app-subscribe"

		emailFields := make(map[string]interface{})
		emailFields["firstname"] = this.mapAppCache["firstname"]
		emailFields["lastname"] = this.mapAppCache["lastname"]
		emailFields["scheme"] = schemetitle

		emailFields["expirydate"] = xDoc["expirydate"]
		emailFields["username"] = this.mapAppCache["username"]

		go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, "", emailFields)
		//generateSubscriptionMail

	} else {
		sMessage += "Please contact support for Manual Subscription"
	}

	return sMessage
}
