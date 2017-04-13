package database

import (
	"fmt"
	"time"

	"io/ioutil"
	"strings"
	"valued/functions"
)

type TelrCode struct {
	Crud
}

func (this *TelrCode) GetName() string {
	return "telrcode"
}

func (this *TelrCode) GetWorkflows() []string {
	return this.Workflows()
}

func (this *TelrCode) GetFields() map[string]interface{} {

	fields := this.Fields()
	return fields
}

func (this *TelrCode) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)

	csvBytes, _ := ioutil.ReadFile("telrcode.csv")
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
						fieldName = "workflow"
					case 1:
						fieldName = "code"
					case 2:
						fieldName = "title"
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

func (this *TelrCode) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *TelrCode) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *TelrCode) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	code as code, control as control, title as title, 
					description as description, workflow as workflow, 
					createdate as createdate, createdby as createdby, 
					updatedate as updatedate, updatedby as updatedby

			from telrcode

			where %s = '%s' order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *TelrCode) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`workflow = '%s' and`, xDocrequest["workflow"])
	}

	wherecode := ""
	if xDocrequest["code"] != nil && xDocrequest["code"].(string) != "" {
		wherecode = fmt.Sprintf(`code like '%%%s%%' and`, xDocrequest["code"])
	}

	sql := `select 	code as code, control as control, title as title, 
					description as description, workflow as workflow, 
					createdate as createdate, createdby as createdby, 
					updatedate as updatedate, updatedby as updatedby
					
			from telrcode

			where %s %s control != '' order by control desc limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wherecode, limit, offset)
	mapRes, _ = curdb.Query(sql)

	return mapRes
}

func (this *TelrCode) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
