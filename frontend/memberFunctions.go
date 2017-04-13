package frontend

import (
	"valued/database"
	"valued/functions"

	"time"

	"fmt"
	"html"
	"net/http"
	"strconv"

	"strings"
)

func (this *Member) viewSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("member")) == "" {
		sMessage += "Valued Member Profile is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["member"] = html.EscapeString(httpReq.FormValue("member"))
	subSearch["memberview-subscription-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Member) searchSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("member")) == "" {
		sMessage += "Valued Member Profile is missing <br>"
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
	xDocrequest["member"] = functions.TrimEscape(httpReq.FormValue("member"))
	xDocresult := tblSubscription.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		default:
			xDoc["action"] = "activateSubscription"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "deactivateSubscription"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#memberview-subscription-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *Member) newSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("member")) == "" {
		sMessage += "Member is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	formNew := make(map[string]interface{})
	formNewFields := make(map[string]interface{})
	formNewFields["member"] = functions.TrimEscape(httpReq.FormValue("member"))
	formNew["memberview-subscription-edit"] = formNewFields

	viewHTML := strconv.Quote(string(this.Generate(formNew, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Member) editSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
	xDocRequest["member"] = xDocRequest["membercontrol"]

	formView := make(map[string]interface{})
	formView["memberview-subscription-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Member) saveSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
	xDoc["membercontrol"] = functions.TrimEscape(httpReq.FormValue("member"))

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		xDoc["workflow"] = "inactive"
	}

	cFormat := "02/01/2006"
	todayDate, _ := time.Parse(cFormat, functions.TrimEscape(httpReq.FormValue("startdate")))
	oneYear := todayDate.Add(time.Hour * 24 * 365)
	//Get Current Subscription Expiry
	sqlSubscription := `select sub.schemecontrol as schemecontrol, sub.code as code, sub.control as control, sub.expirydate as expirydate, sch.title as schemetitle
								from subscription as sub join scheme as sch on sub.schemecontrol = sch.control
								AND sub.workflow = 'active' AND sch.workflow = 'active' AND sch.code in ('lite','lifestyle')  
								AND '%s'::timestamp between sub.startdate::timestamp and sub.expirydate::timestamp AND sub.schemecontrol = '%s' AND sub.membercontrol = '%s' order by control desc`
	sqlSubscription = fmt.Sprintf(sqlSubscription, todayDate.Format("02/01/2006"), xDoc["schemecontrol"], xDoc["membercontrol"])

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

	xDoc["startdate"] = todayDate.Format("02/01/2006")
	xDoc["expirydate"] = oneYear.Format("02/01/2006")

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		tblSubscription.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		sTitle := "SUB" + functions.RandomString(6)
		xDoc["code"] = sTitle
		xDoc["title"] = sTitle
		xDoc["control"] = tblSubscription.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	httpRes.Write([]byte(`{"error":"Record Saved","triggerSubSearch":true,"toggleSubForm":"true"}`))
	return
}

func (this *Member) activateSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscrption"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Subscription Activated","triggerSubSearch":true}`))
}

func (this *Member) deactivateSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a subscription"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update subscription set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Subscription Deactivated","triggerSubSearch":true}`))
}

func (this *Member) deleteSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Subscription Not Deleted}`))
		return
	}
	//
	curdb.Query(fmt.Sprintf(`delete from subscription where control = '%s'`, functions.TrimEscape(httpReq.FormValue("control"))))
	httpRes.Write([]byte(`{"triggerSubSearch":true}`))
}
