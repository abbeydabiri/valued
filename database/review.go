package database

import (
	"fmt"
)

type Review struct {
	Crud
}

func (this *Review) GetName() string {
	return "review"
}

func (this *Review) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Review) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["reviewcategorycontrol"] = ""
	fields["redemptioncontrol"] = ""
	fields["merchantcontrol"] = ""
	fields["employercontrol"] = ""
	fields["schemecontrol"] = ""
	fields["membercontrol"] = ""
	fields["storecontrol"] = ""
	fields["rating"] = ""
	return fields
}

func (this *Review) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Review) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Review) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Review) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	code as code, control as control, title as title, 
					description as description, workflow as workflow, 
					createdate as createdate, createdby as createdby, 
					updatedate as updatedate, updatedby as updatedby,
					image as image, placement as placement

			from review

					where %s = '%s' order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Review) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`title like '%%%s%%' AND`, xDocrequest["title"])
	}

	sqlFields := `
					code as code, control as control, title as title, 
					description as description, workflow as workflow, 
					createdate as createdate, createdby as createdby, 
					updatedate as updatedate, updatedby as updatedby,
					image as image, placement as placement
				`
	if pagination {
		sqlFields = `count(control) as paginationtotal	`
	}

	sql := `select 	%s

			from review

			where %s %s control != '' order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Review) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
