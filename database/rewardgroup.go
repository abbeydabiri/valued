package database

import (
	"fmt"
)

type RewardGroup struct {
	Crud
}

func (this *RewardGroup) GetName() string {
	return "rewardgroup"
}

func (this *RewardGroup) GetWorkflows() []string {
	return []string{"draft", "inactive", "active", "approved", "rejected", "expired"}
}

func (this *RewardGroup) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["groupcontrol"] = ""
	fields["rewardcontrol"] = ""

	return fields
}

func (this *RewardGroup) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *RewardGroup) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *RewardGroup) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *RewardGroup) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	rewardgroup.code as code, rewardgroup.control as control, rewardgroup.title as title,
					rewardgroup.description as description, rewardgroup.workflow as workflow,
					rewardgroup.createdate as createdate, rewardgroup.createdby as createdby,
					rewardgroup.updatedate as updatedate, rewardgroup.updatedby as updatedby,

					rewardgroup.groupcontrol as groupcontrol, rewardgroup.rewardcontrol as rewardcontrol
					
			from rewardgroup

					where rewardgroup.%s = '%s'

			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *RewardGroup) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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

	wherereward := ""
	if xDocrequest["reward"] != nil && xDocrequest["reward"].(string) != "" {
		wherereward = fmt.Sprintf(`rewardgroup.rewardcontrol = '%s' AND`, xDocrequest["reward"])
	}

	sql := `select 	rewardgroup.code as code, rewardgroup.control as control, rewardgroup.title as title,
					rewardgroup.description as description, rewardgroup.workflow as workflow,
					rewardgroup.createdate as createdate, rewardgroup.createdby as createdby,
					rewardgroup.updatedate as updatedate, rewardgroup.updatedby as updatedby,

					rewardgroup.groupcontrol as groupcontrol, rewardgroup.rewardcontrol as rewardcontrol

			from rewardgroup

			where %s %s control != ''

			order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wherereward, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *RewardGroup) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
