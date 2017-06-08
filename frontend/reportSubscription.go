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
	fUSERSPAIDTOTAL := float64(0)
	if mapUsersPaidTotal["1"] != nil {
		mapUsersPaidTotalxDoc := mapUsersPaidTotal["1"].(map[string]interface{})
		switch mapUsersPaidTotalxDoc["userspaidtotal"].(type) {
		case string:
			mapReport["userspaidtotal"] = float64(0)
		case int64:
			fUSERSPAIDTOTAL = functions.Round(float64(mapUsersPaidTotalxDoc["userspaidtotal"].(int64)))
			mapReport["userspaidtotal"] = functions.ThousandSeperator(fUSERSPAIDTOTAL)
		case float64:
			fUSERSPAIDTOTAL = functions.Round(mapUsersPaidTotalxDoc["userspaidtotal"].(float64))
			mapReport["userspaidtotal"] = functions.ThousandSeperator(fUSERSPAIDTOTAL)
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
	sqlNumRedeemedLiteTotal := `select count(rewardcontrol) as numredeemedlitetotal from redemption where schemecontrol = (select control from scheme where code = 'lite')`

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

	sqlNumRedeemedLifestyleTotal := `select count(rewardcontrol) as numredeemedlifestyletotal from redemption where schemecontrol = (select control from scheme where code = 'lifestyle')`

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

	sqlNumRedeemedBritishMumsTotal := `select count(rewardcontrol) as numredeemedbritishmumstotal from redemption where schemecontrol = (select control from scheme where code = 'britishmums')`

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
	sqlAvgSavingsPerRewardLiteTotal := `select sum(savingsvalue) / count(rewardcontrol) as avgsavingperrewardlitetotal from redemption where 
	schemecontrol = (select control from scheme where code = 'lite')`

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

	sqlAvgSavingsPerRewardLifestyleTotal := `select sum(savingsvalue) / count(rewardcontrol) as avgsavingperrewardlifestyletotal from redemption where 
	schemecontrol = (select control from scheme where code = 'lifestyle')`

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

	sqlAvgSavingsPerRewardBritishMumsTotal := `select sum(savingsvalue) / count(rewardcontrol) as avgsavingperrewardbritishmumstotal from redemption where 
	schemecontrol = (select control from scheme where code = 'britishmums')`

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
	sqlAvgSavingsPerEmployeeLiteTotal := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelitetotal from redemption where 
	schemecontrol = (select control from scheme where code = 'lite')`

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

	sqlAvgSavingsPerEmployeeLifestyleTotal := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelifestyletotal from redemption where 
	schemecontrol = (select control from scheme where code = 'lifestyle')`

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

	sqlAvgSavingsPerEmployeeBritishMumsTotal := `select sum(savingsvalue) / count(distinct membercontrol) as avgsavingperemployeelbritishmumstotal from redemption where
	schemecontrol = (select control from scheme where code = 'britishmums')`

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
	sqlSavingLiteTotal := `select sum(savingsvalue) as savinglitetotal from redemption where schemecontrol = (select control from scheme where code = 'lite')`

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

	sqlSavingLifestyleTotal := `select sum(savingsvalue) as savinglifestyletotal from redemption where schemecontrol = (select control from scheme where code = 'lifestyle')`

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

	sqlSavingBritishMumsTotal := `select sum(savingsvalue) as savingbritishmumstotal from redemption where schemecontrol = (select control from scheme where code = 'britishmums')`

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
	sqlYearSavingLiteTotal := fmt.Sprintf(`select sum(savingsvalue) as yearsavinglitetotal from redemption where schemecontrol = 
		(select control from scheme where code = 'lite') and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

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

	sqlYearSavingLifestyleTotal := fmt.Sprintf(`select sum(savingsvalue) as yearsavinglifestyletotal from redemption where schemecontrol = 
		(select control from scheme where code = 'lifestyle') and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

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

	sqlYearSavingBritishMumsTotal := fmt.Sprintf(`select sum(savingsvalue) as yearsavingbritishmumstotal from redemption where schemecontrol = 
		(select control from scheme where code = 'britishmums')	and substring(createdate from 1 for 20)::timestamp between '%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp`, sStartdate, sStopdate)

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

	mapReport["percentofsubscriberssavingsvscostsemployer"] = float64(0)

	if len(mapSavingsEmployerGreater) > 0 {
		mapReport["percentofsubscriberssavingsvscostsemployer"] = functions.RoundUp(float64(len(mapSavingsEmployerGreater)*100)/float64(len(aPercentSubscribersSavingsEmployerSorted)), 0)
	}

	// % OF SUBSCRIBERS WITH A GREATER SAVINGS AMOUNT THAN THE COST OF MEMBERSHIP
	//

	//
	//
	// Number of Users - Pending / Paid / Expired

	sqlSubscriptionsPendingLiteEmployer := `select count(control) as subscriptionspendingliteemployer from subscription
									where workflow = 'inactive' and expirydate::timestamp > '%s'::timestamp
									and schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapSubscriptionsPendingLiteEmployer, _ := curdb.Query(fmt.Sprintf(sqlSubscriptionsPendingLiteEmployer, time.Now().Format(cFormat)))
	mapReport["subscriptionspendingliteemployer"] = float64(0)
	if mapSubscriptionsPendingLiteEmployer["1"] != nil {
		mapSubscriptionsPendingLiteEmployerxDoc := mapSubscriptionsPendingLiteEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsPendingLiteEmployerxDoc["subscriptionspendingliteemployer"].(type) {
		case string:
			mapReport["subscriptionspendingliteemployer"] = float64(0)
		case int64:
			mapReport["subscriptionspendingliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSubscriptionsPendingLiteEmployerxDoc["subscriptionspendingliteemployer"].(int64))))
		case float64:
			mapReport["subscriptionspendingliteemployer"] = functions.ThousandSeperator(functions.Round(mapSubscriptionsPendingLiteEmployerxDoc["subscriptionspendingliteemployer"].(float64)))
		}
	}

	sqlSubscriptionsPendingLifestyleEmployer := `select count(control) as subscriptionspendinglifestyleemployer from subscription
									where workflow = 'inactive' and expirydate::timestamp > '%s'::timestamp
									and schemecontrol = (select control from scheme where code = 'lifestyle')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapSubscriptionsPendingLifestyleEmployer, _ := curdb.Query(fmt.Sprintf(sqlSubscriptionsPendingLifestyleEmployer, time.Now().Format(cFormat)))
	mapReport["subscriptionspendinglifestyleemployer"] = float64(0)
	if mapSubscriptionsPendingLifestyleEmployer["1"] != nil {
		mapSubscriptionsPendingLifestyleEmployerxDoc := mapSubscriptionsPendingLifestyleEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsPendingLifestyleEmployerxDoc["subscriptionspendinglifestyleemployer"].(type) {
		case string:
			mapReport["subscriptionspendinglifestyleemployer"] = float64(0)
		case int64:
			mapReport["subscriptionspendinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSubscriptionsPendingLifestyleEmployerxDoc["subscriptionspendinglifestyleemployer"].(int64))))
		case float64:
			mapReport["subscriptionspendinglifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapSubscriptionsPendingLifestyleEmployerxDoc["subscriptionspendinglifestyleemployer"].(float64)))
		}
	}

	sqlSubscriptionsPaidLiteEmployer := `select count(control) as subscriptionspaidliteemployer from subscription
									where workflow in ('paid','active' ) and expirydate::timestamp > '%s'::timestamp
									and schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapSubscriptionsPaidLiteEmployer, _ := curdb.Query(fmt.Sprintf(sqlSubscriptionsPaidLiteEmployer, time.Now().Format(cFormat)))
	mapReport["subscriptionspaidliteemployer"] = float64(0)
	if mapSubscriptionsPaidLiteEmployer["1"] != nil {
		mapSubscriptionsPaidLiteEmployerxDoc := mapSubscriptionsPaidLiteEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsPaidLiteEmployerxDoc["subscriptionspaidliteemployer"].(type) {
		case string:
			mapReport["subscriptionspaidliteemployer"] = float64(0)
		case int64:
			mapReport["subscriptionspaidliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSubscriptionsPaidLiteEmployerxDoc["subscriptionspaidliteemployer"].(int64))))
		case float64:
			mapReport["subscriptionspaidliteemployer"] = functions.ThousandSeperator(functions.Round(mapSubscriptionsPaidLiteEmployerxDoc["subscriptionspaidliteemployer"].(float64)))
		}
	}

	sqlSubscriptionsPaidLifestyleEmployer := `select count(control) as subscriptionspaidliteemployer from subscription
									where workflow in ('paid','active' ) and expirydate::timestamp > '%s'::timestamp
									and schemecontrol = (select control from scheme where code = 'lifestyle')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapSubscriptionsPaidLifestyleEmployer, _ := curdb.Query(fmt.Sprintf(sqlSubscriptionsPaidLifestyleEmployer, time.Now().Format(cFormat)))
	mapReport["subscriptionspaidlifestyleemployer"] = float64(0)
	if mapSubscriptionsPaidLifestyleEmployer["1"] != nil {
		mapSubscriptionsPaidLifestyleEmployerxDoc := mapSubscriptionsPaidLifestyleEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsPaidLifestyleEmployerxDoc["subscriptionspaidlifestyleemployer"].(type) {
		case string:
			mapReport["subscriptionspaidlifestyleemployer"] = float64(0)
		case int64:
			mapReport["subscriptionspaidlifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSubscriptionsPaidLifestyleEmployerxDoc["subscriptionspaidlifestyleemployer"].(int64))))
		case float64:
			mapReport["subscriptionspaidlifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapSubscriptionsPaidLifestyleEmployerxDoc["subscriptionspaidlifestyleemployer"].(float64)))
		}
	}

	sqlSubscriptionsExpiredLiteEmployer := `select count(control) as subscriptionsexpiredliteemployer from subscription
									where workflow = 'inactive' and expirydate::timestamp > '%s'::timestamp
									and schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapSubscriptionsExpiredLiteEmployer, _ := curdb.Query(fmt.Sprintf(sqlSubscriptionsExpiredLiteEmployer, time.Now().Format(cFormat)))
	mapReport["subscriptionsexpiredliteemployer"] = float64(0)
	if mapSubscriptionsExpiredLiteEmployer["1"] != nil {
		mapSubscriptionsExpiredLiteEmployerxDoc := mapSubscriptionsExpiredLiteEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsExpiredLiteEmployerxDoc["subscriptionsexpiredliteemployer"].(type) {
		case string:
			mapReport["subscriptionsexpiredliteemployer"] = float64(0)
		case int64:
			mapReport["subscriptionsexpiredliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSubscriptionsExpiredLiteEmployerxDoc["subscriptionsexpiredliteemployer"].(int64))))
		case float64:
			mapReport["subscriptionsexpiredliteemployer"] = functions.ThousandSeperator(functions.Round(mapSubscriptionsExpiredLiteEmployerxDoc["subscriptionsexpiredliteemployer"].(float64)))
		}
	}

	sqlSubscriptionsExpiredLifestyleEmployer := `select count(control) as subscriptionsexpiredlifestyleemployer from subscription
									where workflow = 'inactive' and expirydate::timestamp > '%s'::timestamp
									and schemecontrol = (select control from scheme where code = 'lifestyle')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapSubscriptionsExpiredLifestyleEmployer, _ := curdb.Query(fmt.Sprintf(sqlSubscriptionsExpiredLifestyleEmployer, time.Now().Format(cFormat)))
	mapReport["subscriptionsexpiredlifestyleemployer"] = float64(0)
	if mapSubscriptionsExpiredLifestyleEmployer["1"] != nil {
		mapSubscriptionsExpiredLifestyleEmployerxDoc := mapSubscriptionsExpiredLifestyleEmployer["1"].(map[string]interface{})
		switch mapSubscriptionsExpiredLifestyleEmployerxDoc["subscriptionsexpiredlifestyleemployer"].(type) {
		case string:
			mapReport["subscriptionsexpiredlifestyleemployer"] = float64(0)
		case int64:
			mapReport["subscriptionsexpiredlifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapSubscriptionsExpiredLifestyleEmployerxDoc["subscriptionsexpiredlifestyleemployer"].(int64))))
		case float64:
			mapReport["subscriptionsexpiredlifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapSubscriptionsExpiredLifestyleEmployerxDoc["subscriptionsexpiredlifestyleemployer"].(float64)))
		}
	}

	// Number of Users - Pending / Paid / Expired
	//
	//

	//
	//
	// Active Employees & Percentage of Employees
	//
	//

	sqlUsersActiveLiteEmployer := `select count(distinct membercontrol) as usersactiveliteemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersActiveLiteEmployer, _ := curdb.Query(sqlUsersActiveLiteEmployer)
	mapReport["usersactiveliteemployer"] = float64(0)
	fUSERSACTIVELITEEMPLOYER := float64(0)
	if mapUsersActiveLiteEmployer["1"] != nil {
		mapUsersActiveLiteEmployerxDoc := mapUsersActiveLiteEmployer["1"].(map[string]interface{})
		switch mapUsersActiveLiteEmployerxDoc["usersactiveliteemployer"].(type) {
		case string:
			mapReport["usersactiveliteemployer"] = float64(0)
		case int64:
			fUSERSACTIVELITEEMPLOYER = functions.Round(float64(mapUsersActiveLiteEmployerxDoc["usersactiveliteemployer"].(int64)))
			mapReport["usersactiveliteemployer"] = functions.ThousandSeperator(fUSERSACTIVELITEEMPLOYER)
		case float64:
			fUSERSACTIVELITEEMPLOYER = functions.Round(mapUsersActiveLiteEmployerxDoc["usersactiveliteemployer"].(float64))
			mapReport["usersactiveliteemployer"] = functions.ThousandSeperator(fUSERSACTIVELITEEMPLOYER)
		}
	}

	sqlUsersActiveLifestyleEmployer := `select count(distinct membercontrol) as usersactivelifestyleemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersActiveLifestyleEmployer, _ := curdb.Query(sqlUsersActiveLifestyleEmployer)
	mapReport["usersactivelifestyleemployer"] = float64(0)
	fUSERSACTIVELIFESTYLEEMPLOYER := float64(0)
	if mapUsersActiveLifestyleEmployer["1"] != nil {
		mapUsersActiveLifestyleEmployerxDoc := mapUsersActiveLifestyleEmployer["1"].(map[string]interface{})
		switch mapUsersActiveLifestyleEmployerxDoc["usersactivelifestyleemployer"].(type) {
		case string:
			mapReport["usersactivelifestyleemployer"] = float64(0)
		case int64:
			fUSERSACTIVELIFESTYLEEMPLOYER = functions.Round(float64(mapUsersActiveLifestyleEmployerxDoc["usersactivelifestyleemployer"].(int64)))
			mapReport["usersactivelifestyleemployer"] = functions.ThousandSeperator(fUSERSACTIVELIFESTYLEEMPLOYER)
		case float64:
			fUSERSACTIVELIFESTYLEEMPLOYER = functions.Round(mapUsersActiveLifestyleEmployerxDoc["usersactivelifestyleemployer"].(float64))
			mapReport["usersactivelifestyleemployer"] = functions.ThousandSeperator(fUSERSACTIVELIFESTYLEEMPLOYER)
		}
	}

	mapReport["usersactivepercentliteemployer"] = float64(0)
	mapReport["usersactivepercentlifestyleemployer"] = float64(0)
	if fUSERSPAIDTOTAL > float64(0) {
		if fUSERSACTIVELITEEMPLOYER > float64(0) {
			mapReport["usersactivepercentliteemployer"] = functions.RoundUp((fUSERSACTIVELITEEMPLOYER*float64(100))/fUSERSPAIDTOTAL, 0)
		}

		if fUSERSACTIVELIFESTYLEEMPLOYER > float64(0) {
			mapReport["usersactivepercentlifestyleemployer"] = functions.RoundUp((fUSERSACTIVELIFESTYLEEMPLOYER*float64(100))/fUSERSPAIDTOTAL, 0)
		}
	}

	// Active Employees & Percentage of Employees
	//

	//---

	//
	//Number of REWARDS REDEEMED BY EMPLOYEES
	sqlUsersRedeemedLiteEmployer := `select count(rewardcontrol) as usersredeemedliteemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersRedeemedLiteEmployer, _ := curdb.Query(sqlUsersRedeemedLiteEmployer)
	mapReport["usersredeemedliteemployer"] = float64(0)
	if mapUsersRedeemedLiteEmployer["1"] != nil {
		mapUsersRedeemedLiteEmployerxDoc := mapUsersRedeemedLiteEmployer["1"].(map[string]interface{})
		switch mapUsersRedeemedLiteEmployerxDoc["usersredeemedliteemployer"].(type) {
		case string:
			mapReport["usersredeemedliteemployer"] = float64(0)
		case int64:
			mapReport["usersredeemedliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersRedeemedLiteEmployerxDoc["usersredeemedliteemployer"].(int64))))
		case float64:
			mapReport["usersredeemedliteemployer"] = functions.ThousandSeperator(functions.Round(mapUsersRedeemedLiteEmployerxDoc["usersredeemedliteemployer"].(float64)))
		}
	}

	sqlUsersRedeemedLifestyleEmployer := `select count(rewardcontrol) as usersredeemedlifestyleemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lifestyle')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersRedeemedLifestyleEmployer, _ := curdb.Query(sqlUsersRedeemedLifestyleEmployer)
	mapReport["usersredeemedlifestyleemployer"] = float64(0)
	if mapUsersRedeemedLifestyleEmployer["1"] != nil {
		mapUsersRedeemedLifestyleEmployerxDoc := mapUsersRedeemedLifestyleEmployer["1"].(map[string]interface{})
		switch mapUsersRedeemedLifestyleEmployerxDoc["usersredeemedlifestyleemployer"].(type) {
		case string:
			mapReport["usersredeemedlifestyleemployer"] = float64(0)
		case int64:
			mapReport["usersredeemedlifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersRedeemedLifestyleEmployerxDoc["usersredeemedlifestyleemployer"].(int64))))
		case float64:
			mapReport["usersredeemedlifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapUsersRedeemedLifestyleEmployerxDoc["usersredeemedlifestyleemployer"].(float64)))
		}
	}

	//Number of REWARDS REDEEMED BY EMPLOYEES
	//

	//---

	//
	//AVERAGE SAVINGS per EMPLOYEE
	sqlUsersAvgSavingsLiteEmployer := `select sum(savingsvalue) / count(distinct membercontrol) as usersavgsavingsliteemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersAvgSavingsLiteEmployer, _ := curdb.Query(sqlUsersAvgSavingsLiteEmployer)
	mapReport["usersavgsavingsliteemployer"] = float64(0)
	if mapUsersAvgSavingsLiteEmployer["1"] != nil {
		mapUsersAvgSavingsLiteEmployerxDoc := mapUsersAvgSavingsLiteEmployer["1"].(map[string]interface{})
		switch mapUsersAvgSavingsLiteEmployerxDoc["usersavgsavingsliteemployer"].(type) {
		case string:
			mapReport["usersavgsavingsliteemployer"] = float64(0)
		case int64:
			mapReport["usersavgsavingsliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersAvgSavingsLiteEmployerxDoc["usersavgsavingsliteemployer"].(int64))))
		case float64:
			mapReport["usersavgsavingsliteemployer"] = functions.ThousandSeperator(functions.Round(mapUsersAvgSavingsLiteEmployerxDoc["usersavgsavingsliteemployer"].(float64)))
		}
	}

	sqlUsersAvgSavingsLifestyleEmployer := `select sum(savingsvalue) / count(distinct membercontrol) as usersavgsavingslifestyleemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lifestyle')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersAvgSavingsLifestyleEmployer, _ := curdb.Query(sqlUsersAvgSavingsLifestyleEmployer)
	mapReport["usersavgsavingslifestyleemployer"] = float64(0)
	if mapUsersAvgSavingsLifestyleEmployer["1"] != nil {
		mapUsersAvgSavingsLifestyleEmployerxDoc := mapUsersAvgSavingsLifestyleEmployer["1"].(map[string]interface{})
		switch mapUsersAvgSavingsLifestyleEmployerxDoc["usersavgsavingslifestyleemployer"].(type) {
		case string:
			mapReport["usersavgsavingslifestyleemployer"] = float64(0)
		case int64:
			mapReport["usersavgsavingslifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersAvgSavingsLifestyleEmployerxDoc["usersavgsavingslifestyleemployer"].(int64))))
		case float64:
			mapReport["usersavgsavingslifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapUsersAvgSavingsLifestyleEmployerxDoc["usersavgsavingslifestyleemployer"].(float64)))
		}
	}

	//AVERAGE SAVINGS per EMPLOYEE
	//

	//---

	//
	//TOTAL EMPLOYEE SAVINGS
	sqlUsersSavingsLiteEmployer := `select sum(savingsvalue) as userssavingsliteemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lite')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersSavingsLiteEmployer, _ := curdb.Query(sqlUsersSavingsLiteEmployer)
	mapReport["userssavingsliteemployer"] = float64(0)
	if mapUsersSavingsLiteEmployer["1"] != nil {
		mapUsersSavingsLiteEmployerxDoc := mapUsersSavingsLiteEmployer["1"].(map[string]interface{})
		switch mapUsersSavingsLiteEmployerxDoc["userssavingsliteemployer"].(type) {
		case string:
			mapReport["userssavingsliteemployer"] = float64(0)
		case int64:
			mapReport["userssavingsliteemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersSavingsLiteEmployerxDoc["userssavingsliteemployer"].(int64))))
		case float64:
			mapReport["userssavingsliteemployer"] = functions.ThousandSeperator(functions.Round(mapUsersSavingsLiteEmployerxDoc["userssavingsliteemployer"].(float64)))
		}
	}

	sqlUsersSavingsLifestyleEmployer := `select sum(savingsvalue) as userssavingslifestyleemployer from redemption 
									where schemecontrol = (select control from scheme where code = 'lifestyle')
									and employercontrol not in (select control from profile where code in ('main','britishmums'))`

	mapUsersSavingsLifestyleEmployer, _ := curdb.Query(sqlUsersSavingsLifestyleEmployer)
	mapReport["userssavingslifestyleemployer"] = float64(0)
	if mapUsersSavingsLifestyleEmployer["1"] != nil {
		mapUsersSavingsLifestyleEmployerxDoc := mapUsersSavingsLifestyleEmployer["1"].(map[string]interface{})
		switch mapUsersSavingsLifestyleEmployerxDoc["userssavingslifestyleemployer"].(type) {
		case string:
			mapReport["userssavingslifestyleemployer"] = float64(0)
		case int64:
			mapReport["userssavingslifestyleemployer"] = functions.ThousandSeperator(functions.Round(float64(mapUsersSavingsLifestyleEmployerxDoc["userssavingslifestyleemployer"].(int64))))
		case float64:
			mapReport["userssavingslifestyleemployer"] = functions.ThousandSeperator(functions.Round(mapUsersSavingsLifestyleEmployerxDoc["userssavingslifestyleemployer"].(float64)))
		}
	}

	//TOTAL EMPLOYEE SAVINGS
	//

	//---

	//
	//
	//TOP 5 EMPLOYERS WITH THE MOST SUBSCRIPTIONS

	sqlTop5EmployerSubscriptions := `select  (select title from profile where control = subscription.employercontrol) as employer, count(control) as subscription
										from subscription where schemecontrol = (select control from scheme where code = 'lifestyle')
										and employercontrol not in (select control from profile where code in ('main','britishmums'))
										group by employercontrol order by 2 desc limit 5`

	mapTop5EmployerSubscriptions, _ := curdb.Query(sqlTop5EmployerSubscriptions)
	aTop5EmployerSubscriptionsSorted := functions.SortMap(mapTop5EmployerSubscriptions)

	for _, sNumber := range aTop5EmployerSubscriptionsSorted {
		xDocReward := mapTop5EmployerSubscriptions[sNumber].(map[string]interface{})
		xDocReward["row"] = sNumber
		sTag := fmt.Sprintf(`%v#report-subscription-topfiveredeemed-row`, sNumber)
		mapReport[sTag] = xDocReward
	}

	//TOP 5 EMPLOYERS WITH THE MOST SUBSCRIPTIONS
	//
	//

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-subscription"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Subscription | Admin Reports","mainpanelContent":` + contentHTML + `}`))

}
