package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"

	"strings"
	"time"
)

func (this *Employee) viewSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employee")) == "" {
		sMessage += "Valued Employee Profile is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchStore := make(map[string]interface{})
	subSearchStore["employee"] = html.EscapeString(httpReq.FormValue("employee"))
	subSearch["employeeview-subscription-search"] = subSearchStore

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *Employee) searchSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employee")) == "" {
		sMessage += "Valued Employee Profile is missing <br>"
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
				xDocrequest["offset"] = fmt.Sprintf("%s", offset)
			}
		}
	}

	tblSubscription := new(database.Subscription)
	subSearch := make(map[string]interface{})
	xDocrequest["member"] = functions.TrimEscape(httpReq.FormValue("employee"))
	xDocrequest["employer"] = this.mapCache["control"]
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

		subSearch[cNumber+"#employeeview-subscription-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *Employee) newSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employee")) == "" {
		sMessage += "Employee is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	formNew := make(map[string]interface{})
	formNewFields := make(map[string]interface{})
	formNewFields["employee"] = functions.TrimEscape(httpReq.FormValue("employee"))
	formNew["employeeview-subscription-edit"] = formNewFields

	viewHTML := strconv.Quote(string(this.Generate(formNew, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Employee) editSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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
	xDocRequest["employee"] = xDocRequest["membercontrol"]

	formView := make(map[string]interface{})
	formView["employeeview-subscription-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"subForm":` + viewHTML + `}`))
}

func (this *Employee) saveSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("employee")) == "" {
		sMessage += "Employee is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
		sMessage += "Scheme is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("startdate")) == "" {
		sMessage += "Start Date is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	xDoc["employercontrol"] = this.mapCache["control"]

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
	xDoc["membercontrol"] = functions.TrimEscape(httpReq.FormValue("employee"))

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		xDoc["workflow"] = "unpaid"
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
		xDoc["workflow"] = "unpaid"
		xDoc["control"] = tblSubscription.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	httpRes.Write([]byte(`{"error":"Record Saved","triggerSubSearch":true,"toggleSubForm":"true"}`))
	return
}

func (this *Employee) subscribeEmployee(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
	mapScheme, _ := curdb.Query(sqlScheme)
	if mapScheme["1"] == nil {
		sMessage := "Please Select Scheme <br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	mapIsSubscribed := make(map[string]interface{})
	sqlSubscriptions := fmt.Sprintf(`select membercontrol from subscription where schemecontrol = '%s' and employercontrol = '%s'`,
		functions.TrimEscape(httpReq.FormValue("scheme")), this.mapCache["control"])
	mapSubscriptions, _ := curdb.Query(sqlSubscriptions)
	for _, xDoc := range mapSubscriptions {
		xDoc := xDoc.(map[string]interface{})
		mapIsSubscribed[xDoc["membercontrol"].(string)] = xDoc["membercontrol"]
	}

	nTotal := 0
	httpReq.ParseForm()

	sliceEmployees := make(map[string]string)
	for _, control := range httpReq.Form["control"] {
		sliceEmployees[control] = control
	}

	for _, control := range sliceEmployees {
		// controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)

		//Check if User Aleady Linked
		if mapIsSubscribed[control] != nil {
			continue
		}
		//Check if User Aleady Linked

		//Link Employee to Scheme
		sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
		mapScheme, _ := curdb.Query(sqlScheme)

		if mapScheme["1"] != nil {
			tblSubscription := new(database.Subscription)

			xDocSubscription := make(map[string]interface{})
			xDocSubscription["membercontrol"] = control
			xDocSubscription["employercontrol"] = this.mapCache["control"]

			xDocSubscription["price"] = mapScheme["1"].(map[string]interface{})["price"]
			xDocSubscription["schemecontrol"] = functions.TrimEscape(httpReq.FormValue("scheme"))
			xDocSubscription["workflow"] = "unpaid"

			cFormat := "02/01/2006"
			todayDate, _ := time.Parse(cFormat, functions.GetSystemDate())
			oneYear := todayDate.Add(time.Hour * 24 * 365)

			//Get Current Subscription Expiry
			sqlSubscription := `select sub.schemecontrol as schemecontrol, sub.code as code, sub.control as control, sub.expirydate as expirydate, sch.title as schemetitle
								from subscription as sub join scheme as sch on sub.schemecontrol = sch.control
								AND sub.workflow = 'active' AND sch.workflow = 'active' AND sch.code in ('lite','lifestyle')  
								AND '%s'::timestamp between sub.startdate::timestamp and sub.expirydate::timestamp AND sub.schemecontrol = '%s' AND sub.membercontrol = '%s' order by control desc`
			sqlSubscription = fmt.Sprintf(sqlSubscription, todayDate.Format("02/01/2006"), xDocSubscription["schemecontrol"], xDocSubscription["membercontrol"])

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

			xDocSubscription["startdate"] = todayDate.Format("02/01/2006")
			xDocSubscription["expirydate"] = oneYear.Format("02/01/2006")

			tblSubscription.Create(this.mapCache["username"].(string), xDocSubscription, curdb)
		}
		//Link Employee to Scheme
		nTotal++
	}

	if nTotal == 0 {
		sMessage = "No Employees Subscribed"
	} else {
		sMessage = fmt.Sprintf("%v  Employees Subscribed", nTotal)
	}

	httpRes.Write([]byte(`{"error":"` + sMessage + `","triggerSearch":true}`))
}

func (this *Employee) renewSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
	mapScheme, _ := curdb.Query(sqlScheme)
	if mapScheme["1"] == nil {
		sMessage := "Please Select Scheme <br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	mapIsSubscribed := make(map[string]interface{})
	sqlSubscriptions := fmt.Sprintf(`select membercontrol, expirydate from subscription where schemecontrol = '%s' and employercontrol = '%s'`,
		functions.TrimEscape(httpReq.FormValue("scheme")), this.mapCache["control"])
	mapSubscriptions, _ := curdb.Query(sqlSubscriptions)
	for _, xDoc := range mapSubscriptions {
		xDoc := xDoc.(map[string]interface{})
		mapIsSubscribed[xDoc["membercontrol"].(string)] = xDoc["expirydate"]
	}

	nTotal := 0
	httpReq.ParseForm()

	sliceEmployees := make(map[string]string)
	for _, control := range httpReq.Form["control"] {
		sliceEmployees[control] = control
	}

	for _, control := range sliceEmployees {

		//Check if User Aleady Linked
		if mapIsSubscribed[control] == nil {
			continue
		}
		//Check if User Aleady Linked

		//Link Employee to Scheme
		sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
		mapScheme, _ := curdb.Query(sqlScheme)

		if mapScheme["1"] != nil {
			tblSubscription := new(database.Subscription)

			xDocSubscription := make(map[string]interface{})
			xDocSubscription["membercontrol"] = control
			xDocSubscription["employercontrol"] = this.mapCache["control"]

			xDocSubscription["price"] = mapScheme["1"].(map[string]interface{})["price"]
			xDocSubscription["schemecontrol"] = functions.TrimEscape(httpReq.FormValue("scheme"))
			xDocSubscription["workflow"] = "none"

			cFormat := "02/01/2006"
			todayDate, _ := time.Parse(cFormat, mapIsSubscribed[control].(string))
			oneYear := todayDate.Add(time.Hour * 24 * 365)

			xDocSubscription["startdate"] = todayDate.Format("02/01/2006")
			xDocSubscription["expirydate"] = oneYear.Format("02/01/2006")

			tblSubscription.Create(this.mapCache["username"].(string), xDocSubscription, curdb)
		}
		//Link Employee to Scheme
		nTotal++
	}

	if nTotal == 0 {
		sMessage = "No Subscriptions Renewed"
	} else {
		sMessage = fmt.Sprintf("%v Subscriptions Renewed", nTotal)
	}

	httpRes.Write([]byte(`{"error":"` + sMessage + `","triggerSearch":true}`))
}

func (this *Employee) deleteSubscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
	mapScheme, _ := curdb.Query(sqlScheme)
	if mapScheme["1"] == nil {
		sMessage := "Please Select Scheme <br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	mapIsSubscribed := make(map[string]interface{})
	sqlSubscriptions := fmt.Sprintf(`select membercontrol, control from subscription where schemecontrol = '%s' and employercontrol = '%s'`,
		functions.TrimEscape(httpReq.FormValue("scheme")), this.mapCache["control"])
	mapSubscriptions, _ := curdb.Query(sqlSubscriptions)
	for _, xDoc := range mapSubscriptions {
		xDoc := xDoc.(map[string]interface{})
		mapIsSubscribed[xDoc["membercontrol"].(string)] = xDoc["control"]
	}

	nTotal := 0
	httpReq.ParseForm()

	sliceEmployees := make(map[string]string)
	for _, control := range httpReq.Form["control"] {
		sliceEmployees[control] = control
	}

	for _, control := range sliceEmployees {

		//Check if User Aleady Linked
		if mapIsSubscribed[control] == nil {
			continue
		}
		//Check if User Aleady Linked

		//Link Employee to Scheme
		sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
		mapScheme, _ := curdb.Query(sqlScheme)

		if mapScheme["1"] != nil {
			tblSubscription := new(database.Subscription)

			xDocSubscription := make(map[string]interface{})
			xDocSubscription["control"] = mapIsSubscribed[control]
			xDocSubscription["workflow"] = "expired"
			tblSubscription.Update(this.mapCache["username"].(string), xDocSubscription, curdb)
		}
		//Link Employee to Scheme
		nTotal++
	}

	if nTotal == 0 {
		sMessage = "No Subscriptions Deleted"
	} else {
		sMessage = fmt.Sprintf("%v Subscriptions Deleted", nTotal)
	}

	httpRes.Write([]byte(`{"error":"` + sMessage + `","triggerSearch":true}`))
}
