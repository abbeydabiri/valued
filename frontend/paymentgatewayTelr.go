// {
// 	"method":"create", "store":17555, "authkey":"Dkwt3@RDJk^kn5jV",
// 	"order":{
// 		"cartid":"103", "test":1, "amount":190.00, "currency":"AED",
// 		"description":"Valued Member Subscription for Scheme %s cost: %d"
// 	},
// 	"customer":{
// 		"email":"stearartois@home.com",
// 		"name":{
// 			"title":"Mrs", "forenames":"Abbey", "surname":"Stella Marie"
// 		},
// 		"address":{	"country":"AE"},
// 		"ref":"testoder1s"
// 	},
// 	"return": {
// 		"authorised": "http://localhost/app-subscribe?action=authorised",
// 		"declined": "http://localhost/app-subscribe?action=declined",
// 		"cancelled": "http://localhost/app-subscribe?action=cancelled"
// 	}
// }

package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"

	"encoding/json"
	"fmt"
	"net/http"
	// "net/url"
	// "strconv"
	"strings"
)

type PaymentGatewayTELR struct{}

func (this *PaymentGatewayTELR) CreateOrder(redirectToPage string, httpReq *http.Request, curdb database.Database) string {

	// if redirectToPage != "app-subscribe" || redirectToPage != "app-gift" {
	// 	return redirectToPage
	// }

	jsonUrl := `https://secure.telr.com/gateway/order.json`
	jsonStr := `{
		"method":"create", "store":17555, "authkey":"Dkwt3@RDJk^kn5jV",
		"order":{
			"cartid":"%s", "test":1, "amount":%.2f, "currency":"AED",
			"description":"Valued Member Subscription for Scheme %s cost: %.2f %s"
		},
		"customer":{
			"email":"%s",
			"name":{
				"title":"%s", "forenames":"%s", "surname":"%s"
			},
			"address":{	"country":"%s"},
			"ref":"%s"
		},
		"return": {
			"authorised": "http://%s/?a=%s&action=telr&telr=%s&status=authorised",
			"declined": "http://%s/?a=%s&action=telr&telr=%s&status=declined",
			"cancelled": "http://%s/?a=%s&action=telr&telr=%s&status=cancelled"
		}
	}`

	sCouponDetails := ""
	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	mapAppCache := curdb.GetSession(GOSESSID.Value, "mapAppCache")
	mapAppTelr := make(map[string]interface{})

	tblTelrOrder := new(database.TelrOrder)
	xDocTelrOrder := make(map[string]interface{})

	sUsername := ""
	if mapAppCache["username"] != nil {
		sUsername = mapAppCache["username"].(string)
	}

	switch redirectToPage {
	case "app-subscribe":
		mapAppTelr = curdb.GetSession(GOSESSID.Value, "mapAppSubscribe")
		xDocTelrOrder["profilecontrol"] = mapAppCache["control"]
		xDocTelrOrder["logincontrol"] = mapAppCache["control"]

		//Create an Inactive Subscription for User before redirecting
		xDocSub := make(map[string]interface{})
		xDocSub["workflow"] = "inactive"
		xDocSub["price"] = mapAppTelr["totalprice"]
		xDocSub["schemecontrol"] = mapAppTelr["schemecontrol"]
		xDocSub["membercontrol"] = mapAppCache["control"]
		xDocSub["employercontrol"] = mapAppCache["employercontrol"]
		xDocTelrOrder["subscriptioncontrol"] = new(database.Subscription).Create(mapAppCache["username"].(string), xDocSub, curdb)
		//Create an Inactive Subscription for User before redirecting

	case "app-gift":
		mapAppTelr = curdb.GetSession(GOSESSID.Value, "mapAppGift")
		sCouponDetails = "GIFT a Membership - "

		xDocTelrOrder["profilecontrol"] = mapAppTelr["profilecontrol"]
		xDocTelrOrder["logincontrol"] = mapAppTelr["logincontrol"]

		sUsername = mapAppTelr["sendersemail"].(string)

	default:
		return ""
	}

	xDocTelrOrder["code"] = redirectToPage

	if mapAppTelr["couponcode"] != nil && mapAppTelr["coupondiscount"] != nil {
		xDocTelrOrder["couponcontrol"] = mapAppTelr["couponcontrol"]
		sCouponDetails = fmt.Sprintf(`Coupon %s (discount: %s) applied`, mapAppTelr["couponcode"], mapAppTelr["coupondiscount"])
	}

	if mapAppTelr["totalprice"].(float64) < mapAppTelr["schemeprice"].(float64) {
		sCouponDetails += fmt.Sprintf(` You saved %.2f`, (mapAppTelr["schemeprice"].(float64) - mapAppTelr["totalprice"].(float64)))
	}

	xDocTelrOrder["schemecontrol"] = mapAppTelr["schemecontrol"]
	xDocTelrOrder["amount"] = mapAppTelr["totalprice"]
	xDocTelrOrder["currency"] = "AED"
	sTelrControl := tblTelrOrder.Create(sUsername, xDocTelrOrder, curdb)
	// sTelrControlEncrypted := base64.StdEncoding.EncodeToString([]byte(sTelrControl))

	sTelrCartId := sTelrControl + "-" + functions.RandomString(6)
	sTelrControlEncrypted := base64.StdEncoding.EncodeToString([]byte(sTelrCartId))

	jsonStrRequest := ""
	switch redirectToPage {
	case "app-subscribe":
		jsonStrRequest = fmt.Sprintf(jsonStr, sTelrControlEncrypted, mapAppTelr["totalprice"], mapAppTelr["schemetitle"], mapAppTelr["schemeprice"], sCouponDetails,
			mapAppTelr["profileemail"], mapAppTelr["title"], mapAppTelr["firstname"], mapAppTelr["lastname"],
			"AE", mapAppTelr["profileemail"],
			httpReq.Host, redirectToPage, sTelrControlEncrypted,
			httpReq.Host, redirectToPage, sTelrControlEncrypted,
			httpReq.Host, redirectToPage, sTelrControlEncrypted)

	case "app-gift":
		jsonStrRequest = fmt.Sprintf(jsonStr, sTelrControlEncrypted, mapAppTelr["totalprice"], mapAppTelr["schemetitle"], mapAppTelr["schemeprice"], sCouponDetails,
			mapAppTelr["sendersemail"], "", mapAppTelr["sendersname"], "",
			"AE", mapAppTelr["profileemail"],
			httpReq.Host, redirectToPage, sTelrControlEncrypted,
			httpReq.Host, redirectToPage, sTelrControlEncrypted,
			httpReq.Host, redirectToPage, sTelrControlEncrypted)
	}

	// AE = mapAppTelr["countrycode"]

	xDocTelrOrder = make(map[string]interface{})
	xDocTelrOrder["control"] = sTelrControl

	jsonStrRequestCleaned := functions.TrimEscape(jsonStrRequest)
	jsonStrRequestCleaned = strings.Replace(jsonStrRequestCleaned, "\r", "", -1)
	xDocTelrOrder["createrequest"] = strings.Replace(jsonStrRequestCleaned, "\n", "<br>", -1)

	jsonByteResult := functions.HttpPostJSON(jsonUrl, []byte(jsonStrRequest))

	jsonByteResultCleaned := functions.TrimEscape(string(jsonByteResult))
	jsonByteResultCleaned = strings.Replace(jsonByteResultCleaned, "\r", "", -1)
	xDocTelrOrder["createresponse"] = strings.Replace(jsonByteResultCleaned, "\n", "<br>", -1)

	mapJsonResult := make(map[string]interface{})
	json.Unmarshal(jsonByteResult, &mapJsonResult)

	if mapJsonResult["trace"] != nil {
		xDocTelrOrder["workflow"] = "pending"

		mapOrderInterface := mapJsonResult["order"]
		mapOrder := mapOrderInterface.(map[string]interface{})
		xDocTelrOrder["telrurl"] = mapOrder["url"]
		xDocTelrOrder["telrref"] = mapOrder["ref"]
	} else {
		xDocTelrOrder["workflow"] = "failed"
	}

	tblTelrOrder.Update(sUsername, xDocTelrOrder, curdb)

	if xDocTelrOrder["telrurl"] != nil {
		return xDocTelrOrder["telrurl"].(string)
	}

	return ""
}

