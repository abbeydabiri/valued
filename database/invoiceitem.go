package database

import (
	"fmt"
)

type InvoiceItem struct {
	Crud
}

func (this *InvoiceItem) GetName() string {
	return "invoiceitem"
}

func (this *InvoiceItem) GetWorkflows() []string {
	return this.Workflows()
}

func (this *InvoiceItem) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["lineno"] = 0
	fields["invoicecontrol"] = ""

	fields["profilecontrol"] = ""

	fields["shiptocode"] = ""
	fields["toprofilecontrol"] = ""

	fields["shipfromcode"] = ""
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

	fields["quantity"] = 0
	fields["amount"] = 0.0

	fields["itemcontrol"] = ""
	fields["model"] = ""
	fields["brand"] = ""
	fields["serial"] = ""
	fields["barcode"] = ""
	fields["partcode"] = ""

	fields["uom"] = ""
	fields["unit"] = ""
	fields["size"] = ""
	fields["sizeuom"] = ""
	fields["pack"] = ""
	fields["packuom"] = ""
	fields["height"] = ""
	fields["heightuom"] = ""
	fields["weight"] = ""
	fields["weightuom"] = ""
	fields["length"] = ""
	fields["lengthuom"] = ""
	fields["width"] = ""
	fields["widthuom"] = ""
	fields["capacity"] = ""
	fields["capacityuom"] = ""
	return fields
}

