package database

import (
	"fmt"
)

type Industry struct {
	Crud
}

func (this *Industry) GetName() string {
	return "industry"
}

func (this *Industry) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Industry) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["image"] = ""
	fields["placement"] = ""
	fields["industrycontrol"] = ""
	return fields
}

func (this *Industry) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	xDoc := make(map[string]interface{})
	xDoc["code"] = "main"
	xDoc["title"] = "None"
	xDoc["image"] = "None"
	xDoc["placement"] = "None"
	xDoc["industrycontrol"] = ""
	xDoc["workflow"] = "active"
	this.Create("root", xDoc, curdb)
}

func (this *Industry) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Industry) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Industry) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	industry.code as code, industry.control as control, industry.title as title, 
					industry.description as description, industry.workflow as workflow, 
					industry.createdate as createdate, industry.createdby as createdby, 
					industry.updatedate as updatedate, industry.updatedby as updatedby,
					industry.image as image, industry.placement as placement,

					industry.industrycontrol as industrycontrol,
					parentindustry.title as industrytitle	

			from industry left join industry as parentindustry
			on parentindustry.control = industry.industrycontrol

			where industry.%s = '%s' order by placement`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Industry) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`industry.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`industry.title like '%%%s%%' AND`, xDocrequest["title"])
	}

	whereindustry := ""
	if xDocrequest["industry"] != nil {
		whereindustry = fmt.Sprintf(`industry.industrycontrol = '%s' AND`, xDocrequest["industry"])
	}

	wheresub := ""
	if xDocrequest["sub"] != nil {
		switch xDocrequest["sub"].(string) {
		case "true":
			wheresub = `industry.industrycontrol != '' AND`
		default:
			wheresub = `industry.industrycontrol = '' AND`
		}
	}

	sql := `select 	industry.code as code, industry.control as control, industry.title as title, 
					industry.description as description, industry.workflow as workflow, 
					industry.createdate as createdate, industry.createdby as createdby, 
					industry.updatedate as updatedate, industry.updatedby as updatedby,
					industry.image as image, industry.placement as placement,

					industry.industrycontrol as industrycontrol,
					parentindustry.title as industrytitle	

			from industry left join industry as parentindustry
			on parentindustry.control = industry.industrycontrol

			where %s %s %s %s industry.control != '' order by placement limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wheretitle, whereindustry, wheresub, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Industry) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
