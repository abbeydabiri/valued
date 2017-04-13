package database

import (
	"fmt"
)

type Permission struct {
	Crud
}

func (this *Permission) GetName() string {
	return "permission"
}

func (this *Permission) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Permission) GetFields() map[string]interface{} {

	// menu-code, (title)
	// code, (owner, group, others)
	// role, (admin, agent, partner, player)
	// action, (search,new,view,edit,activate,deactivate etc)

	fields := this.Fields()
	fields["role"] = ""
	fields["action"] = ""
	fields["menucontrol"] = ""
	fields["logincontrol"] = ""
	return fields
}

func (this *Permission) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Permission) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Permission) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Permission) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	permission.code as code, permission.control as control, permission.title as title, 
					permission.description as description, permission.workflow as workflow, 
					permission.createdate as createdate, permission.createdby as createdby, 
					permission.updatedate as updatedate, permission.updatedby as updatedby,

					permission.role as role, permission.action as action,
					permission.logincontrol as logincontrol, login.username as loginusername, 

					login.profilecontrol as profilecontrol, profile.title as profiletitle,
					profile.firstname as profilefirstname, profile.lastname as profilelastname,
					profile.image as profileimage,

					permission.menucontrol as menucontrol, menu.code as menucode, menu.title as menutitle

			from permission

					left join menu on menu.control = permission.menucontrol
					left join login on login.control = permission.logincontrol
					left join profile on login.profilecontrol = profile.control
					
					where permission.%s = '%s' order by logincontrol, title, code, role`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Permission) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	pagination := false
	if xDocrequest["pagination"] != nil {
		pagination = xDocrequest["pagination"].(bool)
	}

	if xDocrequest["searchtext"] == nil {
		xDocrequest["searchtext"] = ""
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
		whereworkflow = fmt.Sprintf(`permission.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	whererole := ""
	if xDocrequest["role"] != nil && xDocrequest["role"].(string) != "" {
		whererole = fmt.Sprintf(`permission.role = '%s' AND`, xDocrequest["role"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`permission.code = '%s' AND`, xDocrequest["code"])
	}

	whereaction := ""
	if xDocrequest["action"] != nil && xDocrequest["action"].(string) != "" {
		whereaction = fmt.Sprintf(`permission.action = '%s' AND`, xDocrequest["action"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`permission.title like '%%%s%%' AND`, xDocrequest["title"])
	}

	wheremenu := ""
	if xDocrequest["menu"] != nil && xDocrequest["menu"].(string) != "" {
		wheremenu = fmt.Sprintf(`permission.menucontrol = '%s' AND`, xDocrequest["menu"])
	}

	wherelogin := ""
	if xDocrequest["login"] != nil && xDocrequest["login"].(string) != "" {
		wherelogin = fmt.Sprintf(`permission.logincontrol = '%s' AND`, xDocrequest["login"])
	}

	whereprofile := ""
	if xDocrequest["profile"] != nil && xDocrequest["profile"].(string) != "" {
		whereprofile = fmt.Sprintf(`login.profilecontrol = '%s' AND`, xDocrequest["profile"])
	}

	sqlFields := `
					permission.code as code, permission.control as control, permission.title as title, 
					permission.description as description, permission.workflow as workflow, 
					permission.createdate as createdate, permission.createdby as createdby, 
					permission.updatedate as updatedate, permission.updatedby as updatedby,

					permission.role as role, permission.action as action,
					permission.logincontrol as logincontrol, login.username as loginusername, 
					
					login.profilecontrol as profilecontrol, profile.title as profiletitle,
					profile.firstname as profilefirstname, profile.lastname as profilelastname,
					profile.image as profileimage,

					permission.menucontrol as menucontrol, menu.code as menucode, menu.title as menutitle
				`
	if pagination {
		sqlFields = `count(permission.control) as paginationtotal	`
	}

	sql := `select 	%s

			from permission

					left join menu on menu.control = permission.menucontrol
					left join login on login.control = permission.logincontrol
					left join profile on login.profilecontrol = profile.control


			where %s %s %s %s %s %s %s %s  (permission.code like '%%%s%%' or permission.role like '%%%s%%' or permission.action like '%%%s%%') 
			order by title, role, code, action, logincontrol limit %s offset %s`

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wherecode, whererole, wheretitle, whereaction, wheremenu, wherelogin, whereprofile, searchtext, searchtext, searchtext, limit, offset)

	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Permission) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
