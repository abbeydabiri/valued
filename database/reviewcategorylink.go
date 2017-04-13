package database

import (
	"fmt"
)

type ReviewCategoryLink struct {
	Crud
}

func (this *ReviewCategoryLink) GetName() string {
	return "reviewcategorylink"
}

func (this *ReviewCategoryLink) GetWorkflows() []string {
	return this.Workflows()
}

func (this *ReviewCategoryLink) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["reviewcategorycontrol"] = ""
	fields["merchantcontrol"] = ""
	return fields
}

func (this *ReviewCategoryLink) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *ReviewCategoryLink) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ReviewCategoryLink) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *ReviewCategoryLink) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	reviewcategorylink.code as code, reviewcategorylink.control as control, reviewcategorylink.title as title, 
					reviewcategorylink.description as description, reviewcategorylink.workflow as workflow, 
					reviewcategorylink.createdate as createdate, reviewcategorylink.createdby as createdby, 
					reviewcategorylink.updatedate as updatedate, reviewcategorylink.updatedby as updatedby,
				
					merchant.control as merchantcontrol, merchant.title as merchanttitle,

					reviewcategory.control as reviewcategorycontrol, reviewcategory.title as reviewcategorytitle, 
					reviewcategory.placement as reviewcategoryplacement

			from reviewcategorylink, reviewcategory, profile as merchant
			
			where reviewcategorylink.%s = '%s' AND 
			merchant.control = reviewcategorylink.merchantcontrol AND
			reviewcategory.control = reviewcategorylink.reviewcategorycontrol
			
			order by reviewcategoryplacement`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ReviewCategoryLink) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`reviewcategorylink.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wheremerchant := ""
	if xDocrequest["merchant"] != nil && xDocrequest["merchant"].(string) != "" {
		wheremerchant = fmt.Sprintf(`reviewcategorylink.merchantcontrol like '%%%s%%' AND`, xDocrequest["merchant"])
	}

	wherereviewcategory := ""
	if xDocrequest["reviewcategory"] != nil && xDocrequest["reviewcategory"].(string) != "" {
		wherereviewcategory = fmt.Sprintf(`reviewcategorylink.reviewcategorycontrol like '%%%s%%' AND`, xDocrequest["reviewcategory"])
	}

	sql := `select 	reviewcategorylink.code as code, reviewcategorylink.control as control, reviewcategorylink.title as title, 
					reviewcategorylink.description as description, reviewcategorylink.workflow as workflow, 
					reviewcategorylink.createdate as createdate, reviewcategorylink.createdby as createdby, 
					reviewcategorylink.updatedate as updatedate, reviewcategorylink.updatedby as updatedby,
				
					merchant.control as merchantcontrol, merchant.title as merchanttitle,

					reviewcategory.control as reviewcategorycontrol, reviewcategory.title as reviewcategorytitle, 
					reviewcategory.placement as reviewcategoryplacement

			from reviewcategorylink, reviewcategory, profile as merchant
			
			where %s %s %s

			merchant.control = reviewcategorylink.merchantcontrol AND
			reviewcategory.control = reviewcategorylink.reviewcategorycontrol
			
			order by reviewcategoryplacement limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wheremerchant, wherereviewcategory, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *ReviewCategoryLink) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
