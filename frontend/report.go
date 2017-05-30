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

	case "subscription":
		this.subscription(httpRes, httpReq, curdb)
		return

	case "merchant":
		this.merchant(httpRes, httpReq, curdb)
		return

	}
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
