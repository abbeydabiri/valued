package database

import (
	"fmt"
)

type Category struct {
	Crud
}

func (this *Category) GetName() string {
	return "category"
}

func (this *Category) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Category) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["image"] = ""
	fields["placement"] = ""
	fields["categorycontrol"] = ""
	return fields
}

func (this *Category) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	xDoc := make(map[string]interface{})
	xDoc["code"] = "main"
	xDoc["title"] = "None"
	xDoc["image"] = "None"
	xDoc["placement"] = "None"
	xDoc["categorycontrol"] = ""
	xDoc["workflow"] = "active"
	this.Create("root", xDoc, curdb)
}

func (this *Category) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Category) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Category) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	category.code as code, category.control as control, category.title as title, 
					category.description as description, category.workflow as workflow, 
					category.createdate as createdate, category.createdby as createdby, 
					category.updatedate as updatedate, category.updatedby as updatedby,
					category.image as image, category.placement as placement,

					category.categorycontrol as categorycontrol,
					parentcategory.title as categorytitle	

			from category left join category as parentcategory
			on parentcategory.control = category.categorycontrol

			where category.%s = '%s' order by placement`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Category) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`category.workflow = '%s' and`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`lower(category.title) like lower('%%%s%%') and`, xDocrequest["title"])
	}

	wherecategory := ""
	if xDocrequest["category"] != nil {
		wherecategory = fmt.Sprintf(`category.categorycontrol = '%s' and`, xDocrequest["category"])
	}

	wheresub := ""
	if xDocrequest["sub"] != nil {
		switch xDocrequest["sub"].(string) {
		case "true":
			wheresub = `category.categorycontrol != '' and`
		default:
			wheresub = `category.categorycontrol = '' and`
		}
	}

	sqlFields := `
					category.code as code, category.control as control, category.title as title, 
					category.description as description, category.workflow as workflow, 
					category.createdate as createdate, category.createdby as createdby, 
					category.updatedate as updatedate, category.updatedby as updatedby,
					category.image as image, category.placement as placement,

					category.categorycontrol as categorycontrol,
					parentcategory.title as categorytitle
				`

	if pagination {
		sqlFields = `count(category.control) as paginationtotal	`
	}

	sql := `select 	%s	

			from category left join category as parentcategory
			on parentcategory.control = category.categorycontrol

			where %s %s %s %s category.control != '' order by categorycontrol, placement limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, wherecategory, wheresub, limit, offset)
	mapRes, _ = curdb.Query(sql)

	return mapRes
}

func (this *Category) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
