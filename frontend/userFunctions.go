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

func (this *User) viewPermission(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("login")) == "" {
		sMessage += "User is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	subSearch := make(map[string]interface{})
	subSearchPermission := make(map[string]interface{})
	subSearchPermission["login"] = html.EscapeString(httpReq.FormValue("login"))
	subSearch["userview-permission-search"] = subSearchPermission

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))

	httpRes.Write([]byte(`{"subview":` + viewHTML + `}`))
}

func (this *User) searchPermission(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("login")) == "" {
		sMessage += "User is missing <br>"
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

	tblUserPermission := new(database.Permission)
	subSearch := make(map[string]interface{})
	xDocrequest["login"] = html.EscapeString(httpReq.FormValue("login"))
	xDocresult := tblUserPermission.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		switch xDoc["workflow"].(string) {
		case "inactive":
			xDoc["action"] = "linkedPermissionActivate"
			xDoc["actionColor"] = "success"
			xDoc["actionLabel"] = "Activate"

		case "active":
			xDoc["action"] = "linkedPermissionDeactivate"
			xDoc["actionColor"] = "danger"
			xDoc["actionLabel"] = "De-Activate"
		}

		subSearch[cNumber+"#userview-permission-search-result"] = xDoc
	}

	viewHTML := strconv.Quote(string(this.Generate(subSearch, nil)))
	httpRes.Write([]byte(`{"subsearchresult":` + viewHTML + `}`))
}

func (this *User) linkPermission(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("login")) == "" {
		sMessage += "User is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("store")) == "" {
		sMessage += "Permission is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("merchant")) == "" {
		sMessage += "Merchant is missing <br>"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblUserPermission := new(database.Permission)
	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "inactive"
	xDoc["logincontrol"] = html.EscapeString(httpReq.FormValue("login"))
	xDoc["storecontrol"] = html.EscapeString(httpReq.FormValue("store"))

	if html.EscapeString(httpReq.FormValue("control")) != "" {
		xDoc["control"] = html.EscapeString(httpReq.FormValue("control"))
		tblUserPermission.Update(this.mapCache["username"].(string), xDoc, curdb)
	} else {
		xDoc["control"] = tblUserPermission.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	this.searchPermission(httpRes, httpReq, curdb)
}

func (this *User) linkedPermissionDelete(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a permission"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`delete from loginstore where control = '%s'`,
		html.EscapeString(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Permissiond Deleted","triggerSubSearch":true}`))
}

func (this *User) linkedPermissionActivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a permission"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update loginstore set workflow = 'active', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Permissiond Activated","triggerSubSearch":true}`))
}

func (this *User) linkedPermissionDeactivate(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	if html.EscapeString(httpReq.FormValue("control")) == "" {
		httpRes.Write([]byte(`{"error":"Please select a permission"}`))
		return
	}

	curdb.Query(fmt.Sprintf(`update loginstore set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control = '%s'`,
		this.mapCache["username"], functions.GetSystemTime(), functions.TrimEscape(httpReq.FormValue("control"))))

	httpRes.Write([]byte(`{"error":"Linked Permissiond Deactivated","triggerSubSearch":true}`))
}

func (this *User) linkedPermissionDeactivateAll(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	controlList := ""
	httpReq.ParseForm()
	for _, control := range httpReq.Form["control"] {
		controlList = fmt.Sprintf(`%s,'%s'`, controlList, control)
	}

	if controlList == "" {
		httpRes.Write([]byte(`{"error":"Please select a permission"}`))
		return
	}

	sqlQuery := fmt.Sprintf(`update loginstore set workflow = 'inactive', updatedby = '%s', updatedate = '%s' where control in ('0'%s)`,
		this.mapCache["username"], functions.GetSystemTime(), controlList)
	curdb.Query(sqlQuery)

	httpRes.Write([]byte(`{"error":"Linked Permissions Deactivated","triggerSubSearch":true}`))
}
