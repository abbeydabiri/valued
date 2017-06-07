package frontend

import (
	"fmt"
	"time"
	"valued/database"
	"valued/functions"

	"net/http"
	"strconv"
)

func (this *Report) subscription(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

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

	sqlUsersPendingTotal := `select count(control) as userspendingtotal from profile where company != 'Yes' and workflow = 'pending'`
	mapUsersPendingTotal, _ := curdb.Query(sqlUsersPendingTotal)
	mapReport["userspendingtotal"] = float64(0)
	if mapUsersPendingTotal["1"] != nil {
		mapUsersPendingTotalxDoc := mapUsersPendingTotal["1"].(map[string]interface{})
		switch mapUsersPendingTotalxDoc["userspendingtotal"].(type) {
		case string:
			mapReport["userspendingtotal"] = float64(0)
		case int64:
			mapReport["userspendingtotal"] = functions.ThousandSeperator(functions.Round(float64(mapUsersPendingTotalxDoc["userspendingtotal"].(int64))))
		case float64:
			mapReport["userspendingtotal"] = functions.ThousandSeperator(functions.Round(mapUsersPendingTotalxDoc["userspendingtotal"].(float64)))
		}
	}

	sqlUsersPaidTotal := `select count(control) as userspaidtotal from profile where company != 'Yes' and workflow in ('paid','subscribed','subscribed-pending' )`
	mapUsersPaidTotal, _ := curdb.Query(sqlUsersPaidTotal)
	mapReport["userspaidtotal"] = float64(0)
	if mapUsersPaidTotal["1"] != nil {
		mapUsersPaidTotalxDoc := mapUsersPaidTotal["1"].(map[string]interface{})
		switch mapUsersPaidTotalxDoc["userspaidtotal"].(type) {
		case string:
			mapReport["userspaidtotal"] = float64(0)
		case int64:
			mapReport["userspaidtotal"] = functions.ThousandSeperator(functions.Round(float64(mapUsersPaidTotalxDoc["userspaidtotal"].(int64))))
		case float64:
			mapReport["userspaidtotal"] = functions.ThousandSeperator(functions.Round(mapUsersPaidTotalxDoc["userspaidtotal"].(float64)))
		}
	}

	sqlUsersExpiredTotal := `select count(control) as usersexpiredtotal from profile where company != 'Yes' and workflow in ('expired' )`
	mapUsersExpiredTotal, _ := curdb.Query(sqlUsersExpiredTotal)
	mapReport["usersexpiredtotal"] = float64(0)
	if mapUsersExpiredTotal["1"] != nil {
		mapUsersExpiredTotalxDoc := mapUsersExpiredTotal["1"].(map[string]interface{})
		switch mapUsersExpiredTotalxDoc["usersexpiredtotal"].(type) {
		case string:
			mapReport["usersexpiredtotal"] = float64(0)
		case int64:
			mapReport["usersexpiredtotal"] = functions.ThousandSeperator(functions.Round(float64(mapUsersExpiredTotalxDoc["usersexpiredtotal"].(int64))))
		case float64:
			mapReport["usersexpiredtotal"] = functions.ThousandSeperator(functions.Round(mapUsersExpiredTotalxDoc["usersexpiredtotal"].(float64)))
		}
	}

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
	sqlNumRedeemedLiteTotal := `select count(distinct rewardcontrol) as numredeemedlitetotal from redemption where schemecontrol = (select control from scheme where code = 'lite')`

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

	sqlNumRedeemedLifestyleTotal := `select count(distinct rewardcontrol) as numredeemedlifestyletotal from redemption where schemecontrol = (select control from scheme where code = 'lifestyle')`

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

	sqlNumRedeemedBritishMumsTotal := `select count(distinct rewardcontrol) as numredeemedbritishmumstotal from redemption where schemecontrol = (select control from scheme where code = 'britishmums')`

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

	//

	//

	//
	// Employers Subscriptions

	//
	// % OF SUBSCRIBERS WITH A GREATER SAVINGS AMOUNT THAN THE COST OF MEMBERSHIP
	sqlPercentSubscribersSavingsEmployer := `select subscription.membercontrol as membercontrol, sum(subscription.price) as subscription, 
								(select sum(redemption.savingsvalue)  from redemption where membercontrol = subscription.membercontrol) as savings
								from subscription where subscription.employercontrol not in (select control from profile where code in ('main','britishmums'))
								group by subscription.membercontrol order by 3 desc nulls last`

	mapPercentSubscribersSavingsEmployer, _ := curdb.Query(sqlPercentSubscribersSavingsEmployer)
	aPercentSubscribersSavingsEmployerSorted := functions.SortMap(mapPercentSubscribersSavingsEmployer)

	var mapSavingsEmployerGreater []string
	for _, sNumber := range aPercentSubscribersSavingsEmployerSorted {
		xDocSubscribersSavingsEmployer := mapPercentSubscribersSavingsEmployer[sNumber].(map[string]interface{})

		switch xDocSubscribersSavingsEmployer["savings"].(type) {
		case float64:
			if xDocSubscribersSavingsEmployer["savings"].(float64) >
				xDocSubscribersSavingsEmployer["subscription"].(float64) {
				mapSavingsEmployerGreater = append(mapSavingsEmployerGreater, xDocSubscribersSavingsEmployer["membercontrol"].(string))
			}

		}
	}

	mapReport["percentofsubscriberssavingsvscostsemployer"] = functions.RoundUp(float64(len(mapSavingsEmployerGreater)*100)/float64(len(aPercentSubscribersSavingsEmployerSorted)), 0)

	// % OF SUBSCRIBERS WITH A GREATER SAVINGS AMOUNT THAN THE COST OF MEMBERSHIP
	//

	//
	// Number of Users - Pending / Paid / Expired

	sqlUsersPendingLiteEmployer := `select count(control) as userspendingliteemployer from profile where company != 'Yes' and workflow = 'pending' 
									and employercontrol not in (select control from profile where code in ('main','britishmums'))
									and control in (
										select distinct membercontrol from subscription
										where workflow = 'inactive' and expirydate::timestamp > '%s'::timestamp
										and schemecontrol = (select control from scheme where code = 'lite')
									)`
	mapUsersPendingLiteEmployer, _ := curdb.Query(fmt.Sprintf(sqlUsersPendingLiteEmployer, time.Now().Format(cFormat)))
	mapReport["userspendingliteemployer"] = float64(0)
	if mapUsersPendingLiteEmployer["1"] != nil {
		mapUsersPendingLiteEmployerxDoc := mapUsersPendingLiteEmployer["1"].(map[string]interface{})
		switch mapUsersPendingLiteEmployerxDoc["userspendingliteemployer"].(type) {
		case string:
			mapReport["userspendingliteemployer"] = float64(0)
		case int64:
			mapReport["userspendingliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersPendingLiteEmployerxDoc["userspendingliteemployer"].(int64))))
		case float64:
			mapReport["userspendingliteemployer"] = functions.ThousandSeperator(functions.Round(mapUsersPendingLiteEmployerxDoc["userspendingliteemployer"].(float64)))
		}
	}

	sqlUsersPendingLifestyleEmployer := `select count(control) as userspendinglifestyleemployer from profile where company != 'Yes' and workflow = 'pending' 
									and employercontrol not in (select control from profile where code in ('main','britishmums'))
									and control in (
										select distinct membercontrol from subscription
										where workflow = 'inactive' and expirydate::timestamp > '%s'::timestamp
										and schemecontrol = (select control from scheme where code = 'lifestyle')
									)`
	mapUsersPendingLifestyleEmployer, _ := curdb.Query(fmt.Sprintf(sqlUsersPendingLifestyleEmployer, time.Now().Format(cFormat)))
	mapReport["userspendinglifestyleemployer"] = float64(0)
	if mapUsersPendingLifestyleEmployer["1"] != nil {
		mapUsersPendingLifestyleEmployerxDoc := mapUsersPendingLifestyleEmployer["1"].(map[string]interface{})
		switch mapUsersPendingLifestyleEmployerxDoc["userspendinglifestyleemployer"].(type) {
		case string:
			mapReport["userspendinglifestyleemployer"] = float64(0)
		case int64:
			mapReport["userspendinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersPendingLifestyleEmployerxDoc["userspendinglifestyleemployer"].(int64))))
		case float64:
			mapReport["userspendinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapUsersPendingLifestyleEmployerxDoc["userspendinglifestyleemployer"].(float64)))
		}
	}

	sqlUsersPaidLiteEmployer := `select count(control) as userspaidliteemployer from profile where company != 'Yes' and workflow in ('paid','subscribed','subscribed-paid' ) 
									and employercontrol not in (select control from profile where code in ('main','britishmums'))
									and control in (
										select distinct membercontrol from subscription
										where workflow = 'active' and expirydate::timestamp > '%s'::timestamp
										and schemecontrol = (select control from scheme where code = 'lite')
									)`
	mapUsersPaidLiteEmployer, _ := curdb.Query(fmt.Sprintf(sqlUsersPaidLiteEmployer, time.Now().Format(cFormat)))
	mapReport["userspaidliteemployer"] = float64(0)
	if mapUsersPaidLiteEmployer["1"] != nil {
		mapUsersPaidLiteEmployerxDoc := mapUsersPaidLiteEmployer["1"].(map[string]interface{})
		switch mapUsersPaidLiteEmployerxDoc["userspaidliteemployer"].(type) {
		case string:
			mapReport["userspaidliteemployer"] = float64(0)
		case int64:
			mapReport["userspaidliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersPaidLiteEmployerxDoc["userspaidliteemployer"].(int64))))
		case float64:
			mapReport["userspaidliteemployer"] = functions.ThousandSeperator(functions.Round(mapUsersPaidLiteEmployerxDoc["userspaidliteemployer"].(float64)))
		}
	}

	sqlUsersPaidLifestyleEmployer := `select count(control) as userspaidlifestyleemployer from profile where company != 'Yes' and workflow in ('paid','subscribed','subscribed-paid' ) 
									and employercontrol not in (select control from profile where code in ('main','britishmums'))
									and control in (
										select distinct membercontrol from subscription
										where workflow = 'active' and expirydate::timestamp > '%s'::timestamp
										and schemecontrol = (select control from scheme where code = 'lifestyle')
									)`
	mapUsersPaidLifestyleEmployer, _ := curdb.Query(fmt.Sprintf(sqlUsersPaidLifestyleEmployer, time.Now().Format(cFormat)))
	mapReport["userspaidlifestyleemployer"] = float64(0)
	if mapUsersPaidLifestyleEmployer["1"] != nil {
		mapUsersPaidLifestyleEmployerxDoc := mapUsersPaidLifestyleEmployer["1"].(map[string]interface{})
		switch mapUsersPaidLifestyleEmployerxDoc["userspaidlifestyleemployer"].(type) {
		case string:
			mapReport["userspaidlifestyleemployer"] = float64(0)
		case int64:
			mapReport["userspaidlifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersPaidLifestyleEmployerxDoc["userspaidlifestyleemployer"].(int64))))
		case float64:
			mapReport["userspaidlifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapUsersPaidLifestyleEmployerxDoc["userspaidlifestyleemployer"].(float64)))
		}
	}

	sqlUsersExpiredLiteEmployer := `select count(control) as usersexpiredliteemployer from profile where company != 'Yes' and workflow in ('expired') 
									and employercontrol not in (select control from profile where code in ('main','britishmums'))
									and control in (
										select distinct membercontrol from subscription
										where workflow = 'expired' and expirydate::timestamp < '%s'::timestamp
										and schemecontrol = (select control from scheme where code = 'lite')
									)`
	mapUsersExpiredLiteEmployer, _ := curdb.Query(fmt.Sprintf(sqlUsersExpiredLiteEmployer, time.Now().Format(cFormat)))
	mapReport["usersexpiredliteemployer"] = float64(0)
	if mapUsersExpiredLiteEmployer["1"] != nil {
		mapUsersExpiredLiteEmployerxDoc := mapUsersExpiredLiteEmployer["1"].(map[string]interface{})
		switch mapUsersExpiredLiteEmployerxDoc["usersexpiredliteemployer"].(type) {
		case string:
			mapReport["usersexpiredliteemployer"] = float64(0)
		case int64:
			mapReport["usersexpiredliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersExpiredLiteEmployerxDoc["usersexpiredliteemployer"].(int64))))
		case float64:
			mapReport["usersexpiredliteemployer"] = functions.ThousandSeperator(functions.Round(mapUsersExpiredLiteEmployerxDoc["usersexpiredliteemployer"].(float64)))
		}
	}

	sqlUsersExpiredLifestyleEmployer := `select count(control) as usersexpiredlifestyleemployer from profile where company != 'Yes' and workflow in ('expired') 
									and employercontrol not in (select control from profile where code in ('main','britishmums'))
									and control in (
										select distinct membercontrol from subscription
										where workflow = 'expired' and expirydate::timestamp < '%s'::timestamp
										and schemecontrol = (select control from scheme where code = 'lifestyle')
									)`
	mapUsersExpiredLifestyleEmployer, _ := curdb.Query(fmt.Sprintf(sqlUsersExpiredLifestyleEmployer, time.Now().Format(cFormat)))
	mapReport["usersexpiredlifestyleemployer"] = float64(0)
	if mapUsersExpiredLifestyleEmployer["1"] != nil {
		mapUsersExpiredLifestyleEmployerxDoc := mapUsersExpiredLifestyleEmployer["1"].(map[string]interface{})
		switch mapUsersExpiredLifestyleEmployerxDoc["usersexpiredlifestyleemployer"].(type) {
		case string:
			mapReport["usersexpiredlifestyleemployer"] = float64(0)
		case int64:
			mapReport["usersexpiredlifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersExpiredLifestyleEmployerxDoc["usersexpiredlifestyleemployer"].(int64))))
		case float64:
			mapReport["usersexpiredlifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapUsersExpiredLifestyleEmployerxDoc["usersexpiredlifestyleemployer"].(float64)))
		}
	}

	//
	// Number of Registrations
	sqlRegistrationLiteEmployer := `select count(control) as registrationsvaluedemployer from profile where company != 'Yes' and workflow = 'registered' 
										and employercontrol in (select control from profile where code in ('main'))`
	mapRegistrationLiteEmployer, _ := curdb.Query(sqlRegistrationLiteEmployer)
	mapReport["registrationsvaluedemployer"] = float64(0)
	if mapRegistrationLiteEmployer["1"] != nil {
		mapRegistrationLiteEmployerxDoc := mapRegistrationLiteEmployer["1"].(map[string]interface{})
		switch mapRegistrationLiteEmployerxDoc["registrationsvaluedemployer"].(type) {
		case string:
			mapReport["registrationsvaluedemployer"] = float64(0)
		case int64:
			mapReport["registrationsvaluedemployer"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationLiteEmployerxDoc["registrationsvaluedemployer"].(int64))))
		case float64:
			mapReport["registrationsvaluedemployer"] = functions.ThousandSeperator(functions.Round(mapRegistrationLiteEmployerxDoc["registrationsvaluedemployer"].(float64)))
		}
	}

	sqlRegistrationLifestyleEmployer := `select count(control) as registrationsemployersemployer from profile where company != 'Yes' and workflow = 'registered' 
										and employercontrol not in (select control from profile where code in ('main','britishmums'))`
	mapRegistrationLifestyleEmployer, _ := curdb.Query(sqlRegistrationLifestyleEmployer)
	mapReport["registrationsemployersemployer"] = float64(0)
	if mapRegistrationLifestyleEmployer["1"] != nil {
		mapRegistrationLifestyleEmployerxDoc := mapRegistrationLifestyleEmployer["1"].(map[string]interface{})
		switch mapRegistrationLifestyleEmployerxDoc["registrationsemployersemployer"].(type) {
		case string:
			mapReport["registrationsemployersemployer"] = float64(0)
		case int64:
			mapReport["registrationsemployersemployer"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationLifestyleEmployerxDoc["registrationsemployersemployer"].(int64))))
		case float64:
			mapReport["registrationsemployersemployer"] = functions.ThousandSeperator(functions.Round(mapRegistrationLifestyleEmployerxDoc["registrationsemployersemployer"].(float64)))
		}
	}

	sqlRegistrationBritishMumsEmployer := `select count(control) as registrationsbritishmumsemployer from profile where company != 'Yes' and workflow = 'registered' 
										and employercontrol in (select control from profile where code in ('britishmums'))`
	mapRegistrationBritishMumsEmployer, _ := curdb.Query(sqlRegistrationBritishMumsEmployer)
	mapReport["registrationsbritishmumsemployer"] = float64(0)
	if mapRegistrationBritishMumsEmployer["1"] != nil {
		mapRegistrationBritishMumsEmployerxDoc := mapRegistrationBritishMumsEmployer["1"].(map[string]interface{})
		switch mapRegistrationBritishMumsEmployerxDoc["registrationsbritishmumsemployer"].(type) {
		case string:
			mapReport["registrationsbritishmumsemployer"] = float64(0)
		case int64:
			mapReport["registrationsbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationBritishMumsEmployerxDoc["registrationsbritishmumsemployer"].(int64))))
		case float64:
			mapReport["registrationsbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(mapRegistrationBritishMumsEmployerxDoc["registrationsbritishmumsemployer"].(float64)))
		}
	}

	// Number of Registrations
	//

	//
	// Number of Subscriptions
	fSubscriptionsLiteEmployer := float64(0)
	sqlSubscriptionsLiteEmployer := `select count(control) as subscriptionsliteemployer from subscription  where schemecontrol in (select control from scheme where code in ('lite'))`
	mapSubscriptionsLiteEmployer, _ := curdb.Query(sqlSubscriptionsLiteEmployer)
	mapReport["subscriptionsliteemployer"] = float64(0)
	if mapSubscriptionsLiteEmployer["1"] != nil {
		mapSubscriptionsLiteEmployerxDoc := mapSubscriptionsLiteEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsLiteEmployerxDoc["subscriptionsliteemployer"].(type) {
		case string:
			mapReport["subscriptionsliteemployer"] = float64(0)
		case int64:
			fSubscriptionsLiteEmployer = functions.Round(float64(mapSubscriptionsLiteEmployerxDoc["subscriptionsliteemployer"].(int64)))
			mapReport["subscriptionsliteemployer"] = functions.ThousandSeperator(fSubscriptionsLiteEmployer)
		case float64:
			fSubscriptionsLiteEmployer = functions.Round(float64(mapSubscriptionsLiteEmployerxDoc["subscriptionsliteemployer"].(float64)))
			mapReport["subscriptionsliteemployer"] = functions.ThousandSeperator(fSubscriptionsLiteEmployer)
		}
	}

	fSubscriptionsLifestyleEmployer := float64(0)
	sqlSubscriptionsLifestyleEmployer := `select count(control) as subscriptionslifestyleemployer from subscription  where schemecontrol in (select control from scheme where code in ('lifestyle'))`
	mapSubscriptionsLifestyleEmployer, _ := curdb.Query(sqlSubscriptionsLifestyleEmployer)
	mapReport["subscriptionslifestyleemployer"] = float64(0)
	if mapSubscriptionsLifestyleEmployer["1"] != nil {
		mapSubscriptionsLifestyleEmployerxDoc := mapSubscriptionsLifestyleEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsLifestyleEmployerxDoc["subscriptionslifestyleemployer"].(type) {
		case string:
			mapReport["subscriptionslifestyleemployer"] = float64(0)
		case int64:
			fSubscriptionsLifestyleEmployer = functions.Round(float64(mapSubscriptionsLifestyleEmployerxDoc["subscriptionslifestyleemployer"].(int64)))
			mapReport["subscriptionslifestyleemployer"] = functions.ThousandSeperator(fSubscriptionsLifestyleEmployer)
		case float64:
			fSubscriptionsLifestyleEmployer = functions.Round(float64(mapSubscriptionsLifestyleEmployerxDoc["subscriptionslifestyleemployer"].(float64)))
			mapReport["subscriptionslifestyleemployer"] = functions.ThousandSeperator(fSubscriptionsLifestyleEmployer)
		}
	}

	fSubscriptionsBritishMumsEmployer := float64(0)
	sqlSubscriptionsBritishMumsEmployer := `select count(control) as subscriptionsbritishmumsemployer from subscription  where schemecontrol in (select control from scheme where code in ('britishmums'))`
	mapSubscriptionsBritishMumsEmployer, _ := curdb.Query(sqlSubscriptionsBritishMumsEmployer)
	mapReport["subscriptionsbritishmumsemployer"] = float64(0)
	if mapSubscriptionsBritishMumsEmployer["1"] != nil {
		mapSubscriptionsBritishMumsEmployerxDoc := mapSubscriptionsBritishMumsEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsBritishMumsEmployerxDoc["subscriptionsbritishmumsemployer"].(type) {
		case string:
			mapReport["subscriptionsbritishmumsemployer"] = float64(0)
		case int64:
			fSubscriptionsBritishMumsEmployer = functions.Round(float64(mapSubscriptionsBritishMumsEmployerxDoc["subscriptionsbritishmumsemployer"].(int64)))
			mapReport["subscriptionsbritishmumsemployer"] = functions.ThousandSeperator(fSubscriptionsBritishMumsEmployer)
		case float64:
			fSubscriptionsBritishMumsEmployer = functions.Round(float64(mapSubscriptionsBritishMumsEmployerxDoc["subscriptionsbritishmumsemployer"].(float64)))
			mapReport["subscriptionsbritishmumsemployer"] = functions.ThousandSeperator(fSubscriptionsBritishMumsEmployer)
		}
	}
	// Number of Subscriptions
	//

	//---

	//
	//Number of Redeemed Rewards
	sqlNumRedeemedLiteEmployer := `select count(distinct rewardcontrol) as numredeemedliteemployer from redemption where schemecontrol = (select control from scheme where code = 'lite')`

	mapNumRedeemedLiteEmployer, _ := curdb.Query(sqlNumRedeemedLiteEmployer)
	mapReport["numredeemedliteemployer"] = float64(0)
	if mapNumRedeemedLiteEmployer["1"] != nil {
		mapNumRedeemedLiteEmployerxDoc := mapNumRedeemedLiteEmployer["1"].(map[string]interface{})
		switch mapNumRedeemedLiteEmployerxDoc["numredeemedliteemployer"].(type) {
		case string:
			mapReport["numredeemedliteemployer"] = float64(0)
		case int64:
			mapReport["numredeemedliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapNumRedeemedLiteEmployerxDoc["numredeemedliteemployer"].(int64))))
		case float64:
			mapReport["numredeemedliteemployer"] = functions.ThousandSeperator(functions.Round(mapNumRedeemedLiteEmployerxDoc["numredeemedliteemployer"].(float64)))
		}
	}

	sqlNumRedeemedLifestyleEmployer := `select count(distinct rewardcontrol) as numredeemedlifestyleemployer from redemption where schemecontrol = (select control from scheme where code = 'lifestyle')`

	mapNumRedeemedLifestyleEmployer, _ := curdb.Query(sqlNumRedeemedLifestyleEmployer)
	mapReport["numredeemedlifestyleemployer"] = float64(0)
	if mapNumRedeemedLifestyleEmployer["1"] != nil {
		mapNumRedeemedLifestyleEmployerxDoc := mapNumRedeemedLifestyleEmployer["1"].(map[string]interface{})
		switch mapNumRedeemedLifestyleEmployerxDoc["numredeemedlifestyleemployer"].(type) {
		case string:
			mapReport["numredeemedlifestyleemployer"] = float64(0)
		case int64:
			mapReport["numredeemedlifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapNumRedeemedLifestyleEmployerxDoc["numredeemedlifestyleemployer"].(int64))))
		case float64:
			mapReport["numredeemedlifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapNumRedeemedLifestyleEmployerxDoc["numredeemedlifestyleemployer"].(float64)))
		}
	}

	sqlNumRedeemedBritishMumsEmployer := `select count(distinct rewardcontrol) as numredeemedbritishmumsemployer from redemption where schemecontrol = (select control from scheme where code = 'britishmums')`

	mapNumRedeemedBritishMumsEmployer, _ := curdb.Query(sqlNumRedeemedBritishMumsEmployer)
	mapReport["numredeemedbritishmumsemployer"] = float64(0)
	if mapNumRedeemedBritishMumsEmployer["1"] != nil {
		mapNumRedeemedBritishMumsEmployerxDoc := mapNumRedeemedBritishMumsEmployer["1"].(map[string]interface{})
		switch mapNumRedeemedBritishMumsEmployerxDoc["numredeemedbritishmumsemployer"].(type) {
		case string:
			mapReport["numredeemedbritishmumsemployer"] = float64(0)
		case int64:
			mapReport["numredeemedbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(float64(mapNumRedeemedBritishMumsEmployerxDoc["numredeemedbritishmumsemployer"].(int64))))
		case float64:
			mapReport["numredeemedbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(mapNumRedeemedBritishMumsEmployerxDoc["numredeemedbritishmumsemployer"].(float64)))
		}
	}

	//Number of Redeemed Rewards
	//

	//---

	//
	//AVERAGE SAVINGs per Redeemed Rewards
	sqlAvgSavingsPerRewardLiteEmployer := `select sum(savingsvalue) / count(distinct rewardcontrol) as avgsavingperrewardliteemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))`

	mapAvgSavingsPerRewardLiteEmployer, _ := curdb.Query(sqlAvgSavingsPerRewardLiteEmployer)
	mapReport["avgsavingperrewardliteemployer"] = float64(0)
	if mapAvgSavingsPerRewardLiteEmployer["1"] != nil {
		mapAvgSavingsPerRewardLiteEmployerxDoc := mapAvgSavingsPerRewardLiteEmployer["1"].(map[string]interface{})
		switch mapAvgSavingsPerRewardLiteEmployerxDoc["avgsavingperrewardliteemployer"].(type) {
		case string:
			mapReport["avgsavingperrewardliteemployer"] = float64(0)
		case int64:
			mapReport["avgsavingperrewardliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerRewardLiteEmployerxDoc["avgsavingperrewardliteemployer"].(int64))))
		case float64:
			mapReport["avgsavingperrewardliteemployer"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerRewardLiteEmployerxDoc["avgsavingperrewardliteemployer"].(float64)))
		}
	}

	sqlAvgSavingsPerRewardLifestyleEmployer := `select sum(savingsvalue) / count(distinct rewardcontrol) as avgsavingperrewardlifestyleemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))`

	mapAvgSavingsPerRewardLifestyleEmployer, _ := curdb.Query(sqlAvgSavingsPerRewardLifestyleEmployer)
	mapReport["avgsavingperrewardlifestyleemployer"] = float64(0)
	if mapAvgSavingsPerRewardLifestyleEmployer["1"] != nil {
		mapAvgSavingsPerRewardLifestyleEmployerxDoc := mapAvgSavingsPerRewardLifestyleEmployer["1"].(map[string]interface{})
		switch mapAvgSavingsPerRewardLifestyleEmployerxDoc["avgsavingperrewardlifestyleemployer"].(type) {
		case string:
			mapReport["avgsavingperrewardlifestyleemployer"] = float64(0)
		case int64:
			mapReport["avgsavingperrewardlifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerRewardLifestyleEmployerxDoc["avgsavingperrewardlifestyleemployer"].(int64))))
		case float64:
			mapReport["avgsavingperrewardlifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerRewardLifestyleEmployerxDoc["avgsavingperrewardlifestyleemployer"].(float64)))
		}
	}

	sqlAvgSavingsPerRewardBritishMumsEmployer := `select sum(savingsvalue) / count(distinct rewardcontrol) as avgsavingperrewardbritishmumsemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))`

	mapAvgSavingsPerRewardBritishMumsEmployer, _ := curdb.Query(sqlAvgSavingsPerRewardBritishMumsEmployer)
	mapReport["avgsavingperrewardbritishmumsemployer"] = float64(0)
	if mapAvgSavingsPerRewardBritishMumsEmployer["1"] != nil {
		mapAvgSavingsPerRewardBritishMumsEmployerxDoc := mapAvgSavingsPerRewardBritishMumsEmployer["1"].(map[string]interface{})
		switch mapAvgSavingsPerRewardBritishMumsEmployerxDoc["avgsavingperrewardbritishmumsemployer"].(type) {
		case string:
			mapReport["avgsavingperrewardbritishmumsemployer"] = float64(0)
		case int64:
			mapReport["avgsavingperrewardbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerRewardBritishMumsEmployerxDoc["avgsavingperrewardbritishmumsemployer"].(int64))))
		case float64:
			mapReport["avgsavingperrewardbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerRewardBritishMumsEmployerxDoc["avgsavingperrewardbritishmumsemployer"].(float64)))
		}
	}

	//AVERAGE SAVINGs per Redeemed Rewards
	//

	//---

	//
	//AVERAGE SAVINGs per Employee
	sqlAvgSavingsPerEmployeeLiteEmployer := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeeliteemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))`

	mapAvgSavingsPerEmployeeLiteEmployer, _ := curdb.Query(sqlAvgSavingsPerEmployeeLiteEmployer)
	mapReport["avgsavingperemployeeliteemployer"] = float64(0)
	if mapAvgSavingsPerEmployeeLiteEmployer["1"] != nil {
		mapAvgSavingsPerEmployeeLiteEmployerxDoc := mapAvgSavingsPerEmployeeLiteEmployer["1"].(map[string]interface{})
		switch mapAvgSavingsPerEmployeeLiteEmployerxDoc["avgsavingperemployeeliteemployer"].(type) {
		case string:
			mapReport["avgsavingperemployeeliteemployer"] = float64(0)
		case int64:
			mapReport["avgsavingperemployeeliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerEmployeeLiteEmployerxDoc["avgsavingperemployeeliteemployer"].(int64))))
		case float64:
			mapReport["avgsavingperemployeeliteemployer"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerEmployeeLiteEmployerxDoc["avgsavingperemployeeliteemployer"].(float64)))
		}
	}

	sqlAvgSavingsPerEmployeeLifestyleEmployer := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelifestyleemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))`

	mapAvgSavingsPerEmployeeLifestyleEmployer, _ := curdb.Query(sqlAvgSavingsPerEmployeeLifestyleEmployer)
	mapReport["avgsavingperemployeelifestyleemployer"] = float64(0)
	if mapAvgSavingsPerEmployeeLifestyleEmployer["1"] != nil {
		mapAvgSavingsPerEmployeeLifestyleEmployerxDoc := mapAvgSavingsPerEmployeeLifestyleEmployer["1"].(map[string]interface{})
		switch mapAvgSavingsPerEmployeeLifestyleEmployerxDoc["avgsavingperemployeelifestyleemployer"].(type) {
		case string:
			mapReport["avgsavingperemployeelifestyleemployer"] = float64(0)
		case int64:
			mapReport["avgsavingperemployeelifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerEmployeeLifestyleEmployerxDoc["avgsavingperemployeelifestyleemployer"].(int64))))
		case float64:
			mapReport["avgsavingperemployeelifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerEmployeeLifestyleEmployerxDoc["avgsavingperemployeelifestyleemployer"].(float64)))
		}
	}

	sqlAvgSavingsPerEmployeeBritishMumsEmployer := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelbritishmumsemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))`

	mapAvgSavingsPerEmployeeBritishMumsEmployer, _ := curdb.Query(sqlAvgSavingsPerEmployeeBritishMumsEmployer)
	mapReport["avgsavingperemployeelbritishmumsemployer"] = float64(0)
	if mapAvgSavingsPerEmployeeBritishMumsEmployer["1"] != nil {
		mapAvgSavingsPerEmployeeBritishMumsEmployerxDoc := mapAvgSavingsPerEmployeeBritishMumsEmployer["1"].(map[string]interface{})
		switch mapAvgSavingsPerEmployeeBritishMumsEmployerxDoc["avgsavingperemployeelbritishmumsemployer"].(type) {
		case string:
			mapReport["avgsavingperemployeelbritishmumsemployer"] = float64(0)
		case int64:
			mapReport["avgsavingperemployeelbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(float64(mapAvgSavingsPerEmployeeBritishMumsEmployerxDoc["avgsavingperemployeelbritishmumsemployer"].(int64))))
		case float64:
			mapReport["avgsavingperemployeelbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(mapAvgSavingsPerEmployeeBritishMumsEmployerxDoc["avgsavingperemployeelbritishmumsemployer"].(float64)))
		}
	}

	//AVERAGE SAVINGs per Employee
	//

	//---

	//
	//Employer Saving
	sqlSavingLiteEmployer := `select sum(savingsvalue) as savingliteemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))`

	mapSavingLiteEmployer, _ := curdb.Query(sqlSavingLiteEmployer)
	mapReport["savingliteemployer"] = float64(0)
	if mapSavingLiteEmployer["1"] != nil {
		mapSavingLiteEmployerxDoc := mapSavingLiteEmployer["1"].(map[string]interface{})
		switch mapSavingLiteEmployerxDoc["savingliteemployer"].(type) {
		case string:
			mapReport["savingliteemployer"] = float64(0)
		case int64:
			mapReport["savingliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSavingLiteEmployerxDoc["savingliteemployer"].(int64))))
		case float64:
			mapReport["savingliteemployer"] = functions.ThousandSeperator(functions.Round(mapSavingLiteEmployerxDoc["savingliteemployer"].(float64)))
		}
	}

	sqlSavingLifestyleEmployer := `select sum(savingsvalue) as savinglifestyleemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))`

	mapSavingLifestyleEmployer, _ := curdb.Query(sqlSavingLifestyleEmployer)
	mapReport["savinglifestyleemployer"] = float64(0)
	if mapSavingLifestyleEmployer["1"] != nil {
		mapSavingLifestyleEmployerxDoc := mapSavingLifestyleEmployer["1"].(map[string]interface{})
		switch mapSavingLifestyleEmployerxDoc["savinglifestyleemployer"].(type) {
		case string:
			mapReport["savinglifestyleemployer"] = float64(0)
		case int64:
			mapReport["savinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSavingLifestyleEmployerxDoc["savinglifestyleemployer"].(int64))))
		case float64:
			mapReport["savinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapSavingLifestyleEmployerxDoc["savinglifestyleemployer"].(float64)))
		}
	}

	sqlSavingBritishMumsEmployer := `select sum(savingsvalue) as savingbritishmumsemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))`

	mapSavingBritishMumsEmployer, _ := curdb.Query(sqlSavingBritishMumsEmployer)
	mapReport["savingbritishmumsemployer"] = float64(0)
	if mapSavingBritishMumsEmployer["1"] != nil {
		mapSavingBritishMumsEmployerxDoc := mapSavingBritishMumsEmployer["1"].(map[string]interface{})
		switch mapSavingBritishMumsEmployerxDoc["savingbritishmumsemployer"].(type) {
		case string:
			mapReport["savingbritishmumsemployer"] = float64(0)
		case int64:
			mapReport["savingbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSavingBritishMumsEmployerxDoc["savingbritishmumsemployer"].(int64))))
		case float64:
			mapReport["savingbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(mapSavingBritishMumsEmployerxDoc["savingbritishmumsemployer"].(float64)))
		}
	}

	//Employer Saving
	//

	//---

	//Last Year Saving
	sqlYearSavingLiteEmployer := fmt.Sprintf(`select sum(savingsvalue) as yearsavingliteemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lite'))
	and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapYearSavingLiteEmployer, _ := curdb.Query(sqlYearSavingLiteEmployer)
	mapReport["yearsavingliteemployer"] = float64(0)
	if mapYearSavingLiteEmployer["1"] != nil {
		mapYearSavingLiteEmployerxDoc := mapYearSavingLiteEmployer["1"].(map[string]interface{})
		switch mapYearSavingLiteEmployerxDoc["yearsavingliteemployer"].(type) {
		case string:
			mapReport["yearsavingliteemployer"] = float64(0)
		case int64:
			mapReport["yearsavingliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapYearSavingLiteEmployerxDoc["yearsavingliteemployer"].(int64))))
		case float64:
			mapReport["yearsavingliteemployer"] = functions.ThousandSeperator(functions.Round(mapYearSavingLiteEmployerxDoc["yearsavingliteemployer"].(float64)))
		}
	}

	sqlYearSavingLifestyleEmployer := fmt.Sprintf(`select sum(savingsvalue) as yearsavinglifestyleemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'lifestyle'))
	and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapYearSavingLifestyleEmployer, _ := curdb.Query(sqlYearSavingLifestyleEmployer)
	mapReport["yearsavinglifestyleemployer"] = float64(0)
	if mapYearSavingLifestyleEmployer["1"] != nil {
		mapYearSavingLifestyleEmployerxDoc := mapYearSavingLifestyleEmployer["1"].(map[string]interface{})
		switch mapYearSavingLifestyleEmployerxDoc["yearsavinglifestyleemployer"].(type) {
		case string:
			mapReport["yearsavinglifestyleemployer"] = float64(0)
		case int64:
			mapReport["yearsavinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapYearSavingLifestyleEmployerxDoc["yearsavinglifestyleemployer"].(int64))))
		case float64:
			mapReport["yearsavinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapYearSavingLifestyleEmployerxDoc["yearsavinglifestyleemployer"].(float64)))
		}
	}

	sqlYearSavingBritishMumsEmployer := fmt.Sprintf(`select sum(savingsvalue) as yearsavingbritishmumsemployer from redemption where rewardcontrol in 
	(select distinct rewardcontrol from rewardscheme where schemecontrol = (select control from scheme where code = 'britishmums'))
	and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

	mapYearSavingBritishMumsEmployer, _ := curdb.Query(sqlYearSavingBritishMumsEmployer)
	mapReport["yearsavingbritishmumsemployer"] = float64(0)
	if mapYearSavingBritishMumsEmployer["1"] != nil {
		mapYearSavingBritishMumsEmployerxDoc := mapYearSavingBritishMumsEmployer["1"].(map[string]interface{})
		switch mapYearSavingBritishMumsEmployerxDoc["yearsavingbritishmumsemployer"].(type) {
		case string:
			mapReport["yearsavingbritishmumsemployer"] = float64(0)
		case int64:
			mapReport["yearsavingbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(float64(mapYearSavingBritishMumsEmployerxDoc["yearsavingbritishmumsemployer"].(int64))))
		case float64:
			mapReport["yearsavingbritishmumsemployer"] = functions.ThousandSeperator(functions.Round(mapYearSavingBritishMumsEmployerxDoc["yearsavingbritishmumsemployer"].(float64)))
		}
	}

	//Last Year Saving
	//

	//---

	// Employers Subscriptions
	//

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-subscription"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Subscription | Admin Reports","mainpanelContent":` + contentHTML + `}`))

}
