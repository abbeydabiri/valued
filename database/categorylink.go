package database

import (
	"fmt"
)

type CategoryLink struct {
	Crud
}

func (this *CategoryLink) GetName() string {
	return "categorylink"
}

func (this *CategoryLink) GetWorkflows() []string {
	return this.Workflows()
}

func (this *CategoryLink) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["categorycontrol"] = ""
	fields["merchantcontrol"] = ""
	fields["rewardcontrol"] = ""
	fields["storecontrol"] = ""
	return fields
}

func (this *CategoryLink) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *CategoryLink) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *CategoryLink) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *CategoryLink) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	refTable := ""
	if xDocrequest["reftable"] != nil && xDocrequest["reftable"].(string) != "" {
		refTable = xDocrequest["reftable"].(string)
	}

	sql := `select 	categorylink.code as code, categorylink.control as control, categorylink.title as title, 
					categorylink.description as description, categorylink.workflow as workflow, 
					categorylink.createdate as createdate, categorylink.createdby as createdby, 
					categorylink.updatedate as updatedate, categorylink.updatedby as updatedby,
					category.control as categorycontrol, category.title as categorytitle, 
					category.placement as categoryplacement`

	switch refTable {
	case "store":
		sql += ` ,store.control as storecontrol, store.title as storetitle %s , store %s %s %s %s %s store.control = categorylink.storecontrol AND`
	case "reward":
		sql += ` ,reward.control as rewardcontrol, reward.title as rewardtitle %s , reward %s %s %s %s %s reward.control = categorylink.rewardcontrol AND`
	default: //merchant and others
		sql += ` ,merchant.control as merchantcontrol, merchant.title as merchanttitle %s , profile as merchant %s %s %s %s %s merchant.control = categorylink.merchantcontrol AND`
	case "":
		//donothing
	}

	sql += ` categorylink.%s = '%s' AND category.control = categorylink.categorycontrol
			
					order by categoryplacement`

	sql = fmt.Sprintf(sql, `from category, categorylink`, `where`, searchfield, searchvalue)

	// sql := `select 	categorylink.code as code, categorylink.control as control, categorylink.title as title,
	// 				categorylink.description as description, categorylink.workflow as workflow,
	// 				categorylink.createdate as createdate, categorylink.createdby as createdby,
	// 				categorylink.updatedate as updatedate, categorylink.updatedby as updatedby,

	// 				merchant.control as merchantcontrol, merchant.title as merchanttitle,

	// 				category.control as categorycontrol, category.title as categorytitle,
	// 				category.placement as categoryplacement

	// 		from categorylink, category, profile as merchant

	// 		where categorylink.%s = '%s' &&
	// 		merchant.control = categorylink.merchantcontrol &&
	// 		category.control = categorylink.categorycontrol

	// 		order by categoryplacement`

	// sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *CategoryLink) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	offset := "0"
	if xDocrequest["offset"] != nil && xDocrequest["offset"].(string) != "" {
		offset = xDocrequest["offset"].(string)
	}

	refTable := ""
	wherestore := ""
	if xDocrequest["store"] != nil && xDocrequest["store"].(string) != "" {
		wherestore = fmt.Sprintf(`storelink.storecontrol like '%%%s%%' AND`, xDocrequest["store"])
		refTable = "store"
	}

	wherereward := ""
	if xDocrequest["reward"] != nil && xDocrequest["reward"].(string) != "" {
		wherereward = fmt.Sprintf(`rewardlink.rewardcontrol like '%%%s%%' AND`, xDocrequest["reward"])
		refTable = "reward"
	}

	wheremerchant := ""
	if xDocrequest["merchant"] != nil && xDocrequest["merchant"].(string) != "" {
		wheremerchant = fmt.Sprintf(`categorylink.merchantcontrol like '%%%s%%' AND`, xDocrequest["merchant"])
		refTable = "merchant"
	}

	wherecategory := ""
	if xDocrequest["category"] != nil && xDocrequest["category"].(string) != "" {
		wherecategory = fmt.Sprintf(`categorylink.categorycontrol like '%%%s%%' AND`, xDocrequest["category"])
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`categorylink.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	sql := `select 	categorylink.code as code, categorylink.control as control, categorylink.title as title, 
					categorylink.description as description, categorylink.workflow as workflow, 
					categorylink.createdate as createdate, categorylink.createdby as createdby, 
					categorylink.updatedate as updatedate, categorylink.updatedby as updatedby,
					category.control as categorycontrol, category.title as categorytitle, 
					category.placement as categoryplacement,`

	switch refTable {
	case "store":
		sql += `store.control as storecontrol, store.title as storetitle %s , store %s %s %s %s %s %s store.control = categorylink.storecontrol AND `
	case "reward":
		sql += `reward.control as rewardcontrol, reward.title as rewardtitle %s , reward %s %s %s %s %s %s reward.control = categorylink.rewardcontrol AND `
	case "merchant":
		sql += `merchant.control as merchantcontrol, merchant.title as merchanttitle %s , profile as merchant %s %s %s %s %s %s merchant.control = categorylink.merchantcontrol AND `
	}

	//merchant.control as merchantcontrol, merchant.title as merchanttitle
	// %s , profile as merchant

	// 		%s %s %s %s
	// 		merchant.control = categorylink.merchantcontrol &&

	sql += ` category.control = categorylink.categorycontrol
			
					order by categoryplacement limit %s offset %s`

	sql = fmt.Sprintf(sql, `from category, categorylink`, `where`, whereworkflow, wherecategory, wheremerchant, wherereward, wherestore, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *CategoryLink) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
