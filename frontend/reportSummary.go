package frontend

import (
	"strings"
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

	fTotalSubscription := float64(0)
	sqlTotalSubscription := `select count(control) as totalsubscription from subscription`
	mapTotalSubscription, _ := curdb.Query(sqlTotalSubscription)
	mapReport["totalsubscription"] = float64(0)
	if mapTotalSubscription["1"] != nil {
		mapTotalSubscriptionxDoc := mapTotalSubscription["1"].(map[string]interface{})
		switch mapTotalSubscriptionxDoc["totalsubscription"].(type) {
		case string:
			mapReport["totalsubscription"] = float64(0)
		case int64:
			fTotalSubscription = functions.Round(float64(mapTotalSubscriptionxDoc["totalsubscription"].(int64)))
			mapReport["totalsubscription"] = functions.ThousandSeperator(fTotalSubscription)
			// mapReport["totalsubscription"] = functions.ThousandSeperator(functions.Round(float64(mapTotalSubscriptionxDoc["totalsubscription"].(int64))))
		case float64:
			fTotalSubscription = functions.Round(float64(mapTotalSubscriptionxDoc["totalsubscription"].(float64)))
			mapReport["totalsubscription"] = functions.ThousandSeperator(fTotalSubscription)
			// mapReport["totalsubscription"] = functions.ThousandSeperator(functions.Round(mapTotalSubscriptionxDoc["totalsubscription"].(float64)))
		}
	}

	//Employer
	fEmployerSubscription := float64(0)
	sqlEmployerSubscription := `select count(control) as employersubscription from subscription where employercontrol not in (select control from profile where code in ('main','britishmums')) `
	mapEmployerSubscription, _ := curdb.Query(sqlEmployerSubscription)
	mapReport["employersubscription"] = float64(0)
	if mapEmployerSubscription["1"] != nil {
		mapEmployerSubscriptionxDoc := mapEmployerSubscription["1"].(map[string]interface{})
		switch mapEmployerSubscriptionxDoc["employersubscription"].(type) {
		case string:
			mapReport["employersubscription"] = float64(0)
		case int64:
			fEmployerSubscription = functions.Round(float64(mapEmployerSubscriptionxDoc["employersubscription"].(int64)))
			mapReport["employersubscription"] = functions.ThousandSeperator(fEmployerSubscription)
			// mapReport["employersubscription"] = functions.ThousandSeperator(functions.Round(float64(mapEmployerSubscriptionxDoc["employersubscription"].(int64))))
		case float64:
			fEmployerSubscription = functions.Round(float64(mapEmployerSubscriptionxDoc["employersubscription"].(float64)))
			mapReport["employersubscription"] = functions.ThousandSeperator(fEmployerSubscription)
			// mapReport["employersubscription"] = functions.ThousandSeperator(functions.Round(mapEmployerSubscriptionxDoc["employersubscription"].(float64)))
		}
	}

	//BritishMum
	fBritishMumSubscription := float64(0)
	sqlBritishMumSubscription := `select count(control) as britishmumssubscription from subscription where employercontrol in (select control from profile where code in ('britishmums')) `
	mapBritishMumSubscription, _ := curdb.Query(sqlBritishMumSubscription)
	mapReport["britishmumssubscription"] = float64(0)
	if mapBritishMumSubscription["1"] != nil {
		mapBritishMumSubscriptionxDoc := mapBritishMumSubscription["1"].(map[string]interface{})
		switch mapBritishMumSubscriptionxDoc["britishmumssubscription"].(type) {
		case string:
			mapReport["britishmumssubscription"] = float64(0)
		case int64:
			fBritishMumSubscription = functions.Round(float64(mapBritishMumSubscriptionxDoc["britishmumssubscription"].(int64)))
			mapReport["britishmumssubscription"] = functions.ThousandSeperator(fBritishMumSubscription)
			// mapReport["britishmumssubscription"] = functions.ThousandSeperator(functions.Round(float64(mapBritishMumSubscriptionxDoc["britishmumssubscription"].(int64))))
		case float64:
			fBritishMumSubscription = functions.Round(float64(mapBritishMumSubscriptionxDoc["britishmumssubscription"].(float64)))
			mapReport["britishmumssubscription"] = functions.ThousandSeperator(fBritishMumSubscription)
			// mapReport["britishmumssubscription"] = functions.ThousandSeperator(functions.Round(mapBritishMumSubscriptionxDoc["britishmumssubscription"].(float64)))
		}
	}

	//Other
	fOtherSubscription := float64(0)
	sqlOtherSubscription := `select count(control) as othersubscription from subscription where employercontrol in (select control from profile where code in ('main')) `
	mapOtherSubscription, _ := curdb.Query(sqlOtherSubscription)
	mapReport["othersubscription"] = float64(0)
	if mapOtherSubscription["1"] != nil {
		mapOtherSubscriptionxDoc := mapOtherSubscription["1"].(map[string]interface{})
		switch mapOtherSubscriptionxDoc["othersubscription"].(type) {
		case string:
			mapReport["othersubscription"] = float64(0)
		case int64:
			fOtherSubscription = functions.Round(float64(mapOtherSubscriptionxDoc["othersubscription"].(int64)))
			mapReport["othersubscription"] = functions.ThousandSeperator(fOtherSubscription)
			// mapReport["othersubscription"] = functions.ThousandSeperator(functions.Round(float64(mapOtherSubscriptionxDoc["othersubscription"].(int64))))
		case float64:
			fOtherSubscription = functions.Round(float64(mapOtherSubscriptionxDoc["othersubscription"].(float64)))
			mapReport["othersubscription"] = functions.ThousandSeperator(fOtherSubscription)
			// mapReport["othersubscription"] = functions.ThousandSeperator(functions.Round(mapOtherSubscriptionxDoc["othersubscription"].(float64)))
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

	//Demographics of Report
	//

	curdb.Query("set datestyle = dmy")

	sqlDemographic := `select  subscription.code, profile.control, profile.dob as age, profile.title as gender, profile.nationality as nationality from subscription join profile on profile.control = subscription.membercontrol`
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
	mapPieChartAge["title"] = "Age"

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
		if iSeries == 0 {
			iSeriesPercentage = float64(0)
		}
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-summary-demographics-row`, iNumber)
		mapPieChartAge[sTag] = pieChartRow

	}
	mapReport["1#report-summary-demographics"] = mapPieChartAge

	//

	mapPieChartGender := make(map[string]interface{})
	mapPieChartGender["title"] = "Gender"

	mapLegendGender := []string{"UnKnown", "Female", "Male"}
	for iNumber, sLabel := range mapLegendGender {

		iSeries := mapGender[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(genderPieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		if iSeries == 0 {
			iSeriesPercentage = float64(0)
		}
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-summary-demographics-row`, iNumber)
		mapPieChartGender[sTag] = pieChartRow

	}

	mapReport["2#report-summary-demographics"] = mapPieChartGender

	//

	mapPieChartNationality := make(map[string]interface{})
	mapPieChartNationality["title"] = "Nationality"

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
		if posCounter <= 5 {
			mapLegendNationality = append(mapLegendNationality, sLabel)
		}
		posCounter++
	}

	for iNumber, sLabel := range mapLegendNationality {

		iSeries := mapNationality[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(nationalityPieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		if iSeries == 0 {
			iSeriesPercentage = float64(0)
		}
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-summary-demographics-row`, iNumber)
		mapPieChartNationality[sTag] = pieChartRow

	}

	mapReport["3#report-summary-demographics"] = mapPieChartNationality

	//Engagement Funnel

	//Total Registrations
	sqlTotalRegistration := `select count(control) as totalregistered from profile where company != 'Yes' and workflow = 'registered'`
	mapTotalRegistration, _ := curdb.Query(sqlTotalRegistration)
	mapReport["totalregistered"] = float64(0)
	if mapTotalRegistration["1"] != nil {
		mapTotalRegistrationxDoc := mapTotalRegistration["1"].(map[string]interface{})
		switch mapTotalRegistrationxDoc["totalregistered"].(type) {
		case string:
			mapReport["totalregistered"] = float64(0)
		case int64:
			mapReport["totalregistered"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRegistrationxDoc["totalregistered"].(int64))))
		case float64:
			mapReport["totalregistered"] = functions.ThousandSeperator(functions.Round(mapTotalRegistrationxDoc["totalregistered"].(float64)))
		}
	}

	//Total Subscriptions
	sqlTotalSubscribed := `select count(control) as totalsubscribed from profile where company != 'Yes' and workflow = 'subscribed'`
	mapTotalSubscribed, _ := curdb.Query(sqlTotalSubscribed)
	mapReport["totalsubscribed"] = float64(0)
	if mapTotalSubscribed["1"] != nil {
		mapTotalSubscribedxDoc := mapTotalSubscribed["1"].(map[string]interface{})
		switch mapTotalSubscribedxDoc["totalsubscribed"].(type) {
		case string:
			mapReport["totalsubscribed"] = float64(0)
		case int64:
			mapReport["totalsubscribed"] = functions.ThousandSeperator(functions.Round(float64(mapTotalSubscribedxDoc["totalsubscribed"].(int64))))
		case float64:
			mapReport["totalsubscribed"] = functions.ThousandSeperator(functions.Round(mapTotalSubscribedxDoc["totalsubscribed"].(float64)))
		}
	}

	//Total Active Members
	sqlTotalActiveMember := `select count(control) as totalactivemember from profile where company != 'Yes' and status = 'active'`
	mapTotalActiveMember, _ := curdb.Query(sqlTotalActiveMember)
	mapReport["totalactivemember"] = float64(0)
	if mapTotalActiveMember["1"] != nil {
		mapTotalActiveMemberxDoc := mapTotalActiveMember["1"].(map[string]interface{})
		switch mapTotalActiveMemberxDoc["totalactivemember"].(type) {
		case string:
			mapReport["totalactivemember"] = float64(0)
		case int64:
			mapReport["totalactivemember"] = functions.ThousandSeperator(functions.Round(float64(mapTotalActiveMemberxDoc["totalactivemember"].(int64))))
		case float64:
			mapReport["totalactivemember"] = functions.ThousandSeperator(functions.Round(mapTotalActiveMemberxDoc["totalactivemember"].(float64)))
		}
	}

	//Engagement Funnel

	//

	//Savings

	//Total Savings
	sqlTotalSavings := `select sum(savingsvalue) as totalsavings from redemption`
	mapTotalSavings, _ := curdb.Query(sqlTotalSavings)
	mapReport["totalsavings"] = float64(0)
	if mapTotalSavings["1"] != nil {
		mapTotalSavingsxDoc := mapTotalSavings["1"].(map[string]interface{})
		switch mapTotalSavingsxDoc["totalsavings"].(type) {
		case string:
			mapReport["totalsavings"] = float64(0)
		case int64:
			mapReport["totalsavings"] = functions.ThousandSeperator(functions.Round(float64(mapTotalSavingsxDoc["totalsavings"].(int64))))
		case float64:
			mapReport["totalsavings"] = functions.ThousandSeperator(functions.Round(mapTotalSavingsxDoc["totalsavings"].(float64)))
		}
	}

	//Total Reward Savings
	sqlTotalRewardSavings := `select (sum(savingsvalue) / count(distinct rewardcontrol)) as totalrewardsavings from redemption`
	mapTotalRewardSavings, _ := curdb.Query(sqlTotalRewardSavings)
	mapReport["totalrewardsavings"] = float64(0)
	if mapTotalRewardSavings["1"] != nil {
		mapTotalRewardSavingsxDoc := mapTotalRewardSavings["1"].(map[string]interface{})
		switch mapTotalRewardSavingsxDoc["totalrewardsavings"].(type) {
		case string:
			mapReport["totalrewardsavings"] = float64(0)
		case int64:
			mapReport["totalrewardsavings"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRewardSavingsxDoc["totalrewardsavings"].(int64))))
		case float64:
			mapReport["totalrewardsavings"] = functions.ThousandSeperator(functions.Round(mapTotalRewardSavingsxDoc["totalrewardsavings"].(float64)))
		}
	}

	//Total Active Members Savings
	sqlTotalActiveMemberSavings := `select (sum(savingsvalue) / count(distinct membercontrol)) as totalmembersavings from redemption`
	mapTotalActiveMemberSavings, _ := curdb.Query(sqlTotalActiveMemberSavings)
	mapReport["totalmembersavings"] = float64(0)
	if mapTotalActiveMemberSavings["1"] != nil {
		mapTotalActiveMemberSavingsxDoc := mapTotalActiveMemberSavings["1"].(map[string]interface{})
		switch mapTotalActiveMemberSavingsxDoc["totalmembersavings"].(type) {
		case string:
			mapReport["totalmembersavings"] = float64(0)
		case int64:
			mapReport["totalmembersavings"] = functions.ThousandSeperator(functions.Round(float64(mapTotalActiveMemberSavingsxDoc["totalmembersavings"].(int64))))
		case float64:
			mapReport["totalmembersavings"] = functions.ThousandSeperator(functions.Round(mapTotalActiveMemberSavingsxDoc["totalmembersavings"].(float64)))
		}
	}

	//Savings

	//Revenue

	//Total Revenue
	sqlTotalRevenue := `select sum(transactionvalue) as totalrevenue from redemption`
	mapTotalRevenue, _ := curdb.Query(sqlTotalRevenue)
	mapReport["totalrevenue"] = float64(0)
	if mapTotalRevenue["1"] != nil {
		mapTotalRevenuexDoc := mapTotalRevenue["1"].(map[string]interface{})
		switch mapTotalRevenuexDoc["totalrevenue"].(type) {
		case string:
			mapReport["totalrevenue"] = float64(0)
		case int64:
			mapReport["totalrevenue"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRevenuexDoc["totalrevenue"].(int64))))
		case float64:
			mapReport["totalrevenue"] = functions.ThousandSeperator(functions.Round(mapTotalRevenuexDoc["totalrevenue"].(float64)))
		}
	}

	//Total Reward Revenue
	sqlTotalRewardRevenue := `select (sum(transactionvalue) / count(distinct rewardcontrol)) as totalrewardrevenue from redemption`
	mapTotalRewardRevenue, _ := curdb.Query(sqlTotalRewardRevenue)
	mapReport["totalrewardrevenue"] = float64(0)
	if mapTotalRewardRevenue["1"] != nil {
		mapTotalRewardRevenuexDoc := mapTotalRewardRevenue["1"].(map[string]interface{})
		switch mapTotalRewardRevenuexDoc["totalrewardrevenue"].(type) {
		case string:
			mapReport["totalrewardrevenue"] = float64(0)
		case int64:
			mapReport["totalrewardrevenue"] = functions.ThousandSeperator(functions.Round(float64(mapTotalRewardRevenuexDoc["totalrewardrevenue"].(int64))))
		case float64:
			mapReport["totalrewardrevenue"] = functions.ThousandSeperator(functions.Round(mapTotalRewardRevenuexDoc["totalrewardrevenue"].(float64)))
		}
	}

	//Total Active Members Revenue
	sqlTotalActiveMemberRevenue := `select (sum(transactionvalue) / count(distinct membercontrol)) as totalmemberrevenue from redemption`
	mapTotalActiveMemberRevenue, _ := curdb.Query(sqlTotalActiveMemberRevenue)
	mapReport["totalmemberrevenue"] = float64(0)
	if mapTotalActiveMemberRevenue["1"] != nil {
		mapTotalActiveMemberRevenuexDoc := mapTotalActiveMemberRevenue["1"].(map[string]interface{})
		switch mapTotalActiveMemberRevenuexDoc["totalmemberrevenue"].(type) {
		case string:
			mapReport["totalmemberrevenue"] = float64(0)
		case int64:
			mapReport["totalmemberrevenue"] = functions.ThousandSeperator(functions.Round(float64(mapTotalActiveMemberRevenuexDoc["totalmemberrevenue"].(int64))))
		case float64:
			mapReport["totalmemberrevenue"] = functions.ThousandSeperator(functions.Round(mapTotalActiveMemberRevenuexDoc["totalmemberrevenue"].(float64)))
		}
	}

	//Revenue

	//

	// Top 10 redeemed rewards
	sqlTop10Rewards := `select (select title from profile where control=reward.merchantcontrol ) as merchant, title as reward, 
		(select count(control) from redemption where rewardcontrol = reward.control) as redemption,
		reward.control from reward
		order by 3 desc , 4 `
	mapTop10Rewards, _ := curdb.Query(sqlTop10Rewards)
	aTop10RewardsSorted := functions.SortMap(mapTop10Rewards)

	for nRowCounter, sNumber := range aTop10RewardsSorted {
		if nRowCounter > 9 {
			break
		}

		xDocReward := mapTop10Rewards[sNumber].(map[string]interface{})
		xDocReward["row"] = sNumber
		sTag := fmt.Sprintf(`%v#report-summary-toptenredeemed-row`, sNumber)
		mapReport[sTag] = xDocReward
	}
	// Top 10 redeemed rewards
	//

	// % OF SUBSCRIBERS WITH A GREATER SAVINGS AMOUNT THAN THE COST OF MEMBERSHIP
	sqlPercentSubscribersSavings := `select subscription.membercontrol as membercontrol, sum(subscription.price) as subscription, 
						(select sum(redemption.savingsvalue)  from redemption where membercontrol = subscription.membercontrol) as savings
						from subscription group by subscription.membercontrol order by 3 desc nulls last`

	mapPercentSubscribersSavings, _ := curdb.Query(sqlPercentSubscribersSavings)
	aPercentSubscribersSavingsSorted := functions.SortMap(mapPercentSubscribersSavings)

	var mapSavingsGreater []string
	for _, sNumber := range aPercentSubscribersSavingsSorted {
		xDocSubscribersSavings := mapPercentSubscribersSavings[sNumber].(map[string]interface{})

		switch xDocSubscribersSavings["savings"].(type) {
		case float64:
			if xDocSubscribersSavings["savings"].(float64) >
				xDocSubscribersSavings["subscription"].(float64) {
				mapSavingsGreater = append(mapSavingsGreater, xDocSubscribersSavings["membercontrol"].(string))
			}

		}
	}

	mapReport["percentofsubscriberssavingsvscosts"] = functions.RoundUp(float64(len(mapSavingsGreater)*100)/float64(len(aPercentSubscribersSavingsSorted)), 0)

	// % OF SUBSCRIBERS WITH A GREATER SAVINGS AMOUNT THAN THE COST OF MEMBERSHIP
	//

	//
	//
	//GRAPHS: - NUMBER OF MONTHLY REDEEMED REWARDS FOR LAST 12 MONTHS
	sLabel := "yyyy-Mon"
	sOrderBy := "yyyymm"

	curMonth := time.Now().Add(-(time.Hour * 24 * 365))
	sStartdateChart := curMonth.Format(cFormat)

	oneYearChart := time.Now()
	sStopdateChart := oneYearChart.Format(cFormat)

	redeemedReportGenerator := make(map[string]interface{})

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

		redeemedReportGenerator[sLabelIndex] = fmt.Sprintf(`"%s",`, sLabel)
		redeemedReportGenerator[sSeriesIndex] = functions.ThousandSeperator(functions.Round(float64(0))) + ","
	}

	sqlRedeemedChart := `select to_char(substring(createdate from 1 for 20)::timestamp,'%s') as orderby, to_char(substring(createdate from 1 for 20)::timestamp,'%s') as label, 
	count(control) as redeemed from redemption where substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp 
					group by 1,2 order by 1
					`
	sqlRedeemedChart = fmt.Sprintf(sqlRedeemedChart, sOrderBy, sLabel, sStartdateChart, sStopdateChart)
	mapRedeemedChart, _ := curdb.Query(sqlRedeemedChart)

	redeemedReportGenerator["id"] = "redeemed"
	redeemedReportGenerator["title"] = "NUMBER OF MONTHLY REDEEMED REWARDS FOR LAST 12 MONTHS"
	redeemedHigh := int64(100)
	for _, xDoc := range mapRedeemedChart {
		xDoc := xDoc.(map[string]interface{})

		sLabel := fmt.Sprintf("%s#label", xDoc["orderby"])
		sSeries := fmt.Sprintf("%s#series", xDoc["orderby"])

		redeemedReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["label"])

		redeemedReportGenerator[sSeries] = fmt.Sprintf("%v,", xDoc["redeemed"])
		if xDoc["redeemed"].(int64) > redeemedHigh {
			redeemedHigh = xDoc["redeemed"].(int64)
			redeemedHigh += int64(2)
		}

	}
	redeemedReportGenerator["high"] = redeemedHigh
	mapReport["1#report-generator-barchart-admin"] = redeemedReportGenerator
	//GRAPHS: - NUMBER OF MONTHLY REDEEMED REWARDS FOR LAST 12 MONTHS
	//
	//

	//
	//
	//GRAPHS: - NUMBER OF NEW SUBSCRIPTIONS BY MONTH FOR LAST 12 MONTHS

	subscriptionReportGenerator := make(map[string]interface{})

	for sOrderby, sLabel := range monthLabelSeries {
		sLabelIndex := fmt.Sprintf("%s#label", sOrderby)
		sSeriesIndex := fmt.Sprintf("%s#series", sOrderby)

		subscriptionReportGenerator[sLabelIndex] = fmt.Sprintf(`"%s",`, sLabel)
		subscriptionReportGenerator[sSeriesIndex] = functions.ThousandSeperator(functions.Round(float64(0))) + ","
	}

	sqlSubscriptionChart := `select to_char(substring(createdate from 1 for 20)::timestamp,'%s') as orderby, to_char(substring(createdate from 1 for 20)::timestamp,'%s') as label, 
	count(control) as subscription from subscription where substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp 
					group by 1,2 order by 1
					`
	sqlSubscriptionChart = fmt.Sprintf(sqlSubscriptionChart, sOrderBy, sLabel, sStartdateChart, sStopdateChart)
	mapSubscriptionChart, _ := curdb.Query(sqlSubscriptionChart)

	subscriptionReportGenerator["id"] = "subscription"
	subscriptionReportGenerator["title"] = "NUMBER OF NEW SUBSCRIPTIONS BY MONTH FOR LAST 12 MONTHS"
	subscriptionHigh := int64(100)
	for _, xDoc := range mapSubscriptionChart {
		xDoc := xDoc.(map[string]interface{})

		sLabel := fmt.Sprintf("%s#label", xDoc["orderby"])
		sSeries := fmt.Sprintf("%s#series", xDoc["orderby"])

		subscriptionReportGenerator[sLabel] = fmt.Sprintf(`"%s",`, xDoc["label"])

		subscriptionReportGenerator[sSeries] = fmt.Sprintf("%v,", xDoc["subscription"])
		if xDoc["subscription"].(int64) > subscriptionHigh {
			subscriptionHigh = xDoc["subscription"].(int64)
			subscriptionHigh += int64(2)
		}

	}
	subscriptionReportGenerator["high"] = subscriptionHigh
	mapReport["2#report-generator-barchart-admin"] = subscriptionReportGenerator
	//GRAPHS: - NUMBER OF NEW SUBSCRIPTIONS BY MONTH FOR LAST 12 MONTHS
	//
	//

	//
	//
	//- PIE CHART SHOWING % OF SUBSCRIPTIONS BY EMPLOYERS, BM, OTHER

	mapPieChartSubscription := make(map[string]interface{})
	mapPieChartSubscription["id"] = "subscriptionPie"
	mapPieChartSubscription["title"] = "NUMBER OF SUBSCRIPTIONS PIE CHART"
	mapPieChartSubscription["label"] = ""
	mapPieChartSubscription["series"] = ""

	subscriptionPieTotal := int64(fTotalSubscription)
	mapSubscriptionPie := map[string]int64{
		"Employers":    int64(fEmployerSubscription),
		"British Mums": int64(fBritishMumSubscription),
		"Others":       int64(fOtherSubscription),
	}

	if len(mapSubscriptionPie) == 0 {
		mapPieChartSubscription["label"] = "'No Records'"
		mapPieChartSubscription["series"] = "100"
	}

	mapLegendSubscription := []string{"Employers", "British Mums", "Others"}
	for _, iSeries := range mapSubscriptionPie {
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(subscriptionPieTotal) * 100)

		mapPieChartSubscription["label"] = fmt.Sprintf(`%v'%v%%',`, mapPieChartSubscription["label"], iSeriesPercentage)
		mapPieChartSubscription["series"] = fmt.Sprintf(`%v%v,`, mapPieChartSubscription["series"], iSeriesPercentage)
	}
	for iNumber, sLabel := range mapLegendSubscription {

		iSeries := mapSubscriptionPie[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(subscriptionPieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-generator-piechart-row`, iNumber)
		mapPieChartSubscription[sTag] = pieChartRow

	}
	mapReport["1#report-generator-piechart-admin"] = mapPieChartSubscription

	//- PIE CHART SHOWING % OF SUBSCRIPTIONS BY EMPLOYERS, BM, OTHER
	//
	//

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-summary"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Summary | Admin Reports","mainpanelContent":` + contentHTML + `}`))

}
