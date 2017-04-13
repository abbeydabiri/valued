package database

import (
	"fmt"
)

type Media struct {
	Crud
}

func (this *Media) GetName() string {
	return "media"
}

func (this *Media) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Media) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["public"] = ""
	fields["merchantcontrol"] = ""
	return fields
}

func (this *Media) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Media) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Media) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Media) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	media.code as code, media.control as control, media.title as title, 
					media.description as description, media.workflow as workflow, 
					media.createdate as createdate, media.createdby as createdby, 
					media.updatedate as updatedate, media.updatedby as updatedby,
				
					media.public as public,
					merchant.control as merchantcontrol, merchant.title as merchanttitle
					
			from media, profile as merchant
			
			where media.%s = '%s' AND 
			merchant.control = media.merchantcontrol 
			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Media) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`media.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`media.title like '%%%s%%' AND`, xDocrequest["title"])
	}

	wherepublic := ""
	if xDocrequest["public"] != nil && xDocrequest["public"].(string) != "" {
		wherepublic = fmt.Sprintf(`media.public like '%%%s%%' AND`, xDocrequest["public"])
	}

	wheresource := ""
	if xDocrequest["source"] != nil && xDocrequest["source"].(string) != "" {
		wheresource = fmt.Sprintf(`media.sourcecontrol like '%%%s%%' AND`, xDocrequest["source"])
	}

	wheresourcecontrol := ""
	if xDocrequest["sourcecontrol"] != nil && xDocrequest["sourcecontrol"].(string) != "" {
		wheresourcecontrol = fmt.Sprintf(`media.sourcecontrolcontrol like '%%%s%%' AND`, xDocrequest["sourcecontrol"])
	}

	sql := `select 	media.code as code, media.control as control, media.title as title, 
					media.description as description, media.workflow as workflow, 
					media.createdate as createdate, media.createdby as createdby, 
					media.updatedate as updatedate, media.updatedby as updatedby,
				
					media.public as public,	media.source as source, 
					media.sourcecontrol as sourcecontrol
					
			from media
			
			where %s %s %s %s %s

			order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wheretitle, wherepublic, wheresource, wheresourcecontrol, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Media) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