func (this *InvoiceItem) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *InvoiceItem) Create(cUsername string, lMain bool, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc("system", this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *InvoiceItem) Update(cUsername string, lMain bool, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *InvoiceItem) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {
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
					
					xdoc.shiptocode, xdoc.shipfromcode,
					xdoc.validfrom, xdoc.validtill, xdoc.deliverydate, xdoc.shippingdate, 
					xdoc.shippingvehicle, xdoc.latestshippingdate, xdoc.latestdeliverydate,
					xdoc.discountname, xdoc.discountcode, xdoc.discountamount, xdoc.discountpercent,
					xdoc.taxamount, xdoc.taxpercent, xdoc.totalexcltax, xdoc.totalincltax,
					xdoc.currency, xdoc.currencyrate, 
					
					xdoc.lineno, xdoc.quantity, xdoc.amount,

					xdoc.model, xdoc.brand, xdoc.serial, xdoc.barcode, xdoc.partcode, 
					xdoc.uom, xdoc.unit, xdoc.sizeuom, xdoc.pack, xdoc.packuom, 
					xdoc.height, xdoc.heightuom, xdoc.weight, xdoc.weightuom,
					xdoc.length, xdoc.lengthuom, xdoc.width, xdoc.widthuom, 
					xdoc.capacity, xdoc.capacityuom,

					order.docno as invoicedocno, order.docdate as invoicedocdate,
					xdoc.invoicecontrol, order.title  as invoice, order.code as invoicecode, 
					order.currency as invoicecurrency, order.currencyrate as invoicecurrencyrate, 

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
					
					left join invoice on invoice.control = xdoc.invoicecontrol

					left join profile on profile.control = xdoc.profilecontrol

					left join profile as toprofile on profile.control = xdoc.toprofilecontrol
					left join profile as fromprofile on profile.control = xdoc.fromprofilecontrol

					left join contact as tocontact on contact.control = xdoc.tocontactcontrol
					left join contact as fromcontact on contact.control = xdoc.fromcontactcontrol

					left join types on types.control = xdoc.typecontrol
					left join groups on groups.control = xdoc.groupcontrol
					left join category on category.control = xdoc.categorycontrol

				where xdoc.%s = '%s' order by xdoc.lineno, xdoc.title
			`
	sql = fmt.Sprintf(sql, tablename, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *InvoiceItem) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchtext"] == nil {
		xDocrequest["searchtext"] = "%"
	}
	searchtext := xDocrequest["searchtext"].(string)

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"] != "" {
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
	if xDocrequest["auditwfcode"] != nil && xDocrequest["auditwfcode"] != "" {
		whereauditwfcode = "xdoc.auditwf in (" + xDocrequest["auditwfcode"].(string) + ") AND"
	}

	whereworkflowcode := ""
	if xDocrequest["workflowcode"] != nil && xDocrequest["workflowcode"] != "" {
		whereworkflowcode = "xdoc.workflow in (" + xDocrequest["workflowcode"].(string) + ") AND"
	}

	whereprofilecontrol := ""
	if xDocrequest["profilecontrol"] != nil && xDocrequest["profilecontrol"] != "" {
		whereprofilecontrol = "xdoc.profilecontrol in (" + xDocrequest["profilecontrol"].(string) + ") AND"
	}

	wheretoprofilecontrol := ""
	if xDocrequest["toprofilecontrol"] != nil && xDocrequest["toprofilecontrol"] != "" {
		wheretoprofilecontrol = "xdoc.toprofilecontrol in (" + xDocrequest["toprofilecontrol"].(string) + ") AND"
	}

	wherefromprofilecontrol := ""
	if xDocrequest["fromprofilecontrol"] != nil && xDocrequest["fromprofilecontrol"] != "" {
		wherefromprofilecontrol = "xdoc.fromprofilecontrol in (" + xDocrequest["fromprofilecontrol"].(string) + ") AND"
	}

	whereinvoicecontrol := ""
	if xDocrequest["invoicecontrol"] != nil && xDocrequest["invoicecontrol"] != "" {
		whereinvoicecontrol = "xdoc.invoicecontrol in (" + xDocrequest["invoicecontrol"].(string) + ") AND"
	}

	wherevalidate := ""
	if xDocrequest["validate"] != nil && xDocrequest["validate"] != "" {
		wherevalidate = "'" + xDocrequest["validate"].(string) + "' between xdoc.validfrom::timestamp AND xdoc.validtill::timestamp AND"
	}

	wheredeliverydate := ""
	if xDocrequest["deliverydate"] != nil && xDocrequest["deliverydate"] != "" {
		wheredeliverydate = "'" + xDocrequest["deliverydate"].(string) + "' between xdoc.deliverydate::timestamp AND xdoc.latestdeliverydate::timestamp AND"
	}

	whereshippingdate := ""
	if xDocrequest["shippingdate"] != nil && xDocrequest["shippingdate"] != "" {
		whereshippingdate = "'" + xDocrequest["shippingdate"].(string) + "' between xdoc.shippingdate::timestamp AND xdoc.latestshippingdate::timestamp AND"
	}

	wheretypecode := ""
	if xDocrequest["typecode"] != nil && xDocrequest["typecode"] != "" {
		wheretypecode = "types.code in (" + xDocrequest["typecode"].(string) + ") AND"
	}

	wheregroupcode := ""
	if xDocrequest["groupcode"] != nil && xDocrequest["groupcode"] != "" {
		wheregroupcode = "groups.code in (" + xDocrequest["groupcode"].(string) + ") AND"
	}

	wherecategorycode := ""
	if xDocrequest["categorycode"] != nil && xDocrequest["categorycode"] != "" {
		wherecategorycode = "category.code in (" + xDocrequest["categorycode"].(string) + ") AND"
	}

	sql := `select 	xdoc.control, xdoc.title, xdoc.code, xdoc.description, 
					
					xdoc.shiptocode, xdoc.shipfromcode,
					xdoc.validfrom, xdoc.validtill, xdoc.deliverydate, xdoc.shippingdate, 
					xdoc.shippingvehicle, xdoc.latestshippingdate, xdoc.latestdeliverydate,
					xdoc.discountname, xdoc.discountcode, xdoc.discountamount, xdoc.discountpercent,
					xdoc.taxamount, xdoc.taxpercent, xdoc.totalexcltax, xdoc.totalincltax,
					xdoc.currency, xdoc.currencyrate, 
					
					xdoc.lineno, xdoc.quantity, xdoc.amount,

					xdoc.model, xdoc.brand, xdoc.serial, xdoc.barcode, xdoc.partcode, 
					xdoc.uom, xdoc.unit, xdoc.sizeuom, xdoc.pack, xdoc.packuom, 
					xdoc.height, xdoc.heightuom, xdoc.weight, xdoc.weightuom,
					xdoc.length, xdoc.lengthuom, xdoc.width, xdoc.widthuom, 
					xdoc.capacity, xdoc.capacityuom,

					order.docno as invoicedocno, order.docdate as invoicedocdate,
					xdoc.invoicecontrol, order.title  as invoice, order.code as invoicecode, 
					order.currency as invoicecurrency, order.currencyrate as invoicecurrencyrate, 

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
					
					left join invoice on invoice.control = xdoc.invoicecontrol

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
					) order by xdoc.lineno, xdoc.title limit %s
			`
	sql = fmt.Sprintf(sql, tablename, whereaudited, whereauditwfcode, whereworkflowcode, whereprofilecontrol, wheretoprofilecontrol, wherefromprofilecontrol, whereinvoicecontrol, wherevalidate, wheredeliverydate, whereshippingdate, wheretypecode, wheregroupcode, wherecategorycode, searchtext, searchtext, limit)
	// log.Println(sql)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}
