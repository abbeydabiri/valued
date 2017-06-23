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

	//
	//
	//AVERAGE NPS SCORE Based on Feedback Rating

	sqlAverageNPSscore := `select redemption.merchantcontrol, 
	(select title from profile where control = redemption.merchantcontrol) as merchant, 
	feedback.title as question, feedback.answer as answer, feedback.redemptioncontrol as redemptioncontrol
	from redemption left join feedback on feedback.redemptioncontrol = redemption.control order by merchantcontrol`

	resAverageNPSscore, _ := curdb.Query(sqlAverageNPSscore)
	mapAverageNPSscore := this.calculateNPS(resAverageNPSscore)

	totalNPSscore, averageNPSscore := float64(0), float64(0)
	for _, NPSscore := range mapAverageNPSscore {
		totalNPSscore += NPSscore
	}

	if len(mapAverageNPSscore) > 0 {
		averageNPSscore = functions.RoundUp(totalNPSscore/float64(len(mapAverageNPSscore)), 0)
	}
	mapReport["averagenpsscore"] = averageNPSscore

	//AVERAGE NPS SCORE Based on Feedback Rating
	//
	//

	//---

	//
	//
	//TOP 5 MERCHANTS NPS (Best & Worst) on Feedback Rating

	mapNPSscoreList := make(map[string]interface{})
	for sMerchant, NPSscore := range mapAverageNPSscore {
		sNPSscore := fmt.Sprintf("%v", NPSscore)

		mapNPSscore := make(map[string]interface{})
		mapNPSscore["merchant"] = sMerchant
		mapNPSscore["npsscore"] = sNPSscore
		mapNPSscoreList[sNPSscore] = mapNPSscore
	}

	mapSortedNPS := functions.SortMap(mapNPSscoreList)
	functions.SortDesc(mapSortedNPS)

	iBestFive := 1
	for iNumber, npsscore := range mapSortedNPS {
		if iBestFive <= 5 {
			sTag := fmt.Sprintf(`%v#report-merchant-npsscore-bestfive-row`, iNumber)
			mapReport[sTag] = mapNPSscoreList[npsscore].(map[string]interface{})
		}
		iBestFive++
	}

	iWorstFive := 1
	functions.SortAsc(mapSortedNPS)
	for iNumber, npsscore := range mapSortedNPS {
		if iWorstFive <= 5 {
			sTag := fmt.Sprintf(`%v#report-merchant-npsscore-worstfive-row`, iNumber)
			mapReport[sTag] = mapNPSscoreList[npsscore].(map[string]interface{})
		}
		iWorstFive++
	}

	//TOP 5 MERCHANTS NPS (Best & Worst) on Feedback Rating
	//
	//

	//
	//
	//GRAPHS: - MONTHLY REVENUE FOR LAST 12 MONTHS
	sLabel := "yyyy-Mon"
	sOrderBy := "yyyymm"

	curMonth := time.Now().Add(-(time.Hour * 24 * 365))
	sStartdateChart := curMonth.Format(cFormat)

	oneYearChart := time.Now()
	sStopdateChart := oneYearChart.Format(cFormat)

	revenueMonthlyReportGenerator := make(map[string]interface{})

	counter := 1
	monthLabelSeries := make(map[string]interface{})
	for oneYearChart.After(curMonth) {
		monthLabelSeries[curMonth.Format("200601")] = curMonth.Format("2006-Jan")
		curMonth = curMonth.Add(time.Hour * 24 * 30)
		counter++
	}

	for sOrderby, sLabel := range monthLabelSeries {
		sLabelIndex := fmt.Sprintf("%s#label", sOrderby)
		sSeriesIndex := fmt.Sprintf("%s#series", sOrderby)

		revenueMonthlyReportGenerator[sLabelIndex] = fmt.Sprintf(`"%s",`, sLabel)
		revenueMonthlyReportGenerator[sSeriesIndex] = functions.ThousandSeperator(functions.Round(float64(0))) + ","
	}

	sqlRevenueMonthlyChart := `select to_char(substring(createdate from 1 for 20)::timestamp,'%s') as orderby, 
	to_char(substring(createdate from 1 for 20)::timestamp,'%s') as label, sum(transactionvalue) as revenuemonthly from redemption 
	where substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp group by 1,2 order by 1`

	sqlRevenueMonthlyChart = fmt.Sprintf(sqlRevenueMonthlyChart, sOrderBy, sLabel, sStartdateChart, sStopdateChart)
	mapRevenueMonthlyChart, _ := curdb.Query(sqlRevenueMonthlyChart)

	revenueMonthlyReportGenerator["id"] = "revenuemonthly"
	revenueMonthlyHigh := float64(100)
	for _, xDoc := range mapRevenueMonthlyChart {
		xDoc := xDoc.(map[string]interface{})

		sLabel := fmt.Sprintf("%s#label", xDoc["orderby"])
		sSeries := fmt.Sprintf("%s#series", xDoc["orderby"])

		revenueMonthlyReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["label"])

		revenueMonthlyReportGenerator[sSeries] = fmt.Sprintf("%v,", xDoc["revenuemonthly"])
		if xDoc["revenuemonthly"].(float64) > revenueMonthlyHigh {
			revenueMonthlyHigh = xDoc["revenuemonthly"].(float64)
			revenueMonthlyHigh += float64(2)
		}

	}
	revenueMonthlyReportGenerator["high"] = revenueMonthlyHigh
	mapReport["report-merchant-barchart-monthly-revenue"] = revenueMonthlyReportGenerator
	//GRAPHS: -MONTHLY REVENUE FOR LAST 12 MONTHS
	//
	//

	//
	//
	//GRAPHS: - MONTHLY ACTIVE MERCHANT FOR LAST 12 MONTHS

	activemerchantMonthlyReportGenerator := make(map[string]interface{})

	for sOrderby, sLabel := range monthLabelSeries {
		sLabelIndex := fmt.Sprintf("%s#label", sOrderby)
		sSeriesIndex := fmt.Sprintf("%s#series", sOrderby)

		activemerchantMonthlyReportGenerator[sLabelIndex] = fmt.Sprintf(`"%s",`, sLabel)
		activemerchantMonthlyReportGenerator[sSeriesIndex] = functions.ThousandSeperator(functions.Round(float64(0))) + ","
	}

	sqlActiveMerchantMonthlyChart := `select to_char(substring(createdate from 1 for 20)::timestamp,'%s') as orderby, 
	to_char(substring(createdate from 1 for 20)::timestamp,'%s') as label, count(distinct merchantcontrol) as activemerchant,
	(select count(distinct control) from profile where control in (select distinct merchantcontrol from reward)) as totalmerchants
	from redemption where substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp group by 1,2 order by 1
	`

	sqlActiveMerchantMonthlyChart = fmt.Sprintf(sqlActiveMerchantMonthlyChart, sOrderBy, sLabel, sStartdateChart, sStopdateChart)
	mapActiveMerchantMonthlyChart, _ := curdb.Query(sqlActiveMerchantMonthlyChart)

	activemerchantMonthlyHigh := float64(1)
	activemerchantMonthlyReportGenerator["id"] = "activemerchantmonthly"
	for _, xDoc := range mapActiveMerchantMonthlyChart {
		xDoc := xDoc.(map[string]interface{})

		sLabel := fmt.Sprintf("%s#label", xDoc["orderby"])
		sSeries := fmt.Sprintf("%s#series", xDoc["orderby"])
		activemerchantMonthlyReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["label"])

		xDoc["activemerchantmonthly"] =
			functions.RoundUp(float64(xDoc["activemerchant"].(int64)*100)/float64(xDoc["totalmerchants"].(int64)), 0)

		activemerchantMonthlyReportGenerator[sSeries] = fmt.Sprintf("%v,", xDoc["activemerchantmonthly"])
		if xDoc["activemerchantmonthly"].(float64) > activemerchantMonthlyHigh {
			activemerchantMonthlyHigh = xDoc["activemerchantmonthly"].(float64)
			activemerchantMonthlyHigh += float64(2)
		}

	}
	activemerchantMonthlyReportGenerator["high"] = activemerchantMonthlyHigh
	mapReport["report-merchant-barchart-monthly-activemerchant"] = activemerchantMonthlyReportGenerator
	//GRAPHS: - MONTHLY ACTIVE MERCHANT FOR LAST 12 MONTHS
	//
	//

	//
	//
	//GRAPHS: - MONTHLY AVERAGE NPS SCORE (for the last 12 months)

	sqlMonthlyAverageNPSscore := `select to_char(substring(feedback.createdate from 1 for 20)::timestamp,'%s') as orderby, 	
	to_char(substring(feedback.createdate from 1 for 20)::timestamp,'%s') as label, redemption.merchantcontrol as merchantcontrol, 
	feedback.title as question, feedback.answer as answer, (select title from profile where control = redemption.merchantcontrol) as merchant
	from redemption  left join feedback on feedback.redemptioncontrol = redemption.control where feedback.title like '%%RECOMMEND%%'
	and substring(feedback.createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp
	order by 1`
	sqlMonthlyAverageNPSscore = fmt.Sprintf(sqlMonthlyAverageNPSscore, sOrderBy, sLabel, sStartdateChart, sStopdateChart)
	resMonthlyAverageNPSscore, _ := curdb.Query(sqlMonthlyAverageNPSscore)

	aSortedMonthlyAverageNPSscore := functions.SortMap(resMonthlyAverageNPSscore)
	mapMonthlyNPSFeedbackChart := make(map[string]interface{})

	sCurOrderby, sCurLabel := "", ""
	mapMontlyAverageNPSscore := make(map[string]interface{})
	for _, sNumber := range aSortedMonthlyAverageNPSscore {
		xDocFeedback := resMonthlyAverageNPSscore[sNumber].(map[string]interface{})

		if sCurOrderby == "" {
			sCurOrderby = xDocFeedback["orderby"].(string)
			sCurLabel = xDocFeedback["label"].(string)
			mapMontlyAverageNPSscore = make(map[string]interface{})
		}

		if sCurOrderby != xDocFeedback["orderby"].(string) {

			mapMonthlyNPSFeedbackValues := make(map[string]interface{})
			mapMonthlyNPSFeedbackValues["orderby"] = sCurOrderby
			mapMonthlyNPSFeedbackValues["label"] = sCurLabel

			//Calculate NPS score
			mapAverageNPSscore := this.calculateNPS(mapMontlyAverageNPSscore)
			totalNPSscore, averageNPSscore := float64(0), float64(0)
			for _, NPSscore := range mapAverageNPSscore {
				totalNPSscore += NPSscore
			}
			if len(mapAverageNPSscore) > 0 {
				averageNPSscore = functions.RoundUp(totalNPSscore/float64(len(mapAverageNPSscore)), 0)
			}
			mapMonthlyNPSFeedbackValues["averagenpsscoremonthly"] = averageNPSscore
			mapMonthlyNPSFeedbackChart[sCurOrderby] = mapMonthlyNPSFeedbackValues
			//Calculate NPS score
		}

		mapMontlyAverageNPSscore[sNumber] = xDocFeedback

		sCurOrderby = xDocFeedback["orderby"].(string)
		sCurLabel = xDocFeedback["label"].(string)
	}

	if sCurOrderby != "" {
		mapMonthlyNPSFeedbackValues := make(map[string]interface{})
		mapMonthlyNPSFeedbackValues["orderby"] = sCurOrderby
		mapMonthlyNPSFeedbackValues["label"] = sCurLabel

		//Calculate NPS score
		mapAverageNPSscore := this.calculateNPS(mapMontlyAverageNPSscore)
		totalNPSscore, averageNPSscore := float64(0), float64(0)
		for _, NPSscore := range mapAverageNPSscore {
			totalNPSscore += NPSscore
		}
		if len(mapAverageNPSscore) > 0 {
			averageNPSscore = functions.RoundUp(totalNPSscore/float64(len(mapAverageNPSscore)), 0)
		}
		mapMonthlyNPSFeedbackValues["averagenpsscoremonthly"] = averageNPSscore
		mapMonthlyNPSFeedbackChart[sCurOrderby] = mapMonthlyNPSFeedbackValues
		//Calculate NPS score
	}

	averagenpsscoreMonthlyReportGenerator := make(map[string]interface{})

	for sOrderby, sLabel := range monthLabelSeries {
		sLabelIndex := fmt.Sprintf("%s#label", sOrderby)
		sSeriesIndex := fmt.Sprintf("%s#series", sOrderby)

		averagenpsscoreMonthlyReportGenerator[sLabelIndex] = fmt.Sprintf(`"%s",`, sLabel)
		averagenpsscoreMonthlyReportGenerator[sSeriesIndex] = functions.ThousandSeperator(functions.Round(float64(0))) + ","
	}

	averagenpsscoreMonthlyReportGenerator["id"] = "averagenpsscoremonthly"
	averagenpsscoreMonthlyHigh := float64(10)
	for _, xDoc := range mapMonthlyNPSFeedbackChart {
		xDoc := xDoc.(map[string]interface{})

		sLabel := fmt.Sprintf("%s#label", xDoc["orderby"])
		sSeries := fmt.Sprintf("%s#series", xDoc["orderby"])

		averagenpsscoreMonthlyReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["label"])

		averagenpsscoreMonthlyReportGenerator[sSeries] = fmt.Sprintf("%v,", xDoc["averagenpsscoremonthly"])
		if xDoc["averagenpsscoremonthly"].(float64) > averagenpsscoreMonthlyHigh {
			averagenpsscoreMonthlyHigh = xDoc["averagenpsscoremonthly"].(float64)
			averagenpsscoreMonthlyHigh += float64(2)
		}

	}
	averagenpsscoreMonthlyReportGenerator["high"] = averagenpsscoreMonthlyHigh
	mapReport["report-merchant-barchart-monthly-averagenpsscore"] = averagenpsscoreMonthlyReportGenerator

	//GRAPHS: - MONTHLY AVERAGE NPS SCORE (for the last 12 months)
	//
	//

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-merchant"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Merchant | Admin Reports","mainpanelContent":` + contentHTML + `}`))

}
