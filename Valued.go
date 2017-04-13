// package valued
package main

import (
	"flag"
	"valued/data"
	"valued/frontend"
	"valued/functions"
)

var server frontend.Server

func AssetRequest() string {
	return <-data.AssetRequest
}

func AssetResponse(response []byte) {
	data.AssetResponse <- response
}

func Start(Port string) {
	server.OS = "android"
	data.OS = server.OS
	data.AssetRequest = make(chan string)
	data.AssetResponse = make(chan []byte)

	go server.Start(Port, "", false, false, false)
}

func main() {
	var cPort string
	flag.StringVar(&cPort, "port", "80", "This is the Default Port")

	var cRedirect string
	flag.StringVar(&cRedirect, "redirect", "", "Redirect this Port to Default Port if default port not specified it redirects to port 443")

	var lSSL bool
	flag.BoolVar(&lSSL, "ssl", false, "Use HTTPS Protocol on Default Port")

	var lPGSQL bool
	flag.BoolVar(&lPGSQL, "pgsql", true, "Use PostgreSQL or github.com/cznic/ql")

	var lInit bool
	flag.BoolVar(&lInit, "init", false, "Initialize Database")

	var lLog bool
	flag.BoolVar(&lLog, "log", false, "Enable Logging")

	flag.Parse()

	if lLog {
		functions.LoadLogFile()
	}

	server.Start(cPort, cRedirect, lSSL, lPGSQL, lInit)
}
