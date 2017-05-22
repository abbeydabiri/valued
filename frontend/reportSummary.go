package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (this *Report) summary(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	curdb.Query("set datestyle = dmy")
	mapReport := make(map[string]interface{})

	//Number of Subscriptions
	//Total
	sqlTotalSubscription := `select count(control) as totalsubscription from subscription`
	mapTotalSubscription, _ := curdb.Query(sqlTotalSubscription)
	mapReport["totalsubscription"] = float64(0)
	if mapTotalSubscription["1"] != nil {
		mapTotalSubscriptionxDoc := mapTotalSubscription["1"].(map[string]interface{})
		switch mapTotalSubscriptionxDoc["totalsubscription"].(type) {
		case string:
			mapReport["totalsubscription"] = float64(0)
		case int64:
			mapReport["totalsubscription"] = functions.ThousandSeperator(functions.Round(float64(mapTotalSubscriptionxDoc["totalsubscription"].(int64))))
		case float64:
			mapReport["totalsubscription"] = functions.ThousandSeperator(functions.Round(mapTotalSubscriptionxDoc["totalsubscription"].(float64)))
		}
	}

	//Employer
	sqlEmployerSubscription := `select count(control) as employersubscription from subscription where employercontrol not in (select control from profile where code in ('main','britishmums')) `
	mapEmployerSubscription, _ := curdb.Query(sqlEmployerSubscription)
	mapReport["employersubscription"] = float64(0)
	if mapEmployerSubscription["1"] != nil {
		mapEmployerSubscriptionxDoc := mapEmployerSubscription["1"].(map[string]interface{})
		switch mapEmployerSubscriptionxDoc["employersubscription"].(type) {
		case string:
			mapReport["employersubscription"] = float64(0)
		case int64:
			mapReport["employersubscription"] = functions.ThousandSeperator(functions.Round(float64(mapEmployerSubscriptionxDoc["employersubscription"].(int64))))
		case float64:
			mapReport["employersubscription"] = functions.ThousandSeperator(functions.Round(mapEmployerSubscriptionxDoc["employersubscription"].(float64)))
		}
	}

	//BritishMum
	sqlBritishMumSubscription := `select count(control) as britishmumssubscription from subscription where employercontrol in (select control from profile where code in ('britishmums')) `
	mapBritishMumSubscription, _ := curdb.Query(sqlBritishMumSubscription)
	mapReport["britishmumssubscription"] = float64(0)
	if mapBritishMumSubscription["1"] != nil {
		mapBritishMumSubscriptionxDoc := mapBritishMumSubscription["1"].(map[string]interface{})
		switch mapBritishMumSubscriptionxDoc["britishmumssubscription"].(type) {
		case string:
			mapReport["britishmumssubscription"] = float64(0)
		case int64:
			mapReport["britishmumssubscription"] = functions.ThousandSeperator(functions.Round(float64(mapBritishMumSubscriptionxDoc["britishmumssubscription"].(int64))))
		case float64:
			mapReport["britishmumssubscription"] = functions.ThousandSeperator(functions.Round(mapBritishMumSubscriptionxDoc["britishmumssubscription"].(float64)))
		}
	}

	//Other
	sqlOtherSubscription := `select count(control) as othersubscription from subscription where employercontrol in (select control from profile where code in ('main')) `
	mapOtherSubscription, _ := curdb.Query(sqlOtherSubscription)
	mapReport["othersubscription"] = float64(0)
	if mapOtherSubscription["1"] != nil {
		mapOtherSubscriptionxDoc := mapOtherSubscription["1"].(map[string]interface{})
		switch mapOtherSubscriptionxDoc["othersubscription"].(type) {
		case string:
			mapReport["othersubscription"] = float64(0)
		case int64:
			mapReport["othersubscription"] = functions.ThousandSeperator(functions.Round(float64(mapOtherSubscriptionxDoc["othersubscription"].(int64))))
		case float64:
			mapReport["othersubscription"] = functions.ThousandSeperator(functions.Round(mapOtherSubscriptionxDoc["othersubscription"].(float64)))
		}
	}
	//Number of Subscriptions

	//

	cFormat := "02/01/2006"
	iMonth, _ := strconv.Atoi(time.Now().Format(cFormat)[3:5])
	iMonth++
	if iMonth == 13 {
		iMonth = 1
	}

	sMonth := fmt.Sprintf("%v", iMonth)
	if iMonth < 10 {
		sMonth = "0" + sMonth
	}

	startDate, _ := time.Parse(cFormat, "01/"+sMonth+"/"+time.Now().Format(cFormat)[6:])
	oneMonth := startDate.Add(time.Hour * 24 * 30)

	//Number of Subscriptions Due
	//Total
	sqlTotalSubscriptionDue := fmt.Sprintf(`select count(control) as totalsubscriptiondue from subscription where expirydate::timestamp between '%s'::timestamp and '%s'::timestamp `, startDate.Format(cFormat), oneMonth.Format(cFormat))
	mapTotalSubscriptionDue, _ := curdb.Query(sqlTotalSubscriptionDue)
	mapReport["totalsubscriptiondue"] = float64(0)
	if mapTotalSubscriptionDue["1"] != nil {
		mapTotalSubscriptionDuexDoc := mapTotalSubscriptionDue["1"].(map[string]interface{})
		switch mapTotalSubscriptionDuexDoc["totalsubscriptiondue"].(type) {
		case string:
			mapReport["totalsubscriptiondue"] = float64(0)
		case int64:
			mapReport["totalsubscriptiondue"] = functions.ThousandSeperator(functions.Round(float64(mapTotalSubscriptionDuexDoc["totalsubscriptiondue"].(int64))))
		case float64:
			mapReport["totalsubscriptiondue"] = functions.ThousandSeperator(functions.Round(mapTotalSubscriptionDuexDoc["totalsubscriptiondue"].(float64)))
		}
	}

	//Employer

	sqlEmployerSubscriptionDue := fmt.Sprintf(`select count(control) as employersubscriptiondue from subscription where employercontrol not in (select control from profile where code in ('main','britishmums')) and expirydate::timestamp between '%s'::timestamp and '%s'::timestamp `, startDate.Format(cFormat), oneMonth.Format(cFormat))
	mapEmployerSubscriptionDue, _ := curdb.Query(sqlEmployerSubscriptionDue)
	mapReport["employersubscriptiondue"] = float64(0)
	if mapEmployerSubscriptionDue["1"] != nil {
		mapEmployerSubscriptionDuexDoc := mapEmployerSubscriptionDue["1"].(map[string]interface{})
		switch mapEmployerSubscriptionDuexDoc["employersubscriptiondue"].(type) {
		case string:
			mapReport["employersubscriptiondue"] = float64(0)
		case int64:
			mapReport["employersubscriptiondue"] = functions.ThousandSeperator(functions.Round(float64(mapEmployerSubscriptionDuexDoc["employersubscriptiondue"].(int64))))
		case float64:
			mapReport["employersubscriptiondue"] = functions.ThousandSeperator(functions.Round(mapEmployerSubscriptionDuexDoc["employersubscriptiondue"].(float64)))
		}
	}

	//BritishMum

	sqlBritishMumSubscriptionDue := fmt.Sprintf(`select count(control) as britishmumssubscriptiondue from subscription where employercontrol in (select control from profile where code in ('britishmums')) and expirydate::timestamp between '%s'::timestamp and '%s'::timestamp `, startDate.Format(cFormat), oneMonth.Format(cFormat))
	mapBritishMumSubscriptionDue, _ := curdb.Query(sqlBritishMumSubscriptionDue)
	mapReport["britishmumssubscriptiondue"] = float64(0)
	if mapBritishMumSubscriptionDue["1"] != nil {
		mapBritishMumSubscriptionDuexDoc := mapBritishMumSubscriptionDue["1"].(map[string]interface{})
		switch mapBritishMumSubscriptionDuexDoc["britishmumssubscriptiondue"].(type) {
		case string:
			mapReport["britishmumssubscriptiondue"] = float64(0)
		case int64:
			mapReport["britishmumssubscriptiondue"] = functions.ThousandSeperator(functions.Round(float64(mapBritishMumSubscriptionDuexDoc["britishmumssubscriptiondue"].(int64))))
		case float64:
			mapReport["britishmumssubscriptiondue"] = functions.ThousandSeperator(functions.Round(mapBritishMumSubscriptionDuexDoc["britishmumssubscriptiondue"].(float64)))
		}
	}

	//Other

	sqlOtherSubscriptionDue := fmt.Sprintf(`select count(control) as othersubscriptiondue from subscription where employercontrol in (select control from profile where code in ('main')) and expirydate::timestamp between '%s'::timestamp and '%s'::timestamp `, startDate.Format(cFormat), oneMonth.Format(cFormat))
	mapOtherSubscriptionDue, _ := curdb.Query(sqlOtherSubscriptionDue)
	mapReport["othersubscriptiondue"] = float64(0)
	if mapOtherSubscriptionDue["1"] != nil {
		mapOtherSubscriptionDuexDoc := mapOtherSubscriptionDue["1"].(map[string]interface{})
		switch mapOtherSubscriptionDuexDoc["othersubscriptiondue"].(type) {
		case string:
			mapReport["othersubscriptiondue"] = float64(0)
		case int64:
			mapReport["othersubscriptiondue"] = functions.ThousandSeperator(functions.Round(float64(mapOtherSubscriptionDuexDoc["othersubscriptiondue"].(int64))))
		case float64:
			mapReport["othersubscriptiondue"] = functions.ThousandSeperator(functions.Round(mapOtherSubscriptionDuexDoc["othersubscriptiondue"].(float64)))
		}
	}
	//Number of Subscriptions Due

	//Number of Redeemed Rewards
	//Total
	sqlTotalRedeemed := `select count(control) as totalredeemed from redemption`
	mapTotalRedeemed, _ := curdb.Query(sqlTotalRedeemed)
	mapReport["totalredeemed"] = float64(0)
	if mapTotalRedeemed["1"] != nil {
		mapTotalRedeemedxDoc := mapTotalRedeemed["1"].(map[string]interface{})
		switch mapTotalRedeemedxDoc["totalredeemed"].(type) {
		case string:
			mapReport["totalredeemed"] = float64(0)
		case int64:
			mapReport["totalredeemed"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRedeemedxDoc["totalredeemed"].(int64))))
		case float64:
			mapReport["totalredeemed"] = functions.ThousandSeperator(functions.Round(mapTotalRedeemedxDoc["totalredeemed"].(float64)))
		}
	}

	//Employer
	sqlEmployerRedeemed := `select count(control) as employerredeemed from redemption where employercontrol not in (select control from profile where code in ('main','britishmums')) `
	mapEmployerRedeemed, _ := curdb.Query(sqlEmployerRedeemed)
	mapReport["employerredeemed"] = float64(0)
	if mapEmployerRedeemed["1"] != nil {
		mapEmployerRedeemedxDoc := mapEmployerRedeemed["1"].(map[string]interface{})
		switch mapEmployerRedeemedxDoc["employerredeemed"].(type) {
		case string:
			mapReport["employerredeemed"] = float64(0)
		case int64:
			mapReport["employerredeemed"] = functions.ThousandSeperator(functions.Round(float64(mapEmployerRedeemedxDoc["employerredeemed"].(int64))))
		case float64:
			mapReport["employerredeemed"] = functions.ThousandSeperator(functions.Round(mapEmployerRedeemedxDoc["employerredeemed"].(float64)))
		}
	}

	//BritishMum
	sqlBritishMumRedeemed := `select count(control) as britishmumsredeemed from redemption where employercontrol in (select control from profile where code in ('britishmums')) `
	mapBritishMumRedeemed, _ := curdb.Query(sqlBritishMumRedeemed)
	mapReport["britishmumsredeemed"] = float64(0)
	if mapBritishMumRedeemed["1"] != nil {
		mapBritishMumRedeemedxDoc := mapBritishMumRedeemed["1"].(map[string]interface{})
		switch mapBritishMumRedeemedxDoc["britishmumsredeemed"].(type) {
		case string:
			mapReport["britishmumsredeemed"] = float64(0)
		case int64:
			mapReport["britishmumsredeemed"] = functions.ThousandSeperator(functions.Round(float64(mapBritishMumRedeemedxDoc["britishmumsredeemed"].(int64))))
		case float64:
			mapReport["britishmumsredeemed"] = functions.ThousandSeperator(functions.Round(mapBritishMumRedeemedxDoc["britishmumsredeemed"].(float64)))
		}
	}

	//Other
	sqlOtherRedeemed := `select count(control) as otherredeemed from redemption where employercontrol in (select control from profile where code in ('main')) `
	mapOtherRedeemed, _ := curdb.Query(sqlOtherRedeemed)
	mapReport["otherredeemed"] = float64(0)
	if mapOtherRedeemed["1"] != nil {
		mapOtherRedeemedxDoc := mapOtherRedeemed["1"].(map[string]interface{})
		switch mapOtherRedeemedxDoc["otherredeemed"].(type) {
		case string:
			mapReport["otherredeemed"] = float64(0)
		case int64:
			mapReport["otherredeemed"] = functions.ThousandSeperator(functions.Round(float64(mapOtherRedeemedxDoc["otherredeemed"].(int64))))
		case float64:
			mapReport["otherredeemed"] = functions.ThousandSeperator(functions.Round(mapOtherRedeemedxDoc["otherredeemed"].(float64)))
		}
	}
	//Number of Redeemed Rewards

	//

	mapReport["totalredeemed"] = 0
	mapReport["employerredeemed"] = 0
	mapReport["britishmumsredeemed"] = 0
	mapReport["otherredeemed"] = 0

	// sqlRevenueTotal := `select sum(transactionvalue) as revenue
	// 					from redemption where merchantcontrol = '%s'
	// 					and substring(createdate from 1 for 20)::timestamp between
	// 					'%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp
	// 				`
	// cFormat := "02/01/2006"
	// startDate, _ := time.Parse(cFormat, this.AdminCreatedate)
	// diffYears := int64(functions.GetDifferenceInSeconds(functions.GetSystemDate(), this.AdminCreatedate)) / int64(time.Hour*24*365)
	// if diffYears > 0 {
	// 	startDate = startDate.Add(time.Hour * 24 * 365 * time.Duration(diffYears))
	// }
	// oneYear := startDate.Add(time.Hour * 24 * 365)

	// curdb.Query("set datestyle = dmy")

	// sqlRevenueYear := fmt.Sprintf(sqlRevenueTotal, this.AdminControl, startDate.Format(cFormat), oneYear.Format(cFormat))
	// mapRevenueYear, _ := curdb.Query(sqlRevenueYear)

	// mapReport["revenueyear"] = float64(0)
	// if mapRevenueYear["1"] != nil {
	// 	mapRevenue := mapRevenueYear["1"].(map[string]interface{})
	// 	switch mapRevenue["revenue"].(type) {
	// 	case string:
	// 		mapReport["revenueyear"] = float64(0)
	// 	case int64:
	// 		mapReport["revenueyear"] = functions.ThousandSeperator(functions.Round(float64(mapRevenue["revenue"].(int64))))
	// 	case float64:
	// 		mapReport["revenueyear"] = functions.ThousandSeperator(functions.Round(mapRevenue["revenue"].(float64)))
	// 	}
	// }

	// //

	// sqlRedeemedTotal := `select count(control) as redeemedyear
	// 					from redemption where merchantcontrol = '%s'
	// 					and substring(createdate from 1 for 20)::timestamp between
	// 					'%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp
	// 				`

	// sqlRedeemedYear := fmt.Sprintf(sqlRedeemedTotal, this.AdminControl, startDate.Format(cFormat), oneYear.Format(cFormat))
	// mapRedeemedYear, _ := curdb.Query(sqlRedeemedYear)

	// mapReport["redeemedyear"] = float64(0)
	// if mapRedeemedYear["1"] != nil {
	// 	mapRedeemed := mapRedeemedYear["1"].(map[string]interface{})
	// 	switch mapRedeemed["redeemedyear"].(type) {
	// 	case string:
	// 		mapReport["redeemedyear"] = float64(0)
	// 	case int64:
	// 		mapReport["redeemedyear"] = functions.ThousandSeperator(functions.Round(float64(mapRedeemed["redeemedyear"].(int64))))
	// 	case float64:
	// 		mapReport["redeemedyear"] = functions.ThousandSeperator(functions.Round(mapRedeemed["redeemedyear"].(float64)))
	// 	}
	// }

	// //--

	// //

	// //Get NPS Score Based on Feedback Rating
	// sqlFeedback := `select merchant.title as question, merchant.answer as answer, merchant.redemptioncontrol as redemptioncontrol
	// 					from redemption
	// 						left join merchant on merchant.redemptioncontrol = redemption.control
	// 					where redemption.merchantcontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp
	// 				`
	// sqlFeedback = fmt.Sprintf(sqlFeedback, this.AdminControl, startDate.Format(cFormat), oneYear.Format(cFormat))

	// curdb.Query("set datestyle = dmy")
	// mapFeedback, _ := curdb.Query(sqlFeedback)

	// iNPSTotal := float64(0)
	// iNPSPositive := float64(0)
	// iNPSNegative := float64(0)
	// //100/(iNPSPositive+iNPSNegative)*(iNPSPositive-iNPSNegative)
	// //100/5*(2-3)

	// ratingCategory := make([]int, 11)
	// improveCategory := make(map[string]int)

	// for _, xDoc := range mapFeedback {
	// 	xDoc := xDoc.(map[string]interface{})

	// 	switch {
	// 	case strings.Contains(xDoc["question"].(string), "IMPROVEMENT"):
	// 		improveCategory[xDoc["answer"].(string)]++

	// 	case strings.Contains(xDoc["question"].(string), "RECOMMEND"):
	// 		score, _ := strconv.Atoi(xDoc["answer"].(string))
	// 		switch {
	// 		case score <= 6:
	// 			iNPSNegative++
	// 			break
	// 		case score >= 9:
	// 			iNPSPositive++
	// 			break
	// 		}
	// 		iNPSTotal++
	// 		ratingCategory[score]++
	// 	}
	// }

	// //Get NPS Score Based on Feedback Rating

	// iNPSNegativePercentage := float64(0)
	// iNPSPositivePercentage := float64(0)
	// if iNPSTotal > 0 {
	// 	iNPSPositivePercentage = (iNPSPositive / iNPSTotal) * 100
	// 	iNPSNegativePercentage = (iNPSNegative / iNPSTotal) * 100
	// }

	// mapReport["npsscore"] = functions.RoundUp(iNPSPositivePercentage-iNPSNegativePercentage, 0)

	// //
	// //

	// //BarChart: 12 Months Revenue & Redemption
	// sLabel := "yyyy-Mon"
	// sOrderBy := "yyyymm"

	// sStartdate := startDate.Format(cFormat)
	// sStopdate := oneYear.Format(cFormat)
	// mapReport["startdate"] = sStartdate
	// mapReport["stopdate"] = sStopdate

	// revenueReportGenerator := make(map[string]interface{})
	// redemptionReportGenerator := make(map[string]interface{})

	// counter := 1
	// curMonth := startDate
	// monthLabelSeries := make(map[string]interface{})
	// for oneYear.After(curMonth) {
	// 	monthLabelSeries[curMonth.Format("200601")] = curMonth.Format("2006-Jan")
	// 	curMonth = curMonth.Add(time.Hour * 24 * 30)
	// 	counter++
	// }

	// for sOrderby, sLabel := range monthLabelSeries {
	// 	sLabelIndex := fmt.Sprintf("%s#label", sOrderby)
	// 	sSeriesIndex := fmt.Sprintf("%s#series", sOrderby)

	// 	revenueReportGenerator[sLabelIndex] = fmt.Sprintf(`"%s",`, sLabel)
	// 	revenueReportGenerator[sSeriesIndex] = functions.ThousandSeperator(functions.Round(float64(0))) + ","

	// 	redemptionReportGenerator[sLabelIndex] = fmt.Sprintf(`"%s",`, sLabel)
	// 	redemptionReportGenerator[sSeriesIndex] = functions.ThousandSeperator(functions.Round(float64(0))) + ","
	// }

	// sqlRevenue := `select to_char(substring(createdate from 1 for 20)::timestamp,'%s') as orderby, to_char(substring(createdate from 1 for 20)::timestamp,'%s') as label, sum(transactionvalue) as revenue
	// 				from redemption where merchantcontrol = '%s' and substring(createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp
	// 				group by 1,2 order by 1
	// 				`
	// sqlRevenue = fmt.Sprintf(sqlRevenue, sOrderBy, sLabel, this.AdminControl, sStartdate, sStopdate)

	// sqlRedemption := `select to_char(substring(createdate from 1 for 20)::timestamp,'%s') as orderby, to_char(substring(createdate from 1 for 20)::timestamp,'%s') as label, count(control) as redemption
	// 				from redemption where merchantcontrol = '%s' and substring(createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp
	// 				group by 1,2 order by 1
	// 				`
	// sqlRedemption = fmt.Sprintf(sqlRedemption, sOrderBy, sLabel, this.AdminControl, sStartdate, sStopdate)

	// curdb.Query("set datestyle = dmy")
	// mapRedemption, _ := curdb.Query(sqlRedemption)
	// mapRevenue, _ := curdb.Query(sqlRevenue)

	// revenueReportGenerator["id"] = "revenue"
	// revenueReportGenerator["title"] = "REVENUE"
	// revenueHigh := float64(100)
	// for _, xDoc := range mapRevenue {
	// 	xDoc := xDoc.(map[string]interface{})
	// 	// fmt.Printf("[%v]\n", xDoc["revenue"])
	// 	// if fmt.Sprintf("%v", xDoc["revenue"]) == "" || fmt.Sprintf("%v", xDoc["revenue"]) == "0" {
	// 	// 	continue
	// 	// }

	// 	sLabel := fmt.Sprintf("%s#label", xDoc["orderby"])
	// 	sSeries := fmt.Sprintf("%s#series", xDoc["orderby"])

	// 	// if revenueReportGenerator[sLabel] != nil {
	// 	revenueReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["label"])
	// 	// }

	// 	// if revenueReportGenerator[sSeries] != nil {
	// 	revenueReportGenerator[sSeries] = fmt.Sprintf("%v,", xDoc["revenue"])
	// 	if xDoc["revenue"].(float64) > revenueHigh {
	// 		revenueHigh = xDoc["revenue"].(float64)
	// 		revenueHigh += float64(2)
	// 	}
	// 	// }
	// }
	// revenueReportGenerator["high"] = revenueHigh
	// mapReport["1#report-generator-barchart"] = revenueReportGenerator

	// redemptionReportGenerator["id"] = "redemption"
	// redemptionReportGenerator["title"] = "REDEMPTION"
	// redemptionHigh := int64(5)
	// for _, xDoc := range mapRedemption {
	// 	xDoc := xDoc.(map[string]interface{})
	// 	// fmt.Printf("[%v]\n", xDoc["redemption"])
	// 	// if fmt.Sprintf("%v,", xDoc["redemption"]) == "" || fmt.Sprintf("%v,", xDoc["redemption"]) == "0" {
	// 	// 	continue
	// 	// }

	// 	sLabel := fmt.Sprintf("%s#label", xDoc["orderby"])
	// 	sSeries := fmt.Sprintf("%s#series", xDoc["orderby"])

	// 	if redemptionReportGenerator[sLabel] != nil {
	// 		redemptionReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["label"])
	// 	}

	// 	if redemptionReportGenerator[sSeries] != nil {
	// 		redemptionReportGenerator[sSeries] = fmt.Sprintf("%v,", xDoc["redemption"])
	// 		if xDoc["redemption"].(int64) > redemptionHigh {
	// 			redemptionHigh = xDoc["redemption"].(int64)
	// 			redemptionHigh += int64(2)
	// 		}
	// 	}
	// }
	// redemptionReportGenerator["high"] = redemptionHigh
	// mapReport["2#report-generator-barchart"] = redemptionReportGenerator
	//BarChart: 12 Months Revenue & Redemption

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-summary"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Summary | Merchant Reports","mainpanelContent":` + contentHTML + `}`))

}
