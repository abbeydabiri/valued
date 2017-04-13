package database

import (
	"fmt"
)

type ReviewCategory struct {
	Crud
}

func (this *ReviewCategory) GetName() string {
	return "reviewcategory"
}

func (this *ReviewCategory) GetWorkflows() []string {
	return this.Workflows()
}

func (this *ReviewCategory) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["image"] = ""
	fields["placement"] = ""
	return fields
}

func (this *ReviewCategory) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *ReviewCategory) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ReviewCategory) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ReviewCategory) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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

			from reviewcategory

					where %s = '%s' order by placement`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ReviewCategory) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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

			from reviewcategory

			where %s %s control != '' order by placement limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ReviewCategory) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
