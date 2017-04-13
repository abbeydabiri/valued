package database

import (
	"fmt"
)

type Role struct {
	Crud
}

func (this *Role) GetName() string {
	return "role"
}

func (this *Role) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Role) GetFields() map[string]interface{} {

	fields := this.Fields()
	return fields
}

func (this *Role) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	xDoc := make(map[string]interface{})
	xDoc["code"] = "admin"
	xDoc["title"] = "Admin"
	xDoc["workflow"] = "active"
	this.Create("root", xDoc, curdb)

	xDoc = make(map[string]interface{})
	xDoc["code"] = "merchant"
	xDoc["title"] = "Merchant"
	xDoc["workflow"] = "active"
	this.Create("root", xDoc, curdb)

	xDoc = make(map[string]interface{})
	xDoc["code"] = "employer"
	xDoc["title"] = "Employer"
	xDoc["workflow"] = "active"
	this.Create("root", xDoc, curdb)

	xDoc = make(map[string]interface{})
	xDoc["code"] = "member"
	xDoc["title"] = "Member"
	xDoc["workflow"] = "active"
	this.Create("root", xDoc, curdb)
}

func (this *Role) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Role) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Role) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	role.code as code, role.control as control, role.title as title, 
					role.description as description, role.workflow as workflow, 
					role.createdate as createdate, role.createdby as createdby, 
					role.updatedate as updatedate, role.updatedby as updatedby
					

			from role

			where role.%s = '%s' order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Role) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	pagination := false
	if xDocrequest["pagination"] != nil {
		pagination = xDocrequest["pagination"].(bool)
	}

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	offset := "0"
	if xDocrequest["offset"] != nil && xDocrequest["offset"].(string) != "" {
		offset = xDocrequest["offset"].(string)
	}

	whereprofilecontrol := ""
	if xDocrequest["profilecontrol"] != nil && xDocrequest["profilecontrol"].(string) != "" {
		whereprofilecontrol = fmt.Sprintf(`role.profilecontrol = '%s' and`, xDocrequest["profilecontrol"])
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`role.workflow = '%s' and`, xDocrequest["workflow"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`role.code like '%%%s%%' and`, xDocrequest["code"])
	}

	sqlFields := `
					role.code as code, role.control as control, role.title as title, 
					role.description as description, role.workflow as workflow, 
					role.createdate as createdate, role.createdby as createdby, 
					role.updatedate as updatedate, role.updatedby as updatedby
					
					
				`
	if pagination {
		sqlFields = `count(role.control) as paginationtotal	`
	}

	sql := `select 	%s

			from role

			where %s %s %s role.control != '' order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wherecode, whereprofilecontrol, limit, offset)

	if xDocrequest["pagination"] != nil {

	} else {
		mapRes, _ = curdb.Query(sql)
	}
	return mapRes
}

func (this *Role) Delete(xDocrequest map[string]interface{}, curdb Database) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `delete from %s where %s = '%s'`

	sql = fmt.Sprintf(sql, this.GetName(), searchfield, searchvalue)
	curdb.Query(sql)
}
