package frontend

import (
	"valued/database"
	"valued/functions"

	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Report struct {
	AdminCreatedate string
	AdminControl    string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *Report) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.AdminControl = this.mapCache["control"].(string)
	this.AdminCreatedate = this.mapCache["createdate"].(string)[:10]
	if this.mapCache["company"].(string) != "Yes" {
		this.AdminControl = this.mapCache["employercontrol"].(string)
		this.AdminCreatedate = this.mapCache["employercreatedate"].(string)[:10]
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "", "summary":
		this.summary(httpRes, httpReq, curdb)
		return

	case "downloadReport":
		this.downloadReport(httpRes, httpReq, curdb)
		return

	case "subscription":
		this.subscription(httpRes, httpReq, curdb)
		return

	case "merchant":
		this.merchant(httpRes, httpReq, curdb)
		return

	}
}

func (this *Report) downloadReport(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	cFormat := "02/01/2006"
	startDate, _ := time.Parse(cFormat, this.AdminCreatedate)
	diffYears := int64(functions.GetDifferenceInSeconds(functions.GetSystemDate(), this.AdminCreatedate)) / int64(time.Hour*24*365)
	if diffYears > 0 {
		startDate = startDate.Add(time.Hour * 24 * 365 * time.Duration(diffYears))
	}
	oneYear := startDate.Add(time.Hour * 24 * 365)

	sStartdate := startDate.Format(cFormat)
	sStopdate := oneYear.Format(cFormat)

	if functions.TrimEscape(httpReq.FormValue("startdate")) != "" {
		sStartdate = functions.TrimEscape(httpReq.FormValue("startdate"))
	}

	if functions.TrimEscape(httpReq.FormValue("stopdate")) != "" {
		sStopdate = functions.TrimEscape(httpReq.FormValue("stopdate"))
	}

	sqlFeedback := `select merchant.title as question, merchant.answer as answer, merchant.redemptioncontrol as redemptioncontrol
						from redemption 
							left join merchant on merchant.redemptioncontrol = redemption.control
						where redemption.merchantcontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp 
					`
	sqlFeedback = fmt.Sprintf(sqlFeedback, this.AdminControl, sStartdate, sStopdate)

	/*
		sqlRedemption := `select redemption.control as control, redemption.createdate as date, store.address as location, redemption.transactionvalue as revenue, redemption.schemecontrol as schemecontrol,
							reward.discount as discount, reward.title as reward, member.dob as age, member.title as gender, member.nationality as nationality, coupon.code as coupon
						from redemption
								left join store on store.control = redemption.storecontrol
								left join reward on reward.control = redemption.rewardcontrol
								left join coupon on coupon.control = redemption.couponcontrol
								left join profile as member on member.control = redemption.membercontrol
						where redemption.merchantcontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp order by substring(redemption.createdate from 1 for 20)::timestamp desc
						`
	*/

	sqlRedemption := `select redemption.control as control, redemption.createdate as date, redemption.location as location, redemption.transactionvalue as revenue, redemption.schemecontrol as schemecontrol,
						redemption.discount as discount, redemption.reward as reward, redemption.dob as age, redemption.gender as gender, redemption.nationality as nationality, redemption.code as coupon
					from redemption where redemption.merchantcontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp order by substring(redemption.createdate from 1 for 20)::timestamp desc
					`
	sqlRedemption = fmt.Sprintf(sqlRedemption, this.AdminControl, sStartdate, sStopdate)

	curdb.Query("set datestyle = dmy")
	mapFeedback, _ := curdb.Query(sqlFeedback)
	mapRedemption, _ := curdb.Query(sqlRedemption)

	merchantReportGenerator := make(map[string]interface{})
	for _, xDoc := range mapFeedback {
		xDoc := xDoc.(map[string]interface{})

		sKeyType := ""
		switch {
		case strings.Contains(xDoc["question"].(string), "RECOMMEND"):
			sKeyType = "-RECOMMEND"
		case strings.Contains(xDoc["question"].(string), "IMPROVEMENT"):
			sKeyType = "-IMPROVEMENT"
		}

		sKey := fmt.Sprintf("%s%s", xDoc["redemptioncontrol"], sKeyType)
		merchantReportGenerator[sKey] = xDoc["answer"]
	}

	sLine := `"%v","%v","%v","%v","%v","%v","%v","%v","%v","%v","%v","%v"`
	redemptionReportGenerator := make([]string, len(mapRedemption)+1)
	redemptionReportGenerator[0] = fmt.Sprintf(sLine, "Transaction Number", "Date", "Discount", "Reward", "Coupon",
		"Revenue", "Nationality", "Gender", "Age", "Requested Improvement", "NPS Scoring", "Location")

	for cNumber, xDoc := range mapRedemption {

		iNumber, _ := strconv.Atoi(cNumber)
		xDoc := xDoc.(map[string]interface{})

		switch xDoc["gender"].(string) {
		case "Mrs", "Miss":
			xDoc["gender"] = "Female"
		case "Mr":
			xDoc["gender"] = "Male"
		default:
			xDoc["gender"] = "Other"
		}

		sKey := fmt.Sprintf("%s-IMPROVEMENT", xDoc["control"])

		if merchantReportGenerator[sKey] != nil {
			xDoc["improve"] = merchantReportGenerator[sKey]
		} else {
			xDoc["improve"] = ""
		}

		sKey = fmt.Sprintf("%s-RECOMMEND", xDoc["control"])
		if merchantReportGenerator[sKey] != nil {
			xDoc["rating"] = merchantReportGenerator[sKey]
		} else {
			xDoc["rating"] = ""
		}

		xDoc["age"] = functions.GetDifferenceInYears("", xDoc["age"].(string))
		xDoc["date"] = xDoc["date"].(string)[0:10]

		redemptionReportGenerator[iNumber] = fmt.Sprintf(sLine, xDoc["control"], xDoc["date"], xDoc["discount"], xDoc["reward"], xDoc["coupon"],
			xDoc["revenue"], xDoc["nationality"], xDoc["gender"], xDoc["age"], xDoc["improve"], xDoc["rating"], xDoc["location"])
	}

	sFilename := fmt.Sprintf("Valued-Report-%s-%s.csv", sStartdate, sStopdate)
	httpRes.Header().Set("Content-Type", "text/csv")
	httpRes.Header().Set("Content-Disposition", "attachment;filename="+sFilename)
	httpRes.Write([]byte(strings.Join(redemptionReportGenerator, "\r\n")))
}

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
	startDate, _ := time.Parse(cFormat, "01"+time.Now().Format(cFormat)[2:])
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

	sqlEmployerSubscriptionDue := fmt.Sprintf(`select count(control) as employersubscriptiondue from subscription where employercontrol not in (select control from profile where code in ('main','britishmums')) where expirydate::timestamp between '%s'::timestamp and '%s'::timestamp `, startDate.Format(cFormat), oneMonth.Format(cFormat))
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