func (this *PaymentGatewayTELR) VerifyPayment(sTelrControlEncrypted string, httpReq *http.Request, curdb database.Database) (telrRecord map[string]interface{}) {

	telrRecord = make(map[string]interface{})
	jsonUrl := `https://secure.telr.com/gateway/order.json`
	jsonStr := `{
		"method":"check", "store":17555, "authkey":"Dkwt3@RDJk^kn5jV",
		"order":{
			"ref":"%s"
		}
	}`

	/*
		checkrequest
		{
			"method":"check", "store":17555, "authkey":"Dkwt3@RDJk^kn5jV",
			"order":{
				"ref":"%s"
			}
		}


		checkresponse
		{
			"method": "check",
			"trace": "4000/3475/58492b1d",
			"order": {
			"ref": "734313B52DB69921EEBDA32732289E76AD7A9737CAE563E8652937B7FF072225",
			"cartid": "1.00000000002",
			"test": 1,
			"amount": "530.23",
			"currency": "AED",
			"description": "Valued Member Subscription for Scheme VALUED LIFESTYLE cost:...",
				"status": {
					"code": -2,
					"text": "Cancelled"
					}
			}
		}
	*/

	if sTelrControlEncrypted == "" {
		telrRecord["error"] = "No CartID received from Telr"
		return
	}

	sTelrControlbase64Bytes, err := base64.StdEncoding.DecodeString(sTelrControlEncrypted)
	if err != nil {
		telrRecord["error"] = "Telr Record not found"
		return
	}
	sTelrCartId := string(sTelrControlbase64Bytes)
	sTelrControl := strings.Split(sTelrCartId, "-")[0]

	tblTelrOrder := new(database.TelrOrder)
	xDocRequest := make(map[string]interface{})
	xDocRequest["workflow"] = "pending"
	xDocRequest["control"] = sTelrControl

	xDcocResult := tblTelrOrder.Search(xDocRequest, curdb)
	if xDcocResult["1"] == nil {
		telrRecord["error"] = "TelrOrder ID is Invalid "
		return
	}

	//Lets Get the TelrOrder Details for Verification

	xDocTelrOrder := xDcocResult["1"].(map[string]interface{})
	jsonStrRequest := fmt.Sprintf(jsonStr, xDocTelrOrder["telrref"])

	jsonStrRequestCleaned := functions.TrimEscape(jsonStrRequest)
	jsonStrRequestCleaned = strings.Replace(jsonStrRequestCleaned, "\r", "", -1)
	xDocTelrOrder["checkrequest"] = strings.Replace(jsonStrRequestCleaned, "\n", "<br>", -1)

	jsonByteResult := functions.HttpPostJSON(jsonUrl, []byte(jsonStrRequest))

	jsonByteResultCleaned := functions.TrimEscape(string(jsonByteResult))
	jsonByteResultCleaned = strings.Replace(jsonByteResultCleaned, "\r", "", -1)
	xDocTelrOrder["checkresponse"] = strings.Replace(jsonByteResultCleaned, "\n", "<br>", -1)

	mapJsonResult := make(map[string]interface{})
	json.Unmarshal(jsonByteResult, &mapJsonResult)

	sMessage := ""
	if mapJsonResult["order"] != nil {
		mapOrderInterface := mapJsonResult["order"]
		mapOrder := mapOrderInterface.(map[string]interface{})

		//validate amount, currency,
		if mapOrder["currency"].(string) != xDocTelrOrder["currency"].(string) {
			sMessage += fmt.Sprintf("Currency Mismatch!! expected <b>%s</b>, recevied <b>%s</b> - Possible Fraud Attempt <br>",
				xDocTelrOrder["currency"], mapOrder["currency"])
		}

		sAmount := fmt.Sprintf("%.2f", xDocTelrOrder["amount"])
		if mapOrder["amount"].(string) != sAmount {
			sMessage += fmt.Sprintf("Amount Mismatch!! expected <b>%s</b>, recevied <b>%s</b> - Possible Fraud Attempt <br>",
				sAmount, mapOrder["amount"])
		}

		//retrieve status [] code, text
		mapOrderStatusInterface := mapOrder["status"]
		mapOrderStatus := mapOrderStatusInterface.(map[string]interface{})

		xDocTelrOrder["telrcode"] = mapOrderStatus["code"]
		xDocTelrOrder["telrtext"] = mapOrderStatus["text"].(string)
		xDocTelrOrder["workflow"] = strings.ToLower(mapOrderStatus["text"].(string))

	} else {
		xDocTelrOrder["workflow"] = "failed"
		if mapJsonResult["error"] != nil {
			mapErrorInterface := mapJsonResult["error"]
			mapError := mapErrorInterface.(map[string]interface{})
			sMessage += mapError["note"].(string)
		} else {
			sMessage += "Transaction Validation Failed"
		}
	}

	delete(xDocTelrOrder, "schemetitle")
	xDocTelrOrder["description"] = sMessage
	tblTelrOrder.Update(xDocTelrOrder["createdby"].(string), xDocTelrOrder, curdb)

	if len(sMessage) > 0 {
		telrRecord["error"] = sMessage
		return
	}

	telrRecord = xDocTelrOrder

	/*
		checkresponse
		{
			"method": "check",
			"trace": "4000/3475/58492b1d",
			"order": {
			"ref": "734313B52DB69921EEBDA32732289E76AD7A9737CAE563E8652937B7FF072225",
			"cartid": "1.00000000002",
			"test": 1,
			"amount": "530.23",
			"currency": "AED",
			"description": "Valued Member Subscription for Scheme VALUED LIFESTYLE cost:...",
			"status": {
					"code": -2,
					"text": "Cancelled"
				}
			}
		}
	*/

	return
}
