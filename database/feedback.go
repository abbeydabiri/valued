package database

import (
	"fmt"
)

type Feedback struct {
	Crud
}

func (this *Feedback) GetName() string {
	return "feedback"
}

func (this *Feedback) GetWorkflows() []string {
	return []string{"draft", "inactive", "active", "approved", "rejected", "expired"}
}

func (this *Feedback) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["redemptioncontrol"] = ""
	fields["answer"] = ""
	return fields
}

func (this *Feedback) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Feedback) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Feedback) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Feedback) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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

					answer as answer, redemptioncontrol as redemptioncontrol

			from feedback

					where %s = '%s' order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Feedback) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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

	whereanswer := ""
	if xDocrequest["answer"] != nil && xDocrequest["answer"].(string) != "" {
		whereanswer = fmt.Sprintf(`answer like '%%%s%%' AND`, xDocrequest["answer"])
	}

	whereredemption := ""
	if xDocrequest["redemption"] != nil && xDocrequest["redemption"].(string) != "" {
		whereredemption = fmt.Sprintf(`redemptioncontrol = '%s' AND`, xDocrequest["redemption"])
	}

	sqlFields := `
					code as code, control as control, title as title, 
					description as description, workflow as workflow, 
					createdate as createdate, createdby as createdby, 
					updatedate as updatedate, updatedby as updatedby,

					answer as answer, redemptioncontrol as redemptioncontrol
				`
	if pagination {
		sqlFields = `count(control) as paginationtotal	`
	}

	sql := `select %s

			from feedback

			where %s %s  %s %s control != '' order by control desc limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, whereredemption, whereanswer, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Feedback) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
