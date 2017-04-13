package database

import (
	"fmt"
)

type Redemption struct {
	Crud
}

func (this *Redemption) GetName() string {
	return "redemption"
}

func (this *Redemption) GetWorkflows() []string {
	return []string{"draft", "inactive", "active", "approved", "rejected", "expired"}
}

func (this *Redemption) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["merchantcontrol"] = ""
	fields["employercontrol"] = ""
	fields["membercontrol"] = ""

	fields["schemecontrol"] = ""
	fields["rewardcontrol"] = ""
	fields["couponcontrol"] = ""
	fields["storecontrol"] = ""

	fields["dob"] = ""
	fields["reward"] = ""
	fields["gender"] = ""
	fields["location"] = ""
	fields["discount"] = ""
	fields["nationality"] = ""

	fields["savingsvalue"] = 0.00
	fields["transactionvalue"] = 0.00

	return fields
}

func (this *Redemption) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Redemption) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Redemption) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Redemption) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	redemption.code as code, redemption.control as control, redemption.title as title,
					redemption.description as description, redemption.workflow as workflow,
					redemption.createdate as createdate, redemption.createdby as createdby,
					redemption.updatedate as updatedate, redemption.updatedby as updatedby,

					redemption.savingsvalue as savingsvalue, redemption.transactionvalue as transactionvalue,

					redemption.merchantcontrol as merchantcontrol, redemption.employercontrol as employercontrol,
					redemption.schemecontrol as schemecontrol,
					redemption.membercontrol as membercontrol, redemption.rewardcontrol as rewardcontrol,
					redemption.couponcontrol as couponcontrol, redemption.storecontrol as storecontrol,

					redemption.dob as dob, redemption.reward as reward, redemption.gender as gender,
					redemption.discount as discount, redemption.nationality as nationality

			from redemption

					where redemption.%s = '%s'

			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Redemption) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
					redemption.code as code, redemption.control as control, redemption.title as title,
					redemption.description as description, redemption.workflow as workflow,
					redemption.createdate as createdate, redemption.createdby as createdby,
					redemption.updatedate as updatedate, redemption.updatedby as updatedby,

					redemption.savingsvalue as savingsvalue, redemption.transactionvalue as transactionvalue,

					redemption.merchantcontrol as merchantcontrol, redemption.employercontrol as employercontrol,
					redemption.schemecontrol as schemecontrol,
					redemption.membercontrol as membercontrol, redemption.rewardcontrol as rewardcontrol,
					redemption.couponcontrol as couponcontrol, redemption.storecontrol as storecontrol,

					redemption.dob as dob, redemption.reward as reward, redemption.gender as gender,
					redemption.discount as discount, redemption.nationality as nationality
				`
	if pagination {
		sqlFields = `count(redemption.control) as paginationtotal	`
	}

	sql := `select 	%s

			from redemption

			where %s %s control != ''

			order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Redemption) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
