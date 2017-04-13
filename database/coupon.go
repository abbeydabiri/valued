package database

import (
	"fmt"
)

type Coupon struct {
	Crud
}

func (this *Coupon) GetName() string {
	return "coupon"
}

func (this *Coupon) GetWorkflows() []string {
	return []string{"active", "inactive", "approved", "rejected"}
}

func (this *Coupon) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["rewardcontrol"] = ""
	return fields
}

func (this *Coupon) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Coupon) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Coupon) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Coupon) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	coupon.code as code, coupon.control as control, coupon.title as title,
					coupon.description as description, coupon.workflow as workflow,
					coupon.createdate as createdate, coupon.createdby as createdby,
					coupon.updatedate as updatedate, coupon.updatedby as updatedby,

					coupon.rewardcontrol as rewardcontrol

			from coupon left join reward on reward.control = coupon.rewardcontrol

			where coupon.%s = '%s' order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Coupon) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`coupon.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`coupon.code like '%%%s%%' AND`, xDocrequest["code"])
	}

	wherereward := ""
	if xDocrequest["reward"] != nil && xDocrequest["reward"].(string) != "" {
		wherereward = fmt.Sprintf(`coupon.rewardcontrol = '%s' AND`, xDocrequest["reward"])
	}

	sqlFields := `
					coupon.code as code, coupon.control as control, coupon.title as title,
					coupon.description as description, coupon.workflow as workflow,
					coupon.createdate as createdate, coupon.createdby as createdby,
					coupon.updatedate as updatedate, coupon.updatedby as updatedby,

					coupon.rewardcontrol as rewardcontrol
				`

	if pagination {
		sqlFields = `count(coupon.control) as paginationtotal	`
	}

	sql := `select 	%s

			from coupon left join coupon as parentcoupon
			on parentcoupon.control = coupon.rewardcontrol

			where %s %s %s coupon.control != '' order by control desc limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wherecode, wherereward, limit, offset)
	mapRes, _ = curdb.Query(sql)

	return mapRes
}

func (this *Coupon) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
