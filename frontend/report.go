package frontend

import (
	"strconv"
	"strings"
	"valued/database"
	"valued/functions"

	"net/http"
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

func (this *Report) calculateNPS(mapNPSFeedback map[string]interface{}) (merchantNPS map[string]float64) {

	sCurMechant := ""
	merchantNPS = make(map[string]float64)
	iNPSTotal, iNPSPositive, iNPSNegative := float64(0), float64(0), float64(0)

	aSorted := functions.SortMap(mapNPSFeedback)
	for _, sNumber := range aSorted {
		xDocFeedback := mapNPSFeedback[sNumber].(map[string]interface{})

		if sCurMechant == "" {
			sCurMechant = xDocFeedback["merchant"].(string)
		}

		if sCurMechant != xDocFeedback["merchant"].(string) {

			iNPSNegativePercentage := float64(0)
			iNPSPositivePercentage := float64(0)
			if iNPSTotal > 0 {
				iNPSPositivePercentage = (iNPSPositive / iNPSTotal) * 100
				iNPSNegativePercentage = (iNPSNegative / iNPSTotal) * 100
			}
			merchantNPS[sCurMechant] = functions.RoundUp(iNPSPositivePercentage-iNPSNegativePercentage, 0)

			iNPSTotal, iNPSPositive, iNPSNegative = float64(0), float64(0), float64(0)
		}

		switch {
		case strings.Contains(xDocFeedback["question"].(string), "RECOMMEND"):
			score, _ := strconv.Atoi(xDocFeedback["answer"].(string))
			switch {
			case score <= 6:
				iNPSNegative++
				break
			case score >= 9:
				iNPSPositive++
				break
			}
			iNPSTotal++
		}

		sCurMechant = xDocFeedback["merchant"].(string)
	}

	iNPSNegativePercentage := float64(0)
	iNPSPositivePercentage := float64(0)
	if iNPSTotal > 0 {
		iNPSPositivePercentage = (iNPSPositive / iNPSTotal) * 100
		iNPSNegativePercentage = (iNPSNegative / iNPSTotal) * 100
	}

	if sCurMechant != "" {
		merchantNPS[sCurMechant] = functions.RoundUp(iNPSPositivePercentage-iNPSNegativePercentage, 0)
	}

	return
}
