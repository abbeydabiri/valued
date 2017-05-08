package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"

	"time"

	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type AppProfile struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppProfile) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	this.pageMap = make(map[string]interface{})

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		httpRes.Write([]byte(`{` + this.View(httpRes, httpReq, curdb) + `}`))
		return

	case "edit":
		httpRes.Write([]byte(`{` + this.edit(httpRes, httpReq, curdb) + `}`))
		return

	case "save":
		this.save(httpRes, httpReq, curdb)
		return

	case "redemption":
		this.redemption(httpRes, httpReq, curdb)
		return

	case "mySavings":
		this.mySavings(httpRes, httpReq, curdb)
		return

	case "savePic":
		this.saveProfilePic(httpRes, httpReq, curdb)
		return

	case "forgotPin":
		this.forgotPin(httpRes, httpReq, curdb)
		return

	case "changePassword":
		this.changePassword(httpRes, httpReq, curdb)
		return

	case "savePassword":
		this.savePassword(httpRes, httpReq, curdb)
		return

	case "setPin":
		this.setPin(httpRes, httpReq, curdb)
		return

	case "savePin":
		this.savePin(httpRes, httpReq, curdb)
		return

	case "complete":
		pageMap := make(map[string]interface{})
		pageMap["app-complete"] = make(map[string]interface{})
		pageHTML := strconv.Quote(string(this.Generate(pageMap, nil)))

		httpRes.Write([]byte(`{"pageTitle":"Complete Profile","pageContent":` + pageHTML + `}`))
		return
	}
}

