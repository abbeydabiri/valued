package frontend

import (
	"valued/database"
	"valued/functions"

	"net/http"
	"strconv"
)

type Support struct {
	functions.Templates
	pageMap map[string]interface{}
}

func (this *Support) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	switch httpReq.FormValue("action") {
	default:
		fallthrough
	case "":
		this.pageMap = make(map[string]interface{})
		this.pageMap["support"] = make(map[string]interface{})

		contentHTML := strconv.Quote(string(this.Generate(this.pageMap, nil)))
		httpRes.Write([]byte(`{"pageTitle":"Valued Support","pageContent":` + contentHTML + `}`))
		return
	}
}
