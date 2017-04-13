package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"time"

	"net/http"
)

type AppActivate struct {
	functions.Templates
	pageMap     map[string]interface{}
	htmlContent map[string]interface{}
}

func (this *AppActivate) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if httpReq.Method == "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	if functions.TrimEscape(httpReq.FormValue("code")) == "" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	this.pageMap = make(map[string]interface{})
	this.htmlContent = make(map[string]interface{})

	this.activate(httpRes, httpReq, curdb)

	this.pageMap["html"] = this.htmlContent
	httpRes.Write(this.Generate(this.pageMap, nil))
	return

}

func (this *AppActivate) activate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	//--> Check if Code exists in ActivationLink Table
	xDocrequest := make(map[string]interface{})
	xDocrequest["searchfield"] = "code"
	xDocrequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDocresult := new(database.ActivationLink).Read(xDocrequest, curdb)

	sPayloadLogin := `function redirect(){window.location.replace('/?a=app-login')}; setTimeout(redirect,3000);`

	if xDocresult["1"] == nil {
		sMessage += "Invalid Activation Link <br>"
		this.htmlContent["payload"] = sPayloadLogin
		this.htmlContent["notify"] = `$.notify({message: '` + sMessage + `'},{type: 'danger', timer: 5000});`
		return
	}
	//--> Check if Code exists in ActivationLink Table
	xDoc := xDocresult["1"].(map[string]interface{})

	xDocLink := make(map[string]interface{})
	xDocLink["control"] = xDoc["control"]
	xDocLink["workflow"] = "active"
	xDocLink["description"] = xDoc["description"]

	xDocLink["vieweddate"] = functions.GetSystemTime()
	xDocLink["expirydate"] = functions.GetTimeFormat(time.Now().Add(time.Hour * 24 * 15))

	if functions.TrimEscape(httpReq.FormValue("action")) == "email" {
		new(database.ActivationLink).Update("activation", xDocLink, curdb)

		//Force Logout of current user if any
		// GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		// curdb.SetSession(GOSESSID.Value, "mapAppCache", make(map[string]interface{}), false)

		sEmailNew := xDoc["description"].(string)
		if sEmailNew == "" {
			sMessage += "Valued Email was not Updated. <br>"
		} else {

			xDocProfile := make(map[string]interface{})
			xDocProfile["email"] = xDoc["description"]
			xDocProfile["control"] = xDoc["logincontrol"]
			new(database.Profile).Update("root", xDocProfile, curdb)

			sPayloadEmail := `function redirect(){window.location.replace('/?a=app-profile')}; setTimeout(redirect,3000);`
			sMessage += "Your Valued Email has been Updated. <br>"
			this.htmlContent["payload"] = sPayloadEmail
		}

		this.htmlContent["notify"] = `$.notify({message: '` + sMessage + `'},{type: 'success', timer: 5000});`
		return
	}

	if functions.TrimEscape(httpReq.FormValue("action")) == "username" {
		new(database.ActivationLink).Update("activation", xDocLink, curdb)

		sEmailNew := xDoc["description"].(string)
		if sEmailNew == "" {
			sMessage += "Valued Email was not Updated. <br>"
		} else {

			xDocLogin := make(map[string]interface{})
			xDocLogin["username"] = xDoc["description"]
			xDocLogin["control"] = xDoc["logincontrol"]
			new(database.Profile).Update("root", xDocLogin, curdb)

			sPayloadUsername := `function redirect(){window.location.replace('/?a=app-login')}; setTimeout(redirect,3000);`
			sMessage += "Your Valued Email has been Updated. <br>"
			this.htmlContent["payload"] = sPayloadUsername

			// Force Logout of current user if any
			GOSESSID, _ := httpReq.Cookie(_COOKIE_)
			curdb.SetSession(GOSESSID.Value, "mapAppCache", make(map[string]interface{}), false)
		}

		this.htmlContent["notify"] = `$.notify({message: '` + sMessage + `'},{type: 'success', timer: 5000});`
		return
	}

	if functions.TrimEscape(httpReq.FormValue("action")) == "view" {
		new(database.ActivationLink).Update("activation", xDocLink, curdb)
		httpRes.Write([]byte(``))
		return
	}

	if functions.TrimEscape(httpReq.FormValue("action")) == "reset" {

		// Force Logout of current user if any
		GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		curdb.SetSession(GOSESSID.Value, "mapAppReset", xDocLink, false)
		curdb.SetSession(GOSESSID.Value, "mapAppCache", make(map[string]interface{}), false)

		sMessage += "Resetting Your Valued Member account Password. <br>"
		sPayloadReset := `function redirect(){window.location.replace('/?a=app-login&action=reset')}; setTimeout(redirect,3000);`
		this.htmlContent["payload"] = sPayloadReset
		this.htmlContent["notify"] = `$.notify({message: '` + sMessage + `'},{type: 'success', timer: 5000});`
		return
	}

	if functions.TrimEscape(httpReq.FormValue("action")) == "resetpin" {

		// Force Logout of current user if any
		// GOSESSID, _ := httpReq.Cookie(_COOKIE_)
		// curdb.SetSession(GOSESSID.Value, "mapAppReset", xDocLink, false)
		// curdb.SetSession(GOSESSID.Value, "mapAppCache", make(map[string]interface{}), false)

		sqlResetPin := fmt.Sprintf(`update profile set pincode = '' where control = '%s'`, xDoc["logincontrol"])
		curdb.Query(sqlResetPin)

		sMessage += "Pin has been Reset, Please login to create a new 4-Digit Pin. <br>"
		sPayloadReset := `function redirect(){window.location.replace('/?a=app-profile')}; setTimeout(redirect,3000);`
		this.htmlContent["payload"] = sPayloadReset
		this.htmlContent["notify"] = `$.notify({message: '` + sMessage + `'},{type: 'success', timer: 5000});`
		return
	}

	if xDoc["workflow"].(string) == "inactive" {
		sMessage += "Expired Activation Link <br>"
		this.htmlContent["payload"] = sPayloadLogin
		this.htmlContent["notify"] = `$.notify({message: '` + sMessage + `'},{type: 'danger', timer: 5000});`
		return
	}

	if xDoc["control"] != nil && xDoc["control"].(string) != "" {
		xDocLink := make(map[string]interface{})
		xDocLink["workflow"] = "inactive"
		xDocLink["control"] = xDoc["control"]
		new(database.ActivationLink).Update("activation", xDocLink, curdb)
	}

	//Logic to activate login & profile
	if xDoc["logincontrol"] != nil && xDoc["logincontrol"].(string) != "" {
		xDocProfile := make(map[string]interface{})
		xDocProfile["status"] = "active"
		xDocProfile["control"] = xDoc["logincontrol"]
		new(database.Profile).Update("activation", xDocProfile, curdb)
	}

	// if xDoc["profilecontrol"] != nil && xDoc["profilecontrol"].(string) != "" {
	// 	xDocProfile := make(map[string]interface{})
	// 	xDocProfile["status"] = "inactive"
	// 	xDocProfile["control"] = xDoc["profilecontrol"]
	// 	new(database.Profile).Update("activation", xDocProfile, curdb)
	// }
	//Logic to activate login & profile

	//Force Logout of current user if any
	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	curdb.SetSession(GOSESSID.Value, "mapAppCache", make(map[string]interface{}), false)

	sMessage += "Your Valued Member account has been activated. Please login to continue. <br>"
	this.htmlContent["payload"] = sPayloadLogin
	this.htmlContent["notify"] = `$.notify({message: '` + sMessage + `'},{type: 'success', timer: 5000});`
	return
}
