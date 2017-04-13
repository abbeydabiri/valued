package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"

	"strings"
)

func (this *Employer) viewUser(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("employer")) == "" {
		sMessage += "Employer is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchUser := make(map[string]interface{})
	subSearchUser["employer"] = html.EscapeString(httpReq.FormValue("employer"))
	subSearch["employerview-user-search"] = subSearchUser

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Employer) searchUser(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("employer")) == "" {
		sMessage += "Employer is missing <br>"
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

	tblUser := new(database.Profile)
	subSearch := make(map[string]interface{})
	xDocrequest["employercontrol"] = html.EscapeString(httpReq.FormValue("employer"))
	xDocresult := tblUser.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "linkedUserActivate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "linkedUserDeactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#employerview-user-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *Employer) linkedUserActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update login set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked User Activated","triggerSubSearch":true}`))
}

func (this *Employer) linkedUserDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a user"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update login set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked User Deactivated","triggerSubSearch":true}`))
}

func (this *Employer) viewEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		sMessage += "Valued Employer Profile is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["employer"] = html.EscapeString(httpReq.FormValue("employer"))
	subSearch["employerview-member-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Employer) searchEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		sMessage += "Valued Employer Profile is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocrequest := make(map[string]interface{})
	xDocrequest["limit"] = "10"
	if functions.TrimEscape(httpReq.FormValue("limit")) == "" {
		xDocrequest["limit"] = functions.TrimEscape(httpReq.FormValue("limit"))
	}

	xDocrequest["offset"] = "0"
	if functions.TrimEscape(httpReq.FormValue("offset")) == "" {
		offset, err := strconv.Atoi(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("offset"))))
		if err == nil {
			if offset > 1 {
				xDocrequest["offset"] = fmt.Sprintf("%s", offset-1)
			}
		}
	}

	tblMember := new(database.Profile)
	subSearch := make(map[string]interface{})
	xDocrequest["employercontrol"] = functions.TrimEscape(httpReq.FormValue("employer"))
	xDocresult := tblMember.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "activateMember"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "deactivateMember"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#employerview-member-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *Employer) newEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		sMessage += "Employer is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	formNew := make(map[string]interface{})
	formNewFields := make(map[string]interface{})
	formNewFields["employer"] = functions.TrimEscape(httpReq.FormValue("employer"))
	formNew["employerview-member-edit"] = formNewFields

	viewHTML := strconv.Quote(string(this.Generate(formNew, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Employer) editEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		sMessage += "Employer is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblMember := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblMember.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})
	xDocRequest["employer"] = xDocRequest["employercontrol"]

	formView := make(map[string]interface{})
	formView["employerview-member-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Employer) saveEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		sMessage += "Employer is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
		sMessage += "Scheme is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("startdate")) == "" {
		sMessage += "Start Date is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("expirydate")) == "" {
		sMessage += "Expiry Date is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})

	sqlEmployer := fmt.Sprintf(`select employercontrol from profile where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("employer")))
	mapEmployer, _ := curdb.Query(sqlEmployer)
	if mapEmployer["1"] == nil {
		sMessage += "Employer is Missing! <br>"
	} else {
		xDoc["employercontrol"] = mapEmployer["1"].(map[string]interface{})["employercontrol"]
	}

	sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("scheme")))
	mapScheme, _ := curdb.Query(sqlScheme)
	if mapScheme["1"] == nil {
		sMessage += "Scheme is Missing! <br>"
	} else {
		xDoc["price"] = mapScheme["1"].(map[string]interface{})["price"]
		xDoc["schemecontrol"] = functions.TrimEscape(httpReq.FormValue("scheme"))
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblMember := new(database.Profile)
	xDoc["startdate"] = functions.TrimEscape(httpReq.FormValue("startdate"))
	xDoc["expirydate"] = functions.TrimEscape(httpReq.FormValue("expirydate"))
	xDoc["membercontrol"] = functions.TrimEscape(httpReq.FormValue("member"))
	xDoc["company"] = ""

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblMember.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		sTitle := "SUB" + functions.RandomString(6)
		xDoc["code"] = sTitle
		xDoc["title"] = sTitle
		xDoc["workflow"] = "inactive"
		xDoc["control"] = tblMember.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	httpRes.Write([]byte(`{"error":"Record Saved","triggerSubSearch":true,"toggleAppSidebar":"subForm"}`))
	return
}

