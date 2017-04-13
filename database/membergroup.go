package database

import (
	"fmt"
)

type MemberGroup struct {
	Crud
}

func (this *MemberGroup) GetName() string {
	return "membergroup"
}

func (this *MemberGroup) GetWorkflows() []string {
	return []string{"draft", "inactive", "active", "approved", "rejected", "expired"}
}

func (this *MemberGroup) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["groupcontrol"] = ""
	fields["membercontrol"] = ""

	return fields
}

func (this *MemberGroup) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *MemberGroup) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *MemberGroup) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *MemberGroup) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	membergroup.code as code, membergroup.control as control, membergroup.title as title,
					membergroup.description as description, membergroup.workflow as workflow,
					membergroup.createdate as createdate, membergroup.createdby as createdby,
					membergroup.updatedate as updatedate, membergroup.updatedby as updatedby,

					membergroup.groupcontrol as groupcontrol, membergroup.membercontrol as membercontrol
					
			from membergroup

					where membergroup.%s = '%s'

			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *MemberGroup) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`membergroup.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheremember := ""
	if xDocrequest["member"] != nil && xDocrequest["member"].(string) != "" {
		wheremember = fmt.Sprintf(`membergroup.membercontrol = '%s' AND`, xDocrequest["member"])
	}

	sql := `select 	membergroup.code as code, membergroup.control as control, membergroup.title as title,
					membergroup.description as description, membergroup.workflow as workflow,
					membergroup.createdate as createdate, membergroup.createdby as createdby,
					membergroup.updatedate as updatedate, membergroup.updatedby as updatedby,

					membergroup.groupcontrol as groupcontrol, membergroup.membercontrol as membercontrol

			from membergroup

			where %s %s control != ''

			order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wheremember, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *MemberGroup) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
