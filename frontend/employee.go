package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bufio"
	"bytes"

	"github.com/jung-kurt/gofpdf"
	// "io"
)

type Employee struct {
	EmployerControl string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Employee) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if !strings.HasSuffix(httpReq.FormValue("action"), "download") {
		if httpReq.Method != "POST" {
			http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
		}
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.EmployerControl = this.mapCache["control"].(string)
	if this.mapCache["company"].(string) != "Yes" {
		this.EmployerControl = this.mapCache["employercontrol"].(string)
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})

		searchResult, searchPagination := this.search(httpReq, curdb)
		for sKey, iPagination := range searchPagination {
			searchResult[sKey] = iPagination
		}

		schemeMap, _ := curdb.Query(`select title, price, control from scheme order by title`)
		for cNumber, xDoc := range schemeMap {
			xDoc := xDoc.(map[string]interface{})
			searchResult[cNumber+"#employee-search-select"] = xDoc
		}

		//Search Scheme
		schemeMapResult, _ := curdb.Query(`select control, title from scheme order by title`)
		for cNumber, schemeXdoc := range schemeMapResult {
			schemeXdoc := schemeXdoc.(map[string]interface{})
			searchResult[cNumber+"#select-scheme"] = schemeXdoc
		}
		//Search Scheme

		this.pageMap["employee-search"] = searchResult

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Search Employees","mainpanelContentSearch":` + contentHTML + `}`))
		return

	case "quicksearch":
		this.quicksearch(httpRes, httpReq, curdb)
		return

	case "search":
		searchResult, searchPagination := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		paginationHTML := strconv.Quote(string(this.Generate(searchPagination, nil)))
		httpRes.Write([]byte(`{"searchresult":` + contentHTML + `,"searchPage":` + paginationHTML + `}`))
		return

	case "new":
		newHtml := this.new(httpReq, curdb)
		httpRes.Write([]byte(`{"mainpanelContent":` + newHtml + `}`))
		return

	case "importcsv":
		this.importcsv(httpRes, httpReq, curdb)
		return

	case "importcsvdownload":
		this.importcsvdownload(httpRes, httpReq, curdb)
		return

	case "importcsvsave":
		this.importcsvsave(httpRes, httpReq, curdb)
		return

	case "edit":
		this.edit(httpRes, httpReq, curdb)
		return

	case "save":
		this.save(httpRes, httpReq, curdb)
		return

	case "view":
		viewHTML := this.view(functions.TrimEscape(httpReq.FormValue("control")), curdb)
		httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
		return

	case "deleteSubscription":
		this.deleteSubscription(httpRes, httpReq, curdb)
		return

	case "renewSubscription":
		this.renewSubscription(httpRes, httpReq, curdb)
		return

	case "subscribeEmployee":
		this.subscribeEmployee(httpRes, httpReq, curdb)
		return

	case "viewSubscription":
		this.viewSubscription(httpRes, httpReq, curdb)
		return
	case "searchSubscription":
		this.searchSubscription(httpRes, httpReq, curdb)
		return
	case "saveSubscription":
		this.saveSubscription(httpRes, httpReq, curdb)
		return
	case "newSubscription":
		this.newSubscription(httpRes, httpReq, curdb)
		return
	case "editSubscription":
		this.editSubscription(httpRes, httpReq, curdb)
		return

	case "invoicedownload":
		this.invoicedownload(httpRes, httpReq, curdb)
		return

	case "requestInvoice":
		this.requestInvoice(httpRes, httpReq, curdb)
		return

	}
}

func (this *Employee) quicksearch(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	tblEmployee := new(database.Profile)
	quickSearch := make(map[string]interface{})
	xDocrequest := make(map[string]interface{})

	xDocrequest["limit"] = "10"
	xDocrequest["member"] = "Yes"
	xDocrequest["employercontrol"] = this.EmployerControl
	xDocrequest["title"] = functions.TrimEscape(httpReq.FormValue("title"))
	xDocresult := tblEmployee.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		xDoc["title"] = fmt.Sprintf("%s %s %s", xDoc["title"], xDoc["firstname"], xDoc["lastname"])

		quickSearch[cNumber+"#quick-search-result"] = xDoc
	}

	if len(quickSearch) == 0 {
		xDoc := make(map[string]interface{})
		xDoc["tag"] = functions.TrimEscape(httpReq.FormValue("tag"))
		xDoc["title"] = "No Employees Found"
		quickSearch["0#quick-search-result"] = xDoc
	}

	viewDropdownHtml := strconv.Quote(string(this.Generate(quickSearch, nil)))

	httpRes.Write([]byte(`{"quicksearch":` + viewDropdownHtml + `}`))
}

