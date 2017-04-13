package database

import (
	"fmt"
)

type Invoice struct {
	Crud
}

func (this *Invoice) GetName() string {
	return "invoice"
}

func (this *Invoice) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Invoice) GetFields() map[string]interface{} {

	fields := this.Fields()

	fields["docdate"] = ""
	fields["profilecontrol"] = ""

	fields["shiptocode"] = ""
	fields["tocontactcontrol"] = ""
	fields["toprofilecontrol"] = ""

	fields["shipfromcode"] = ""
	fields["fromcontactcontrol"] = ""
	fields["fromprofilecontrol"] = ""

	fields["validfrom"] = ""
	fields["validtill"] = ""
	fields["deliverydate"] = ""
	fields["shippingdate"] = ""
	fields["shippingvehicle"] = ""
	fields["latestshippingdate"] = ""
	fields["latestdeliverydate"] = ""

	fields["discountname"] = ""
	fields["discountcode"] = ""
	fields["discountamount"] = 0.0
	fields["discountpercent"] = 0.0

	fields["taxamount"] = 0.0
	fields["taxpercent"] = 0.0

	fields["totalexcltax"] = 0.0
	fields["totalincltax"] = 0.0

	fields["currency"] = ""
	fields["currencyrate"] = 0.0

	return fields
}

