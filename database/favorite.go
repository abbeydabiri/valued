package database

import (
	"fmt"
)

type Favorite struct {
	Crud
}

func (this *Favorite) GetName() string {
	return "favorite"
}

func (this *Favorite) GetWorkflows() []string {
	return this.Workflows()
}

func (this *Favorite) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["profilecontrol"] = ""
	fields["merchantcontrol"] = ""
	fields["rewardcontrol"] = ""
	fields["storecontrol"] = ""
	return fields
}

func (this *Favorite) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *Favorite) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Favorite) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *Favorite) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	refTable := ""
	if xDocrequest["reftable"] != nil && xDocrequest["reftable"] != "" {
		refTable = xDocrequest["reftable"].(string)
	}

	sql := `select 	favorite.code as code, favorite.control as control, favorite.title as title, 
					favorite.description as description, favorite.workflow as workflow, 
					favorite.createdate as createdate, favorite.createdby as createdby, 
					favorite.updatedate as updatedate, favorite.updatedby as updatedby,
					profile.control as profilecontrol, profile.title as profiletitle `

	switch refTable {
	case "store":
		sql += ` ,store.control as storecontrol, store.title as storetitle %s , store %s %s %s %s %s store.control = favorite.storecontrol AND`
	case "reward":
		sql += ` ,reward.control as rewardcontrol, reward.title as rewardtitle %s , reward %s %s %s %s %s reward.control = favorite.rewardcontrol AND`
	default: //merchant and others
		sql += ` ,merchant.control as merchantcontrol, merchant.title as merchanttitle %s , profile as merchant %s %s %s %s %s merchant.control = favorite.merchantcontrol AND`
	case "":
		//donothing
	}

	sql += ` favorite.%s = '%s' && profile.control = favorite.profilecontrol
			
					order by control desc`

	sql = fmt.Sprintf(sql, `from profile, favorite`, `where`, searchfield, searchvalue)

	// sql := `select 	favorite.code as code, favorite.control as control, favorite.title as title,
	// 				favorite.description as description, favorite.workflow as workflow,
	// 				favorite.createdate as createdate, favorite.createdby as createdby,
	// 				favorite.updatedate as updatedate, favorite.updatedby as updatedby,

	// 				merchant.control as merchantcontrol, merchant.title as merchanttitle,

	// 				profile.control as profilecontrol, profile.title as profiletitle,
	// 				profile.placement as profileplacement

	// 		from favorite, profile, profile as merchant

	// 		where favorite.%s = '%s' &&
	// 		merchant.control = favorite.merchantcontrol &&
	// 		profile.control = favorite.profilecontrol

	// 		order by profileplacement`

	// sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Favorite) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	limit := "100"
	if xDocrequest["limit"] != nil && xDocrequest["limit"].(string) != "" {
		limit = xDocrequest["limit"].(string)
	}

	offset := "0"
	if xDocrequest["offset"] != nil && xDocrequest["offset"].(string) != "" {
		offset = xDocrequest["offset"].(string)
	}

	refTable := ""
	if xDocrequest["reftable"] != nil && xDocrequest["reftable"].(string) != "" {
		refTable = xDocrequest["reftable"].(string)
	}

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
		wheremerchant = fmt.Sprintf(`favorite.merchantcontrol like '%%%s%%' AND`, xDocrequest["merchant"])
		refTable = "merchant"
	}

	whereprofile := ""
	if xDocrequest["profile"] != nil && xDocrequest["profile"].(string) != "" {
		whereprofile = fmt.Sprintf(`favorite.profilecontrol like '%%%s%%' AND`, xDocrequest["profile"])
	}

	whereworkflow := ""
	if xDocrequest["workflow"] != nil && xDocrequest["workflow"].(string) != "" {
		whereworkflow = fmt.Sprintf(`favorite.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherecategory := ""
	if xDocrequest["category"] != nil && xDocrequest["category"].(string) != "" {
		switch refTable {
		case "reward", "merchant":
			wherecategory = fmt.Sprintf(`%s.categorycontrol like '%%%s%%' AND`, refTable, xDocrequest["category"])
		}
	}
	wheresubcategory := ""
	if xDocrequest["merchant"] != nil && xDocrequest["merchant"].(string) != "" {
		switch refTable {
		case "reward", "merchant":
			wheresubcategory = fmt.Sprintf(`%s.subcategorycontrol like '%%%s%%' AND`, refTable, xDocrequest["category"])
		}
	}

	sql := `select 	favorite.code as code, favorite.control as control, favorite.title as title, 
					favorite.description as description, favorite.workflow as workflow, 
					favorite.createdate as createdate, favorite.createdby as createdby, 
					favorite.updatedate as updatedate, favorite.updatedby as updatedby,
					profile.control as profilecontrol, profile.title as profiletitle `

	switch refTable {
	case "store":
		sql += `,store.control as storecontrol, store.title as storetitle %s , store %s %s %s %s %s %s %s %s store.control = favorite.storecontrol AND`
	case "reward":
		sql += `,reward.control as rewardcontrol, reward.title as rewardtitle %s , reward %s %s %s %s %s %s %s %s reward.control = favorite.rewardcontrol AND`
	case "merchant":
		sql += `,merchant.control as merchantcontrol, merchant.title as merchanttitle %s , profile as merchant %s %s %s %s %s %s %s %s merchant.control = favorite.merchantcontrol AND`
	}

	//merchant.control as merchantcontrol, merchant.title as merchanttitle
	// %s , profile as merchant

	// 		%s %s %s %s
	// 		merchant.control = favorite.merchantcontrol &&

	sql += `profile.control = favorite.profilecontrol
			
					order by control desc limit %s offset %s`

	sql = fmt.Sprintf(sql, `from profile, favorite`, `where`, whereworkflow, whereprofile, wheremerchant, wherereward, wherestore,
		wherecategory, wheresubcategory, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *Favorite) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
