package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"

	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type AppLogin struct {
	functions.Templates
	mapAppCache map[string]interface{}
}

func (this *AppLogin) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	if httpReq.Method == "GET" {

		sPayload := "getForm('/app-login')"
		pageMap := make(map[string]interface{})
		htmlContent := make(map[string]interface{})

		if this.mapAppCache["control"] != nil {
			sPayload = "getForm('/app-home')"
		}

		htmlContent["payload"] = sPayload
		pageMap["html"] = htmlContent

		httpRes.Write(this.Generate(pageMap, nil))
		return
	}

	if this.mapAppCache["control"] != nil {
		new(AppProfile).Process(httpRes, httpReq, curdb)
		return
	}

	httpRes.Header().Set("content-type", "application/json")

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":

		pageMap := make(map[string]interface{})
		loginPage := make(map[string]interface{})

		appFooter := make(map[string]interface{})
		loginPage["app-footer"] = appFooter
		appFooter["home"] = "white"

		loginPage["formId"] = strings.TrimSpace(httpReq.FormValue("p"))

		loginPage["app-navbar"] = new(AppNavbar).GetNavBar(this.mapAppCache)

		pageMap["app-login"] = loginPage
		pageHTML := strconv.Quote(string(this.Generate(pageMap, nil)))

		httpRes.Write([]byte(`{"pageTitle":"Valued Login","pageContent":` + pageHTML + `}`))
		return

	case "forgotpassword":
		this.forgotpassword(httpRes, httpReq, curdb)
		return

	case "reset":
		this.reset(httpRes, httpReq, curdb)
		return

	case "savepassword":
		this.savepassword(httpRes, httpReq, curdb)
		return

	case "login":
		cUsername := strings.TrimSpace(httpReq.FormValue("user"))
		cPassword := strings.TrimSpace(httpReq.FormValue("pass"))
		if cUsername != "" && cPassword != "" {

			tblProfile := new(database.Profile)
			mapRes := tblProfile.VerifyLogin("member", cUsername, cPassword, curdb)

			if mapRes["1"] != nil {
				GOSESSID, _ := httpReq.Cookie(_COOKIE_)
				xDoc := mapRes["1"].(map[string]interface{})

				if xDoc["employercode"] != nil && xDoc["employercode"].(string) == "britishmums" {
					httpRes.Write([]byte(`{"error":"Invalid Login Details"}`))
					return
				}

				curdb.SetSession(GOSESSID.Value, "mapAppCache", xDoc, false)
				new(AppProfile).Process(httpRes, httpReq, curdb)
				return
			}
		}

		httpRes.Write([]byte(`{"error":"Invalid Login Details"}`))
		return

	case "signup":
		this.signup(httpRes, httpReq, curdb)
		return
	}
}

func (this *AppLogin) savepassword(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	xDocUser := curdb.GetSession(GOSESSID.Value, "mapAppResetPassword")

	if xDocUser["control"] == nil {
		sMessage := "Invalid Activation Link <br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `","getform":"/"}`))
		return
	}

	sMessage := ""
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

	xDocNew := make(map[string]interface{})
	xDocNew["control"] = xDocUser["control"].(string)
	xDocNew["password"] = functions.TrimEscape(httpReq.FormValue("password"))
	new(database.Profile).Update("reset", xDocNew, curdb)

	httpRes.Write([]byte(`{"alertSuccess":"Password Changed Successfully - <br> Please Login With New Password","getform":"/app-login"}`))

}

func (this *AppLogin) reset(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	xDocLink := curdb.GetSession(GOSESSID.Value, "mapAppReset")
	curdb.SetSession(GOSESSID.Value, "mapAppReset", make(map[string]interface{}), false)

	if xDocLink["workflow"] == nil {
		sMessage := "Invalid Activation Link <br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `","getform":"/"}`))
		return
	}

	if xDocLink["workflow"].(string) != "active" {
		sMessage := "Expired Activation Link <br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `","getform":"/"}`))
		return
	}

	xDocRequest := make(map[string]interface{})
	xDocRequest["searchfield"] = "email"
	xDocRequest["searchvalue"] = xDocLink["description"]
	xDocResult := new(database.Profile).Read(xDocRequest, curdb)

	xDocLink["workflow"] = "inactive"
	new(database.ActivationLink).Update("reset", xDocLink, curdb)

	if xDocResult["1"] == nil {
		sMessage := "Invalid User Activation Link <br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `","getform":"/"}`))
		return
	}

	xDoc := xDocResult["1"].(map[string]interface{})
	curdb.SetSession(GOSESSID.Value, "mapAppResetPassword", xDoc, false)

	pageMap := make(map[string]interface{})
	pageMap["app-password-reset"] = xDoc
	pageHTML := strconv.Quote(string(this.Generate(pageMap, nil)))

	httpRes.Write([]byte(`{"pageTitle":"Set New Password","pageContent":` + pageHTML + `}`))

}