func (this *Invoice) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Invoice) Create(cUsername string, lMain bool, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc("system", this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Invoice) Update(cUsername string, lMain bool, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Invoice) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {
	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = "%"
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	lMain := true
	if xDocrequest["main"] != nil {
		lMain = xDocrequest["main"].(bool)
	}

	tablename := this.GetName()
	if !lMain {
		tablename = tablename + "audit"
	}

	sql := `select 	xdoc.control, xdoc.title, xdoc.code, xdoc.description, 
					
					xdoc.docno, xdoc.docdate, xdoc.shiptocode, xdoc.shipfromcode,
					xdoc.validfrom, xdoc.validtill, xdoc.deliverydate, xdoc.shippingdate, 
					xdoc.shippingvehicle, xdoc.latestshippingdate, xdoc.latestdeliverydate,
					xdoc.discountname, xdoc.discountcode, xdoc.discountamount, xdoc.discountpercent,
					xdoc.taxamount, xdoc.taxpercent, xdoc.totalexcltax, xdoc.totalincltax,
					xdoc.currency, xdoc.currencyrate, 

					xdoc.profilecontrol, profile.title  as profile, profile.code  as profilecode,
					xdoc.toprofilecontrol, toprofile.title  as toprofile, toprofile.code  as toprofilecode,
					xdoc.fromprofilecontrol, fromprofile.title  as fromprofile, fromprofile.code  as fromprofilecode,

					xdoc.tocontactcontrol, tocontact.title  as tocontact, tocontact.code  as tocontactcode,
					xdoc.fromcontactcontrol, fromcontact.title  as fromcontact, fromcontact.code  as fromcontactcode,
					
					xdoc.typecontrol, types.title  as type, types.code as typecode,
					xdoc.groupcontrol, groups.title  as group, groups.code as groupcode,
					xdoc.categorycontrol, category.title  as category, category.code as categorycode,
					
					xdoc.workflow, xdoc.auditwf, xdoc.createdate, xdoc.createdby,  
					xdoc.updatedate, xdoc.updatedby, xdoc.auditdate, xdoc.auditedby 

				from  %s as xdoc
					
					left join profile on profile.control = xdoc.profilecontrol

					left join profile as toprofile on profile.control = xdoc.toprofilecontrol
					left join profile as fromprofile on profile.control = xdoc.fromprofilecontrol

					left join contact as tocontact on contact.control = xdoc.tocontactcontrol
					left join contact as fromcontact on contact.control = xdoc.fromcontactcontrol

					left join types on types.control = xdoc.typecontrol
					left join groups on groups.control = xdoc.groupcontrol
					left join category on category.control = xdoc.categorycontrol

				where xdoc.%s = '%s' order by xdoc.title
			`
	sql = fmt.Sprintf(sql, tablename, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Invoice) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchtext"] == nil {
		xDocrequest["searchtext"] = "%"
	}
	searchtext := xDocrequest["searchtext"].(string)

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	lMain := true
	if xDocrequest["main"] != nil {
		lMain = xDocrequest["main"].(bool)
	}

	tablename := this.GetName()
	if !lMain {
		tablename = tablename + "audit"
	}

	lAudited := false
	whereaudited := "xdoc.auditdate != '' AND xdoc.auditedby != '' AND"
	if xDocrequest["audited"] != nil {
		lAudited = xDocrequest["audited"].(bool)
	}
	if lAudited {
		whereaudited = "xdoc.auditdate = '' AND xdoc.auditedby = '' AND"
	}

	whereauditwfcode := ""
	if xDocrequest["auditwfcode"] != nil && xDocrequest["auditwfcode"].(string) != "" {
		whereauditwfcode = "xdoc.auditwf in (" + xDocrequest["auditwfcode"].(string) + ") AND"
	}

	whereworkflowcode := ""
	if xDocrequest["workflowcode"] != nil && xDocrequest["workflowcode"].(string) != "" {
		whereworkflowcode = "xdoc.workflow in (" + xDocrequest["workflowcode"].(string) + ") AND"
	}

	whereprofilecontrol := ""
	if xDocrequest["profilecontrol"] != nil && xDocrequest["profilecontrol"].(string) != "" {
		whereprofilecontrol = "xdoc.profilecontrol in (" + xDocrequest["profilecontrol"].(string) + ") AND"
	}

	wheretoprofilecontrol := ""
	if xDocrequest["toprofilecontrol"] != nil && xDocrequest["toprofilecontrol"].(string) != "" {
		wheretoprofilecontrol = "xdoc.toprofilecontrol in (" + xDocrequest["toprofilecontrol"].(string) + ") AND"
	}

	wherefromprofilecontrol := ""
	if xDocrequest["fromprofilecontrol"] != nil && xDocrequest["fromprofilecontrol"].(string) != "" {
		wherefromprofilecontrol = "xdoc.fromprofilecontrol in (" + xDocrequest["fromprofilecontrol"].(string) + ") AND"
	}

	wherevalidate := ""
	if xDocrequest["validate"] != nil && xDocrequest["validate"].(string) != "" {
		wherevalidate = "'" + xDocrequest["validate"].(string) + "' between xdoc.validfrom::timestamp AND xdoc.validtill::timestamp AND"
	}

	wheredeliverydate := ""
	if xDocrequest["deliverydate"] != nil && xDocrequest["deliverydate"].(string) != "" {
		wheredeliverydate = "'" + xDocrequest["deliverydate"].(string) + "' between xdoc.deliverydate::timestamp AND xdoc.latestdeliverydate::timestamp AND"
	}

	whereshippingdate := ""
	if xDocrequest["shippingdate"] != nil && xDocrequest["shippingdate"].(string) != "" {
		whereshippingdate = "'" + xDocrequest["shippingdate"].(string) + "' between xdoc.shippingdate::timestamp AND xdoc.latestshippingdate::timestamp AND"
	}

	wheretypecode := ""
	if xDocrequest["typecode"] != nil && xDocrequest["typecode"].(string) != "" {
		wheretypecode = "types.code in (" + xDocrequest["typecode"].(string) + ") AND"
	}

	wheregroupcode := ""
	if xDocrequest["groupcode"] != nil && xDocrequest["groupcode"].(string) != "" {
		wheregroupcode = "groups.code in (" + xDocrequest["groupcode"].(string) + ") AND"
	}

	wherecategorycode := ""
	if xDocrequest["categorycode"] != nil && xDocrequest["categorycode"].(string) != "" {
		wherecategorycode = "category.code in (" + xDocrequest["categorycode"].(string) + ") AND"
	}

	sql := `select 	xdoc.control, xdoc.title, xdoc.code, xdoc.description, 
					
					xdoc.docno, xdoc.docdate, xdoc.shiptocode, xdoc.shipfromcode,
					xdoc.validfrom, xdoc.validtill, xdoc.deliverydate, xdoc.shippingdate, 
					xdoc.shippingvehicle, xdoc.latestshippingdate, xdoc.latestdeliverydate,
					xdoc.discountname, xdoc.discountcode, xdoc.discountamount, xdoc.discountpercent,
					xdoc.taxamount, xdoc.taxpercent, xdoc.totalexcltax, xdoc.totalincltax,
					xdoc.currency, xdoc.currencyrate, 

					xdoc.profilecontrol, profile.title  as profile, profile.code  as profilecode,
					xdoc.toprofilecontrol, toprofile.title  as toprofile, toprofile.code  as toprofilecode,
					xdoc.fromprofilecontrol, fromprofile.title  as fromprofile, fromprofile.code  as fromprofilecode,

					xdoc.tocontactcontrol, tocontact.title  as tocontact, tocontact.code  as tocontactcode,
					xdoc.fromcontactcontrol, fromcontact.title  as fromcontact, fromcontact.code  as fromcontactcode,
					
					xdoc.typecontrol, types.title  as type, types.code as typecode,
					xdoc.groupcontrol, groups.title  as group, groups.code as groupcode,
					xdoc.categorycontrol, category.title  as category, category.code as categorycode,
					
					xdoc.workflow, xdoc.auditwf, xdoc.createdate, xdoc.createdby,  
					xdoc.updatedate, xdoc.updatedby, xdoc.auditdate, xdoc.auditedby 

				from  %s as xdoc
					
					left join profile on profile.control = xdoc.profilecontrol

					left join profile as toprofile on profile.control = xdoc.toprofilecontrol
					left join profile as fromprofile on profile.control = xdoc.fromprofilecontrol

					left join contact as tocontact on contact.control = xdoc.tocontactcontrol
					left join contact as fromcontact on contact.control = xdoc.fromcontactcontrol

					left join types on types.control = xdoc.typecontrol
					left join groups on groups.control = xdoc.groupcontrol
					left join category on category.control = xdoc.categorycontrol

				where %s %s %s %s %s %s %s %s %s %s %s %s %s %s (
						lower( xdoc.code) like lower('%%%s%%') or 
						lower( xdoc.title) like lower('%%%s%%') 
					)  order by  xdoc.auditdate::timestamp desc limit %s
			`
	sql = fmt.Sprintf(sql, tablename, whereaudited, whereauditwfcode, whereworkflowcode, whereprofilecontrol, wheretoprofilecontrol, wherefromprofilecontrol, wherevalidate, wheredeliverydate, whereshippingdate, wheretypecode, wheregroupcode, wherecategorycode, searchtext, searchtext, limit)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}
