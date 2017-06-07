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

	//
	// Total Subscriptions

	//
	// % OF SUBSCRIBERS WITH A GREATER SAVINGS AMOUNT THAN THE COST OF MEMBERSHIP
	sqlPercentSubscribersSavingsTotal := `select subscription.membercontrol as membercontrol, sum(subscription.price) as subscription, 
								(select sum(redemption.savingsvalue)  from redemption where membercontrol = subscription.membercontrol) as savings
								from subscription group by subscription.membercontrol order by 3 desc nulls last`

	mapPercentSubscribersSavingsTotal, _ := curdb.Query(sqlPercentSubscribersSavingsTotal)
	aPercentSubscribersSavingsTotalSorted := functions.SortMap(mapPercentSubscribersSavingsTotal)

	var mapSavingsTotalGreater []string
	for _, sNumber := range aPercentSubscribersSavingsTotalSorted {
		xDocSubscribersSavingsTotal := mapPercentSubscribersSavingsTotal[sNumber].(map[string]interface{})

		switch xDocSubscribersSavingsTotal["savings"].(type) {
		case float64:
			if xDocSubscribersSavingsTotal["savings"].(float64) >
				xDocSubscribersSavingsTotal["subscription"].(float64) {
				mapSavingsTotalGreater = append(mapSavingsTotalGreater, xDocSubscribersSavingsTotal["membercontrol"].(string))
			}

		}
	}

	mapReport["percentofsubscriberssavingsvscoststotal"] = functions.RoundUp(float64(len(mapSavingsTotalGreater)*100)/float64(len(aPercentSubscribersSavingsTotalSorted)), 0)

	// % OF SUBSCRIBERS WITH A GREATER SAVINGS AMOUNT THAN THE COST OF MEMBERSHIP
	//

	//
	// Number of Users - Pending / Paid / Expired

	mapReport["subscriptionspendingtotal"] = float64(0)
	mapReport["subscriptionspaidtotal"] = float64(0)
	mapReport["subscriptionsexpiredtotal"] = float64(0)

	// fActiveEmployeesLifestyleTotal := float64(0)
	// sqlActiveEmployeesLifestyleTotal := `select count(control) as numactivelifestyletotal from profile where control in
	// (select membercontrol from subscription  where schemecontrol in (select control from scheme where code in ('lifestyle'))) and status = 'active'`

	// mapActiveEmployeesLifestyleTotal, _ := curdb.Query(sqlActiveEmployeesLifestyleTotal)
	// mapReport["numactivelifestyletotal"] = float64(0)
	// if mapActiveEmployeesLifestyleTotal["1"] != nil {
	// 	mapActiveEmployeesLifestyleTotalxDoc := mapActiveEmployeesLifestyleTotal["1"].(map[string]interface{})
	// 	switch mapActiveEmployeesLifestyleTotalxDoc["numactivelifestyletotal"].(type) {
	// 	case string:
	// 		mapReport["numactivelifestyletotal"] = float64(0)
	// 	case int64:
	// 		fActiveEmployeesLifestyleTotal = functions.Round(float64(mapActiveEmployeesLifestyleTotalxDoc["numactivelifestyletotal"].(int64)))
	// 		mapReport["numactivelifestyletotal"] = functions.ThousandSeperator(fActiveEmployeesLifestyleTotal)
	// 	case float64:
	// 		fActiveEmployeesLifestyleTotal = functions.Round(float64(mapActiveEmployeesLifestyleTotalxDoc["numactivelifestyletotal"].(float64)))
	// 		mapReport["numactivelifestyletotal"] = functions.ThousandSeperator(fActiveEmployeesLifestyleTotal)
	// 	}
	// }

	// fActiveEmployeesLiteTotal := float64(0)
	// sqlActiveEmployeesLiteTotal := `select count(control) as numactivelitetotal from profile where control in
	// (select membercontrol from subscription  where schemecontrol in (select control from scheme where code in ('lite'))) and status = 'active'`

	// mapActiveEmployeesLiteTotal, _ := curdb.Query(sqlActiveEmployeesLiteTotal)
	// mapReport["numactivelitetotal"] = float64(0)
	// if mapActiveEmployeesLiteTotal["1"] != nil {
	// 	mapActiveEmployeesLiteTotalxDoc := mapActiveEmployeesLiteTotal["1"].(map[string]interface{})
	// 	switch mapActiveEmployeesLiteTotalxDoc["numactivelitetotal"].(type) {
	// 	case string:
	// 		mapReport["numactivelitetotal"] = float64(0)

	// 	case int64:
	// 		fActiveEmployeesLiteTotal = functions.Round(float64(mapActiveEmployeesLiteTotalxDoc["numactivelitetotal"].(int64)))
	// 		mapReport["numactivelitetotal"] = functions.ThousandSeperator(fActiveEmployeesLiteTotal)
	// 	case float64:
	// 		fActiveEmployeesLiteTotal = functions.Round(float64(mapActiveEmployeesLiteTotalxDoc["numactivelitetotal"].(float64)))
	// 		mapReport["numactivelitetotal"] = functions.ThousandSeperator(fActiveEmployeesLiteTotal)
	// 	}
	// }

	// Number of Users - Pending / Paid / Expired
	//

	/*
		//
		// Total & Percent of Active Employees

		fActiveEmployeesLifestyleTotal := float64(0)
		sqlActiveEmployeesLifestyleTotal := `select count(control) as numactivelifestyletotal from profile where control in
		(select membercontrol from subscription  where schemecontrol in (select control from scheme where code in ('lifestyle'))) and status = 'active'`

		mapActiveEmployeesLifestyleTotal, _ := curdb.Query(sqlActiveEmployeesLifestyleTotal)
		mapReport["numactivelifestyletotal"] = float64(0)
		if mapActiveEmployeesLifestyleTotal["1"] != nil {
			mapActiveEmployeesLifestyleTotalxDoc := mapActiveEmployeesLifestyleTotal["1"].(map[string]interface{})
			switch mapActiveEmployeesLifestyleTotalxDoc["numactivelifestyletotal"].(type) {
			case string:
				mapReport["numactivelifestyletotal"] = float64(0)
			case int64:
				fActiveEmployeesLifestyleTotal = functions.Round(float64(mapActiveEmployeesLifestyleTotalxDoc["numactivelifestyletotal"].(int64)))
				mapReport["numactivelifestyletotal"] = functions.ThousandSeperator(fActiveEmployeesLifestyleTotal)
			case float64:
				fActiveEmployeesLifestyleTotal = functions.Round(float64(mapActiveEmployeesLifestyleTotalxDoc["numactivelifestyletotal"].(float64)))
				mapReport["numactivelifestyletotal"] = functions.ThousandSeperator(fActiveEmployeesLifestyleTotal)
			}
		}

		fActiveEmployeesLiteTotal := float64(0)
		sqlActiveEmployeesLiteTotal := `select count(control) as numactivelitetotal from profile where control in
		(select membercontrol from subscription  where schemecontrol in (select control from scheme where code in ('lite'))) and status = 'active'`

		mapActiveEmployeesLiteTotal, _ := curdb.Query(sqlActiveEmployeesLiteTotal)
		mapReport["numactivelitetotal"] = float64(0)
		if mapActiveEmployeesLiteTotal["1"] != nil {
			mapActiveEmployeesLiteTotalxDoc := mapActiveEmployeesLiteTotal["1"].(map[string]interface{})
			switch mapActiveEmployeesLiteTotalxDoc["numactivelitetotal"].(type) {
			case string:
				mapReport["numactivelitetotal"] = float64(0)

			case int64:
				fActiveEmployeesLiteTotal = functions.Round(float64(mapActiveEmployeesLiteTotalxDoc["numactivelitetotal"].(int64)))
				mapReport["numactivelitetotal"] = functions.ThousandSeperator(fActiveEmployeesLiteTotal)
			case float64:
				fActiveEmployeesLiteTotal = functions.Round(float64(mapActiveEmployeesLiteTotalxDoc["numactivelitetotal"].(float64)))
				mapReport["numactivelitetotal"] = functions.ThousandSeperator(fActiveEmployeesLiteTotal)
			}
		}

		fActiveEmployeesTotal := float64(0)
		sqlActiveEmployeesTotal := `select count(control) as numactivetotal from profile where control in (select membercontrol from subscription) and status = 'active'`

		mapActiveEmployeesTotal, _ := curdb.Query(sqlActiveEmployeesTotal)
		mapReport["percentactivelitetotal"] = float64(0)
		mapReport["percentactivelifestyletotal"] = float64(0)

		if mapActiveEmployeesTotal["1"] != nil {
			mapActiveEmployeesTotalxDoc := mapActiveEmployeesTotal["1"].(map[string]interface{})
			switch mapActiveEmployeesTotalxDoc["numactivetotal"].(type) {
			case int64:
				fActiveEmployeesTotal = functions.Round(float64(mapActiveEmployeesTotalxDoc["numactivetotal"].(int64)))

			case float64:
				fActiveEmployeesTotal = functions.Round(float64(mapActiveEmployeesTotalxDoc["numactivetotal"].(float64)))
			}

			if fActiveEmployeesTotal > float64(0) && fActiveEmployeesLiteTotal > float64(0) {
				mapReport["percentactivelitetotal"] = functions.RoundDown((fActiveEmployeesLiteTotal*float64(100))/fActiveEmployeesTotal, 0)
			}

			if fActiveEmployeesTotal > float64(0) && fActiveEmployeesLifestyleTotal > float64(0) {
				mapReport["percentactivelifestyletotal"] = functions.RoundDown((fActiveEmployeesLifestyleTotal*float64(100))/fActiveEmployeesTotal, 0)
			}
		}

		// Total & Percent of Active Employees
		//
	*/

	//
	// Number of Registrations
	sqlRegistrationLiteTotal := `select count(control) as registrationsvaluedtotal from profile where company != 'Yes' and workflow = 'registered' 
										and employercontrol in (select control from profile where code in ('main'))`
	mapRegistrationLiteTotal, _ := curdb.Query(sqlRegistrationLiteTotal)
	mapReport["registrationsvaluedtotal"] = float64(0)
	if mapRegistrationLiteTotal["1"] != nil {
		mapRegistrationLiteTotalxDoc := mapRegistrationLiteTotal["1"].(map[string]interface{})
		switch mapRegistrationLiteTotalxDoc["registrationsvaluedtotal"].(type) {
		case string:
			mapReport["registrationsvaluedtotal"] = float64(0)
		case int64:
			mapReport["registrationsvaluedtotal"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationLiteTotalxDoc["registrationsvaluedtotal"].(int64))))
		case float64:
			mapReport["registrationsvaluedtotal"] = functions.ThousandSeperator(functions.Round(mapRegistrationLiteTotalxDoc["registrationsvaluedtotal"].(float64)))
		}
	}

	sqlRegistrationLifestyleTotal := `select count(control) as registrationsemployerstotal from profile where company != 'Yes' and workflow = 'registered' 
										and employercontrol not in (select control from profile where code in ('main','britishmums'))`
	mapRegistrationLifestyleTotal, _ := curdb.Query(sqlRegistrationLifestyleTotal)
	mapReport["registrationsemployerstotal"] = float64(0)
	if mapRegistrationLifestyleTotal["1"] != nil {
		mapRegistrationLifestyleTotalxDoc := mapRegistrationLifestyleTotal["1"].(map[string]interface{})
		switch mapRegistrationLifestyleTotalxDoc["registrationsemployerstotal"].(type) {
		case string:
			mapReport["registrationsemployerstotal"] = float64(0)
		case int64:
			mapReport["registrationsemployerstotal"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationLifestyleTotalxDoc["registrationsemployerstotal"].(int64))))
		case float64:
			mapReport["registrationsemployerstotal"] = functions.ThousandSeperator(functions.Round(mapRegistrationLifestyleTotalxDoc["registrationsemployerstotal"].(float64)))
		}
	}

	sqlRegistrationBritishMumsTotal := `select count(control) as registrationsbritishmumstotal from profile where company != 'Yes' and workflow = 'registered' 
										and employercontrol in (select control from profile where code in ('britishmums'))`
	mapRegistrationBritishMumsTotal, _ := curdb.Query(sqlRegistrationBritishMumsTotal)
	mapReport["registrationsbritishmumstotal"] = float64(0)
	if mapRegistrationBritishMumsTotal["1"] != nil {
		mapRegistrationBritishMumsTotalxDoc := mapRegistrationBritishMumsTotal["1"].(map[string]interface{})
		switch mapRegistrationBritishMumsTotalxDoc["registrationsbritishmumstotal"].(type) {
		case string:
			mapReport["registrationsbritishmumstotal"] = float64(0)
		case int64:
			mapReport["registrationsbritishmumstotal"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationBritishMumsTotalxDoc["registrationsbritishmumstotal"].(int64))))
		case float64:
			mapReport["registrationsbritishmumstotal"] = functions.ThousandSeperator(functions.Round(mapRegistrationBritishMumsTotalxDoc["registrationsbritishmumstotal"].(float64)))
		}
	}

	// Number of Registrations
	//

	//
	// Number of Subscriptions
	fSubscriptionsLiteTotal := float64(0)
	sqlSubscriptionsLiteTotal := `select count(control) as subscriptionslitetotal from subscription  where schemecontrol in (select control from scheme where code in ('lite'))`
	mapSubscriptionsLiteTotal, _ := curdb.Query(sqlSubscriptionsLiteTotal)
	mapReport["subscriptionslitetotal"] = float64(0)
	if mapSubscriptionsLiteTotal["1"] != nil {
		mapSubscriptionsLiteTotalxDoc := mapSubscriptionsLiteTotal["1"].(map[string]interface{})
		switch mapSubscriptionsLiteTotalxDoc["subscriptionslitetotal"].(type) {
		case string:
			mapReport["subscriptionslitetotal"] = float64(0)
		case int64:
			fSubscriptionsLiteTotal = functions.Round(float64(mapSubscriptionsLiteTotalxDoc["subscriptionslitetotal"].(int64)))
			mapReport["subscriptionslitetotal"] = functions.ThousandSeperator(fSubscriptionsLiteTotal)
		case float64:
			fSubscriptionsLiteTotal = functions.Round(float64(mapSubscriptionsLiteTotalxDoc["subscriptionslitetotal"].(float64)))
			mapReport["subscriptionslitetotal"] = functions.ThousandSeperator(fSubscriptionsLiteTotal)
		}
	}

	fSubscriptionsLifestyleTotal := float64(0)
	sqlSubscriptionsLifestyleTotal := `select count(control) as subscriptionslifestyletotal from subscription  where schemecontrol in (select control from scheme where code in ('lifestyle'))`
	mapSubscriptionsLifestyleTotal, _ := curdb.Query(sqlSubscriptionsLifestyleTotal)
	mapReport["subscriptionslifestyletotal"] = float64(0)
	if mapSubscriptionsLifestyleTotal["1"] != nil {
		mapSubscriptionsLifestyleTotalxDoc := mapSubscriptionsLifestyleTotal["1"].(map[string]interface{})
		switch mapSubscriptionsLifestyleTotalxDoc["subscriptionslifestyletotal"].(type) {
		case string:
			mapReport["subscriptionslifestyletotal"] = float64(0)
		case int64:
			fSubscriptionsLifestyleTotal = functions.Round(float64(mapSubscriptionsLifestyleTotalxDoc["subscriptionslifestyletotal"].(int64)))
			mapReport["subscriptionslifestyletotal"] = functions.ThousandSeperator(fSubscriptionsLifestyleTotal)
		case float64:
			fSubscriptionsLifestyleTotal = functions.Round(float64(mapSubscriptionsLifestyleTotalxDoc["subscriptionslifestyletotal"].(float64)))
			mapReport["subscriptionslifestyletotal"] = functions.ThousandSeperator(fSubscriptionsLifestyleTotal)
		}
	}

	fSubscriptionsBritishMumsTotal := float64(0)
	sqlSubscriptionsBritishMumsTotal := `select count(control) as subscriptionsbritishmumstotal from subscription  where schemecontrol in (select control from scheme where code in ('britishmums'))`
	mapSubscriptionsBritishMumsTotal, _ := curdb.Query(sqlSubscriptionsBritishMumsTotal)
	mapReport["subscriptionsbritishmumstotal"] = float64(0)
	if mapSubscriptionsBritishMumsTotal["1"] != nil {
		mapSubscriptionsBritishMumsTotalxDoc := mapSubscriptionsBritishMumsTotal["1"].(map[string]interface{})
		switch mapSubscriptionsBritishMumsTotalxDoc["subscriptionsbritishmumstotal"].(type) {
		case string:
			mapReport["subscriptionsbritishmumstotal"] = float64(0)
		case int64:
			fSubscriptionsBritishMumsTotal = functions.Round(float64(mapSubscriptionsBritishMumsTotalxDoc["subscriptionsbritishmumstotal"].(int64)))
			mapReport["subscriptionsbritishmumstotal"] = functions.ThousandSeperator(fSubscriptionsBritishMumsTotal)
		case float64:
			fSubscriptionsBritishMumsTotal = functions.Round(float64(mapSubscriptionsBritishMumsTotalxDoc["subscriptionsbritishmumstotal"].(float64)))
			mapReport["subscriptionsbritishmumstotal"] = functions.ThousandSeperator(fSubscriptionsBritishMumsTotal)
		}
	}
	// Number of Subscriptions
	//

	//---

	//
	//Number of Redeemed Rewards
	sqlNumRedeemedLiteTotal := `select count(distinct rewardcontrol) as numredeemedlitetotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))`

	mapNumRedeemedLiteTotal, _ := curdb.Query(sqlNumRedeemedLiteTotal)
	mapReport["numredeemedlitetotal"] = float64(0)
	if mapNumRedeemedLiteTotal["1"] != nil {
		mapNumRedeemedLiteTotalxDoc := mapNumRedeemedLiteTotal["1"].(map[string]interface{})
		switch mapNumRedeemedLiteTotalxDoc["numredeemedlitetotal"].(type) {
		case string:
			mapReport["numredeemedlitetotal"] = float64(0)
		case int64:
			mapReport["numredeemedlitetotal"] = functions.ThousandSeperator(functions.Round(float64(mapNumRedeemedLiteTotalxDoc["numredeemedlitetotal"].(int64))))
		case float64:
			mapReport["numredeemedlitetotal"] = functions.ThousandSeperator(functions.Round(mapNumRedeemedLiteTotalxDoc["numredeemedlitetotal"].(float64)))
		}
	}

	sqlNumRedeemedLifestyleTotal := `select count(distinct rewardcontrol) as numredeemedlifestyletotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))`

	mapNumRedeemedLifestyleTotal, _ := curdb.Query(sqlNumRedeemedLifestyleTotal)
	mapReport["numredeemedlifestyletotal"] = float64(0)
	if mapNumRedeemedLifestyleTotal["1"] != nil {
		mapNumRedeemedLifestyleTotalxDoc := mapNumRedeemedLifestyleTotal["1"].(map[string]interface{})
		switch mapNumRedeemedLifestyleTotalxDoc["numredeemedlifestyletotal"].(type) {
		case string:
			mapReport["numredeemedlifestyletotal"] = float64(0)
		case int64:
			mapReport["numredeemedlifestyletotal"] = functions.ThousandSeperator(functions.Round(float64(mapNumRedeemedLifestyleTotalxDoc["numredeemedlifestyletotal"].(int64))))
		case float64:
			mapReport["numredeemedlifestyletotal"] = functions.ThousandSeperator(functions.Round(mapNumRedeemedLifestyleTotalxDoc["numredeemedlifestyletotal"].(float64)))
		}
	}

	sqlNumRedeemedBritishMumsTotal := `select count(distinct rewardcontrol) as numredeemedbritishmumstotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))`

	mapNumRedeemedBritishMumsTotal, _ := curdb.Query(sqlNumRedeemedBritishMumsTotal)
	mapReport["numredeemedbritishmumstotal"] = float64(0)
	if mapNumRedeemedBritishMumsTotal["1"] != nil {
		mapNumRedeemedBritishMumsTotalxDoc := mapNumRedeemedBritishMumsTotal["1"].(map[string]interface{})
		switch mapNumRedeemedBritishMumsTotalxDoc["numredeemedbritishmumstotal"].(type) {
		case string:
			mapReport["numredeemedbritishmumstotal"] = float64(0)
		case int64:
			mapReport["numredeemedbritishmumstotal"] = functions.ThousandSeperator(functions.Round(float64(mapNumRedeemedBritishMumsTotalxDoc["numredeemedbritishmumstotal"].(int64))))
		case float64:
			mapReport["numredeemedbritishmumstotal"] = functions.ThousandSeperator(functions.Round(mapNumRedeemedBritishMumsTotalxDoc["numredeemedbritishmumstotal"].(float64)))
		}
	}

	//Number of Redeemed Rewards
	//

	//---

	//
	//AVERAGE SAVINGs per Redeemed Rewards
	sqlAvgSavingsPerRewardLiteTotal := `select sum(savingsvalue) / count(distinct rewardcontrol) as avgsavingperrewardlitetotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))`

	mapAvgSavingsPerRewardLiteTotal, _ := curdb.Query(sqlAvgSavingsPerRewardLiteTotal)
	mapReport["avgsavingperrewardlitetotal"] = float64(0)
	if mapAvgSavingsPerRewardLiteTotal["1"] != nil {
		mapAvgSavingsPerRewardLiteTotalxDoc := mapAvgSavingsPerRewardLiteTotal["1"].(map[string]interface{})
		switch mapAvgSavingsPerRewardLiteTotalxDoc["avgsavingperrewardlitetotal"].(type) {
		case string:
			mapReport["avgsavingperrewardlitetotal"] = float64(0)
		case int64:
			mapReport["avgsavingperrewardlitetotal"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerRewardLiteTotalxDoc["avgsavingperrewardlitetotal"].(int64))))
		case float64:
			mapReport["avgsavingperrewardlitetotal"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerRewardLiteTotalxDoc["avgsavingperrewardlitetotal"].(float64)))
		}
	}

	sqlAvgSavingsPerRewardLifestyleTotal := `select sum(savingsvalue) / count(distinct rewardcontrol) as avgsavingperrewardlifestyletotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))`

	mapAvgSavingsPerRewardLifestyleTotal, _ := curdb.Query(sqlAvgSavingsPerRewardLifestyleTotal)
	mapReport["avgsavingperrewardlifestyletotal"] = float64(0)
	if mapAvgSavingsPerRewardLifestyleTotal["1"] != nil {
		mapAvgSavingsPerRewardLifestyleTotalxDoc := mapAvgSavingsPerRewardLifestyleTotal["1"].(map[string]interface{})
		switch mapAvgSavingsPerRewardLifestyleTotalxDoc["avgsavingperrewardlifestyletotal"].(type) {
		case string:
			mapReport["avgsavingperrewardlifestyletotal"] = float64(0)
		case int64:
			mapReport["avgsavingperrewardlifestyletotal"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerRewardLifestyleTotalxDoc["avgsavingperrewardlifestyletotal"].(int64))))
		case float64:
			mapReport["avgsavingperrewardlifestyletotal"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerRewardLifestyleTotalxDoc["avgsavingperrewardlifestyletotal"].(float64)))
		}
	}

	sqlAvgSavingsPerRewardBritishMumsTotal := `select sum(savingsvalue) / count(distinct rewardcontrol) as avgsavingperrewardbritishmumstotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))`

	mapAvgSavingsPerRewardBritishMumsTotal, _ := curdb.Query(sqlAvgSavingsPerRewardBritishMumsTotal)
	mapReport["avgsavingperrewardbritishmumstotal"] = float64(0)
	if mapAvgSavingsPerRewardBritishMumsTotal["1"] != nil {
		mapAvgSavingsPerRewardBritishMumsTotalxDoc := mapAvgSavingsPerRewardBritishMumsTotal["1"].(map[string]interface{})
		switch mapAvgSavingsPerRewardBritishMumsTotalxDoc["avgsavingperrewardbritishmumstotal"].(type) {
		case string:
			mapReport["avgsavingperrewardbritishmumstotal"] = float64(0)
		case int64:
			mapReport["avgsavingperrewardbritishmumstotal"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerRewardBritishMumsTotalxDoc["avgsavingperrewardbritishmumstotal"].(int64))))
		case float64:
			mapReport["avgsavingperrewardbritishmumstotal"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerRewardBritishMumsTotalxDoc["avgsavingperrewardbritishmumstotal"].(float64)))
		}
	}

	//AVERAGE SAVINGs per Redeemed Rewards
	//

	//---

	//
	//AVERAGE SAVINGs per Employee
	sqlAvgSavingsPerEmployeeLiteTotal := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelitetotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))`

	mapAvgSavingsPerEmployeeLiteTotal, _ := curdb.Query(sqlAvgSavingsPerEmployeeLiteTotal)
	mapReport["avgsavingperemployeelitetotal"] = float64(0)
	if mapAvgSavingsPerEmployeeLiteTotal["1"] != nil {
		mapAvgSavingsPerEmployeeLiteTotalxDoc := mapAvgSavingsPerEmployeeLiteTotal["1"].(map[string]interface{})
		switch mapAvgSavingsPerEmployeeLiteTotalxDoc["avgsavingperemployeelitetotal"].(type) {
		case string:
			mapReport["avgsavingperemployeelitetotal"] = float64(0)
		case int64:
			mapReport["avgsavingperemployeelitetotal"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerEmployeeLiteTotalxDoc["avgsavingperemployeelitetotal"].(int64))))
		case float64:
			mapReport["avgsavingperemployeelitetotal"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerEmployeeLiteTotalxDoc["avgsavingperemployeelitetotal"].(float64)))
		}
	}

	sqlAvgSavingsPerEmployeeLifestyleTotal := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelifestyletotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))`

	mapAvgSavingsPerEmployeeLifestyleTotal, _ := curdb.Query(sqlAvgSavingsPerEmployeeLifestyleTotal)
	mapReport["avgsavingperemployeelifestyletotal"] = float64(0)
	if mapAvgSavingsPerEmployeeLifestyleTotal["1"] != nil {
		mapAvgSavingsPerEmployeeLifestyleTotalxDoc := mapAvgSavingsPerEmployeeLifestyleTotal["1"].(map[string]interface{})
		switch mapAvgSavingsPerEmployeeLifestyleTotalxDoc["avgsavingperemployeelifestyletotal"].(type) {
		case string:
			mapReport["avgsavingperemployeelifestyletotal"] = float64(0)
		case int64:
			mapReport["avgsavingperemployeelifestyletotal"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerEmployeeLifestyleTotalxDoc["avgsavingperemployeelifestyletotal"].(int64))))
		case float64:
			mapReport["avgsavingperemployeelifestyletotal"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerEmployeeLifestyleTotalxDoc["avgsavingperemployeelifestyletotal"].(float64)))
		}
	}

	sqlAvgSavingsPerEmployeeBritishMumsTotal := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelbritishmumstotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))`

	mapAvgSavingsPerEmployeeBritishMumsTotal, _ := curdb.Query(sqlAvgSavingsPerEmployeeBritishMumsTotal)
	mapReport["avgsavingperemployeelbritishmumstotal"] = float64(0)
	if mapAvgSavingsPerEmployeeBritishMumsTotal["1"] != nil {
		mapAvgSavingsPerEmployeeBritishMumsTotalxDoc := mapAvgSavingsPerEmployeeBritishMumsTotal["1"].(map[string]interface{})
		switch mapAvgSavingsPerEmployeeBritishMumsTotalxDoc["avgsavingperemployeelbritishmumstotal"].(type) {
		case string:
			mapReport["avgsavingperemployeelbritishmumstotal"] = float64(0)
		case int64:
			mapReport["avgsavingperemployeelbritishmumstotal"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerEmployeeBritishMumsTotalxDoc["avgsavingperemployeelbritishmumstotal"].(int64))))
		case float64:
			mapReport["avgsavingperemployeelbritishmumstotal"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerEmployeeBritishMumsTotalxDoc["avgsavingperemployeelbritishmumstotal"].(float64)))
		}
	}

	//AVERAGE SAVINGs per Employee
	//

	//---

	//
	//Total Saving
	sqlSavingLiteTotal := `select sum(savingsvalue) as savinglitetotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))`

	mapSavingLiteTotal, _ := curdb.Query(sqlSavingLiteTotal)
	mapReport["savinglitetotal"] = float64(0)
	if mapSavingLiteTotal["1"] != nil {
		mapSavingLiteTotalxDoc := mapSavingLiteTotal["1"].(map[string]interface{})
		switch mapSavingLiteTotalxDoc["savinglitetotal"].(type) {
		case string:
			mapReport["savinglitetotal"] = float64(0)
		case int64:
			mapReport["savinglitetotal"] = functions.ThousandSeperator(functions.Round(float64(mapSavingLiteTotalxDoc["savinglitetotal"].(int64))))
		case float64:
			mapReport["savinglitetotal"] = functions.ThousandSeperator(functions.Round(mapSavingLiteTotalxDoc["savinglitetotal"].(float64)))
		}
	}

	sqlSavingLifestyleTotal := `select sum(savingsvalue) as savinglifestyletotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))`

	mapSavingLifestyleTotal, _ := curdb.Query(sqlSavingLifestyleTotal)
	mapReport["savinglifestyletotal"] = float64(0)
	if mapSavingLifestyleTotal["1"] != nil {
		mapSavingLifestyleTotalxDoc := mapSavingLifestyleTotal["1"].(map[string]interface{})
		switch mapSavingLifestyleTotalxDoc["savinglifestyletotal"].(type) {
		case string:
			mapReport["savinglifestyletotal"] = float64(0)
		case int64:
			mapReport["savinglifestyletotal"] = functions.ThousandSeperator(functions.Round(float64(mapSavingLifestyleTotalxDoc["savinglifestyletotal"].(int64))))
		case float64:
			mapReport["savinglifestyletotal"] = functions.ThousandSeperator(functions.Round(mapSavingLifestyleTotalxDoc["savinglifestyletotal"].(float64)))
		}
	}

	sqlSavingBritishMumsTotal := `select sum(savingsvalue) as savingbritishmumstotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))`

	mapSavingBritishMumsTotal, _ := curdb.Query(sqlSavingBritishMumsTotal)
	mapReport["savingbritishmumstotal"] = float64(0)
	if mapSavingBritishMumsTotal["1"] != nil {
		mapSavingBritishMumsTotalxDoc := mapSavingBritishMumsTotal["1"].(map[string]interface{})
		switch mapSavingBritishMumsTotalxDoc["savingbritishmumstotal"].(type) {
		case string:
			mapReport["savingbritishmumstotal"] = float64(0)
		case int64:
			mapReport["savingbritishmumstotal"] = functions.ThousandSeperator(functions.Round(float64(mapSavingBritishMumsTotalxDoc["savingbritishmumstotal"].(int64))))
		case float64:
			mapReport["savingbritishmumstotal"] = functions.ThousandSeperator(functions.Round(mapSavingBritishMumsTotalxDoc["savingbritishmumstotal"].(float64)))
		}
	}

	//Total Saving
	//

	//---
	cFormat := "02/01/2006"
	sStopdate := time.Now().Format(cFormat)
	sStartdate := time.Now().Add(-(time.Hour * 24 * 365)).Format(cFormat)

	//Last Year Saving
	sqlYearSavingLiteTotal := fmt.Sprintf(`select sum(savingsvalue) as yearsavinglitetotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))
	and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapYearSavingLiteTotal, _ := curdb.Query(sqlYearSavingLiteTotal)
	mapReport["yearsavinglitetotal"] = float64(0)
	if mapYearSavingLiteTotal["1"] != nil {
		mapYearSavingLiteTotalxDoc := mapYearSavingLiteTotal["1"].(map[string]interface{})
		switch mapYearSavingLiteTotalxDoc["yearsavinglitetotal"].(type) {
		case string:
			mapReport["yearsavinglitetotal"] = float64(0)
		case int64:
			mapReport["yearsavinglitetotal"] = functions.ThousandSeperator(functions.Round(float64(mapYearSavingLiteTotalxDoc["yearsavinglitetotal"].(int64))))
		case float64:
			mapReport["yearsavinglitetotal"] = functions.ThousandSeperator(functions.Round(mapYearSavingLiteTotalxDoc["yearsavinglitetotal"].(float64)))
		}
	}

	sqlYearSavingLifestyleTotal := fmt.Sprintf(`select sum(savingsvalue) as yearsavinglifestyletotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))
	and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapYearSavingLifestyleTotal, _ := curdb.Query(sqlYearSavingLifestyleTotal)
	mapReport["yearsavinglifestyletotal"] = float64(0)
	if mapYearSavingLifestyleTotal["1"] != nil {
		mapYearSavingLifestyleTotalxDoc := mapYearSavingLifestyleTotal["1"].(map[string]interface{})
		switch mapYearSavingLifestyleTotalxDoc["yearsavinglifestyletotal"].(type) {
		case string:
			mapReport["yearsavinglifestyletotal"] = float64(0)
		case int64:
			mapReport["yearsavinglifestyletotal"] = functions.ThousandSeperator(functions.Round(float64(mapYearSavingLifestyleTotalxDoc["yearsavinglifestyletotal"].(int64))))
		case float64:
			mapReport["yearsavinglifestyletotal"] = functions.ThousandSeperator(functions.Round(mapYearSavingLifestyleTotalxDoc["yearsavinglifestyletotal"].(float64)))
		}
	}

	sqlYearSavingBritishMumsTotal := fmt.Sprintf(`select sum(savingsvalue) as yearsavingbritishmumstotal from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))
	and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapYearSavingBritishMumsTotal, _ := curdb.Query(sqlYearSavingBritishMumsTotal)
	mapReport["yearsavingbritishmumstotal"] = float64(0)
	if mapYearSavingBritishMumsTotal["1"] != nil {
		mapYearSavingBritishMumsTotalxDoc := mapYearSavingBritishMumsTotal["1"].(map[string]interface{})
		switch mapYearSavingBritishMumsTotalxDoc["yearsavingbritishmumstotal"].(type) {
		case string:
			mapReport["yearsavingbritishmumstotal"] = float64(0)
		case int64:
			mapReport["yearsavingbritishmumstotal"] = functions.ThousandSeperator(functions.Round(float64(mapYearSavingBritishMumsTotalxDoc["yearsavingbritishmumstotal"].(int64))))
		case float64:
			mapReport["yearsavingbritishmumstotal"] = functions.ThousandSeperator(functions.Round(mapYearSavingBritishMumsTotalxDoc["yearsavingbritishmumstotal"].(float64)))
		}
	}

	//Last Year Saving
	//

	//---

	// Total Subscriptions
	//

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-merchant"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Subscription | Admin Reports","mainpanelContent":` + contentHTML + `}`))

}
