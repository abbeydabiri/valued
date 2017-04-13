package database

import (
	"fmt"
)

type Referral struct {
	Crud
}

func (this *Referral) GetName() string {
	return "referral"
}

func (this *Referral) GetWorkflows() []string {
	return []string{"referred", "subscribed"}
}

func (this *Referral) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["firstname"] = ""
	fields["lastname"] = ""
	fields["email"] = ""
	fields["profilecontrol"] = ""
	fields["friendfirstname"] = ""
	fields["friendlastname"] = ""
	fields["friendemail"] = ""
	return fields
}

func (this *Referral) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Referral) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Referral) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Referral) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	referral.code as code, referral.control as control, referral.title as title, 
					referral.description as description, referral.workflow as workflow, 
					referral.createdate as createdate, referral.createdby as createdby, 
					referral.updatedate as updatedate, referral.updatedby as updatedby,
					
					referral.firstname as firstname, referral.lastname as lastname,
					referral.email as email, 

					referral.friendfirstname as friendfirstname, referral.friendlastname as friendlastname,
					referral.friendemail as friendemail,

					referral.profilecontrol as profilecontrol, profile.title as profiletitle, 
					profile.firstname as profilefirstname, profile.lastname as profilelastname

			from referral left join profile as profile
			on profile.control = referral.profilecontrol

			where referral.%s = '%s' order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Referral) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`referral.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`referral.code like '%%%s%%' AND`, xDocrequest["code"])
	}

	whereemail := ""
	if xDocrequest["email"] != nil && xDocrequest["email"].(string) != "" {
		whereemail = fmt.Sprintf(`referral.email like '%%%s%%' AND`, xDocrequest["email"])
	}

	whereprofile := ""
	if xDocrequest["profile"] != nil {
		whereprofile = fmt.Sprintf(`referral.profilecontrol = '%s' AND`, xDocrequest["profile"])
	}

	sqlFields := `
					referral.code as code, referral.control as control, referral.title as title, 
					referral.description as description, referral.workflow as workflow, 
					referral.createdate as createdate, referral.createdby as createdby, 
					referral.updatedate as updatedate, referral.updatedby as updatedby,
					
					referral.firstname as firstname, referral.lastname as lastname,
					referral.email as email, 

					referral.friendfirstname as friendfirstname, referral.friendlastname as friendlastname,
					referral.friendemail as friendemail,

					referral.profilecontrol as profilecontrol, profile.title as profiletitle, 
					profile.firstname as profilefirstname, profile.lastname as profilelastname
				`
	sqlOrderByLimit := ` order by referral.control desc limit %s offset %s`
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(referral.control) as paginationtotal	`
	}

	sql := `select 	%s

			from referral left join profile as profile
			on profile.control = referral.profilecontrol

			where %s %s %s %s referral.control != '' ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wherecode, whereemail, whereprofile, limit, offset)

	if xDocrequest["pagination"] != nil {

	} else {
		mapRes, _ = curdb.Query(sql)
	}
	return mapRes
}

func (this *Referral) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
