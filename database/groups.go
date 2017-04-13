package database

import (
	"fmt"
)

type Groups struct {
	Crud
}

func (this *Groups) GetName() string {
	return "groups"
}

func (this *Groups) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Groups) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["groupscontrol"] = ""
	return fields
}

func (this *Groups) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Groups) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Groups) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Groups) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	groups.code as code, groups.control as control, groups.title as title, 
					groups.description as description, groups.workflow as workflow, 
					groups.createdate as createdate, groups.createdby as createdby, 
					groups.updatedate as updatedate, groups.updatedby as updatedby,
					
					groups.groupscontrol as groupscontrol,
					parentgroups.title as groupstitle	

			from groups left join groups as parentgroups
			on parentgroups.control = groups.groupscontrol

			where groups.%s = '%s' order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Groups) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`groups.workflow = '%s' and`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`groups.title like '%%%s%%' and`, xDocrequest["title"])
	}

	wheregroups := ""
	if xDocrequest["groups"] != nil {
		wheregroups = fmt.Sprintf(`groups.groupscontrol = '%s' and`, xDocrequest["groups"])
	}

	wheresub := ""
	if xDocrequest["sub"] != nil {
		switch xDocrequest["sub"].(string) {
		case "true":
			wheresub = `groups.groupscontrol != '' and`
		default:
			wheresub = `groups.groupscontrol = '' and`
		}
	}

	sqlFields := `
					groups.code as code, groups.control as control, groups.title as title, 
					groups.description as description, groups.workflow as workflow, 
					groups.createdate as createdate, groups.createdby as createdby, 
					groups.updatedate as updatedate, groups.updatedby as updatedby,
					
					groups.groupscontrol as groupscontrol,
					parentgroups.title as groupstitle	
				`
	if pagination {
		sqlFields = `count(groups.control) as paginationtotal	`
	}

	sql := `select 	%s

			from groups left join groups as parentgroups
			on parentgroups.control = groups.groupscontrol

			where %s %s %s %s groups.control != '' order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, wheregroups, wheresub, limit, offset)

	if xDocrequest["pagination"] != nil {

	} else {
		mapRes, _ = curdb.Query(sql)
	}
	return mapRes
}

func (this *Groups) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