func (this *Employee) search(httpReq *http.Request, curdb database.Database) (formSearch, searchPagination map[string]interface{}) {

	formSearch = make(map[string]interface{})
	searchPagination = make(map[string]interface{})

	tblEmployee := new(database.Profile)
	xDocrequest := make(map[string]interface{})

	//Get Pagination Limit & Offset
	sLimit := "10"
	if functions.TrimEscape(httpReq.FormValue("limit")) != "" {
		sLimit = functions.TrimEscape(httpReq.FormValue("limit"))
	}

	sOffset := "0"
	if functions.TrimEscape(httpReq.FormValue("offset")) != "" {
		sOffset = functions.TrimEscape(httpReq.FormValue("offset"))
	}

	intLimit, _ := strconv.Atoi(sLimit)
	intOffset, _ := strconv.Atoi(sOffset)

	if intLimit > 0 && intOffset > 0 {
		sOffset = fmt.Sprintf("%v", (intOffset-1)*intLimit)
	}

	xDocrequest["limit"] = sLimit
	xDocrequest["offset"] = sOffset
	//Get Pagination Limit & Offset

	xDocrequest["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDocrequest["firstname"] = functions.TrimEscape(httpReq.FormValue("firstname"))
	xDocrequest["lastname"] = functions.TrimEscape(httpReq.FormValue("lastname"))
	xDocrequest["email"] = functions.TrimEscape(httpReq.FormValue("email"))
	xDocrequest["workflow"] = functions.TrimEscape(httpReq.FormValue("status"))
	xDocrequest["scheme"] = functions.TrimEscape(httpReq.FormValue("scheme"))
	xDocrequest["employer"] = this.EmployerControl
	xDocrequest["role"] = "member"

	//Set Pagination Limit & Offset
	nTotal := int64(0)
	xDocrequest["pagination"] = true
	xDocPagination := tblEmployee.SearchEmployee(xDocrequest, curdb)
	if xDocPagination["1"] != nil {
		xDocPagination := xDocPagination["1"].(map[string]interface{})
		if xDocPagination["paginationtotal"] != nil {
			nTotal = xDocPagination["paginationtotal"].(int64)
		}
	}
	delete(xDocrequest, "pagination")

	if nTotal > int64(intLimit) {
		nPage := int64(1)
		nPageMax := int64(nTotal/int64(intLimit)) + 1

		for nPage <= nPageMax {
			sPage := fmt.Sprintf("%v", nPage)
			mapPage := make(map[string]interface{})

			if intOffset > 0 && int64(intOffset) == nPage {
				mapPage["state"] = "selected"
			}
			mapPage["page"] = sPage
			searchPagination[sPage+"#select-page"] = mapPage
			nPage++
		}
	} else {
		searchPagination["1#select-page"] = "1"
	}
	//Set Pagination Limit & Offset

	xDocresult := tblEmployee.SearchEmployee(xDocrequest, curdb)
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["status"].(string) {
		case "inactive":
			xDoc["action"] = "activate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "deactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}
		formSearch[cNumber+"#employee-search-result"] = xDoc
	}

	return
}

func (this *Employee) new(httpReq *http.Request, curdb database.Database) string {
	formNew := make(map[string]interface{})

	formSelection := make(map[string]interface{})
	formSelection["formtitle"] = "Add"
	formNew["employee-new"] = formSelection
	return strconv.Quote(string(this.Generate(formNew, nil)))
}