func (this *AppProfile) View(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) string {

	if this.pageMap == nil {
		this.pageMap = make(map[string]interface{})
	}

	if this.mapAppCache == nil {
		GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")
	}

	todayDate := time.Now()
	curdb.Query("set datestyle = dmy")
	sqlSubscriptionPaid := `select sub.control as control
						from subscription as sub join scheme as sch on sub.schemecontrol = sch.control
						AND sub.workflow = 'paid' AND sch.workflow = 'active' AND sch.code in ('lite','lifestyle')  
						AND '%s'::timestamp between sub.startdate::timestamp and sub.expirydate::timestamp AND sub.membercontrol = '%s' order by control desc`
	sqlSubscriptionPaid = fmt.Sprintf(sqlSubscriptionPaid, todayDate.Format("02/01/2006"), this.mapAppCache["control"])
	xDocresultPaid, _ := curdb.Query(sqlSubscriptionPaid)

	for _, xDoc := range xDocresultPaid {
		xDoc := xDoc.(map[string]interface{})

		todayDate = time.Now()
		xDoc["startdate"] = todayDate.Format("02/01/2006")

		oneYear := todayDate.Add(time.Hour * 24 * 365)
		xDoc["expirydate"] = oneYear.Format("02/01/2006")

		xDoc["workflow"] = "active"
		new(database.Subscription).Update(this.mapAppCache["username"].(string), xDoc, curdb)
	}

	formProfile := make(map[string]interface{})
	mapSubscription := make(map[string]interface{})
	sqlSubscription := `select sub.schemecontrol as schemecontrol, sub.code as code, sub.control as control, sub.expirydate as expirydate, sch.title as schemetitle
						from subscription as sub join scheme as sch on sub.schemecontrol = sch.control
						AND sub.workflow = 'active' AND sch.workflow = 'active' AND sch.code in ('lite','lifestyle')  
						AND '%s'::timestamp between sub.startdate::timestamp and sub.expirydate::timestamp AND sub.membercontrol = '%s' order by control desc`
	sqlSubscription = fmt.Sprintf(sqlSubscription, todayDate.Format("02/01/2006"), this.mapAppCache["control"])

	curdb.Query("set datestyle = dmy")
	xDocresult, _ := curdb.Query(sqlSubscription)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		mapSubscription[xDoc["code"].(string)] = xDoc
		formProfile[cNumber+"#app-profile-subscription"] = xDoc
	}

	if len(mapSubscription) > 0 {
		this.mapAppCache["subscription"] = mapSubscription
	} else {
		xDoc := make(map[string]interface{})
		xDoc["code"] = `<button class="profileEditBtn" onclick="getForm('\/app-subscribe')">BUY NOW</button>`
		formProfile["0#app-profile-subscription"] = xDoc
		delete(this.mapAppCache, "subscription")
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	curdb.SetSession(GOSESSID.Value, "mapAppCache", this.mapAppCache, false)

	// -

	//Check if subscription exists and Pin code == "" or nil "" force Pin Code Creation
	if this.mapAppCache["subscription"] != nil {
		formProfile["0#app-profile-btn-savings"] = make(map[string]interface{})

		if this.mapAppCache["company"] == nil || this.mapAppCache["company"].(string) != "Yes" {
			formProfile["0#app-profile-btn-changepin"] = make(map[string]interface{})
		}

		if (this.mapAppCache["pincode"] == nil || this.mapAppCache["pincode"].(string) == "") &&
			(this.mapAppCache["company"] == nil || this.mapAppCache["company"].(string) != "Yes") {

			formProfile := make(map[string]interface{})

			appFooter := make(map[string]interface{})
			formProfile["app-footer"] = appFooter
			appFooter["profile"] = "white"

			sPageTemplate := "app-profile-createpin"

			this.pageMap[sPageTemplate] = formProfile
			contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
			contentHTML = `"pageTitle":"Create Pin","pageContent":` + contentHTML
			return contentHTML
		}
	}
	//Check if subscription exists and Pin code == "" or nil "" force Pin Code Creation

	// -

	lForceProfileCompletion := false
	sqlProfile := `select image, title, firstname, lastname, phonecode, phone, email, dob, nationality from profile where control = '%s'`
	sqlProfile = fmt.Sprintf(sqlProfile, this.mapAppCache["control"])
	xDocresult, _ = curdb.Query(sqlProfile)
	formProfilexDoc := make(map[string]interface{})
	if xDocresult["1"] != nil {

		formProfilexDoc = xDocresult["1"].(map[string]interface{})
		for key, value := range formProfilexDoc {
			formProfile[key] = value
		}

		if formProfile["phone"] == nil || formProfile["phone"].(string) == "" {
			lForceProfileCompletion = true
		}

		if formProfile["phone"] == nil || formProfile["phone"].(string) == "" {
			formProfile["phone"] = "<small><i>**Mobile Number Required**</i></small>"
			lForceProfileCompletion = true
		}

		if formProfile["dob"] == nil || formProfile["dob"].(string) == "" {
			formProfile["dob"] = "<small><i>**Date of Birth Required**</i></small>"
			lForceProfileCompletion = true
		}

		if formProfile["nationality"] == nil || formProfile["nationality"].(string) == "" {
			formProfile["nationality"] = "<small><i>**Nationality Required**</i></small>"
			lForceProfileCompletion = true
		}
	}

	//lForceProfileCompletion
	if lForceProfileCompletion && this.mapAppCache["subscription"] != nil {

		appFooter := make(map[string]interface{})
		formProfilexDoc["app-footer"] = appFooter
		appFooter["profile"] = "white"

		sPageTemplate := "app-profile-edit"

		if this.mapAppCache["company"].(string) == "Yes" {
			formProfilexDoc["titleclass"] = "hide"
		}

		this.pageMap[sPageTemplate] = formProfilexDoc
		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		contentHTML = `"pageTitle":"Edit My Profile","pageContent":` + contentHTML
		return contentHTML
	}
	//lForceProfileCompletion

	// -

	//Get Total Friends Referred & Friends Subscribed
	// sFriendsSQL := `select count(distinct(friendemail)) as %s from referral where lower(workflow) = lower('%s') AND lower(email) = lower('%s')`

	// nReferred := int64(0)
	// sReferredSQL := fmt.Sprintf(sFriendsSQL, "referred", "referred", this.mapAppCache["username"])
	// sReferredRes, _ := curdb.Query(sReferredSQL)
	// if sReferredRes["1"] != nil {
	// 	xDocReferredRes := sReferredRes["1"].(map[string]interface{})
	// 	if xDocReferredRes["referred"] != nil {
	// 		nReferred = xDocReferredRes["referred"].(int64)
	// 	}
	// }
	// formProfile["friendsreferred"] = nReferred

	formProfile["friendsreferred"] = this.GetRefferalCount(this.mapAppCache["email"].(string), curdb)

	// nSubscribed := int64(0)
	// f := fmt.Sprintf(sFriendsSQL, "subscribed", "subscribed", this.mapAppCache["email"])
	// sSubscribedRes, _ := curdb.Query(sSubscribedSQL)
	// if sSubscribedRes["1"] != nil {
	// 	xDocSubscribedRes := sSubscribedRes["1"].(map[string]interface{})
	// 	if xDocSubscribedRes["subscribed"] != nil {
	// 		nSubscribed = xDocSubscribedRes["subscribed"].(int64)
	// 	}
	// }
	formProfile["friendssubscribed"] = this.GetSubscribedCount(this.mapAppCache["email"].(string), curdb)
	//Get Total Friends Referred & Friends Subscribed

	// println("sReferredSQL:" + sReferredSQL)
	// println("sSubscribedSQL:" + sSubscribedSQL)

	//-

	formProfile["username"] = this.mapAppCache["username"]

	appFooter := make(map[string]interface{})
	formProfile["app-footer"] = appFooter
	appFooter["profile"] = "white"

	this.pageMap["app-profile"] = formProfile
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	// httpRes.Write([]byte(`{"pageTitle":"My Profile","pageContent":` + contentHTML + `}`))

	contentHTML = `"pageTitle":"My Profile","pageContent":` + contentHTML
	return contentHTML
}

