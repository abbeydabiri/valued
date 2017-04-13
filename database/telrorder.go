package database

import (
	"fmt"
)

type TelrOrder struct {
	Crud
}

func (this *TelrOrder) GetName() string {
	return "telrorder"
}

func (this *TelrOrder) GetWorkflows() []string {
	return []string{"pending", "authorised", "declined", "cancelled", "failed"}
}

func (this *TelrOrder) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["logincontrol"] = ""
	fields["profilecontrol"] = ""
	fields["schemecontrol"] = ""
	fields["subscription"] = ""
	fields["couponcontrol"] = ""

	fields["telrtext"] = ""
	fields["telrcode"] = 0
	fields["telrref"] = ""
	fields["telrurl"] = ""

	fields["amount"] = 0.0
	fields["currency"] = ""
	fields["createrequest"] = ""
	fields["createresponse"] = ""
	fields["checkrequest"] = ""
	fields["checkresponse"] = ""

	return fields
}

func (this *TelrOrder) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *TelrOrder) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *TelrOrder) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *TelrOrder) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	telrorder.code as code, telrorder.control as control, telrorder.title as title, 
					telrorder.description as description, telrorder.workflow as workflow, 
					telrorder.createdate as createdate, telrorder.createdby as createdby, 
					telrorder.updatedate as updatedate, telrorder.updatedby as updatedby,

					telrorder.telrcode as telrcode, telrorder.telrtext as telrtext, 
					telrorder.telrref as telrref, telrorder.telrurl as telrurl,	


					telrorder.amount as amount, telrorder.currency as currency, 
					telrorder.createrequest as createrequest,
					telrorder.createresponse as createresponse,
					telrorder.checkrequest as checkrequest,
					telrorder.checkresponse as checkresponse,

					telrorder.logincontrol as logincontrol, telrorder.profilecontrol as profilecontrol,
					telrorder.couponcontrol as couponcontrol, telrorder.schemecontrol as schemecontrol,
					telrorder.subscriptioncontrol as subscriptioncontrol,

					scheme.title as schemetitle


			from telrorder, scheme

			where %s = '%s' AND	scheme.control = telrorder.schemecontrol

			order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *TelrOrder) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`telrorder.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheretelrref := ""
	if xDocrequest["telrref"] != nil && xDocrequest["telrref"].(string) != "" {
		wheretelrref = fmt.Sprintf(`telrorder.telrref = '%s' AND`, xDocrequest["telrref"])
	}

	whereprofile := ""
	if xDocrequest["profile"] != nil && xDocrequest["profile"].(string) != "" {
		whereprofile = fmt.Sprintf(`telrorder.profilecontrol = '%s' AND`, xDocrequest["profile"])
	}

	wherecontrol := ""
	if xDocrequest["control"] != nil && xDocrequest["control"].(string) != "" {
		wherecontrol = fmt.Sprintf(`telrorder.control = '%s' AND`, xDocrequest["control"])
	}

	sql := `select 	telrorder.code as code, telrorder.control as control, telrorder.title as title, 
					telrorder.description as description, telrorder.workflow as workflow, 
					telrorder.createdate as createdate, telrorder.createdby as createdby, 
					telrorder.updatedate as updatedate, telrorder.updatedby as updatedby,

					telrorder.telrcode as telrcode, telrorder.telrtext as telrtext, 
					telrorder.telrref as telrref, telrorder.telrurl as telrurl,	


					telrorder.amount as amount, telrorder.currency as currency, 
					telrorder.createrequest as createrequest,
					telrorder.createresponse as createresponse,
					telrorder.checkrequest as checkrequest,
					telrorder.checkresponse as checkresponse,

					telrorder.logincontrol as logincontrol, telrorder.profilecontrol as profilecontrol,
					telrorder.couponcontrol as couponcontrol, telrorder.schemecontrol as schemecontrol,
					telrorder.subscriptioncontrol as subscriptioncontrol,

					scheme.title as schemetitle


			from telrorder, scheme where %s %s %s %s scheme.control = telrorder.schemecontrol

			order by control desc limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wheretelrref, whereprofile, wherecontrol, limit, offset)
	mapRes, _ = curdb.Query(sql)

	return mapRes
}

func (this *TelrOrder) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
