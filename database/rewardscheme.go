package database

import (
	"fmt"
)

type RewardScheme struct {
	Crud
}

func (this *RewardScheme) GetName() string {
	return "rewardscheme"
}

func (this *RewardScheme) GetWorkflows() []string {
	return this.Workflows()
}

func (this *RewardScheme) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["rewardcontrol"] = ""
	fields["schemecontrol"] = ""
	return fields
}

func (this *RewardScheme) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *RewardScheme) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *RewardScheme) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *RewardScheme) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	rewardscheme.code as code, rewardscheme.control as control, rewardscheme.title as title, 
					rewardscheme.description as description, rewardscheme.workflow as workflow, 
					rewardscheme.createdate as createdate, rewardscheme.createdby as createdby, 
					rewardscheme.updatedate as updatedate, rewardscheme.updatedby as updatedby,
					
						
					rewardscheme.schemecontrol as schemecontrol, scheme.title as schemetitle, 
					scheme.code as schemecode, scheme.price as schemeprice,

					rewardscheme.rewardcontrol as rewardcontrol,
					reward.title as rewardtitle

			from rewardscheme, reward, scheme

			where rewardscheme.%s = '%s' and

			reward.control = rewardscheme.rewardcontrol 
			AND scheme.control = rewardscheme.schemecontrol

			order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *RewardScheme) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`rewardscheme.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherereward := ""
	if xDocrequest["reward"] != nil && xDocrequest["reward"].(string) != "" {
		wherereward = fmt.Sprintf(`rewardscheme.rewardcontrol like '%%%s%%' AND`, xDocrequest["reward"])
	}

	wherescheme := ""
	if xDocrequest["scheme"] != nil && xDocrequest["scheme"].(string) != "" {
		wherescheme = fmt.Sprintf(`rewardscheme.schemecontrol like '%%%s%%' AND`, xDocrequest["scheme"])
	}

	sql := `select 	rewardscheme.code as code, rewardscheme.control as control, rewardscheme.title as title, 
					rewardscheme.description as description, rewardscheme.workflow as workflow, 
					rewardscheme.createdate as createdate, rewardscheme.createdby as createdby, 
					rewardscheme.updatedate as updatedate, rewardscheme.updatedby as updatedby,
					
						
					rewardscheme.schemecontrol as schemecontrol, scheme.title as schemetitle, 
					scheme.code as schemecode, scheme.price as schemeprice,

					rewardscheme.rewardcontrol as rewardcontrol,
					reward.title as rewardtitle

			from rewardscheme, reward, scheme

			where %s %s %s 

			reward.control = rewardscheme.rewardcontrol 
			AND scheme.control = rewardscheme.schemecontrol

			order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wherereward, wherescheme, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *RewardScheme) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
