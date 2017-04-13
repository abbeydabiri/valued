package frontend

import (
	"fmt"
	"valued/functions"
)

type AppNavbar struct{}

func (this *AppNavbar) GetNavBarButton(mapAppCache map[string]interface{}) map[string]interface{} {

	NavBarButton := make(map[string]interface{})
	if mapAppCache["control"] == nil {
		NavBarButton["title"] = "REGISTER"
		NavBarButton["onclick"] = "app-login?p=signup"
	} else {
		if mapAppCache["subscription"] == nil {
			NavBarButton["title"] = "BUY"
			NavBarButton["onclick"] = "app-subscribe"
		} else {

			nPrevDiff := 0
			lRenew := false
			sRenewAction := ""
			for _, xDocSubscription := range mapAppCache["subscription"].(map[string]interface{}) {
				xDocSubscription := xDocSubscription.(map[string]interface{})
				nCurDiff := functions.GetDifferenceInMonths("", xDocSubscription["expirydate"].(string))

				if nCurDiff < 2 {
					lRenew = true
					if nPrevDiff == 0 || nCurDiff < nPrevDiff {
						sRenewAction = "app-subscribe?action=subscribe&scheme=" + xDocSubscription["schemecontrol"].(string)
					}
					nPrevDiff = nCurDiff
				}
			}

			if lRenew {
				NavBarButton["title"] = "RENEW"
				NavBarButton["onclick"] = sRenewAction
			} else {
				NavBarButton["title"] = "GIFT"
				NavBarButton["onclick"] = "app-gift"
			}
		}
	}

	return NavBarButton
}

func (this *AppNavbar) GetNavBar(mapAppCache map[string]interface{}) map[string]interface{} {

	appNavbarItem := make(map[string]interface{})
	appNavbar := make(map[string]interface{})
	if mapAppCache["control"] != nil {

		if mapAppCache["subscription"] == nil {

			appNavbarItem["title"] = "LOGOUT"
			appNavbarItem["onclick"] = "getForm('/app-logout')"
			appNavbar["1#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "BECOME A MEMBER"
			appNavbarItem["onclick"] = "getForm('/app-subscribe')"
			appNavbar["2#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "GIFT A MEMBERSHIP"
			appNavbarItem["onclick"] = "getForm('/app-gift')"
			appNavbar["3#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "TELL A FRIEND"
			appNavbarItem["onclick"] = "getForm('/app-tell')"
			appNavbar["4#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "FAVOURITES"
			appNavbarItem["onclick"] = "getForm('/app-favorite')"
			appNavbar["5#app-navbar-item"] = appNavbarItem

		} else {

			appNavbarItem["title"] = "LOGOUT"
			appNavbarItem["onclick"] = "getForm('/app-logout')"
			appNavbar["1#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "RENEW MEMBERSHIP"

			nPrevDiff := 0
			for _, xDocSubscription := range mapAppCache["subscription"].(map[string]interface{}) {
				xDocSubscription := xDocSubscription.(map[string]interface{})
				nCurDiff := functions.GetDifferenceInMonths("", xDocSubscription["expirydate"].(string))

				if nPrevDiff == 0 || nCurDiff < nPrevDiff {
					appNavbarItem["onclick"] = fmt.Sprintf("getForm('/app-subscribe?action=subscribe&scheme=%s')", xDocSubscription["schemecontrol"])
				}
				nPrevDiff = nCurDiff
			}

			// appNavbarItem["onclick"] = "getForm('/app-subscribe?action=subscribe')"
			appNavbar["2#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "GIFT A MEMBERSHIP"
			appNavbarItem["onclick"] = "getForm('/app-gift')"
			appNavbar["3#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "TELL A FRIEND"
			appNavbarItem["onclick"] = "getForm('/app-tell')"
			appNavbar["4#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "FAVOURITES"
			appNavbarItem["onclick"] = "getForm('/app-favorite')"
			appNavbar["5#app-navbar-item"] = appNavbarItem

			appNavbarItem = make(map[string]interface{})
			appNavbarItem["title"] = "MY SAVINGS"
			appNavbarItem["onclick"] = "getForm('/app-profile?action=redemption')"
			appNavbar["6#app-navbar-item"] = appNavbarItem

		}
	} else {
		appNavbarItem["title"] = "LOGIN"
		appNavbarItem["onclick"] = "getForm('/app-login')"
		appNavbar["1#app-navbar-item"] = appNavbarItem

		appNavbarItem = make(map[string]interface{})
		appNavbarItem["title"] = "GIFT A MEMBERSHIP"
		appNavbarItem["onclick"] = "getForm('/app-gift')"
		appNavbar["2#app-navbar-item"] = appNavbarItem
	}

	return appNavbar
}