func (this *AppProfile) GetRefferalCount(email string, curdb database.Database) (nReferred int64) {

	nReferred = int64(0)
	//Get Total Friends Referred & Friends Subscribed
	sFriendsSQL := `select count(distinct(friendemail)) as referred from referral where lower(workflow) = lower('referred') AND lower(email) = lower('%s')`

	sReferredSQL := fmt.Sprintf(sFriendsSQL, this.mapAppCache["username"])
	sReferredRes, _ := curdb.Query(sReferredSQL)
	if sReferredRes["1"] != nil {
		xDocReferredRes := sReferredRes["1"].(map[string]interface{})
		if xDocReferredRes["referred"] != nil {
			nReferred = xDocReferredRes["referred"].(int64)
		}
	}
	return

}

func (this *AppProfile) GetSubscribedCount(email string, curdb database.Database) (nSubscribed int64) {

	sReferredSQL := `select friendemail, cast(control as float) as control from referral where 
	lower(workflow) = lower('referred') AND lower(email) = lower('%s') order by control`
	sReferredSQL = fmt.Sprintf(sReferredSQL, email)
	emailList := ""
	nSubscribed = int64(0)
	mapRefferedByMe := make(map[string]float64)

	sReferredByMeRes, _ := curdb.Query(sReferredSQL)
	aReferralSorted := functions.SortMap(sReferredByMeRes)
	for _, sNumber := range aReferralSorted {
		xDoc := sReferredByMeRes[sNumber].(map[string]interface{})
		if mapRefferedByMe[xDoc["friendemail"].(string)] > float64(1) {
			continue
		}
		mapRefferedByMe[xDoc["friendemail"].(string)] = xDoc["control"].(float64)
		emailList = fmt.Sprintf(`%s,'%s'`, emailList, xDoc["friendemail"])
	}

	emailListSubscribed := ""
	sSubscribedSQL := `select friendemail, cast(control as float) as control from referral where 
	lower(workflow) = lower('referred') AND friendemail in ('0'%s) and lower(email) != lower('%s') order by control`
	sSubscribedSQL = fmt.Sprintf(sSubscribedSQL, emailList, email)
	sSubscribedRes, _ := curdb.Query(sSubscribedSQL)
	aReferralSorted = functions.SortMap(sSubscribedRes)

	for _, sNumber := range aReferralSorted {
		xDoc := sSubscribedRes[sNumber].(map[string]interface{})
		if mapRefferedByMe[xDoc["friendemail"].(string)] > float64(1) {
			if mapRefferedByMe[xDoc["friendemail"].(string)] > xDoc["control"].(float64) {
				delete(mapRefferedByMe, xDoc["friendemail"].(string))
			}
		}
	}

	for sEmail := range mapRefferedByMe {
		emailListSubscribed = fmt.Sprintf(`%s,'%s'`, emailListSubscribed, sEmail)
	}

	sSubscribedSQL = fmt.Sprintf(`select count(distinct(membercontrol)) as subscribed from subscription where 
		membercontrol in (select control from profile where email in ('0'%s))`, emailListSubscribed)
	sSubscribedMapRes, _ := curdb.Query(sSubscribedSQL)

	if sSubscribedMapRes["1"] != nil {
		xDocReferredRes := sSubscribedMapRes["1"].(map[string]interface{})
		if xDocReferredRes["subscribed"] != nil {
			nSubscribed = xDocReferredRes["subscribed"].(int64)
		}
	}
	return
}

