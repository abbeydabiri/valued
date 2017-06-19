package frontend

import (
	"fmt"
	"time"
	"valued/database"
	"valued/functions"

	"net/http"
	"strconv"
)

func (this *Report) merchant(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	curdb.Query("set datestyle = dmy")
	mapReport := make(map[string]interface{})

	//---
	cFormat := "02/01/2006"
	sStopdate := time.Now().Format(cFormat)
	sStartdate := time.Now().Add(-(time.Hour * 24 * 365)).Format(cFormat)
	//---

	//---

	//
	//TOTAL GENERATED REVENUE
	sqlTotalRevenueInception := `select sum(transactionvalue) as totalrevenueinception from redemption where schemecontrol in (select control from scheme)`

	mapTotalRevenueInception, _ := curdb.Query(sqlTotalRevenueInception)
	mapReport["totalrevenueinception"] = float64(0)
	if mapTotalRevenueInception["1"] != nil {
		mapTotalRevenueInceptionxDoc := mapTotalRevenueInception["1"].(map[string]interface{})
		switch mapTotalRevenueInceptionxDoc["totalrevenueinception"].(type) {
		case string:
			mapReport["totalrevenueinception"] = float64(0)
		case int64:
			mapReport["totalrevenueinception"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRevenueInceptionxDoc["totalrevenueinception"].(int64))))
		case float64:
			mapReport["totalrevenueinception"] = functions.ThousandSeperator(functions.Round(mapTotalRevenueInceptionxDoc["totalrevenueinception"].(float64)))
		}
	}

	sqlTotalRevenue12months := fmt.Sprintf(`select sum(transactionvalue) as totalrevenue12months from redemption where schemecontrol in (select control from scheme)
		and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapTotalRevenue12months, _ := curdb.Query(sqlTotalRevenue12months)
	mapReport["totalrevenue12months"] = float64(0)
	if mapTotalRevenue12months["1"] != nil {
		mapTotalRevenue12monthsxDoc := mapTotalRevenue12months["1"].(map[string]interface{})
		switch mapTotalRevenue12monthsxDoc["totalrevenue12months"].(type) {
		case string:
			mapReport["totalrevenue12months"] = float64(0)
		case int64:
			mapReport["totalrevenue12months"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRevenue12monthsxDoc["totalrevenue12months"].(int64))))
		case float64:
			mapReport["totalrevenue12months"] = functions.ThousandSeperator(functions.Round(mapTotalRevenue12monthsxDoc["totalrevenue12months"].(float64)))
		}
	}
	//TOTAL GENERATED REVENUE
	//

	//---

	//
	//TOTAL REDEMPTIONS
	sqlTotalRedemptionInception := `select count(control) as totalredemptioninception from redemption where schemecontrol in (select control from scheme)`

	mapTotalRedemptionInception, _ := curdb.Query(sqlTotalRedemptionInception)
	mapReport["totalredemptioninception"] = float64(0)
	if mapTotalRedemptionInception["1"] != nil {
		mapTotalRedemptionInceptionxDoc := mapTotalRedemptionInception["1"].(map[string]interface{})
		switch mapTotalRedemptionInceptionxDoc["totalredemptioninception"].(type) {
		case string:
			mapReport["totalredemptioninception"] = float64(0)
		case int64:
			mapReport["totalredemptioninception"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRedemptionInceptionxDoc["totalredemptioninception"].(int64))))
		case float64:
			mapReport["totalredemptioninception"] = functions.ThousandSeperator(functions.Round(mapTotalRedemptionInceptionxDoc["totalredemptioninception"].(float64)))
		}
	}

	sqlTotalRedemption12months := fmt.Sprintf(`select count(control) as totalredemption12months from redemption where schemecontrol in (select control from scheme)
		and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapTotalRedemption12months, _ := curdb.Query(sqlTotalRedemption12months)
	mapReport["totalredemption12months"] = float64(0)
	if mapTotalRedemption12months["1"] != nil {
		mapTotalRedemption12monthsxDoc := mapTotalRedemption12months["1"].(map[string]interface{})
		switch mapTotalRedemption12monthsxDoc["totalredemption12months"].(type) {
		case string:
			mapReport["totalredemption12months"] = float64(0)
		case int64:
			mapReport["totalredemption12months"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRedemption12monthsxDoc["totalredemption12months"].(int64))))
		case float64:
			mapReport["totalredemption12months"] = functions.ThousandSeperator(functions.Round(mapTotalRedemption12monthsxDoc["totalredemption12months"].(float64)))
		}
	}
	//TOTAL REDEMPTIONS
	//

	//---

	//
	//ACTIVE MERCHANTS
	sqlActiveMerchantsNumber := `select count(distinct merchantcontrol) as activemerchantsnumber from redemption where schemecontrol in (select control from scheme)`

	mapActiveMerchantsNumber, _ := curdb.Query(sqlActiveMerchantsNumber)
	fACTIVEMERCHANTSNUMBER := float64(0)
	mapReport["activemerchantsnumber"] = float64(0)
	if mapActiveMerchantsNumber["1"] != nil {
		mapActiveMerchantsNumberxDoc := mapActiveMerchantsNumber["1"].(map[string]interface{})
		switch mapActiveMerchantsNumberxDoc["activemerchantsnumber"].(type) {
		case string:
			mapReport["activemerchantsnumber"] = float64(0)
		case int64:
			fACTIVEMERCHANTSNUMBER = functions.Round(float64(mapActiveMerchantsNumberxDoc["activemerchantsnumber"].(int64)))
			mapReport["activemerchantsnumber"] = functions.ThousandSeperator(fACTIVEMERCHANTSNUMBER)
		case float64:
			fACTIVEMERCHANTSNUMBER = functions.Round(mapActiveMerchantsNumberxDoc["activemerchantsnumber"].(float64))
			mapReport["activemerchantsnumber"] = functions.ThousandSeperator(fACTIVEMERCHANTSNUMBER)
		}
	}

	sqlActiveMerchantsPercentage := `select count(distinct control) as activemerchantspercentage from profile where control in (select distinct merchantcontrol from reward)`
	mapActiveMerchantsPercentage, _ := curdb.Query(sqlActiveMerchantsPercentage)
	mapReport["activemerchantspercentage"] = float64(0)
	if mapActiveMerchantsPercentage["1"] != nil {
		mapActiveMerchantsPercentagexDoc := mapActiveMerchantsPercentage["1"].(map[string]interface{})
		switch mapActiveMerchantsPercentagexDoc["activemerchantspercentage"].(type) {
		case string:
			mapReport["activemerchantspercentage"] = float64(0)
		case int64:
			mapReport["activemerchantspercentage"] = functions.RoundUp(float64(fACTIVEMERCHANTSNUMBER*100)/float64(mapActiveMerchantsPercentagexDoc["activemerchantspercentage"].(int64)), 0)
		case float64:
			mapReport["activemerchantspercentage"] = functions.RoundUp(float64(fACTIVEMERCHANTSNUMBER*100)/mapActiveMerchantsPercentagexDoc["activemerchantspercentage"].(float64), 0)
		}
	}
	//ACTIVE MERCHANTS
	//

	//---

	//
	//
	//TOP 5 MERCHANTS WITH HIGHEST REVENUE

	sqlTop5MerchantsRevenue := `select  (select title from profile where control = merchantcontrol) as merchant, sum(transactionvalue) as revenue
								from redemption where schemecontrol in  (select control from scheme)
								group by merchantcontrol order by 2 desc limit 5`

	mapTop5MerchantsRevenue, _ := curdb.Query(sqlTop5MerchantsRevenue)
	aTop5MerchantsRevenueSorted := functions.SortMap(mapTop5MerchantsRevenue)

	for _, sNumber := range aTop5MerchantsRevenueSorted {
		xDocReward := mapTop5MerchantsRevenue[sNumber].(map[string]interface{})
		xDocReward["row"] = sNumber

		switch xDocReward["revenue"].(type) {
		case int64:
			xDocReward["revenue"] = functions.ThousandSeperator(functions.Round(float64(xDocReward["revenue"].(int64))))
		case float64:
			xDocReward["revenue"] = functions.ThousandSeperator(functions.Round(xDocReward["revenue"].(float64)))
		}

		sTag := fmt.Sprintf(`%v#report-merchant-topfiverevenue-row`, sNumber)
		mapReport[sTag] = xDocReward
	}

	//TOP 5 MERCHANTS WITH HIGHEST REVENUE
	//
	//

	//---

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-merchant"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Merchant | Admin Reports","mainpanelContent":` + contentHTML + `}`))

}
