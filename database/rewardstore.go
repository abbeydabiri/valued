package database

import (
	"fmt"
)

type RewardStore struct {
	Crud
}

func (this *RewardStore) GetName() string {
	return "rewardstore"
}

func (this *RewardStore) GetWorkflows() []string {
	return this.Workflows()
}

func (this *RewardStore) GetFields() map[string]interface{} {

	fields := this.Fields()
	fields["rewardcontrol"] = ""
	fields["storecontrol"] = ""
	return fields
}

func (this *RewardStore) Initialize(curdb Database) {
	this.InitializeXdoc(this.GetName(), this.GetFields(), curdb)
}

func (this *RewardStore) Create(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["createdby"] = cUsername
	xDoc["createdate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *RewardStore) Update(cUsername string, xDoc map[string]interface{}, curdb Database) string {
	xDoc["updatedby"] = cUsername
	xDoc["updatedate"] = this.GetSystemTime()
	return this.WriteXdoc(cUsername, this.GetName(), this.GetFields(), xDoc, curdb)
}

func (this *RewardStore) Read(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

	if xDocrequest["searchfield"] == nil {
		xDocrequest["searchfield"] = "control"
	}
	searchfield := xDocrequest["searchfield"].(string)

	if xDocrequest["searchvalue"] == nil {
		xDocrequest["searchvalue"] = ""
	}
	searchvalue := xDocrequest["searchvalue"].(string)

	sql := `select 	rewardstore.code as code, rewardstore.control as control, rewardstore.title as title, 
					rewardstore.description as description, rewardstore.workflow as workflow, 
					rewardstore.createdate as createdate, rewardstore.createdby as createdby, 
					rewardstore.updatedate as updatedate, rewardstore.updatedby as updatedby,
						
					rewardstore.storecontrol as storecontrol, store.title as storetitle, 
					store.code as storecode, store.city as storecity, store.address as storeaddress, 
					store.image as storeimage, store.contact as storecontact, 
					store.phone as storephone, store.email as storeemail, 
					store.gpslat as storegpslat, store.gpslong as storegpslong, 
					store.flagship as storeflagship, 

					rewardstore.rewardcontrol as rewardcontrol,
					reward.title as rewardtitle

			from rewardstore, reward, store

			where rewardstore.%s = '%s' and

			reward.control = rewardstore.rewardcontrol 
			AND store.control = rewardstore.storecontrol

			order by control desc`

	sql = fmt.Sprintf(sql, searchfield, searchvalue)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *RewardStore) Search(xDocrequest map[string]interface{}, curdb Database) (mapRes map[string]interface{}) {

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
		whereworkflow = fmt.Sprintf(`rewardstore.workflow = '%s' AND`, xDocrequest["workflow"])
	}

	wherereward := ""
	if xDocrequest["reward"] != nil && xDocrequest["reward"].(string) != "" {
		wherereward = fmt.Sprintf(`rewardstore.rewardcontrol like '%%%s%%' AND`, xDocrequest["reward"])
	}

	wherestore := ""
	if xDocrequest["store"] != nil && xDocrequest["store"].(string) != "" {
		wherestore = fmt.Sprintf(`rewardstore.storecontrol like '%%%s%%' AND`, xDocrequest["store"])
	}

	whereflagship := ""
	if xDocrequest["flagship"] != nil && xDocrequest["flagship"].(string) != "" {
		whereflagship = fmt.Sprintf(`store.flagship = '%s' AND`, xDocrequest["flagship"])
	}

	wheremerchant := ""
	if xDocrequest["merchant"] != nil && xDocrequest["merchant"].(string) != "" {
		wheremerchant = fmt.Sprintf(`store.merchantcontrol = '%s' AND`, xDocrequest["merchant"])
	}

	sql := `select 	rewardstore.code as code, rewardstore.control as control, rewardstore.title as title, 
					rewardstore.description as description, rewardstore.workflow as workflow, 
					rewardstore.createdate as createdate, rewardstore.createdby as createdby, 
					rewardstore.updatedate as updatedate, rewardstore.updatedby as updatedby,
					
					rewardstore.storecontrol as storecontrol, store.title as storetitle, 
					store.code as storecode, store.city as storecity, store.address as storeaddress, 
					store.image as storeimage, store.contact as storecontact, 
					store.phone as storephone, store.email as storeemail, 
					store.gpslat as storegpslat, store.gpslong as storegpslong,
					store.flagship as storeflagship, 

					rewardstore.rewardcontrol as rewardcontrol,
					reward.title as rewardtitle

			from rewardstore, reward, store

			where %s %s %s  %s %s 

			reward.control = rewardstore.rewardcontrol 
			AND store.control = rewardstore.storecontrol

			order by title limit %s offset %s`

	sql = fmt.Sprintf(sql, whereworkflow, wherereward, wherestore, whereflagship, wheremerchant, limit, offset)
	mapRes, _ = curdb.Query(sql)
	return mapRes
}

func (this *RewardStore) Delete(xDocrequest map[string]interface{}, curdb Database) {

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