func (this *AppProfile) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) string {

	formProfile := make(map[string]interface{})
	sqlProfile := `select image, title, firstname, lastname, phonecode, phone, email, dob, nationality from profile where control = '%s'`
	sqlProfile = fmt.Sprintf(sqlProfile, this.mapAppCache["control"])
	xDocresult, _ := curdb.Query(sqlProfile)
	if xDocresult["1"] != nil {
		formProfile = xDocresult["1"].(map[string]interface{})
		if formProfile["title"].(string) != "" {
			formProfile["titleState"] = "disabled"
		}

		if formProfile["dob"].(string) != "" {
			formProfile["dobState"] = "disabled"
		}

		if formProfile["nationality"].(string) != "" {
			formProfile["nationalityState"] = "disabled"
		}

		formProfile[strings.Replace(formProfile["phonecode"].(string), "+", "", 1)] = "selected"

	}
	formProfile["username"] = this.mapAppCache["username"]

	appFooter := make(map[string]interface{})
	formProfile["app-footer"] = appFooter
	appFooter["profile"] = "white"

	if this.mapAppCache["company"].(string) == "Yes" {
		formProfile["titleclass"] = "hide"
	}

	this.pageMap["app-profile-edit"] = formProfile

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	// httpRes.Write([]byte(`{"pageTitle":"Edit My Profile","pageContent":` + contentHTML + `}`))
	contentHTML = `"pageTitle":"Edit My Profile","pageContent":` + contentHTML
	return contentHTML
}

