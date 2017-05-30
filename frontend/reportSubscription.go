package frontend

import (
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
	// Number of Subscriptions
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
	// Number of Subscriptions
	//

	//
	// Number of Registrations
	sqlRegistrationTotal := `select count(control) as totalregistered from profile where company != 'Yes' and workflow = 'registered'`
	mapRegistrationTotal, _ := curdb.Query(sqlRegistrationTotal)
	mapReport["registrationslitetotal"] = float64(0)
	mapReport["registrationslifestyletotal"] = float64(0)
	if mapRegistrationTotal["1"] != nil {
		mapRegistrationTotalxDoc := mapRegistrationTotal["1"].(map[string]interface{})
		switch mapRegistrationTotalxDoc["totalregistered"].(type) {
		case string:
			mapReport["registrationslitetotal"] = float64(0)
			mapReport["registrationslifestyletotal"] = float64(0)
		case int64:
			mapReport["registrationslitetotal"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationTotalxDoc["totalregistered"].(int64))))
			mapReport["registrationslifestyletotal"] = functions.ThousandSeperator(functions.Round(float64(mapRegistrationTotalxDoc["totalregistered"].(int64))))
		case float64:
			mapReport["registrationslitetotal"] = functions.ThousandSeperator(functions.Round(mapRegistrationTotalxDoc["totalregistered"].(float64)))
			mapReport["registrationslifestyletotal"] = functions.ThousandSeperator(functions.Round(mapRegistrationTotalxDoc["totalregistered"].(float64)))
		}
	}
	// Number of Registrations
	//

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

	// Total Subscriptions
	//

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-subscription"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Subscription | Admin Reports","mainpanelContent":` + contentHTML + `}`))

}