func (this *Report) subscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	//Get List of Stores and Their Address
	mapReport := make(map[string]interface{})

	cFormat := "02/01/2006"
	startDate, _ := time.Parse(cFormat, this.AdminCreatedate)
	diffYears := int64(functions.GetDifferenceInSeconds(functions.GetSystemDate(), this.AdminCreatedate)) / int64(time.Hour*24*365)
	if diffYears > 0 {
		startDate = startDate.Add(time.Hour * 24 * 365 * time.Duration(diffYears))
	}
	oneYear := startDate.Add(time.Hour * 24 * 365)

	sStartdate := startDate.Format(cFormat)
	sStopdate := oneYear.Format(cFormat)
	mapReport["startdate"] = sStartdate
	mapReport["stopdate"] = sStopdate

	sqlDemographic := `select redemption.location as location, redemption.transactionvalue as revenue,
								redemption.dob as age, redemption.gender as gender, redemption.nationality as nationality
							from redemption where redemption.merchantcontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp 
					`
	sqlDemographic = fmt.Sprintf(sqlDemographic, this.AdminControl, sStartdate, sStopdate)

	curdb.Query("set datestyle = dmy")
	mapDemographic, _ := curdb.Query(sqlDemographic)

	mapAge := make(map[string]int)
	mapGender := make(map[string]int)
	mapNationality := make(map[string]int)
	var agePieTotal, genderPieTotal, nationalityPieTotal int

	for _, xDocDemographic := range mapDemographic {
		xDocDemographic := xDocDemographic.(map[string]interface{})

		sGender := "UnKnown"
		switch xDocDemographic["gender"].(string) {
		case "Mrs", "Miss":
			sGender = "Female"
		case "Mr":
			sGender = "Male"
		}
		mapGender[sGender] += 1
		genderPieTotal += 1

		sAge := ""
		if xDocDemographic["age"].(string) != "" {
			iAge := functions.GetDifferenceInYears("", xDocDemographic["age"].(string))
			switch {
			case iAge >= 18 && iAge <= 25:
				sAge = "18-25"
			case iAge >= 26 && iAge <= 30:
				sAge = "26-30"
			case iAge >= 31 && iAge <= 40:
				sAge = "31-40"
			case iAge >= 41 && iAge <= 60:
				sAge = "41-60"
			case iAge >= 61:
				sAge = ">61"
			}
		} else {
			sAge = "UnKnown"
		}
		mapAge[sAge] += 1
		agePieTotal += 1

		sNationality := xDocDemographic["nationality"].(string)
		if sNationality == "" {
			sNationality = "UnKnown"
		}

		mapNationality[sNationality] += 1
		nationalityPieTotal += 1

	}

	mapPieChartAge := make(map[string]interface{})
	mapPieChartAge["id"] = "age"
	mapPieChartAge["title"] = "Age"
	mapPieChartAge["label"] = ""
	mapPieChartAge["series"] = ""

	if len(mapAge) == 0 {
		mapPieChartAge["label"] = "'No Records'"
		mapPieChartAge["series"] = "100"

	}

	mapLegendAge := []string{"UnKnown", "18-25", "26-30", "31-40", "41-60", ">61"}
	for _, iSeries := range mapAge {
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(agePieTotal) * 100)

		mapPieChartAge["label"] = fmt.Sprintf(`%v'%v%%',`, mapPieChartAge["label"], iSeriesPercentage)
		mapPieChartAge["series"] = fmt.Sprintf(`%v%v,`, mapPieChartAge["series"], iSeriesPercentage)
	}
	for iNumber, sLabel := range mapLegendAge {

		iSeries := mapAge[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(agePieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-generator-piechart-row`, iNumber)
		mapPieChartAge[sTag] = pieChartRow

	}
	mapReport["1#report-generator-piechart"] = mapPieChartAge

	//

	mapPieChartGender := make(map[string]interface{})
	mapPieChartGender["id"] = "gender"
	mapPieChartGender["title"] = "Gender"
	mapPieChartGender["label"] = ""
	mapPieChartGender["series"] = ""

	if len(mapGender) == 0 {
		mapPieChartGender["label"] = "'No Records'"
		mapPieChartGender["series"] = "100"
		mapPieChartGender["total"] = "100"
	}

	mapLegendGender := []string{"UnKnown", "Female", "Male"}
	for _, iSeries := range mapGender {
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(genderPieTotal) * 100)
		mapPieChartGender["label"] = fmt.Sprintf(`%v'%v%%',`, mapPieChartGender["label"], iSeriesPercentage)
		mapPieChartGender["series"] = fmt.Sprintf(`%v%v,`, mapPieChartGender["series"], iSeriesPercentage)
	}
	for iNumber, sLabel := range mapLegendGender {

		iSeries := mapGender[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(genderPieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-generator-piechart-row`, iNumber)
		mapPieChartGender[sTag] = pieChartRow

	}

	mapReport["2#report-generator-piechart"] = mapPieChartGender

	//

	mapPieChartNationality := make(map[string]interface{})
	mapPieChartNationality["id"] = "nationality"
	mapPieChartNationality["title"] = "Nationality"
	mapPieChartNationality["label"] = ""
	mapPieChartNationality["series"] = ""

	if len(mapNationality) == 0 {
		mapPieChartNationality["label"] = "'No Records'"
		mapPieChartNationality["series"] = "100"
		mapPieChartNationality["total"] = "100"
	}

	var mapNationalityTop5 []string
	for sLabel, iSeries := range mapNationality {
		sTag := fmt.Sprintf(`%v#%v`, iSeries, sLabel)
		mapNationalityTop5 = append(mapNationalityTop5, sTag)
	}
	functions.SortDesc(mapNationalityTop5)

	posCounter := 1
	var mapLegendNationality []string

	for _, iSeriesLabel := range mapNationalityTop5 {

		sLabel := strings.Split(iSeriesLabel, "#")[1]
		iSeries := mapNationality[sLabel]
		delete(mapNationality, sLabel)

		if posCounter > 4 {
			sLabel = "Others"
		}

		mapNationality[sLabel] += iSeries
		mapLegendNationality = append(mapLegendNationality, sLabel)
		posCounter++
	}

	for _, iSeries := range mapNationality {
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(nationalityPieTotal) * 100)
		mapPieChartNationality["label"] = fmt.Sprintf(`%v'%v%%',`, mapPieChartNationality["label"], iSeriesPercentage)
		mapPieChartNationality["series"] = fmt.Sprintf(`%v%v,`, mapPieChartNationality["series"], iSeriesPercentage)
	}
	for iNumber, sLabel := range mapLegendNationality {

		iSeries := mapNationality[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(nationalityPieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-generator-piechart-row`, iNumber)
		mapPieChartNationality[sTag] = pieChartRow

	}

	mapReport["3#report-generator-piechart"] = mapPieChartNationality

	//Get List of Stores and Their Address

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-subscription"] = mapReport

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Users | Demographics Reports","mainpanelContent":` + contentHTML + `}`))
}

