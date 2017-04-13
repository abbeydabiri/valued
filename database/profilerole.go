package database

import (
	"fmt"
)

type ProfileRole struct {
	Crud
}

func (this *ProfileRole) GetName() string {
	return "profilerole"
}

func (this *ProfileRole) GetWorkflows() []string {
	return []string{"draft", "inactive", "active", "approved", "rejected", "expired"}
}

func (this *ProfileRole) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["rolecontrol"] = ""
	fields["profilecontrol"] = ""

	return fields
}

func (this *ProfileRole) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	xDoc := make(map[string]interface{})
	xDoc["workflow"] = "active"
	xDoc["rolecontrol"] = "1.00000000001"
	xDoc["profilecontrol"] = "1.00000000001"
	this.Create("root", xDoc, curdb)

	xDoc = make(map[string]interface{})
	xDoc["workflow"] = "active"
	xDoc["rolecontrol"] = "1.00000000002"
	xDoc["profilecontrol"] = "1.00000000001"
	this.Create("root", xDoc, curdb)

	xDoc = make(map[string]interface{})
	xDoc["workflow"] = "active"
	xDoc["rolecontrol"] = "1.00000000003"
	xDoc["profilecontrol"] = "1.00000000001"
	this.Create("root", xDoc, curdb)

}

func (this *ProfileRole) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ProfileRole) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ProfileRole) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	profilerole.code as code, profilerole.control as control, profilerole.title as title,
					profilerole.description as description, profilerole.workflow as workflow,
					profilerole.createdate as createdate, profilerole.createdby as createdby,
					profilerole.updatedate as updatedate, profilerole.updatedby as updatedby,

					profilerole.rolecontrol as rolecontrol, profilerole.profilecontrol as profilecontrol
					
			from profilerole

					where profilerole.%s = '%s'

			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ProfileRole) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	offset := "0"
	if xDocrequest["offset"] != nil && xDocrequest["offset"].(string) != "" {
		offset = xDocrequest["offset"].(string)
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`profilerole.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	whereprofile := ""
	if xDocrequest["profile"] != nil && xDocrequest["profile"].(string) != "" {
		whereprofile = fmt.Sprintf(`profilerole.profilecontrol = '%s' AND`, xDocrequest["profile"])
	}

	sql := `select 	profilerole.code as code, profilerole.control as control, profilerole.title as title,
					profilerole.description as description, profilerole.workflow as workflow,
					profilerole.createdate as createdate, profilerole.createdby as createdby,
					profilerole.updatedate as updatedate, profilerole.updatedby as updatedby,

					profilerole.rolecontrol as rolecontrol, profilerole.profilecontrol as profilecontrol

			from profilerole

			where %s %s control != ''

			order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, whereprofile, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ProfileRole) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
