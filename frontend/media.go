package frontend

import (
	"valued/database"
	"valued/functions"

	"encoding/base64"
	"fmt"
	"html"

	"net/http"
	"strconv"
	"strings"
)

type Media struct {
	functions.Templates
	mapCache map[string]interface{}
}

func (this *Media) Process(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	httpRes.Header().Set("content-type", "application/json")

	if httpReq.Method != "POST" {
		http.Redirect(httpRes, httpReq, "/", http.StatusMovedPermanently)
	}

	GOSESSID, _ := httpReq.Cookie(_COOKIE_)
	this.mapCache = curdb.GetSession(GOSESSID.Value, "mapCache")

	switch httpReq.FormValue("action") {
	case "search":
		searchResult := this.search(httpReq, curdb)
		contentHTML := strconv.Quote(string(this.Generate(searchResult, nil)))
		httpRes.Write([]byte(`{"medialist":` + contentHTML + `}`))
		return
	case "upload":
		this.upload(httpRes, httpReq, curdb)
		return
	}

	httpRes.Write([]byte(`{}`))
	return
}

func (this *Media) upload(httpRes http.ResponseWriter, httpReq *http.Request, curdb database.Database) {

	sMessage := ""
	if html.EscapeString(httpReq.FormValue("filename")) == "" {
		sMessage += "File Name is missing <br>"
	}

	if html.EscapeString(httpReq.FormValue("encoded")) == "" {
		sMessage += "Image Data is missing"
	}

	if len(sMessage) > 0 {
		httpRes.Write([]byte(`{"error":"` + sMessage + `"}`))
		return
	}

	tblMedia := new(database.Media)
	xDoc := make(map[string]interface{})
	xDoc["title"] = html.EscapeString(httpReq.FormValue("filename"))

	base64String := httpReq.FormValue("encoded")
	base64String = strings.Split(base64String, "base64,")[1]
	base64Bytes, err := base64.StdEncoding.DecodeString(base64String)

	if err != nil {
		println(err.Error())
		httpRes.Write([]byte(`{}`))
		return
	}

	if base64Bytes != nil {
		fileName := fmt.Sprintf("media_%s_%s", functions.RandomString(6), xDoc["title"].(string))
		xDoc["code"] = functions.SaveImage(fileName, curdb.OSfilepath, base64Bytes)
		tblMedia.Create(this.mapCache["username"].(string), xDoc, curdb)
	}

	httpRes.Write([]byte(`{}`))
	return
}

func (this *Media) search(httpReq *http.Request, curdb database.Database) map[string]interface{} {
	formSearch := make(map[string]interface{})

	tblMedia := new(database.Media)
	xDocrequest := make(map[string]interface{})

	xDocrequest["offset"] = html.EscapeString(httpReq.FormValue("offset"))
	xDocrequest["title"] = html.EscapeString(httpReq.FormValue("title"))

	xDocresult := tblMedia.Search(xDocrequest, curdb)

	for cNumber, xDoc := range xDocresult {
		xDoc := xDoc.(map[string]interface{})
		xDoc["number"] = cNumber
		formSearch[cNumber+"#newseditor-media"] = xDoc
	}

	return formSearch
}
