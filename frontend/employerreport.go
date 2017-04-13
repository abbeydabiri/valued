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

type EmployerReport struct {
	EmployerControl string
	functions.Templates
	mapCache map[string]interface{}
	pageMap  map[string]interface{}
}

func (this *EmployerReport) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	this.EmployerControl = this.mapCache["control"].(string)
	if this.mapCache["company"].(string) != "Yes" {
		this.EmployerControl = this.mapCache["employercontrol"].(string)
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

	case "users":
		this.users(httpRes, httpReq, curdb)
		return

	case "rewards":
		this.rewards(httpRes, httpReq, curdb)
		return

	}
}

func (this *EmployerReport) downloadReport(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	cFormatDate := "02/01/2006"
	todayDate, _ := time.Parse(cFormatDate, functions.GetSystemDate())
	oneYear := todayDate.Add(-(time.Hour * 24 * 365))

	sStartdate := oneYear.Format(cFormatDate)
	sStopdate := functions.GetSystemDate()

	if functions.TrimEscape(httpReq.FormValue("startdate")) != "" {
		sStartdate = functions.TrimEscape(httpReq.FormValue("startdate"))
	}

	if functions.TrimEscape(httpReq.FormValue("stopdate")) != "" {
		sStopdate = functions.TrimEscape(httpReq.FormValue("stopdate"))
	}

	sqlFeedback := `select feedback.title as question, feedback.answer as answer, feedback.redemptioncontrol as redemptioncontrol
						from redemption 
							left join feedback on feedback.redemptioncontrol = redemption.control
						where redemption.employercontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp 
					`
	sqlFeedback = fmt.Sprintf(sqlFeedback, this.EmployerControl, sStartdate, sStopdate)

	/*
		sqlRedemption := `select redemption.control as control, redemption.createdate as date, store.address as location, redemption.transactionvalue as revenue, redemption.schemecontrol as schemecontrol,
							reward.discount as discount, reward.title as reward, member.dob as age, member.title as gender, member.nationality as nationality, coupon.code as coupon
						from redemption
								left join store on store.control = redemption.storecontrol
								left join reward on reward.control = redemption.rewardcontrol
								left join coupon on coupon.control = redemption.couponcontrol
								left join profile as member on member.control = redemption.membercontrol
						where redemption.employercontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp order by substring(redemption.createdate from 1 for 20)::timestamp desc
						`
	*/

	sqlRedemption := `select redemption.control as control, redemption.createdate as date, redemption.location as location, redemption.transactionvalue as revenue, redemption.schemecontrol as schemecontrol,
						redemption.discount as discount, redemption.reward as reward, redemption.dob as age, redemption.gender as gender, redemption.nationality as nationality, redemption.code as coupon
					from redemption where redemption.merchantcontrol = '%s' and substring(redemption.createdate from 1 for 20)::timestamp between '%s'::timestamp and '%s 23:59:59'::timestamp order by substring(redemption.createdate from 1 for 20)::timestamp desc
					`
	sqlRedemption = fmt.Sprintf(sqlRedemption, this.EmployerControl, sStartdate, sStopdate)

	curdb.Query("set datestyle = dmy")
	mapFeedback, _ := curdb.Query(sqlFeedback)
	mapRedemption, _ := curdb.Query(sqlRedemption)

	feedbackReportGenerator := make(map[string]interface{})
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
		feedbackReportGenerator[sKey] = xDoc["answer"]
	}

	sLine := `"%v","%v","%v","%v","%v","%v","%v","%v","%v","%v","%v"`
	redemptionReportGenerator := make([]string, len(mapRedemption)+1)
	redemptionReportGenerator[0] = fmt.Sprintf(sLine, "date", "discount", "reward", "coupon",
		"revenue", "nationality", "gender", "age", "improve", "rating", "location")

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

		if feedbackReportGenerator[sKey] != nil {
			xDoc["improve"] = feedbackReportGenerator[sKey]
		} else {
			xDoc["improve"] = ""
		}

		sKey = fmt.Sprintf("%s-RECOMMEND", xDoc["control"])
		if feedbackReportGenerator[sKey] != nil {
			xDoc["rating"] = feedbackReportGenerator[sKey]
		} else {
			xDoc["rating"] = ""
		}

		xDoc["age"] = functions.GetDifferenceInYears("", xDoc["age"].(string))

		redemptionReportGenerator[iNumber] = fmt.Sprintf(sLine, xDoc["date"], xDoc["discount"], xDoc["reward"], xDoc["coupon"],
			xDoc["revenue"], xDoc["nationality"], xDoc["gender"], xDoc["age"], xDoc["improve"], xDoc["rating"], xDoc["location"])
	}

	sFilename := fmt.Sprintf("Valued-Report-%s-%s.csv", sStartdate, sStopdate)
	httpRes.Header().Set("Content-Type", "text/csv")
	httpRes.Header().Set("Content-Disposition", "attachment;filename="+sFilename)
	httpRes.Write([]byte(strings.Join(redemptionReportGenerator, "\r\n")))
}