func (this *AppProfile) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	xDocProfile := make(map[string]interface{})
	xDocProfile["control"] = this.mapAppCache["control"]

	sqlProfile := `select image, title, firstname, lastname,  phone, email, dob, nationality from profile where control = '%s'`
	sqlProfile = fmt.Sprintf(sqlProfile, this.mapAppCache["control"])
	xDocresult, _ := curdb.Query(sqlProfile)

	sProfileEmail := ""
	if xDocresult["1"] != nil {

		formProfile := xDocresult["1"].(map[string]interface{})
		sProfileEmail = formProfile["email"].(string)

		if formProfile["title"].(string) == "" {
			if functions.TrimEscape(httpReq.FormValue("title")) == "" {
				sMessage += "Title is missing <br>"
			} else {
				xDocProfile["title"] = functions.TrimEscape(httpReq.FormValue("title"))
			}
		}

		if formProfile["dob"].(string) == "" {
			if functions.TrimEscape(httpReq.FormValue("dob")) == "" {
				sMessage += "Date of Birth is missing <br>"
			} else {
				xDocProfile["dob"] = functions.TrimEscape(httpReq.FormValue("dob"))

				yearsDifference := functions.GetDifferenceInYears("", functions.TrimEscape(httpReq.FormValue("dob")))
				if yearsDifference < 18 {
					sMessage += "You Must be at least 18yrs old to Subscribe<br>"
				}
			}
		}

		if formProfile["nationality"].(string) == "" {
			if functions.TrimEscape(httpReq.FormValue("nationality")) == "" {
				sMessage += "Nationality is missing <br>"
			} else {
				xDocProfile["nationality"] = functions.TrimEscape(httpReq.FormValue("nationality"))
			}
		}
	}

	if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
		xDocProfile["phonecode"] = "+971"
	} else {
		xDocProfile["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
	}

	if functions.TrimEscape(httpReq.FormValue("phone")) == "" {
		sMessage += "Mobile Number is missing <br>"
	} else {
		xDocProfile["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))
	}

	sProfileEmailNew := ""
	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Primary Email is missing <br>"
	} else {
		sProfileEmailNew = functions.TrimEscape(httpReq.FormValue("email"))
	}

	sUsernameNew := ""
	if functions.TrimEscape(httpReq.FormValue("username")) != "" {
		sUsernameNew = functions.TrimEscape(httpReq.FormValue("username"))
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	sNotify := ""
	//Check if Username has changed and if it exists
	if sUsernameNew != "" && sUsernameNew != this.mapAppCache["username"].(string) {

		//--> Check if Email/Username exists in Login Table
		sqlCheckEmail := fmt.Sprintf(`select control, username, title from profile where username = '%s'`, sUsernameNew)
		mapCheckEmail, _ := curdb.Query(sqlCheckEmail)
		//--> Check if Email/Username exists in Login Table

		if mapCheckEmail["1"] != nil {
			sMessage += "Username exists and cannot be used<br>"
		} else {

			sNotify = "To complete the process please check your email (" + sProfileEmail + ") . <br> <br> Within the email you will find an activation link which you must click in order to change your username."
			sNotify = fmt.Sprintf(`"error":"%s",`, sNotify)

			//generateActivationLink
			tblActivationLink := new(database.ActivationLink)

			xDocLink := make(map[string]interface{})
			xDocLink["code"] = functions.RandomString(32)
			xDocLink["fullname"] = fmt.Sprintf("User [%s] Please confirm your username", this.mapAppCache["username"])
			xDocLink["workflow"] = "draft"
			xDocLink["description"] = sUsernameNew
			xDocLink["logincontrol"] = this.mapAppCache["control"]
			tblActivationLink.Create("reset", xDocLink, curdb)

			activationLink := fmt.Sprintf("%sapp-activate/?action=username&code=%s", httpReq.Referer(), xDocLink["code"])
			//generateActivationLink

			//generateUsernameChangeEmail
			emailFrom := "username@valued.com"
			emailFromName := "VALUED"
			emailTo := sProfileEmail
			emailSubject := "VALUED - Please confirm your username"
			emailTemplate := "app-username-change"

			emailFields := make(map[string]interface{})
			emailFields["fullname"] = fmt.Sprintf(`%v %v %v`, this.mapAppCache["title"], this.mapAppCache["firstname"], this.mapAppCache["lastname"])
			emailFields["link"] = activationLink

			go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, nil, emailFields)
			//generateUsernameChangeEmail

		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}
	//Check if Username has changed and if it exists

	//Check if Email has changed
	if sProfileEmail != sProfileEmailNew {

		sNotify = "To complete the process please check your email (" + sProfileEmailNew + ") . <br> <br> Within the email you will find an activation link which you must click in order to change your email."
		sNotify = fmt.Sprintf(`"error":"%s",`, sNotify)

		//generateActivationLink
		tblActivationLink := new(database.ActivationLink)

		xDocLink := make(map[string]interface{})
		xDocLink["code"] = functions.RandomString(32)
		xDocLink["title"] = fmt.Sprintf("User [%s] Please confirm your email", this.mapAppCache["username"])
		xDocLink["workflow"] = "draft"
		xDocLink["description"] = sProfileEmailNew
		xDocLink["logincontrol"] = this.mapAppCache["control"]
		tblActivationLink.Create("reset", xDocLink, curdb)

		activationLink := fmt.Sprintf("%sapp-activate/?action=email&code=%s", httpReq.Referer(), xDocLink["code"])
		//generateActivationLink

		//generateEmailChangeEmail
		emailFrom := "email@valued.com"
		emailFromName := "VALUED"
		emailTo := sProfileEmail
		emailSubject := "VALUED - Please confirm your email"
		emailTemplate := "app-email-change"

		emailFields := make(map[string]interface{})
		emailFields["fullname"] = fmt.Sprintf(`%v %v %v`, this.mapAppCache["title"], this.mapAppCache["firstname"], this.mapAppCache["lastname"])
		emailFields["link"] = activationLink

		go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, nil, emailFields)
		//generateEmailChangeEmail

	}
	//Check if Email has changed

	tblProfile := new(database.Profile)
	tblProfile.Update(this.mapAppCache["username"].(string), xDocProfile, curdb)
	// this.View(httpRes, httpReq, curdb)
	httpRes.Write([]byte(`{` + sNotify + this.View(httpRes, httpReq, curdb) + `}`))
	return
}

