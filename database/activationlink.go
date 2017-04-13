package database

import (
	"fmt"
)

type ActivationLink struct {
	Crud
}

func (this *ActivationLink) GetName() string {
	return "activationlink"
}

func (this *ActivationLink) GetWorkflows() []string {
	return this.Workflows()
}

func (this *ActivationLink) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["vieweddate"] = ""
	fields["expirydate"] = ""
	fields["logincontrol"] = ""
	return fields
}

func (this *ActivationLink) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *ActivationLink) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ActivationLink) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ActivationLink) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	activationlink.code as code, activationlink.control as control, activationlink.title as title, 
					activationlink.description as description, activationlink.workflow as workflow, 
					activationlink.createdate as createdate, activationlink.createdby as createdby, 
					activationlink.updatedate as updatedate, activationlink.updatedby as updatedby,
					activationlink.expirydate as expirydate, activationlink.vieweddate as vieweddate, 

					activationlink.logincontrol as logincontrol, login.username as loginusername, 
					login.profilecontrol as profilecontrol, profile.title as profiletitle,
					profile.firstname as profilefirstname, profile.lastname as profilelastname,
					profile.image as profileimage

			from activationlink 
					left join login on activationlink.logincontrol = login.control
					left join profile on login.profilecontrol = profile.control

			where activationlink.%s = '%s' order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ActivationLink) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	pagination := false
	if xDocrequest["pagination"] != nil {
		pagination = xDocrequest["pagination"].(bool)
	}

	if xDocrequest["searchtext"] == nil {
		xDocrequest["searchtext"] = "%"
	}
	searchtext := xDocrequest["searchtext"].(string)

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
		whereworkflow = fmt.Sprintf(`activationlink.workflow = '%s' and`, xDocrequest["workflow"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`activationlink.code = '%s' and`, xDocrequest["code"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`activationlink.title like '%%%s%%' and`, xDocrequest["title"])
	}

	wherelogin := ""
	if xDocrequest["login"] != nil && xDocrequest["login"].(string) != "" {
		wherelogin = fmt.Sprintf(`activationlink.logincontrol = '%s' and`, xDocrequest["login"])
	}

	whereprofile := ""
	if xDocrequest["profile"] != nil && xDocrequest["profile"].(string) != "" {
		whereprofile = fmt.Sprintf(`login.profilecontrol = '%s' and`, xDocrequest["profile"])
	}

	sqlFields := `
					activationlink.code as code, activationlink.control as control, activationlink.title as title, 
					activationlink.description as description, activationlink.workflow as workflow, 
					activationlink.createdate as createdate, activationlink.createdby as createdby, 
					activationlink.updatedate as updatedate, activationlink.updatedby as updatedby,
					activationlink.expirydate as expirydate, activationlink.vieweddate as vieweddate, 

					activationlink.logincontrol as logincontrol, login.username as loginusername, 
					login.profilecontrol as profilecontrol, profile.title as profiletitle,
					profile.firstname as profilefirstname, profile.lastname as profilelastname,
					profile.image as profileimage
				`

	sqlOrderByLimit := `order by activationlink.control desc limit %s offset %s `
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(activationlink.control) as paginationtotal	`
	}

	sql := `select 	%s

			from activationlink 
					left join login on activationlink.logincontrol = login.control
					left join profile on login.profilecontrol = profile.control

			where %s %s %s %s %s (
				lower( profile.title) like lower('%%%s%%') or
				lower( profile.firstname) like lower('%%%s%%') or
				lower( profile.lastname) like lower('%%%s%%') or

				lower( login.username) like lower('%%%s%%') or
				lower( activationlink.title) like lower('%%%s%%') or
				lower( activationlink.description) like lower('%%%s%%') 
				
			) ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wherecode, wheretitle, wherelogin, whereprofile,
		searchtext, searchtext, searchtext, searchtext, searchtext, searchtext, limit, offset)

	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ActivationLink) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
