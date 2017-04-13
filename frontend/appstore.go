package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"html"
	"net/http"
	"strconv"
)

type AppStore struct {
	functions.Templates
	mapAppCache map[string]interface{}
	pageMap     map[string]interface{}
}

func (this *AppStore) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapAppCache = curdb.GetSession(GOSESSID.Value, "mapAppCache")

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "search":
		this.pageMap = make(map[string]interface{})
		this.pageMap["app-store"] = this.search(httpReq, curdb)

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Valued Rewards","pageContent":` + contentHTML + `}`))
		return
	}
}

func (this *AppStore) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {

	formSearch := make(map[string]interface{})

	appFooter := make(map[string]interface{})
	formSearch["app-footer"] = appFooter
	appFooter["brand"] = "white"

	//Search Rewards
	tblReward := new(database.Reward)
	xDocrequest := make(map[string]interface{})

	xDocrequest["merchantcontrol"] = html.EscapeString(httpReq.FormValue("merchant"))
	xDocrequest["workflow"] = "active"
	xDocrequest["limit"] = "200"

	xDocresult := tblReward.Search(xDocrequest, curdb)

	favoriteItem := ""
	favoriteList := new(AppFavorite).List(this.mapAppCache, curdb)

	myAppRedeem := new(AppRedeem)
	myAppRedeem.SetAppCache(this.mapAppCache)

	//Logic to Manage Grouping
	//Fetch all Groupings and
	sqlGrouping := `select rewardgroup.rewardcontrol as rewardcontrol, membergroup.membercontrol as membercontrol , 
							groups.control as groupcontrol, groups.title as grouptitle 
					from groups 
							left join rewardgroup on rewardgroup.groupcontrol = groups.control 
							left join membergroup on membergroup.groupcontrol = groups.control
					`

	xDocGrouping, _ := curdb.Query(sqlGrouping)
	mapRewardGroup := make(map[string]interface{})
	mapMemberGroup := make(map[string]interface{})

	for _, xDocGroup := range xDocGrouping {
		xDocGroup := xDocGroup.(map[string]interface{})

		//Create Map with RewardControl as Index
		//Create Map with RewardControl and MemberControl as Index
		mapRewardGroup[xDocGroup["rewardcontrol"].(string)] = xDocGroup

		sMemberIndex := fmt.Sprintf("%s-%s", xDocGroup["rewardcontrol"], xDocGroup["membercontrol"])
		mapMemberGroup[sMemberIndex] = xDocGroup

	}

	//Logic to Manage Grouping

	sGroup := ""
	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber

		sProfilecontrol := ""
		if this.mapAppCache["control"] != nil {
			sProfilecontrol = this.mapAppCache["control"].(string)
			xDoc["signedIn"] = "yes"
			_, xDoc["state"] = myAppRedeem.ValidateEligibility(xDoc["control"].(string), curdb)
		}

		favoriteItem = fmt.Sprintf("reward-%s", xDoc["control"])
		if favoriteList[favoriteItem] == nil {
			xDoc["heart"] = "heart.png"
		} else {
			xDoc["heart"] = "heart_filled.png"
		}

		sGroup = "00000000"
		if mapRewardGroup[xDoc["control"].(string)] != nil {
			sGroup = ""
			sMemberIndex := fmt.Sprintf("%s-%s", xDoc["control"], sProfilecontrol)
			if mapMemberGroup[sMemberIndex] != nil {
				sGroup = "0"
				xDocGroup := mapMemberGroup[sMemberIndex].(map[string]interface{})
				xDoc["grouptitle"] = fmt.Sprintf(" Group: %s", xDocGroup["grouptitle"])
			}
		}

		if sGroup != "" {
			formSearch[sGroup+cNumber+"#app-reward-list"] = xDoc
		}
	}

	//Search Rewards

	xDocMerchant := make(map[string]interface{})
	//Search Stores & Get Flagship Store for Presentation
	tblAppStore := new(database.Store)
	locatorMarker := make(map[string]interface{})
	xDocrequest = make(map[string]interface{})

	xDocrequest["limit"] = "200"
	xDocrequest["merchantcontrol"] = html.EscapeString(httpReq.FormValue("merchant"))
	xDocrequest["workflow"] = "active"

	xDocresult = tblAppStore.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {

		xDoc := xDoc.(map[string]interface{})
		locatorMarker[cNumber+"#app-map-locator-marker"] = xDoc

		if cNumber == "1" || xDoc["flagship"].(string) == "Yes" {
			xDocMerchant = xDoc
		}
	}

	if len(xDocMerchant) > 0 {

		formSearch["merchanttitle"] = xDocMerchant["merchanttitle"]
		formSearch["merchantimage"] = xDocMerchant["merchantimage"]
		formSearch["merchantphone"] = xDocMerchant["merchantphone"]
		formSearch["merchantemail"] = xDocMerchant["merchantemail"]
		formSearch["merchantwebsite"] = xDocMerchant["merchantwebsite"]
		formSearch["merchantdescription"] = xDocMerchant["merchantdescription"]

		formSearch["title"] = xDocMerchant["title"]
		formSearch["address"] = xDocMerchant["address"]

		if functions.TrimEscape(formSearch["address"].(string)) != "" {
			switch len(xDocresult) {
			default:
				locatorMarker["zoom"] = 5
			case 1, 2:
				locatorMarker["zoom"] = 3
			}
			formSearch["app-map-locator"] = locatorMarker
		}
	}

	//Search Stores & Get Flagship Store for Presentation

	return formSearch
}