func (this *AppProfile) redemption(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sqlRedemption := `select r.control as control,
	 r.savingsvalue as savingsvalue, r.transactionvalue as transactionvalue,
	 r.merchantcontrol as merchantcontrol, r.membercontrol as membercontrol,
	 rwd.title as rewardtitle, s.title as storetitle,
	 m.title as merchanttitle, m.image as merchantimage
     
     
	 from redemption as r 
     left join reward as rwd on r.rewardcontrol = rwd.control
     left join profile as m on r.merchantcontrol = m.control
     left join store as s on r.storecontrol = s.control

	 where r.membercontrol = '%s' 
	 
	 order by control desc`

	sqlRedemption = fmt.Sprintf(sqlRedemption, this.mapAppCache["control"])
	xDocresult, _ := curdb.Query(sqlRedemption)

	sumofsavings := 0.00
	formProfile := make(map[string]interface{})

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		savingsValue := xDoc["savingsvalue"].(float64)
		delete(xDoc, "savingsvalue")

		sumofsavings += savingsValue
		xDoc["savingsvalue"] = functions.ThousandSeperator(functions.Round(savingsValue))
		formProfile[cNumber+"#app-profile-redemption-list"] = xDoc
	}

	formProfile["sumofsavings"] = functions.ThousandSeperator(functions.Round(sumofsavings))

	appNavbar := new(AppNavbar)
	formProfile["app-navbar"] = appNavbar.GetNavBar(this.mapAppCache)
	formProfile["app-navbar-button"] = appNavbar.GetNavBarButton(this.mapAppCache)

	appFooter := make(map[string]interface{})
	formProfile["app-footer"] = appFooter
	appFooter["profile"] = "white"

	this.pageMap["app-profile-redemption"] = formProfile
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Redemption History","pageContent":` + contentHTML + `}`))
}

func (this *AppProfile) mySavings(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	formProfile := make(map[string]interface{})
	appFooter := make(map[string]interface{})
	formProfile["app-footer"] = appFooter
	appFooter["profile"] = "white"

	this.pageMap["app-profile-savings"] = formProfile
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"My Savings","pageContent":` + contentHTML + `}`))
}

func (this *AppProfile) changePassword(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	formProfile := make(map[string]interface{})

	appFooter := make(map[string]interface{})
	formProfile["app-footer"] = appFooter
	appFooter["profile"] = "white"

	this.pageMap["app-profile-password"] = formProfile
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Change Password","pageContent":` + contentHTML + `}`))
}

