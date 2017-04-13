package frontend

import (
	"valued/database"
	"valued/functions"

	"net/http"
	"strconv"
)

type ChangePin struct {
	MerchantControl string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *ChangePin) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.MerchantControl = this.mapCache["control"].(string)
	if this.mapCache["company"].(string) != "Yes" {
		this.MerchantControl = this.mapCache["employercontrol"].(string)
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})
		this.pageMap["merchant-changepin"] = make(map[string]interface{})

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Validate Coupon","mainpanelContent":` + contentHTML + `}`))
		return

	case "savePin":
		this.savePin(httpRes, httpReq, curdb)
		return
	}
}

func (this *ChangePin) savePin(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""

	sCurrentPin := functions.TrimEscape(httpReq.FormValue("currentPin1")) + functions.TrimEscape(httpReq.FormValue("currentPin2")) +
		functions.TrimEscape(httpReq.FormValue("currentPin3")) + functions.TrimEscape(httpReq.FormValue("currentPin4"))

	if sCurrentPin == "" || len(sCurrentPin) != 4 {
		sMessage += "Current 4 Digit Pin is missing <br>"
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

	tblMerchant := new(database.Profile)
	xDoc := make(map[string]interface{})

	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = this.MerchantControl

	ResMap := tblMerchant.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		sMessage += "Cannot find Profile to update <br>"
	} else {
		xDoc = ResMap["1"].(map[string]interface{})

		if xDoc["pincode"].(string) != sCurrentPin {
			sMessage += " Current Pin is incorrect <br> "
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocNew := make(map[string]interface{})
	xDocNew["pincode"] = sConfirmPin
	xDocNew["control"] = this.MerchantControl
	tblMerchant.Update(this.mapCache["username"].(string), xDocNew, curdb)

	this.mapCache["pincode"] = sConfirmPin
	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	curdb.SetSession(GOSESSID.Value, "mapCache", this.mapCache, false)

	httpRes.Write([]byte(`{"error":"Your 4-Digit Pin has been Changed","getform":"/dashboard"}`))
}