func (this *Employee) importcsv(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	formImport := make(map[string]interface{})
	formImport["employee-importcsv"] = make(map[string]interface{})

	importcsvHtml := strconv.Quote(string(this.Generate(formImport, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + importcsvHtml + `}`))
}

func (this *Employee) importcsvdownload(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {
	httpRes.Header().Set("Content-Type", "text/csv")
	httpRes.Header().Set("Content-Disposition", "attachment;filename=Valued_ImportEmployeeTemplate.csv")
	httpRes.Write([]byte(strings.Join([]string{"title", "firstname", "lastname", "employeeid", "email", "activationdate"}, ",")))
}

func (this *Employee) importcsvsave(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("iagree")) == "" {
		sMessage += "Please Confirm Employees are Over 18 Years of Age <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
		sMessage += "Scheme is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("activationdate")) == "" {
		sMessage += "Activation Date is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	mapFiletypes := make(map[string]string)
	mapFiletypes["csvfile"] = "text/plain"
	mapFiles, sMessage := functions.UploadFile([]string{"csvfile"}, mapFiletypes, httpReq)

	if mapFiles["csvfile"] != nil {
		mapEmployee := mapFiles["csvfile"].(map[string]interface{})
		employeeList := string(mapEmployee["filebytes"].([]byte))
		employeeList = strings.Replace(employeeList, "\r", "", -1)

		sliceRow := strings.Split(employeeList, "\n")
		tblEmployee := new(database.Profile)
		tblSubscription := new(database.Subscription)

		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// Link Employee to Scheme
		var mapScheme map[string]interface{}
		if functions.TrimEscape(httpReq.FormValue("scheme")) != "" {
			sqlScheme := fmt.Sprintf(`select price,control from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
			mapScheme, _ = curdb.Query(sqlScheme)
		}
		// Link Employee to Scheme
		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

		sMessage := ""
		sNoneIndustry := ""
		defaultMap, _ := curdb.Query(`select control from industry where code = 'main'`)
		if defaultMap["1"] == nil {
			sMessage += "Missing Default Industry <br>"
		} else {
			sNoneIndustry = defaultMap["1"].(map[string]interface{})["control"].(string)
		}

		sNoneCategory := ""
		defaultMap, _ = curdb.Query(`select control from category where code = 'main'`)
		if defaultMap["1"] == nil {
			sMessage += "Missing Default Category <br>"
		} else {
			sNoneCategory = defaultMap["1"].(map[string]interface{})["control"].(string)
		}

		if len(sMessage) > 0 {
			httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
			return
		}

		go func() {
			for index, stringCols := range sliceRow {

				stringCols = strings.TrimSpace(stringCols)
				if index == 0 || len(stringCols) == 0 {
					continue
				}

				sliceCols := strings.Split(stringCols, ",")

				xDoc := make(map[string]interface{})
				xDoc["title"] = strings.TrimSpace(sliceCols[0])

				if len(sliceCols) > 1 {
					xDoc["firstname"] = strings.TrimSpace(sliceCols[1])
				}
				if len(sliceCols) > 2 {
					xDoc["lastname"] = strings.TrimSpace(sliceCols[2])
				}
				if len(sliceCols) > 3 {
					xDoc["employeeid"] = strings.TrimSpace(sliceCols[3])
				}
				if len(sliceCols) > 4 {
					xDoc["email"] = strings.TrimSpace(sliceCols[4])
				}

				sActivationDate := ""
				if len(sliceCols) > 5 {
					sActivationDate = strings.TrimSpace(sliceCols[5])
				}

				xDoc["industrycontrol"] = sNoneIndustry
				xDoc["subindustrycontrol"] = sNoneIndustry
				xDoc["categorycontrol"] = sNoneCategory
				xDoc["subcategorycontrol"] = sNoneCategory

				xDoc["employercontrol"] = this.EmployerControl
				xDoc["workflow"] = "registered"
				xDoc["member"] = "Yes"

				xDoc["control"] = tblEmployee.Create(this.mapCache["username"].(string), xDoc, curdb)

				if mapScheme["1"] != nil {

					xDocSubscription := make(map[string]interface{})
					xDocSubscription["membercontrol"] = xDoc["control"]
					xDocSubscription["employercontrol"] = this.EmployerControl

					xDocSubscription["price"] = mapScheme["1"].(map[string]interface{})["price"]
					xDocSubscription["schemecontrol"] = mapScheme["1"].(map[string]interface{})["control"]

					xDocSubscription["workflow"] = "registered"

					if sActivationDate == "" {
						sActivationDate = functions.GetSystemDate()
					}

					cFormat := "02/01/2006"
					todayDate, _ := time.Parse(cFormat, sActivationDate)
					oneYear := todayDate.Add(time.Hour * 24 * 365)

					xDocSubscription["startdate"] = todayDate.Format("02/01/2006")
					xDocSubscription["expirydate"] = oneYear.Format("02/01/2006")

					tblSubscription.Create(this.mapCache["username"].(string), xDocSubscription, curdb)
				}

				<-time.Tick(time.Millisecond * 25)
			}
		}()

		this.pageMap = make(map[string]interface{})
		searchResult, searchPagination := this.search(httpReq, curdb)
		for sKey, iPagination := range searchPagination {
			searchResult[sKey] = iPagination
		}

		schemeMap, _ := curdb.Query(`select title, price, control from scheme order by title`)
		for cNumber, xDoc := range schemeMap {
			xDoc := xDoc.(map[string]interface{})
			searchResult[cNumber+"#employee-search-select"] = xDoc
		}
		this.pageMap["employee-search"] = searchResult
		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))

		sMessage = fmt.Sprintf("Importing <b>%d</b> Employee Records", len(sliceRow)-1)
		httpRes.Write([]byte(`{"error":"` + sMessage + `","pageTitle":"Search Employees","mainpanelContentSearch":` + contentHTML + `}`))

	} else {
		sMessage += "File is empty "
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
	}

	return
}

func (this *Employee) save(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if functions.TrimEscape(httpReq.FormValue("iagree")) == "" {
		sMessage += "Please Confirm Employee is Over 18 Years of Age <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("code")) == "" {
		sMessage += "Employee ID is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
			sMessage += "Scheme is missing <br>"
		}

		if functions.TrimEscape(httpReq.FormValue("activationdate")) == "" {
			sMessage += "Activation Date is missing <br>"
		}
	}

	if functions.TrimEscape(httpReq.FormValue("firstname")) == "" {
		sMessage += "First Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("lastname")) == "" {
		sMessage += "Last Name is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("email")) == "" {
		sMessage += "Email is missing <br>"
	}

	if functions.TrimEscape(httpReq.FormValue("username")) == "" {
		sMessage += "Username is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	xDoc := make(map[string]interface{})
	defaultMap, _ := curdb.Query(`select control from industry where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Industry <br>"
	} else {
		xDoc["industrycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
		xDoc["subindustrycontrol"] = xDoc["industrycontrol"]
	}

	defaultMap, _ = curdb.Query(`select control from category where code = 'main'`)
	if defaultMap["1"] == nil {
		sMessage += "Missing Default Category <br>"
	} else {
		xDoc["categorycontrol"] = defaultMap["1"].(map[string]interface{})["control"]
		xDoc["subcategorycontrol"] = xDoc["categorycontrol"]
	}

	//--> Check if Email/Username exists in Login Table
	sqlCheckEmail := fmt.Sprintf(`select control from profile where email = '%s'`,
		functions.TrimEscape(httpReq.FormValue("email")))
	mapCheckEmail, _ := curdb.Query(sqlCheckEmail)

	if mapCheckEmail["1"] != nil {
		sMessage += "Email already exists <br>"
	}

	sqlCheckUsername := fmt.Sprintf(`select control from profile where username = '%s'`,
		functions.TrimEscape(httpReq.FormValue("username")))
	mapCheckUsername, _ := curdb.Query(sqlCheckUsername)

	if mapCheckUsername["1"] != nil {
		sMessage += "Username already exists <br>"
	}
	//--> Check if Email/Username exists in Login Table

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblEmployee := new(database.Profile)
	xDoc["code"] = functions.TrimEscape(httpReq.FormValue("code"))
	xDoc["title"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("title")))
	xDoc["firstname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("firstname")))
	xDoc["lastname"] = functions.CamelCase(functions.TrimEscape(httpReq.FormValue("lastname")))
	xDoc["username"] = functions.TrimEscape(httpReq.FormValue("username"))
	xDoc["email"] = functions.TrimEscape(httpReq.FormValue("email"))

	xDoc["phone"] = functions.TrimEscape(httpReq.FormValue("phone"))
	if functions.TrimEscape(httpReq.FormValue("phonecode")) == "" {
		xDoc["phonecode"] = "+971"
	} else {
		xDoc["phonecode"] = functions.TrimEscape(httpReq.FormValue("phonecode"))
	}

	xDoc["nationality"] = functions.TrimEscape(httpReq.FormValue("nationality"))
	xDoc["dob"] = functions.TrimEscape(httpReq.FormValue("dob"))
	xDoc["keywords"] = functions.TrimEscape(httpReq.FormValue("keywords"))
	xDoc["employercontrol"] = this.EmployerControl

	if httpReq.FormValue("image") != "" {
		base64String := httpReq.FormValue("image")
		base64String = strings.Split(base64String, "base64,")[1]
		base64Bytes, err := base64.StdEncoding.DecodeString(base64String)
		if base64Bytes != nil && err == nil {
			fileName := fmt.Sprintf("employee_%s_%s", functions.RandomString(6),
				functions.TrimEscape(httpReq.FormValue("imageName")))
			xDoc["image"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		}
	}

	if functions.TrimEscape(httpReq.FormValue("control")) != "" {
		xDoc["control"] = functions.TrimEscape(httpReq.FormValue("control"))
		// tblEmployee.Update(this.mapCache["username"].(string), xDoc, curdb)
		//Send An Email Alert Requesting A Change

	} else {
		xDoc["status"] = "inactive"
		xDoc["workflow"] = "registered"
		xDoc["control"] = tblEmployee.Create(this.mapCache["username"].(string), xDoc, curdb)

		if functions.TrimEscape(httpReq.FormValue("scheme")) != "" {

			//Link Employee to Scheme
			sqlScheme := fmt.Sprintf(`select price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
			mapScheme, _ := curdb.Query(sqlScheme)

			if mapScheme["1"] != nil {
				tblSubscription := new(database.Subscription)

				xDocSubscription := make(map[string]interface{})
				xDocSubscription["membercontrol"] = xDoc["control"]
				xDocSubscription["employercontrol"] = this.EmployerControl

				xDocSubscription["price"] = mapScheme["1"].(map[string]interface{})["price"]
				xDocSubscription["schemecontrol"] = functions.TrimEscape(httpReq.FormValue("scheme"))
				xDocSubscription["workflow"] = "none"

				cFormat := "02/01/2006"
				todayDate, _ := time.Parse(cFormat, functions.TrimEscape(httpReq.FormValue("activationdate")))
				oneYear := todayDate.Add(time.Hour * 24 * 365)

				xDocSubscription["startdate"] = todayDate.Format("02/01/2006")
				xDocSubscription["expirydate"] = oneYear.Format("02/01/2006")

				tblSubscription.Create(this.mapCache["username"].(string), xDocSubscription, curdb)
			}
			//Link Employee to Scheme
		}
	}

	viewHTML := this.view(xDoc["control"].(string), curdb)
	httpRes.Write([]byte(`{"error":"Record Saved","mainpanelContent":` + viewHTML + `}`))
	return
}

func (this *Employee) view(control string, curdb database.Database) string {

	if control == "" {
		return control
	}

	tblEmployee := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = control

	ResMap := tblEmployee.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		return ""
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	switch xDocRequest["workflow"].(string) {
	case "inactive":
		xDocRequest["actionView"] = "activateView"
		xDocRequest["actionColor"] = "success"
		xDocRequest["actionLabel"] = "Activate"

	case "active":
		xDocRequest["actionView"] = "deactivateView"
		xDocRequest["actionColor"] = "danger"
		xDocRequest["actionLabel"] = "De-Activate"
	}

	xDocRequest["createdate"] = xDocRequest["createdate"].(string)[0:19]

	nCurYear, _ := strconv.Atoi(functions.GetSystemTime()[6:10])
	nDobYear := nCurYear
	if len(xDocRequest["dob"].(string)) == 10 {
		nDobYear, _ = strconv.Atoi(xDocRequest["dob"].(string)[6:10])
	}
	xDocRequest["age"] = nCurYear - nDobYear

	//Get Employer Info via EmployerControl
	tblEmployer := new(database.Profile)
	employerReq := make(map[string]interface{})
	employerReq["searchvalue"] = xDocRequest["employercontrol"]
	employerRes := tblEmployer.Read(employerReq, curdb)
	if employerRes["1"] != nil {
		employerXdoc := employerRes["1"].(map[string]interface{})
		for cFieldname, iFieldvalue := range employerXdoc {
			xDocRequest["employer"+cFieldname] = iFieldvalue
		}
	}
	//Get Employer Info via EmployerControl

	formView := make(map[string]interface{})
	formView["employee-view"] = xDocRequest

	return strconv.Quote(string(this.Generate(formView, nil)))
}

func (this *Employee) edit(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if functions.TrimEscape(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{}`))
		return
	}

	tblEmployee := new(database.Profile)
	xDocRequest := make(map[string]interface{})
	xDocRequest["searchvalue"] = functions.TrimEscape(httpReq.FormValue("control"))

	ResMap := tblEmployee.Read(xDocRequest, curdb)
	if ResMap["1"] == nil {
		httpRes.Write([]byte(`{"error":"No Record Found"}`))
		return
	}
	xDocRequest = ResMap["1"].(map[string]interface{})

	xDocRequest[strings.Replace(xDocRequest["phonecode"].(string), "+", "", 1)] = "selected"

	xDocRequest["formtitle"] = "Edit"
	formView := make(map[string]interface{})
	formView["employee-edit"] = xDocRequest

	viewHTML := strconv.Quote(string(this.Generate(formView, nil)))
	httpRes.Write([]byte(`{"mainpanelContent":` + viewHTML + `}`))
}

func (this *Employee) invoicedownload(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Valued Subscription Invoice!")

	bytesBuffer := bytes.NewBuffer(nil)
	bytesBufferWriter := bufio.NewWriter(bytesBuffer)

	pdf.Output(bytesBufferWriter)
	bytesBufferWriter.Flush()

	httpRes.Header().Set("Content-Type", "application/pdf")
	httpRes.Header().Set("Content-Disposition", "attachment;filename=MyInvoice.pdf")
	httpRes.Write(bytesBuffer.Bytes())
}

func (this *Employee) requestInvoice(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		sMessage = "Please select a User<br>"
	}

	sScheme := ""
	if functions.TrimEscape(httpReq.FormValue("scheme")) == "" {
		sMessage += "Scheme is missing <br>"
	} else {
		sqlScheme := fmt.Sprintf(`select title, price from scheme where control = '%s'`, functions.TrimEscape(httpReq.FormValue("scheme")))
		resScheme, _ := curdb.Query(sqlScheme)
		if resScheme["1"] == nil {
			sMessage += "Selected Scheme not Found<br>"
		} else {
			sScheme = fmt.Sprintf(`%s - AED %v`, resScheme["1"].(map[string]interface{})["title"],
				resScheme["1"].(map[string]interface{})["price"])
		}
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	//SEND AN EMAIL USING TEMPLATE

	sqlEmployee := fmt.Sprintf(`select profile.firstname as firstname, profile.lastname as lastname, profile.email as email 
		from profile where profile.control in ('0'%s)`, controlList)
	resMember, _ := curdb.Query(sqlEmployee)

	if len(resMember) == 0 {
		sMessage += "Selected User not Found<br>"
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	var emailCC []string
	emailCC = append(emailCC, "employers@valued.com")
	emailTo := this.mapCache["email"].(string)

	emailFrom := "rewards@valued.com"
	emailFromName := "VALUED EMPLOYERS"
	emailTemplate := "employer-invoice"
	emailSubject := fmt.Sprintf("VALUED INVOICE - EMPLOYEE ADDITIONS")

	sTableRows := ""
	sRow := `<tr> <td>%v</td> <td>%v</td> <td>%v</td> <td>%v</td></tr>`

	emailFields := make(map[string]interface{})
	for _, xDoc := range resMember {
		xDoc := xDoc.(map[string]interface{})
		sTableRows += fmt.Sprintf(sRow, xDoc["firstnme"], xDoc["lastname"], xDoc["email"], sScheme)
	}

	emailFields["rows"] = sTableRows
	emailFields["employer"] = this.mapCache["title"]

	if this.mapCache["company"].(string) != "Yes" {
		emailFields["employer"] = this.mapCache["employertitle"]
		emailCC = append(emailCC, this.mapCache["employeremail"].(string))
	}

	go functions.GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate, emailCC, emailFields)
	//SEND AN EMAIL USING TEMPLATE

	sMessage = fmt.Sprintf(`{"error":"Invoice Request Email Has Been Sent"}`, len(resMember))
	httpRes.Write([]byte(sMessage))

}
