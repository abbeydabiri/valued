package database

import (
	"fmt"
)

type Smtp struct {
	Crud
}

func (this *Smtp) GetName() string {
	return "smtp"
}

func (this *Smtp) GetWorkflows() []string {
	return []string{"unverified", "invalid", "valid"}
}

func (this *Smtp) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["port"] = ""
	fields["server"] = ""
	fields["username"] = ""
	fields["password"] = ""
	fields["rate"] = "0"
	fields["delay"] = "0"

	return fields
}

func (this *Smtp) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Smtp) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Smtp) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Smtp) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ".*"
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	control, title, code, description, workflow, 
				 createdate, createdby, updatedate, updatedby,
				
				server, port, username, password,
				rate, delay

				from  smtp

				where %s = '%s' order by title`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Smtp) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	limit := "200"
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

	whereusername := ""
	if xDocrequest["username"] != nil && xDocrequest["username"].(string) != "" {
		whereusername = fmt.Sprintf(`username like '%%%s%%' AND`, xDocrequest["username"])
	}

	sql := `select 	control, title, code, description, workflow, 
				 createdate, createdby, updatedate, updatedby,
				
				server, port, username, password,
				rate, delay

			from smtp

			where %s %s %s

			order by server, username limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wheretitle, whereusername, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}
