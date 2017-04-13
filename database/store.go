package database

import (
	"fmt"
)

type Store struct {
	Crud
}

func (this *Store) GetName() string {
	return "store"
}

func (this *Store) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Store) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["merchantcontrol"] = ""

	fields["contact"] = ""
	fields["phone"] = ""
	fields["phonecode"] = ""
	fields["email"] = ""

	fields["address"] = ""
	fields["city"] = ""
	fields["country"] = ""

	fields["hoursholiday"] = ""
	fields["hoursmontofri"] = ""
	fields["hourssat"] = ""
	fields["hourssun"] = ""

	fields["image"] = ""
	fields["flagship"] = ""

	fields["gpslat"] = ""
	fields["gpslong"] = ""

	return fields
}

func (this *Store) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Store) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Store) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Store) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	store.code as code, store.control as control, store.title as title, 
					store.description as description, store.workflow as workflow, 
					store.createdate as createdate, store.createdby as createdby, 
					store.updatedate as updatedate, store.updatedby as updatedby,

					store.image as image, store.flagship as flagship,
					store.contact as contact, store.phone as phone, store.phonecode as phonecode, store.email as email, 
					store.address as address, store.city as city, store.country as country,
					
					store.hoursholiday as hoursholiday, store.hoursmontofri as hoursmontofri, 
					store.hourssat as hourssat, store.hourssun as hourssun, 
					store.gpslat as gpslat, store.gpslong as gpslong, 

					merchant.title as merchanttitle, merchant.description as merchantdescription,
					merchant.image as merchantimage, merchant.email as merchantemail,
					merchant.phone as merchantphone, merchant.website as merchantwebsite,
					
					store.merchantcontrol as merchantcontrol

			from store, profile as merchant

			where store.%s = '%s' and

			merchant.control = store.merchantcontrol
					
			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Store) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`store.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheremerchant := ""
	if xDocrequest["merchant"] != nil && xDocrequest["merchant"].(string) != "" {
		wheremerchant = fmt.Sprintf(`lower(merchant.title) like lower('%%%s%%') AND`, xDocrequest["merchant"])
	}

	wheremerchantcontrol := ""
	if xDocrequest["merchantcontrol"] != nil && xDocrequest["merchantcontrol"].(string) != "" {
		wheremerchant = fmt.Sprintf(`merchant.control = '%s' AND`, xDocrequest["merchantcontrol"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`lower(store.title) like lower('%%%s%%') AND`, xDocrequest["title"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`lower(store.code) like lower('%%%s%%') AND`, xDocrequest["code"])
	}

	wherecity := ""
	if xDocrequest["city"] != nil && xDocrequest["city"].(string) != "" {
		wherecity = fmt.Sprintf(`lower(store.city) like lower('%%%s%%') AND`, xDocrequest["city"])
	}

	wherecountry := ""
	if xDocrequest["country"] != nil && xDocrequest["country"].(string) != "" {
		wherecountry = fmt.Sprintf(`lower(store.country) like lower('%%%s%%') AND`, xDocrequest["country"])
	}

	whereflagship := ""
	if xDocrequest["flagship"] != nil && xDocrequest["flagship"].(string) != "" {
		whereflagship = fmt.Sprintf(`lower(store.flagship) like lower('%%%s%%') AND`, xDocrequest["flagship"])
	}

	sqlFields := `
					store.code as code, store.control as control, store.title as title, 
					store.description as description, store.workflow as workflow, 
					store.createdate as createdate, store.createdby as createdby, 
					store.updatedate as updatedate, store.updatedby as updatedby,

					store.image as image, store.flagship as flagship,
					store.contact as contact, store.phone as phone, store.phonecode as phonecode, store.email as email, 
					store.address as address, store.city as city, store.country as country,
					
					store.hoursholiday as hoursholiday, store.hoursmontofri as hoursmontofri, 
					store.hourssat as hourssat, store.hourssun as hourssun, 
					store.gpslat as gpslat, store.gpslong as gpslong,

					merchant.title as merchanttitle, merchant.description as merchantdescription,
					merchant.image as merchantimage, merchant.email as merchantemail,
					merchant.phone as merchantphone, merchant.website as merchantwebsite,

					store.merchantcontrol as merchantcontrol
				`
	sqlOrderByLimit := `order by title limit %s offset %s`
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(store.control) as paginationtotal	`
	}

	sql := `select 	%s

			from store left join profile as merchant on store.merchantcontrol = merchant.control

			where %s %s %s %s %s %s  %s %s 

			merchant.control = store.merchantcontrol ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, wherecode, wheremerchant,
		wheremerchantcontrol, wherecity, wherecountry, whereflagship, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Store) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
