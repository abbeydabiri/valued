package frontend

import (
	"fmt"
	"net/http"
	"valued/database"
)

func GetFrontend(httpReq *http.Request, curdb database.Database, cFrontend string) frontender {

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	mapCache := curdb.GetSession(GOSESSID.Value, "mapCache")

	// mapNavigation := curdb.GetSession(GOSESSID.Value, "mapNavigation")
	// if mapNavigation == nil {
	// 	mapNavigation = make(map[string]interface{})
	// }
	// mapNavigation["from"] = httpReq.URL.Path
	// curdb.SetSession(GOSESSID.Value, "mapNavigation", mapNavigation, false)

	if mapCache["control"] != nil {

		switch cFrontend {
		case "dashboard":
			return new(Dashboard)

		case "profile":
			return new(Profile)

		case "merchant":
			return new(Merchant)

		case "industry":
			return new(Industry)

		case "employer":
			return new(Employer)

		case "member":
			return new(Member)
		case "employee":
			return new(Employee)

		case "reward":
			return new(Reward)

		case "merchantreward":
			return new(MerchantReward)

		case "redemption":
			return new(Redemption)

		case "category":
			return new(Category)

		case "review":
			return new(Review)

		case "reviewcategory":
			return new(ReviewCategory)

		case "scheme":
			return new(Scheme)

		case "store":
			return new(Store)

		case "merchantstore":
			return new(MerchantStore)

		case "subscription":
			return new(Subscription)

		case "referral":
			return new(Referral)

		case "group":
			return new(Group)

		case "media":
			return new(Media)

		case "report":
			return new(Report)

		case "merchantreport":
			return new(MerchantReport)

		case "employerreport":
			return new(EmployerReport)

		case "pendingapproval":
			return new(PendingApproval)

		case "validatecoupon":
			return new(ValidateCoupon)

		case "changepin":
			return new(ChangePin)

		// case "smtp":
		// 	return new(Smtp)

		// case "report":
		// 	return new(Report)

		case "password":
			return new(Password)

		case "user":
			return new(User)

		case "permission":
			return new(Permission)

		case "logout":
			xDoc := make(map[string]interface{})
			xDoc["role"] = mapCache["role"]
			curdb.SetSession(GOSESSID.Value, "mapCache", xDoc, false)
			return new(Admin)
		}
	}

	//Signed in to Mobile App
	mapAppCache := curdb.GetSession(GOSESSID.Value, "mapAppCache")
	if mapAppCache["control"] != nil {

		//Get User Scheme and Status
		mapRes := new(database.Profile).VerifyLogin("member", mapAppCache["username"].(string), mapAppCache["password"].(string), curdb)
		if mapRes["1"] == nil {
			curdb.SetSession(GOSESSID.Value, "mapAppCache", make(map[string]interface{}), false)
			return new(AppLogin)

		} else {
			mapAppCache = mapRes["1"].(map[string]interface{})

			mapSubscription := make(map[string]interface{})
			sqlSubscription := `select sub.schemecontrol as schemecontrol, sub.code as code, sub.control as control, sub.expirydate as expirydate, sch.title as schemetitle
							from subscription as sub, scheme as sch where sub.schemecontrol = sch.control
							AND sub.workflow = 'active' AND sch.workflow = 'active' AND sub.membercontrol = '%s' order by control desc`
			sqlSubscription = fmt.Sprintf(sqlSubscription, mapAppCache["control"])
			xDocresult, _ := curdb.Query(sqlSubscription)

			for _, xDoc := range xDocresult {
				xDoc := xDoc.(map[string]interface{})
				mapSubscription[xDoc["code"].(string)] = xDoc
			}

			if len(mapSubscription) > 0 {
				mapAppCache["subscription"] = mapSubscription
			}
			curdb.SetSession(GOSESSID.Value, "mapAppCache", mapAppCache, false)
		}
		//Get User Scheme and Status

		switch cFrontend {
		case "app-pin":
			return new(AppPin)

		case "app-redeem":
			return new(AppRedeem)

		//Authenticated User Navbar Links
		case "app-tell":
			return new(AppTell)

		case "app-favorite":
			return new(AppFavorite)

		case "app-notifications":
			return new(AppNotification)

		case "app-profile":
			return new(AppProfile)
		//Authenticated User Navbar Links

		case "app-subscribe":
			return new(AppSubscribe)

		case "app-store":
			return new(AppStore)

		case "app-logout":
			xDoc := make(map[string]interface{})
			curdb.SetSession(GOSESSID.Value, "mapAppCache", xDoc, false)
			return new(AppLogin)
		}
	}

	switch cFrontend {
	// case "app-redeem", "app-redeem-qrcode", "app-redeem-pin", "app-reward-rate", "app-reward-redeemed":
	// 	return new(AppRedeem)

	case "support":
		return new(Support)

	case "admin", "employer", "merchant":
		return new(Admin)

	case "app-search":
		return new(AppSearch)

	case "app-activate":
		return new(AppActivate)

	//Default Navbar Links
	case "app-version":
		return new(AppVersion)

	case "app-terms":
		return new(AppTerms)
	//Default Navbar Links

	//redirect to login frontend
	case "app-login", "app-profile", "app-favorite",
		"app-store", "app-redeem", "app-pin":
		return new(AppLogin)
	//redirect to login frontend

	case "app-home":
		return new(AppHome)

	case "app-reward":
		return new(AppReward)

	case "app-merchant":
		return new(AppMerchant)

	case "app-map":
		return new(AppMap)

	//app-tell referral with coupon code
	case "app-tell":
		return new(AppLogin)
	//app-tell referral with coupon code

	case "app-gift":
		return new(AppGift)

	}

	return new(Html)
}