func (this *AppProfile) savePassword(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("current")) == "" {
		sMessage += "Current Password is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("password")) == "" {
		sMessage += "New Password is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("confirm")) == "" {
		sMessage += "Confirm New Password is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("password")) != functions.TrimEscape(httpReq.FormValue("confirm")) {
		sMessage += "New Password does not match <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblUser := new(database.Profile)
	xDoc := make(map[string]interface{})
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = this.mapAppCache["control"].(string)

	ResMap := tblUser.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		sMessage += "Cannot find user to update <br>"
	} else {
		xDoc = ResMap["1"].(map[string]interface{})
		if xDoc["password"].(string) != functions.TrimEscape(httpReq.FormValue("current")) {
			sMessage += "Current Password is incorrect <br>"
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocNew := make(map[string]interface{})
	xDocNew["control"] = this.mapAppCache["control"].(string)
	xDocNew["password"] = functions.TrimEscape(httpReq.FormValue("password"))
	tblUser.Update(this.mapAppCache["username"].(string), xDocNew, curdb)

	httpRes.Write([]byte(`{"alertSuccess":"Password Changed Successfully","getform":"/app-profile"}`))
}

func (this *AppProfile) setPin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	formProfile := make(map[string]interface{})

	appFooter := make(map[string]interface{})
	formProfile["app-footer"] = appFooter
	appFooter["profile"] = "white"

	sPageTitle := "Change Pin"
	sPageTemplate := "app-profile-pin"
	if this.mapAppCache["pincode"] == nil || this.mapAppCache["pincode"].(string) == "" {
		sPageTitle = "Create Pin"
		sPageTemplate = "app-profile-createpin"
	}

	this.pageMap[sPageTemplate] = formProfile
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"` + sPageTitle + `","pageContent":` + contentHTML + `}`))
}

func (this *AppProfile) savePin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	lChangePin := true
	sPinAlert := "Changed"
	if this.mapAppCache["pincode"] == nil || this.mapAppCache["pincode"].(string) == "" {
		sPinAlert = "Created"
		lChangePin = false
	}

	sCurrentPin := ""
	if lChangePin {
		sCurrentPin = functions.TrimEscape(httpReq.FormValue("currentPin1")) + functions.TrimEscape(httpReq.FormValue("currentPin2")) +
			functions.TrimEscape(httpReq.FormValue("currentPin3")) + functions.TrimEscape(httpReq.FormValue("currentPin4"))

		if sCurrentPin == "" || len(sCurrentPin) != 4 {
			sMessage += "Current 4 Digit Pin is missing <br>"
		}
	}

	sNewPin := functions.TrimEscape(httpReq.FormValue("newPin1")) + functions.TrimEscape(httpReq.FormValue("newPin2")) +
		functions.TrimEscape(httpReq.FormValue("newPin3")) + functions.TrimEscape(httpReq.FormValue("newPin4"))

	if sNewPin == "" || len(sNewPin) != 4 {
		sMessage += "New 4 Digit Pin is missing <br>"
	}

	sConfirmPin := functions.TrimEscape(httpReq.FormValue("confirmPin1")) + functions.TrimEscape(httpReq.FormValue("confirmPin2")) +
		functions.TrimEscape(httpReq.FormValue("confirmPin3")) + functions.TrimEscape(httpReq.FormValue("confirmPin4"))

	if sConfirmPin == "" || len(sConfirmPin) != 4 {
		sMessage += "Confirm New 4 Digit Pin is missing <br>"
	}

	if sNewPin != sConfirmPin {
		sMessage += "New Pin does not match <br>"
	}

	tblUser := new(database.Profile)
	xDoc := make(map[string]interface{})

	if lChangePin {
		xDocRequest := make(map[string]interface{})
		xDocRequest["searchvalue"] = this.mapAppCache["control"].(string)

		ResMap := tblUser.Read(xDocRequest, curdb)
		if ResMap["1"] == nil {
			sMessage += "Cannot find user to update <br>"
		} else {
			xDoc = ResMap["1"].(map[string]interface{})
			if xDoc["pincode"].(string) != sCurrentPin {
				sMessage += "Current Pin is incorrect <br>"
			}
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocNew := make(map[string]interface{})
	xDocNew["pincode"] = sConfirmPin
	xDocNew["control"] = this.mapAppCache["control"].(string)
	tblUser.Update(this.mapAppCache["username"].(string), xDocNew, curdb)

	this.mapAppCache["pincode"] = sConfirmPin
	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	curdb.SetSession(GOSESSID.Value, "mapAppCache", this.mapAppCache, false)

	httpRes.Write([]byte(`{"error":"Your 4-Digit Pin has been ` + sPinAlert + `","getform":"/app-profile"}`))
}

func (this *AppProfile) saveProfilePic(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if httpReq.FormValue("image") == "" {
		sMessage += "Image is missing <br>"
	}

	if httpReq.FormValue("imageName") == "" {
		sMessage += "Image Name is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	base64String := httpReq.FormValue("image")
	base64String = strings.Split(base64String, "base64,")[1]
	base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
	if base64Bytes != nil && err == nil {
		fileName := fmt.Sprintf("profile-%s-%s", functions.RandomString(6),
			functions.TrimEscape(httpReq.FormValue("imageName")))
		xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
	}

	xDoc["control"] = this.mapAppCache["control"].(string)
	new(database.Profile).Update(this.mapAppCache["username"].(string), xDoc, curdb)

	httpRes.Write([]byte(`{"error":"Profile Pic Changed","getform":"/app-profile"}`))

}

func (this *AppProfile) forgotPin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""

	//generateActivationLink
	tblActivationLink := new(database.ActivationLink)

	xDocLink := make(map[string]interface{})
	xDocLink["code"] = functions.RandomString(32)
	xDocLink["title"] = fmt.Sprintf("User [%v] Pin Reset", this.mapAppCache["email"])
	xDocLink["workflow"] = "draft"
	xDocLink["logincontrol"] = this.mapAppCache["control"]
	xDocLink["description"] = this.mapAppCache["email"]
	tblActivationLink.Create("reset", xDocLink, curdb)

	activationLink := fmt.Sprintf("%sapp-activate/?action=resetpin&code=%v", httpReq.Referer(), xDocLink["code"])
	//generateActivationLink

	//generateActivationEmailLink and sendActivationMail
	emailFrom := "pin@valued.com"
	emailFromName := "VALUED"
	emailTo := this.mapAppCache["username"].(string)
	emailSubject := "VALUED - Pin reset"
	emailTemplate := "app-pin-reset"

	emailFields := make(map[string]interface{})
	emailFields["fullname"] = fmt.Sprintf(`%v %v`, this.mapAppCache["firstname"], this.mapAppCache["lastname"])
	emailFields["email"] = this.mapAppCache["email"]
	emailFields["link"] = activationLink

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, nil, emailFields)
	//generateActivationEmailLink and sendActivationMail

	sMessage = "The first stage of your Pin Reset has been successful. <br><br>"
	sMessage += "To complete the process please check your email. Within the email you will find a link which you must click in order to reset your pin."
	httpRes.Write([]byte(`{"alertSuccess":"` + sMessage + `"}`))
}
