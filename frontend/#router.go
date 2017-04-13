package frontend

import (
	"valued/data"
	"valued/database"
	"valued/functions"

	"compress/gzip"
	"fmt"
	"io"

	"net/http"
	"strings"
	"time"
)

const _COOKIE_ = "vlid"

type frontender interface {
	Process(http.ResponseWriter, *http.Request, database.Database)
}

type Router struct {
	database.Database
	apiFrontend frontender
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (this gzipResponseWriter) Write(fileBytes []byte) (int, error) {
	return this.Writer.Write(fileBytes)
}

func (this Router) Route(httpRes http.ResponseWriter, httpReq *http.Request) {

	if this.Database.OS == "android" {
		httpRes.Header().Set("Access-Control-Allow-Origin", "*")
	}

	mbSize := uint(10)
	mbShifter := uint(20)
	maxSize := mbSize << mbShifter
	if httpReq.ContentLength > int64(maxSize) {
		http.Error(httpRes, "request too large", http.StatusExpectationFailed)
		return
	}

	GOSESSID, ERR_GOSESSID := httpReq.Cookie(_COOKIE_)
	if ERR_GOSESSID != nil {
		tExpires := time.Now().Add(14 * 24 * time.Hour)
		goCookie := http.Cookie{Name: _COOKIE_, Value: functions.RandomString(26), Expires: tExpires}
		http.SetCookie(httpRes, &goCookie)
		httpReq.AddCookie(&goCookie)
	}

	if httpReq.URL.Path[1:] == "favicon.ico" {
		return
	}

	frontend := strings.Split(httpReq.URL.Path[1:], "/")

	switch frontend[0] {
	case "robots.txt":
		http.ServeFile(httpRes, httpReq, httpReq.URL.Path[1:])
		return

	case "images":
		if strings.HasSuffix(httpReq.URL.Path[1:], "/") || strings.HasSuffix(httpReq.URL.Path[1:], "main.log") {
			httpRes.Write([]byte("Directory Listing Forbidden!!!"))
			return
		}

		if len(frontend) > 1 {

			if !strings.Contains(httpReq.Header.Get("Accept-Encoding"), "gzip") {
				http.ServeFile(httpRes, httpReq, this.Database.OSfilepath+httpReq.URL.Path[1:])
				return
			}

			httpRes.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(httpRes)
			defer gz.Close()
			gzipRes := gzipResponseWriter{Writer: gz, ResponseWriter: httpRes}
			http.ServeFile(gzipRes, httpReq, this.Database.OSfilepath+httpReq.URL.Path[1:])
		}
		return

	case "docs":
		if ERR_GOSESSID != nil {
			return
		}

		mapCache := this.GetSession(GOSESSID.Value, "mapCache")
		if mapCache["control"] == nil {
			return
		}

		fallthrough
	case "files":
		seconds := 900
		httpRes.Header().Add("Vary", "Accept-Encoding")
		httpRes.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))

		if strings.HasSuffix(httpReq.URL.Path[1:], "/") || strings.HasSuffix(httpReq.URL.Path[1:], "main.log") {
			httpRes.Write([]byte("Directory Listing Forbidden!!!"))
			return
		}

		if len(frontend) > 1 {
			assetFileBytes, _ := data.Asset(httpReq.URL.Path[1:])

			if len(assetFileBytes) == 0 {
				httpRes.Write([]byte(`<b>Forbidden</b>`))
				return
			}

			httpRes.Header().Add("Vary", "Accept-Encoding")

			contentType := ""
			path := httpReq.URL.Path[1:]
			if strings.HasSuffix(path, ".html") {
				contentType = "text/html"
			} else if strings.HasSuffix(path, ".css") {
				contentType = "text/css"
			} else if strings.HasSuffix(path, ".js") {
				contentType = "application/javascript"
			} else if strings.HasSuffix(path, ".png") {
				contentType = "image/png"
			} else if strings.HasSuffix(path, ".jpg") {
				contentType = "image/jpeg"
			} else if strings.HasSuffix(path, ".gif") {
				contentType = "image/gif"
			} else if strings.HasSuffix(path, ".svg") {
				contentType = "image/svg+xml"
			} else if strings.HasSuffix(path, ".mp4") {
				contentType = "video/mp4"
			} else if strings.HasSuffix(path, ".webm") {
				contentType = "video/webm"
			} else if strings.HasSuffix(path, ".ogg") {
				contentType = "video/ogg"
			} else if strings.HasSuffix(path, ".mp3") {
				contentType = "audio/mp3"
			} else if strings.HasSuffix(path, ".wav") {
				contentType = "audio/wav"
			} else {
				contentType = "text/plain"
			}

			httpRes.Header().Add("Content-Type", contentType)

			httpRes.Write(assetFileBytes)
			return

			if !strings.Contains(httpReq.Header.Get("Accept-Encoding"), "gzip") {
				httpRes.Write(assetFileBytes)
				return
			}

			httpRes.Header().Set("Transfer-Encoding", "gzip")
			httpRes.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(httpRes)
			defer gz.Close()
			gzipRes := gzipResponseWriter{Writer: gz, ResponseWriter: httpRes}
			gzipRes.Write(assetFileBytes)

		}
		return
	default:
		if strings.HasSuffix(httpReq.URL.Path[1:], "favicon.ico") ||
			strings.HasSuffix(httpReq.URL.Path[1:], "favicon.png") {
			return
		}
		this.apiFrontend = GetFrontend(httpReq, this.Database, frontend[0])
	}
	this.apiFrontend.Process(httpRes, httpReq, this.Database)
}
