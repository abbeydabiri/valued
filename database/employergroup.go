package database

import (
	"fmt"
)

type EmployerGroup struct {
	Crud
}

func (this *EmployerGroup) GetName() string {
	return "employergroup"
}

func (this *EmployerGroup) GetWorkflows() []string {
	return []string{"draft", "inactive", "active", "approved", "rejected", "expired"}
}

func (this *EmployerGroup) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["groupcontrol"] = ""
	fields["employercontrol"] = ""

	return fields
}

func (this *EmployerGroup) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *EmployerGroup) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *EmployerGroup) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *EmployerGroup) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	employergroup.code as code, employergroup.control as control, employergroup.title as title,
					employergroup.description as description, employergroup.workflow as workflow,
					employergroup.createdate as createdate, employergroup.createdby as createdby,
					employergroup.updatedate as updatedate, employergroup.updatedby as updatedby,

					employergroup.groupcontrol as groupcontrol, employergroup.employercontrol as employercontrol
					
			from employergroup

					where employergroup.%s = '%s'

			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *EmployerGroup) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	offset := "0"
	if xDocrequest["offset"] != nil && xDocrequest["offset"].(string) != "" {
		offset = xDocrequest["offset"].(string)
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`employergroup.code = '%s' AND`, xDocrequest["code"])
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`employergroup.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	whereemployer := ""
	if xDocrequest["employer"] != nil && xDocrequest["employer"].(string) != "" {
		whereemployer = fmt.Sprintf(`employergroup.employercontrol = '%s' AND`, xDocrequest["employer"])
	}

	sql := `select 	employergroup.code as code, employergroup.control as control, employergroup.title as title,
					employergroup.description as description, employergroup.workflow as workflow,
					employergroup.createdate as createdate, employergroup.createdby as createdby,
					employergroup.updatedate as updatedate, employergroup.updatedby as updatedby,

					employergroup.groupcontrol as groupcontrol, employergroup.employercontrol as employercontrol

			from employergroup

			where %s %s %s employergroup.control != ''

			order by employergroup.control limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, whereemployer, wherecode, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *EmployerGroup) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