func (this *EmployerReport) summary(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	mapReport := make(map[string]interface{})

	//Total Redeemed Year and 30 Days
	sqlRedeemedTotal := `select count(control) as redeemed 
						from redemption where employercontrol = '%s'
						and substring(createdate from 1 for 20)::timestamp between
						'%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp
					`
	cFormat := "02/01/2006"
	todayDate, _ := time.Parse(cFormat, functions.GetSystemDate())
	oneYear := todayDate.Add(-(time.Hour * 24 * 365))
	oneMonth := todayDate.Add(-(time.Hour * 24 * 30))

	sqlRedeemedYear := fmt.Sprintf(sqlRedeemedTotal, this.EmployerControl, oneYear.Format(cFormat), todayDate.Format(cFormat))
	sqlRedeemedMonth := fmt.Sprintf(sqlRedeemedTotal, this.EmployerControl, oneMonth.Format(cFormat), todayDate.Format(cFormat))

	curdb.Query("set datestyle = dmy")
	mapRedeemedYear, _ := curdb.Query(sqlRedeemedYear)
	mapRedeemedMonth, _ := curdb.Query(sqlRedeemedMonth)

	mapReport["redeemedyear"] = float64(0)
	if mapRedeemedYear["1"] != nil {
		mapRedeemed := mapRedeemedYear["1"].(map[string]interface{})
		switch mapRedeemed["redeemed"].(type) {
		case string:
			mapReport["redeemedyear"] = float64(0)
		case int64:
			mapReport["redeemedyear"] = functions.ThousandSeperator(functions.Round(float64(mapRedeemed["redeemed"].(int64))))
		case float64:
			mapReport["redeemedyear"] = functions.ThousandSeperator(functions.Round(mapRedeemed["redeemed"].(float64)))
		}
	}

	mapReport["redeemedmonth"] = float64(0)
	if mapRedeemedMonth["1"] != nil {
		mapRedeemed := mapRedeemedMonth["1"].(map[string]interface{})
		switch mapRedeemed["redeemed"].(type) {
		case string:
			mapReport["redeemedmonth"] = float64(0)
		case int64:
			mapReport["redeemedmonth"] = functions.ThousandSeperator(functions.Round(float64(mapRedeemed["redeemed"].(int64))))
		case float64:
			mapReport["redeemedmonth"] = functions.ThousandSeperator(functions.Round(mapRedeemed["redeemed"].(float64)))
		}
	} //Total Redeemed Year and 30 Days

	//Total Employees Year and 30 Days
	sqlEmployeesTotal := `select count(distinct(membercontrol)) as employees 
						from redemption where employercontrol = '%s'
						and substring(createdate from 1 for 20)::timestamp between
						'%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp
					`

	sqlEmployeesYear := fmt.Sprintf(sqlEmployeesTotal, this.EmployerControl, oneYear.Format(cFormat), todayDate.Format(cFormat))
	sqlEmployeesMonth := fmt.Sprintf(sqlEmployeesTotal, this.EmployerControl, oneMonth.Format(cFormat), todayDate.Format(cFormat))

	mapEmployeesYear, _ := curdb.Query(sqlEmployeesYear)
	mapEmployeesMonth, _ := curdb.Query(sqlEmployeesMonth)

	mapReport["employeesyear"] = float64(0)
	if mapEmployeesYear["1"] != nil {
		mapEmployees := mapEmployeesYear["1"].(map[string]interface{})
		switch mapEmployees["employees"].(type) {
		case string:
			mapReport["employeesyear"] = float64(0)
		case int64:
			mapReport["employeesyear"] = functions.ThousandSeperator(functions.Round(float64(mapEmployees["employees"].(int64))))
		case float64:
			mapReport["employeesyear"] = functions.ThousandSeperator(functions.Round(mapEmployees["employees"].(float64)))
		}
	}

	mapReport["employeesmonth"] = float64(0)
	if mapEmployeesMonth["1"] != nil {
		mapEmployees := mapEmployeesMonth["1"].(map[string]interface{})
		switch mapEmployees["employees"].(type) {
		case string:
			mapReport["employeesmonth"] = float64(0)
		case int64:
			mapReport["employeesmonth"] = functions.ThousandSeperator(functions.Round(float64(mapEmployees["employees"].(int64))))
		case float64:
			mapReport["employeesmonth"] = functions.ThousandSeperator(functions.Round(mapEmployees["employees"].(float64)))
		}
	}
	//Total Employees Year and 30 Days

	//Total Savings Year and 30 Days
	sqlSavingsTotal := `select sum(savingsvalue) as savings 
						from redemption where employercontrol = '%s'
						and substring(createdate from 1 for 20)::timestamp between
						'%s 00:00:00'::timestamp and '%s 23:59:59'::timestamp
					`

	sqlSavingsYear := fmt.Sprintf(sqlSavingsTotal, this.EmployerControl, oneYear.Format(cFormat), todayDate.Format(cFormat))
	sqlSavingsMonth := fmt.Sprintf(sqlSavingsTotal, this.EmployerControl, oneMonth.Format(cFormat), todayDate.Format(cFormat))

	mapSavingsYear, _ := curdb.Query(sqlSavingsYear)
	mapSavingsMonth, _ := curdb.Query(sqlSavingsMonth)

	mapReport["savingsyear"] = float64(0)
	if mapSavingsYear["1"] != nil {
		mapSavings := mapSavingsYear["1"].(map[string]interface{})
		switch mapSavings["savings"].(type) {
		case string:
			mapReport["savingsyear"] = float64(0)
		case int64:
			mapReport["savingsyear"] = functions.ThousandSeperator(functions.Round(float64(mapSavings["savings"].(int64))))
		case float64:
			mapReport["savingsyear"] = functions.ThousandSeperator(functions.Round(mapSavings["savings"].(float64)))
		}
	}

	mapReport["savingsmonth"] = float64(0)
	if mapSavingsMonth["1"] != nil {
		mapSavings := mapSavingsMonth["1"].(map[string]interface{})
		switch mapSavings["savings"].(type) {
		case string:
			mapReport["savingsmonth"] = float64(0)
		case int64:
			mapReport["savingsmonth"] = functions.ThousandSeperator(functions.Round(float64(mapSavings["savings"].(int64))))
		case float64:
			mapReport["savingsmonth"] = functions.ThousandSeperator(functions.Round(mapSavings["savings"].(float64)))
		}
	}
	//Total Savings Year and 30 Days

	cFormatDate := "02/01/2006"
	sStartdate := oneYear.Format(cFormatDate)
	sStopdate := todayDate.Format(cFormatDate)
	mapReport["startdate"] = sStartdate
	mapReport["stopdate"] = sStopdate

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-employer-summary"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Summary | Employer Reports","mainpanelContent":` + contentHTML + `}`))

}

func (this *EmployerReport) users(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	mapReport := make(map[string]interface{})

	//Average Savings
	sqlSavingsAverage := fmt.Sprintf(`select sum(savingsvalue)/count(control) as savingsaverage 
						from redemption where employercontrol = '%s'`, this.EmployerControl)

	mapSavingsAverage, _ := curdb.Query(sqlSavingsAverage)

	mapReport["savingsaverage"] = float64(0)
	if mapSavingsAverage["1"] != nil {
		mapSavings := mapSavingsAverage["1"].(map[string]interface{})
		switch mapSavings["savingsaverage"].(type) {
		case string:
			mapReport["savingsaverage"] = float64(0)
		case int64:
			mapReport["savingsaverage"] = functions.ThousandSeperator(functions.Round(float64(mapSavings["savingsaverage"].(int64))))
		case float64:
			mapReport["savingsaverage"] = functions.ThousandSeperator(functions.Round(mapSavings["savingsaverage"].(float64)))
		}
	}
	//Average Savings

	//Active Employees
	sqlActiveEmployees := fmt.Sprintf(`select count(distinct(membercontrol)) as activeemployees 
						from redemption where employercontrol = '%s'`, this.EmployerControl)

	mapRedeemedTotal, _ := curdb.Query(sqlActiveEmployees)

	mapReport["activeemployees"] = float64(0)
	if mapRedeemedTotal["1"] != nil {
		mapRedeemed := mapRedeemedTotal["1"].(map[string]interface{})
		switch mapRedeemed["activeemployees"].(type) {
		case string:
			mapReport["activeemployees"] = float64(0)
		case int64:
			mapReport["activeemployees"] = functions.ThousandSeperator(functions.Round(float64(mapRedeemed["activeemployees"].(int64))))
		case float64:
			mapReport["activeemployees"] = functions.ThousandSeperator(functions.Round(mapRedeemed["activeemployees"].(float64)))
		}
	}
	//Active Employees

	//Scheme Piechart
	sqlScheme := fmt.Sprintf(`select scheme.title as label, count(subscription.control) as series from scheme left join subscription
		on subscription.schemecontrol = scheme.control where subscription.employercontrol = '%s' group by 1`, this.EmployerControl)
	resScheme, _ := curdb.Query(sqlScheme)

	sLabel := ""
	iSeries := int64(0)
	schemePieTotal := int64(0)
	mapScheme := make(map[string]int64)
	for _, xDoc := range resScheme {
		xDoc := xDoc.(map[string]interface{})
		iSeries = xDoc["series"].(int64)
		sLabel = xDoc["label"].(string)
		schemePieTotal += int64(iSeries)
		mapScheme[sLabel] = iSeries
	}

	mapPieChartScheme := make(map[string]interface{})
	mapPieChartScheme["id"] = "scheme"
	mapPieChartScheme["title"] = "Scheme"
	mapPieChartScheme["label"] = ""
	mapPieChartScheme["series"] = ""

	var mapLegendScheme []string
	if len(mapScheme) == 0 {
		mapPieChartScheme["label"] = "'No Records'"
		mapPieChartScheme["series"] = "100"
		mapPieChartScheme["total"] = "100"
	} else {
		for sLabel, iSeries := range mapScheme {
			mapLegendScheme = append(mapLegendScheme, sLabel)
			iSeriesPercentage := functions.Round(float64(iSeries) / float64(schemePieTotal) * 100)
			mapPieChartScheme["label"] = fmt.Sprintf(`%v'%v%%',`, mapPieChartScheme["label"], iSeriesPercentage)
			mapPieChartScheme["series"] = fmt.Sprintf(`%v%v,`, mapPieChartScheme["series"], iSeriesPercentage)
		}
	}

	for iNumber, sLabel := range mapLegendScheme {

		iSeries := mapScheme[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(schemePieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-generator-piechart-row`, iNumber)
		mapPieChartScheme[sTag] = pieChartRow

	}

	mapReport["1#report-employer-users-piechart"] = mapPieChartScheme
	//Scheme Piechart

	//User Piechart

	sqlUser := fmt.Sprintf(`select workflow as label, count(control) as series
						from profile where employercontrol = '%s' and company != 'Yes' group by 1`, this.EmployerControl)
	resUser, _ := curdb.Query(sqlUser)

	sLabel = ""
	iSeries = int64(0)
	userPieTotal := int64(0)
	mapUser := make(map[string]int64)
	for _, xDoc := range resUser {
		xDoc := xDoc.(map[string]interface{})
		iSeries = xDoc["series"].(int64)
		sLabel = xDoc["label"].(string)
		userPieTotal += int64(iSeries)
		mapUser[sLabel] = iSeries
	}

	mapPieChartUser := make(map[string]interface{})
	mapPieChartUser["id"] = "user"
	mapPieChartUser["title"] = "User"
	mapPieChartUser["label"] = ""
	mapPieChartUser["series"] = ""

	var mapLegendUser []string
	if len(mapUser) == 0 {
		mapPieChartUser["label"] = "'No Users'"
		mapPieChartUser["series"] = "100"
		mapPieChartUser["total"] = "100"
	} else {
		for sLabel, iSeries := range mapUser {
			mapLegendUser = append(mapLegendUser, sLabel)
			iSeriesPercentage := functions.Round(float64(iSeries) / float64(userPieTotal) * 100)
			mapPieChartUser["label"] = fmt.Sprintf(`%v'%v%%',`, mapPieChartUser["label"], iSeriesPercentage)
			mapPieChartUser["series"] = fmt.Sprintf(`%v%v,`, mapPieChartUser["series"], iSeriesPercentage)
		}
	}

	for iNumber, sLabel := range mapLegendUser {

		iSeries := mapUser[sLabel]
		iSeriesPercentage := functions.Round(float64(iSeries) / float64(userPieTotal) * 100)

		pieChartRow := make(map[string]interface{})
		pieChartRow["label"] = sLabel
		pieChartRow["value"] = iSeries
		pieChartRow["percentage"] = fmt.Sprintf(`%v%%`, iSeriesPercentage)

		sTag := fmt.Sprintf(`%v#report-generator-piechart-row`, iNumber)
		mapPieChartUser[sTag] = pieChartRow

	}

	mapReport["2#report-employer-users-piechart"] = mapPieChartUser
	//User Piechart

	cFormatDate := "02/01/2006"
	todayDate, _ := time.Parse(cFormatDate, functions.GetSystemDate())
	sStartdate := todayDate.Add(-(time.Hour * 24 * 365)).Format(cFormatDate)
	sStopdate := todayDate.Format(cFormatDate)

	mapReport["startdate"] = sStartdate
	mapReport["stopdate"] = sStopdate

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-employer-users"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Users | Employer Reports","mainpanelContent":` + contentHTML + `}`))

}

func (this *EmployerReport) rewards(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	mapReport := make(map[string]interface{})

	//Total Redeemed
	sqlRedeemedTotal := fmt.Sprintf(`select count(control) as redeemed  from redemption where employercontrol = '%s'`, this.EmployerControl)

	curdb.Query("set datestyle = dmy")
	mapRedeemedTotal, _ := curdb.Query(sqlRedeemedTotal)

	mapReport["redeemedtotal"] = float64(0)
	if mapRedeemedTotal["1"] != nil {
		mapRedeemed := mapRedeemedTotal["1"].(map[string]interface{})
		switch mapRedeemed["redeemed"].(type) {
		case string:
			mapReport["redeemedtotal"] = float64(0)
		case int64:
			mapReport["redeemedtotal"] = functions.ThousandSeperator(functions.Round(float64(mapRedeemed["redeemed"].(int64))))
		case float64:
			mapReport["redeemedtotal"] = functions.ThousandSeperator(functions.Round(mapRedeemed["redeemed"].(float64)))
		}
	}
	//Total Redeemed

	//Average Savings
	sqlSavingsAverage := fmt.Sprintf(`select sum(savingsvalue)/count(control) as savingsaverage 
						from redemption where employercontrol = '%s'`, this.EmployerControl)

	mapSavingsAverage, _ := curdb.Query(sqlSavingsAverage)

	mapReport["savingsaverage"] = float64(0)
	if mapSavingsAverage["1"] != nil {
		mapSavings := mapSavingsAverage["1"].(map[string]interface{})
		switch mapSavings["savingsaverage"].(type) {
		case string:
			mapReport["savingsaverage"] = float64(0)
		case int64:
			mapReport["savingsaverage"] = functions.ThousandSeperator(functions.Round(float64(mapSavings["savingsaverage"].(int64))))
		case float64:
			mapReport["savingsaverage"] = functions.ThousandSeperator(functions.Round(mapSavings["savingsaverage"].(float64)))
		}
	}
	//Average Savings

	//Popular Rewards
	sqlPopularRewards := fmt.Sprintf(`select merchant.title as merchant, reward.title as reward, reward.discount as discount, 
				count(redemption.control) as redemption from redemption 
				left join reward on reward.control = redemption.rewardcontrol
				left join profile as merchant on merchant.control = redemption.merchantcontrol
				where redemption.employercontrol = '%s'
				group by merchant.title, reward.title, reward.discount
				order by redemption desc limit 5 offset 0`, this.EmployerControl)

	mapPopularRewards, _ := curdb.Query(sqlPopularRewards)
	for cNumber, xDoc := range mapPopularRewards {
		xDoc := xDoc.(map[string]interface{})
		mapReport[cNumber+"#report-employer-rewards-popular"] = xDoc
	}
	//Popular Rewards

	cFormatDate := "02/01/2006"
	todayDate, _ := time.Parse(cFormatDate, functions.GetSystemDate())
	sStartdate := todayDate.Add(-(time.Hour * 24 * 365)).Format(cFormatDate)
	sStopdate := todayDate.Format(cFormatDate)

	mapReport["startdate"] = sStartdate
	mapReport["stopdate"] = sStopdate

	this.pageMap = make(map[string]interface{})
	this.pageMap["report-employer-rewards"] = mapReport
	contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
	httpRes.Write([]byte(`{"pageTitle":"Rewards | Employer Reports","mainpanelContent":` + contentHTML + `}`))

}
