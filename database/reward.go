package database

import (
	"fmt"
	"strings"
)

type Reward struct {
	Crud
}

func (this *Reward) GetName() string {
	return "reward"
}

func (this *Reward) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Reward) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["beneficiary"] = ""
	fields["merchantcontrol"] = ""
	fields["categorycontrol"] = ""
	fields["subcategorycontrol"] = ""

	fields["startdate"] = ""
	fields["enddate"] = ""

	fields["image"] = ""

	fields["type"] = ""
	fields["method"] = ""

	fields["restriction"] = ""
	fields["discount"] = ""

	fields["discounttype"] = ""
	fields["discountvalue"] = 0.0
	fields["visibleto"] = ""

	fields["keywords"] = ""

	fields["maxuse"] = 0
	fields["maxperuser"] = 0
	fields["maxpermonth"] = 0

	fields["orderby"] = 0

	return fields
}

func (this *Reward) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Reward) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Reward) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Reward) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	reward.code as code, reward.control as control, reward.title as title, 
					reward.description as description, reward.workflow as workflow, 
					reward.createdate as createdate, reward.createdby as createdby, 
					reward.updatedate as updatedate, reward.updatedby as updatedby,
					
					reward.image as image, reward.type as type, reward.method as method,
					reward.startdate as startdate, reward.enddate as enddate, 
					reward.maxuse as maxuse, reward.maxperuser as maxperuser,
					reward.maxpermonth as maxpermonth,
					reward.orderby as orderby,
					reward.beneficiary as beneficiary, reward.restriction as restriction, 
					reward.discount as discount, reward.discounttype as discounttype, 
					reward.keywords as keywords, reward.visibleto as visibleto, 
					reward.discountvalue as discountvalue,

					reward.merchantcontrol as merchantcontrol, merchant.title as merchanttitle,
					merchant.firstname as merchantfirstname, merchant.lastname as merchantlastname,
					merchant.image as merchantimage, merchant.email as merchantemail,
					merchant.phone as merchantphone, merchant.website as merchantwebsite,

					reward.categorycontrol as categorycontrol, reward.subcategorycontrol as subcategorycontrol

			from reward, profile as merchant

			where reward.%s = '%s' AND

			reward.merchantcontrol = merchant.control 

			order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Reward) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	pagination := false
	if xDocrequest["pagination"] != nil {
		pagination = xDocrequest["pagination"].(bool)
	}

	searchtext := ""
	if xDocrequest["searchtext"] != nil && xDocrequest["searchtext"].(string) != "" {
		searchtext = xDocrequest["searchtext"].(string)
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
		whereworkflow = fmt.Sprintf(`reward.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`lower(reward.title) like lower('%%%s%%') AND`, xDocrequest["title"])
	}

	wheretype := ""
	if xDocrequest["type"] != nil && xDocrequest["type"].(string) != "" {
		wheretype = fmt.Sprintf(`lower(reward.type) like lower('%%%s%%') AND`, xDocrequest["type"])
	}

	whererewardcontrol := ""
	if xDocrequest["rewardcontrol"] != nil && len(xDocrequest["rewardcontrol"].([]string)) > 0 {
		whererewardcontrol = strings.Join(xDocrequest["rewardcontrol"].([]string), `","`)
		whererewardcontrol = fmt.Sprintf(`reward.control in ("%s") AND`, whererewardcontrol)
	}

	wheremerchant := ""
	if xDocrequest["merchant"] != nil && xDocrequest["merchant"].(string) != "" {
		wheremerchant = fmt.Sprintf(`lower(merchant.company) = lower('Yes') and lower(merchant.title) like lower('%%%s%%') AND`, xDocrequest["merchant"])
	}

	wheremerchantcontrol := ""
	if xDocrequest["merchantcontrol"] != nil && xDocrequest["merchantcontrol"].(string) != "" {
		wheremerchantcontrol = fmt.Sprintf(`lower(merchant.company) = lower('Yes') and merchant.control = '%s' AND`, xDocrequest["merchantcontrol"])
	}

	wherecategory := ""
	if xDocrequest["category"] != nil && xDocrequest["category"].(string) != "" {
		wherecategory = fmt.Sprintf(`reward.categorycontrol = '%s' AND`, xDocrequest["category"])
	}
	wheresubcategory := ""
	if xDocrequest["subcategory"] != nil && xDocrequest["subcategory"].(string) != "" {
		wheresubcategory = fmt.Sprintf(`reward.subcategorycontrol = '%s' AND`, xDocrequest["subcategory"])
	}

	wherekeywords := ""
	if xDocrequest["keywords"] != nil && xDocrequest["keywords"].(string) != "" {
		wherekeywords = fmt.Sprintf(`lower(reward.keywords) like lower('%%%s%%') AND`, xDocrequest["keywords"])
	}

	wherestartdate := ""
	if xDocrequest["startdate"] != nil && xDocrequest["startdate"].(string) != "" {
		wherestartdate = fmt.Sprintf(`reward.startdate like '%%%s%%' AND`, xDocrequest["startdate"])
	}

	wherebeneficiary := ""
	if xDocrequest["beneficiary"] != nil && xDocrequest["beneficiary"].(string) != "" {
		wherebeneficiary = fmt.Sprintf(`lower(reward.beneficiary) like lower('%%%s%%') AND`, xDocrequest["beneficiary"])
	}

	sqlFields := `
					reward.code as code, reward.control as control, reward.title as title, 
					reward.description as description, reward.workflow as workflow, 
					reward.createdate as createdate, reward.createdby as createdby, 
					reward.updatedate as updatedate, reward.updatedby as updatedby,
					
					reward.image as image, reward.type as type, reward.method as method,
					reward.startdate as startdate, reward.enddate as enddate, 
					reward.maxuse as maxuse, reward.maxperuser as maxperuser,
					reward.maxpermonth as maxpermonth,
					reward.orderby as orderby,
					reward.beneficiary as beneficiary, reward.restriction as restriction, 
					reward.discount as discount, reward.discounttype as discounttype, 
					reward.keywords as keywords, reward.visibleto as visibleto, 
					reward.discountvalue as discountvalue,

					reward.merchantcontrol as merchantcontrol, merchant.title as merchanttitle,
					merchant.firstname as merchantfirstname, merchant.lastname as merchantlastname,
					merchant.image as merchantimage, merchant.email as merchantemail,
					merchant.phone as merchantphone, merchant.website as merchantwebsite,

					reward.categorycontrol as categorycontrol, reward.subcategorycontrol as subcategorycontrol
				`

	sqlOrderByLimit := `order by reward.title limit %s offset %s`
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(reward.control) as paginationtotal	`
	}

	sql := `select %s

			from reward left join profile as merchant on reward.merchantcontrol = merchant.control

			where %s %s %s %s %s %s %s %s %s %s %s 

			(lower(reward.title) like lower('%%%s%%') OR lower(merchant.title) like lower('%%%s%%')OR reward.discount like lower('%%%s%%') ) AND

			reward.merchantcontrol = merchant.control ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretype, wheretitle, wheremerchant, wheremerchantcontrol,
		wherecategory, wheresubcategory, wherekeywords, wherestartdate, wherebeneficiary, whererewardcontrol,
		searchtext, searchtext, searchtext, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Reward) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
