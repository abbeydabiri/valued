package database

import (
	"fmt"

	"strings"
	"time"
)

type Subscription struct {
	Crud
}

func (this *Subscription) GetName() string {
	return "subscription"
}

func (this *Subscription) GetWorkflows() []string {
	// pending invoice-requested paid active
	return this.Workflows()
}

func (this *Subscription) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["price"] = 0.0
	fields["startdate"] = ""
	fields["expirydate"] = ""
	fields["schemecontrol"] = ""
	fields["membercontrol"] = ""
	fields["employercontrol"] = ""

	fields["sendersname"] = ""
	fields["sendersemail"] = ""
	return fields
}

func (this *Subscription) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Subscription) GenerateCode(curdb Database) string {

	pad := ""
	next := int64(0)

	sliceDate := strings.Split(time.Now().Format("02/01/2006"), "/")
	cSql := fmt.Sprintf("select count(control) as last from subscription where createdate like '%%/%s/%s%%'", sliceDate[1], sliceDate[2])
	mapResult, _ := curdb.Query(cSql)
	if mapResult["1"] != nil {
		next = mapResult["1"].(map[string]interface{})["last"].(int64) + int64(1)
	}

	padLenght := 4 - len(fmt.Sprintf("%d", next))

	for len(pad) < padLenght {
		pad += "0"
	}

	sliceDate = strings.Split(time.Now().Format("02/01/06"), "/")
	return fmt.Sprintf("%s%s%s%d", sliceDate[2], sliceDate[1], pad, next)

}

func (this *Subscription) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {

	//Include Subscription Generation Code in Format : YY-MM-4 digit Running Number e.g 17-01-0001 = First Member in January 17, 17-05-0100 = One hundreds member in May 17
	xDoc["code"] = this.GenerateCode(curdb)
	xDoc["title"] = xDoc["code"]
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Subscription) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Subscription) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	subscription.code as code, subscription.control as control, subscription.title as title, 
					subscription.description as description, subscription.workflow as workflow, 
					subscription.createdate as createdate, subscription.createdby as createdby, 
					subscription.updatedate as updatedate, subscription.updatedby as updatedby,

					subscription.price as price, subscription.expirydate as expirydate, 
					subscription.startdate as startdate, subscription.sendersname as sendersname,
					subscription.sendersemail as sendersemail,


					scheme.title as schemetitle, subscription.schemecontrol as schemecontrol,
					
					member.title as membertitle, member.firstname as memberfirstname, 
					member.lastname as memberlastname, subscription.membercontrol as membercontrol,
					subscription.employercontrol as employercontrol


			from subscription 
			left join profile as member on subscription.membercontrol = member.control
			left join scheme on subscription.schemecontrol = scheme.control

			where %s = '%s'  order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Subscription) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`scheme.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`scheme.title like '%%%s%%' AND`, xDocrequest["code"])
	}

	wherescheme := ""
	if xDocrequest["scheme"] != nil && xDocrequest["scheme"].(string) != "" {
		wherescheme = fmt.Sprintf(`subscription.schemecontrol = '%s' AND`, xDocrequest["scheme"])
	}

	whereschemename := ""
	if xDocrequest["schemename"] != nil && xDocrequest["schemename"].(string) != "" {
		whereschemename = fmt.Sprintf(`scheme.title like '%%%s%%' AND`, xDocrequest["schemename"])
	}

	wheremember := ""
	if xDocrequest["member"] != nil && xDocrequest["member"].(string) != "" {
		wheremember = fmt.Sprintf(`subscription.membercontrol = '%s' AND`, xDocrequest["member"])
	}

	wheremembername := ""
	if xDocrequest["membername"] != nil && xDocrequest["membername"].(string) != "" {
		wheremembername = fmt.Sprintf(`(member.firstname like '%%%s%%' OR member.lastname like '%%%s%%' ) AND`,
			xDocrequest["membername"])
	}

	whereemployer := ""
	if xDocrequest["employer"] != nil && xDocrequest["employer"].(string) != "" {
		whereemployer = fmt.Sprintf(`subscription.employercontrol = '%s' AND`, xDocrequest["employer"])
	}

	sqlFields := `
					subscription.code as code, subscription.control as control, subscription.title as title, 
					subscription.description as description, subscription.workflow as workflow, 
					subscription.createdate as createdate, subscription.createdby as createdby, 
					subscription.updatedate as updatedate, subscription.updatedby as updatedby,

					subscription.price as price, subscription.expirydate as expirydate, 
					subscription.startdate as startdate, subscription.sendersname as sendersname,
					subscription.sendersemail as sendersemail,


					scheme.title as schemetitle, subscription.schemecontrol as schemecontrol,
					
					member.title as membertitle, member.firstname as memberfirstname, 
					member.lastname as memberlastname, subscription.membercontrol as membercontrol,
					subscription.employercontrol as employercontrol
				`
	if pagination {
		sqlFields = `count(subscription.control) as paginationtotal	`
	}

	sql := `select 	%s
	
			from subscription 
			left join profile as member on subscription.membercontrol = member.control
			left join scheme on subscription.schemecontrol = scheme.control

			where %s %s %s %s %s %s %s

			subscription.control != '' 

			order by schemetitle, memberlastname limit %s offset %s
			`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wherecode, wherescheme, whereschemename,
		wheremember, wheremembername, whereemployer, limit, offset)
	mapRes, _ = curdb.Query(sql)

	return mapRes
}

func (this *Subscription) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