func (this *Report) merchant(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	mapReport := make(map[string]interface{})

	cFormat := "02/01/2006"
	startDate, _ := time.Parse(cFormat, this.AdminCreatedate)
	diffYears := int64(functions.GetDifferenceInSeconds(functions.GetSystemDate(), this.AdminCreatedate)) / int64(time.Hour*24*365)
	if diffYears > 0 {
		startDate = startDate.Add(time.Hour * 24 * 365 * time.Duration(diffYears))
	}
	oneYear := startDate.Add(time.Hour * 24 * 365)

	sStartdate := startDate.Format(cFormat)
	sStopdate := oneYear.Format(cFormat)
	mapReport["startdate"] = sStartdate
	mapReport["stopdate"] = sStopdate

	//Get NPS Score Based on Feedback Rating
	sqlFeedback := `select merchant.title as question, merchant.answer as answer, merchant.redemptioncontrol as redemptioncontrol
						from redemption 
							left join merchant on merchant.redemptioncontrol = redemption.control
						where redemption.merchantcontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp 
					`
	sqlFeedback = fmt.Sprintf(sqlFeedback, this.AdminControl, sStartdate, sStopdate)

	curdb.Query("set datestyle = dmy")
	mapFeedback, _ := curdb.Query(sqlFeedback)

	iNPSTotal := float64(0)
	iNPSPositive := float64(0)
	iNPSNegative := float64(0)
	//100/(iNPSPositive+iNPSNegative)*(iNPSPositive-iNPSNegative)
	//100/5*(2-3)

	ratingCategory := make([]int, 11)
	improveCategory := make(map[string]int)

	for _, xDoc := range mapFeedback {
		xDoc := xDoc.(map[string]interface{})

		switch {
		case strings.Contains(xDoc["question"].(string), "IMPROVEMENT"):
			improveCategory[xDoc["answer"].(string)]++

		case strings.Contains(xDoc["question"].(string), "RECOMMEND"):
			score, _ := strconv.Atoi(xDoc["answer"].(string))
			switch {
			case score <= 6:
				iNPSNegative++
				break
			case score >= 9:
				iNPSPositive++
				break
			}
			iNPSTotal++
			ratingCategory[score]++
		}
	}

	//Get NPS Score Based on Feedback Rating

	iNPSNegativePercentage := float64(0)
	iNPSPositivePercentage := float64(0)
	if iNPSTotal > 0 {
		iNPSPositivePercentage = (iNPSPositive / iNPSTotal) * 100
		iNPSNegativePercentage = (iNPSNegative / iNPSTotal) * 100
	}

	mapReport["npsscore"] = functions.RoundUp(iNPSPositivePercentage-iNPSNegativePercentage, 0)

	//Get BarChart of Slider Rating
	ratingReportGenerator := make(map[string]interface{})
	ratingReportGenerator["id"] = "rating"
	ratingReportGenerator["title"] = "NPS SCORING"
	ratingHigh := float64(10)

	for iKey, iValue := range ratingCategory {

		sLabel := fmt.Sprintf("%v#label", iKey)
		sSeries := fmt.Sprintf("%v#series", iKey)

		ratingReportGenerator[sLabel] = fmt.Sprintf(`"%v",`, iKey)
		ratingReportGenerator[sSeries] = fmt.Sprintf("%v,", iValue)
		if float64(iValue) > ratingHigh {
			ratingHigh = float64(iValue)
			ratingHigh += float64(2)
		}
	}
	ratingReportGenerator["high"] = ratingHigh
	mapReport["1#report-generator-barchart"] = ratingReportGenerator
	//Get BarChart of Slider Rating

	//Get BarChart of Improvement Categories
	improvementCategoriesReportGenerator := make(map[string]interface{})
	improvementCategoriesReportGenerator["id"] = "revenue"
	improvementCategoriesReportGenerator["title"] = "REQUESTED IMPROVEMENT AREAS"
	improvementCategoriesHigh := float64(10)

	//Get Categories Allowed
	sqlReviewCategory := `select reviewcategory.description from reviewcategory 
							left join reviewcategorylink on reviewcategorylink.reviewcategorycontrol = reviewcategory.control
							where reviewcategorylink.merchantcontrol = '%s' order by reviewcategory.control`

	sqlReviewCategory = fmt.Sprintf(sqlReviewCategory, this.AdminControl)
	curdb.Query("set datestyle = dmy")
	mapReviewCategory, _ := curdb.Query(sqlReviewCategory)
	//Get Categories Allowed

	for cNumber, xDoc := range mapReviewCategory {
		xDoc := xDoc.(map[string]interface{})

		sSeries := fmt.Sprintf("%v#series", cNumber)
		sLabel := fmt.Sprintf("%v#label", cNumber)

		nSeries := 0
		sDescription := strings.TrimSpace(xDoc["description"].(string))
		if improveCategory[sDescription] != 0 {
			nSeries = improveCategory[sDescription]
		}

		improvementCategoriesReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["description"].(string))
		improvementCategoriesReportGenerator[sSeries] = fmt.Sprintf("%v,", nSeries)
		if float64(nSeries) > improvementCategoriesHigh {
			improvementCategoriesHigh = float64(nSeries)
			improvementCategoriesHigh += float64(2)
		}
	}
	improvementCategoriesReportGenerator["high"] = improvementCategoriesHigh
	mapReport["2#report-generator-barchart"] = improvementCategoriesReportGenerator
	//Get BarChart of Improvement Categories

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-merchant"] = mapReport

	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Users | Feedback Reports","mainpanelContent":` + contentHTML + `}`))
}