func (this *AppLogin) forgotpassword(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Email is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//--> Check if Email/Username exists in Login Table
	sqlCheckEmail := fmt.Sprintf(`select control, username, title, firstname, lastname, email from profile where email = '%s' `,
		functions.TrimEscape(httpReq.FormValue("email")))
	mapCheckEmail, _ := curdb.Query(sqlCheckEmail)

	if mapCheckEmail["1"] == nil {
		sMessage += "Email does not exist<br>"
	}
	//--> Check if Email/Username exists in Login Table

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//generateActivationLink
	tblActivationLink := new(database.ActivationLink)
	xDocProfile := mapCheckEmail["1"].(map[string]interface{})

	xDocLink := make(map[string]interface{})
	xDocLink["code"] = functions.RandomString(32)
	xDocLink["title"] = fmt.Sprintf("User [%s] Password Reset", xDocProfile["email"])
	xDocLink["workflow"] = "draft"
	xDocLink["logincontrol"] = xDocProfile["control"]
	xDocLink["description"] = xDocProfile["email"]
	tblActivationLink.Create("reset", xDocLink, curdb)

	activationLink := fmt.Sprintf("%sapp-activate/?action=reset&code=%s", httpReq.Referer(), xDocLink["code"])
	//generateActivationLink

	//generateActivationEmailLink and sendActivationMail
	emailFrom := "password@valued.com"
	emailFromName := "VALUED"
	emailTo := xDocProfile["username"].(string)
	emailSubject := "VALUED - Password reset"
	emailTemplate := "app-password-reset"

	emailFields := make(map[string]interface{})
	emailFields["fullname"] = fmt.Sprintf(`%v %v`, xDocProfile["firstname"], xDocProfile["lastname"])
	emailFields["email"] = xDocProfile["email"]
	emailFields["link"] = activationLink

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, nil, emailFields)
	//generateActivationEmailLink and sendActivationMail

	sMessage = "The first stage of your Password Reset has been successful. <br><br>"
	sMessage += "To complete the process please check your email. Within the email you will find a link which you must click in order to reset your password."
	httpRes.Write([]byte(`{"alertSuccess":"` + sMessage + `","getform":"/app-login"}`))
}

func (this *AppLogin) signup(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("firstname")) == "" {
		sMessage += "First Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("surname")) == "" {
		sMessage += "Surname is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Email is missing <br>"
	} else {
		regexValidEmail := "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"
		lValidEmail, _ := regexp.MatchString(regexValidEmail, functions.TrimEscape(httpReq.FormValue("email")))
		if !lValidEmail {
			sMessage += "Your Email is Not Valid <br>"
		}
	}

	if functions.TrimEscape(httpReq.FormValue("password")) == "" {
		sMessage += "Password is missing <br>"
	}

	//--> Check if Email/Username exists in Login Table
	sqlCheckEmail := fmt.Sprintf(`select control from profile where username = '%s' or email = '%s'`,
		functions.TrimEscape(httpReq.FormValue("email")), functions.TrimEscape(httpReq.FormValue("email")))
	mapCheckEmail, _ := curdb.Query(sqlCheckEmail)

	if mapCheckEmail["1"] != nil {
		sMessage += "Email already exists <br>"
	}
	//--> Check if Email/Username exists in Login Table

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//Create Profile
	xDocProfile := make(map[string]interface{})
	defaultMap, _ := curdb.Query(`select control from profile where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Employer <b>'None'</b> <br>"
	} else {
		xDocProfile["employercontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	}

	defaultMap, _ = curdb.Query(`select control from category where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Employer <b>'None'</b> <br>"
	} else {
		xDocProfile["categorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
		xDocProfile["subcategorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	}

	defaultMap, _ = curdb.Query(`select control from industry where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Employer <b>'None'</b> <br>"
	} else {
		xDocProfile["industrycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
		xDocProfile["subindustrycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblProfile := new(database.Profile)
	xDocProfile["firstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("firstname")))
	xDocProfile["lastname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("surname")))
	xDocProfile["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDocProfile["password"] = functions.TrimEscape(httpReq.FormValue("password"))
	xDocProfile["username"] = xDocProfile["email"]
	xDocProfile["workflow"] = "registered"
	xDocProfile["status"] = "inactive"
	xDocProfile["company"] = ""

	xDocProfile["control"] = tblProfile.Create("signup", xDocProfile, curdb)

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

	//generateActivationLink
	tblActivationLink := new(database.ActivationLink)
	xDocLink := make(map[string]interface{})
	xDocLink["code"] = functions.RandomString(32)
	xDocLink["title"] = fmt.Sprintf("User [%v] Account Activation", xDocProfile["username"])
	xDocLink["workflow"] = "draft"
	xDocLink["logincontrol"] = xDocProfile["control"]
	tblActivationLink.Create("signup", xDocLink, curdb)

	activationLink := fmt.Sprintf("%sapp-activate/?code=%v", httpReq.Referer(), xDocLink["code"])
	//generateActivationLink

	//generateRegistrationMail
	emailFrom := "membership@valued.com"
	emailFromName := "VALUED membership"
	emailTo := xDocProfile["username"].(string)
	emailSubject := "Welcome to VALUED - Member Registration"
	emailTemplate := "app-registration"

	emailFields := make(map[string]interface{})
	emailFields["fullname"] = fmt.Sprintf(`%v %v %v`, xDocProfile["title"], xDocProfile["firstname"], xDocProfile["lastname"])
	emailFields["email"] = xDocProfile["email"]
	emailFields["username"] = xDocProfile["username"]
	emailFields["link"] = activationLink

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, nil, emailFields)
	//generateRegistrationMail

	sMessage = "The first stage of your sign up has been successful. <br><br>"
	sMessage += "To complete the process please check your email. Within the email you will find a link which you must click in order to activate your account."
	httpRes.Write([]byte(`{"alertSuccess":"` + sMessage + `","getform":"/app-home"}`))
}
