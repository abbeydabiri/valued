package database

import (
	"fmt"
	"time"

	"io/ioutil"
	"strings"
	"valued/functions"
)

type Menu struct {
	Crud
}

func (this *Menu) GetName() string {
	return "menu"
}

func (this *Menu) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Menu) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["placement"] = ""
	fields["parent"] = ""
	fields["action"] = ""
	fields["icon"] = ""
	fields["role"] = ""
	return fields
}

func (this *Menu) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	csvBytes, _ := ioutil.ReadFile("menu.csv")
	if csvBytes != nil {
		dealerList := string(csvBytes)
		dealerList = strings.Replace(dealerList, "\r", "", -1)
		sliceRow := strings.Split(dealerList, "\n")

		go func() {
			for rowNo, stringCols := range sliceRow {

				if rowNo == 0 {
					continue
				}

				sliceCols := strings.Split(stringCols, ",")

				xDoc := make(map[string]interface{})
				for index, value := range sliceCols {
					fieldName := ""
					switch index {
					case 0:
						fieldName = "placement"
					case 1:
						fieldName = "code"
					case 2:
						fieldName = "parent"
					case 3:
						fieldName = "role"
					case 4:
						fieldName = "title"
					case 5:
						fieldName = "icon"
					case 6:
						fieldName = "action"
					}

					if fieldName != "" {
						xDoc[fieldName] = functions.TrimEscape(value)
					}
				}

				if xDoc["title"] != nil && xDoc["title"].(string) != "" {
					xDoc["workflow"] = "active"
					this.Create("root", xDoc, curdb)
				}
				<-time.Tick(time.Millisecond * 25)
			}
		}()
	}
}

func (this *Menu) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Menu) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Menu) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	menu.code as code, menu.control as control, menu.title as title, 
					menu.description as description, menu.workflow as workflow, 
					menu.createdate as createdate, menu.createdby as createdby, 
					menu.updatedate as updatedate, menu.updatedby as updatedby,

					menu.placement as placement, menu.parent as parent,
					menu.action as action, menu.icon as icon, menu.role as role

			from menu

			where menu.%s = '%s' order by placement`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Menu) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`menu.workflow = '%s' and`, xDocrequest["workflow"])
	}

	wheretitle := ""
	if xDocrequest["title"] != nil && xDocrequest["title"].(string) != "" {
		wheretitle = fmt.Sprintf(`menu.title like '%%%s%%' and`, xDocrequest["title"])
	}

	whereparent := ""
	if xDocrequest["parent"] != nil {
		whereparent = fmt.Sprintf(`menu.parent = '%s' and`, xDocrequest["parent"])
	}

	whererole := ""
	if xDocrequest["role"] != nil {
		whererole = fmt.Sprintf(`menu.role = '%s' and`, xDocrequest["role"])
	}

	wheresub := ""
	if xDocrequest["sub"] != nil {
		switch xDocrequest["sub"].(string) {
		case "true":
			wheresub = `menu.parent != '' and`
		default:
			wheresub = `menu.parent = '' and`
		}
	}

	sqlFields := `
					menu.code as code, menu.control as control, menu.title as title, 
					menu.description as description, menu.workflow as workflow, 
					menu.createdate as createdate, menu.createdby as createdby, 
					menu.updatedate as updatedate, menu.updatedby as updatedby,
					
					menu.placement as placement, menu.parent as parent,
					menu.action as action, menu.icon as icon, menu.role as role
				`

	sqlOrderByLimit := `order by menu.placement limit %s offset %s `
	if pagination {
		limit = ``
		offset = ``
		sqlOrderByLimit = `%s%s`
		sqlFields = `count(menu.control) as paginationtotal	`
	}

	sql := `select 	%s

			from menu 

			where %s %s %s %s %s menu.control != '' ` + sqlOrderByLimit

	sql = fmt.Sprintf(sql, sqlFields, whereworkflow, wheretitle, whereparent, whererole, wheresub, limit, offset)
	mapRes, _ = curdb.Query(sql)

	return mapRes
}

func (this *Menu) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