func (this *Employer) activateEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select an employee"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update member set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Employee Activated","triggerSubSearch":true}`))
}

func (this *Employer) deactivateEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select an employee"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update member set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Employee Deactivated","triggerSubSearch":true}`))
}

func (this *Employer) viewSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		sMessage += "Valued Employer Profile is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["employer"] = html.EscapeString(httpReq.FormValue("employer"))
	subSearch["employerview-subscription-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Employer) searchSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		sMessage += "Valued Employer Profile is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDocrequest := make(map[string]interface{})
	xDocrequest["limit"] = "10"
	if functions.TrimEscape(httpReq.FormValue("limit")) == "" {
		xDocrequest["limit"] = functions.TrimEscape(httpReq.FormValue("limit"))
	}

	xDocrequest["offset"] = "0"
	if functions.TrimEscape(httpReq.FormValue("offset")) == "" {
		offset, err := strconv.Atoi(strings.TrimSpace(functions.TrimEscape(httpReq.FormValue("offset"))))
		if err == nil {
			if offset > 1 {
				xDocrequest["offset"] = fmt.Sprintf("%s", offset-1)
			}
		}
	}

	tblSubscription := new(database.Subscription)
	subSearch := make(map[string]interface{})
	xDocrequest["employer"] = functions.TrimEscape(httpReq.FormValue("employer"))
	xDocresult := tblSubscription.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "activateSubscription"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "deactivateSubscription"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#employerview-subscription-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *Employer) newSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employer")) == "" {
		sMessage += "Employer is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	formNew := make(map[string]interface{})
	formNewFields := make(map[string]interface{})
	formNewFields["employer"] = functions.TrimEscape(httpReq.FormValue("employer"))
	formNew["employerview-subscription-edit"] = formNewFields

	viewHTML := strconv.Quote(string(this.Generate(formNew, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Employer) editSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		sMessage += "Subscription is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblSubscription := new(database.Subscription)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblSubscription.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})
	xDocRequest["employer"] = xDocRequest["employercontrol"]

	formView := make(map[string]interface{})
	formView["employerview-subscription-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Employer) saveSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("member")) == "" {
		sMessage += "Member is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
		sMessage += "Scheme is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("startdate")) == "" {
		sMessage += "Start Date is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("expirydate")) == "" {
		sMessage += "Expiry Date is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})

	sqlEmployer := fmt.Sprintf(`select employercontrol from profile where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("member")))
	mapEmployer, _ := curdb.Query(sqlEmployer)
	if mapEmployer["1"] == nil {
		sMessage += "Employer is Missing! <br>"
	} else {
		xDoc["employercontrol"] = mapEmployer["1"].(map[string]interface{})["employercontrol"]
	}

	sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`,
		functions.TrimEscape(httpReq.FormValue("scheme")))
	mapScheme, _ := curdb.Query(sqlScheme)
	if mapScheme["1"] == nil {
		sMessage += "Scheme is Missing! <br>"
	} else {
		xDoc["price"] = mapScheme["1"].(map[string]interface{})["price"]
		xDoc["schemecontrol"] = functions.TrimEscape(httpReq.FormValue("scheme"))
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblSubscription := new(database.Subscription)
	xDoc["startdate"] = functions.TrimEscape(httpReq.FormValue("startdate"))
	xDoc["expirydate"] = functions.TrimEscape(httpReq.FormValue("expirydate"))
	xDoc["membercontrol"] = functions.TrimEscape(httpReq.FormValue("member"))

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblSubscription.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		sTitle := "SUB" + functions.RandomString(6)
		xDoc["code"] = sTitle
		xDoc["title"] = sTitle
		xDoc["workflow"] = "inactive"
		xDoc["control"] = tblSubscription.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	httpRes.Write([]byte(`{"error":"Record Saved","triggerSubSearch":true,"toggleAppSidebar":"subForm"}`))
	return
}

func (this *Employer) activateSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscription"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Subscription Activated","triggerSubSearch":true}`))
}

func (this *Employer) deactivateSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a Subscription"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Subscription Deactivated","triggerSubSearch":true}`))
}
