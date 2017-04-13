package frontend

import (
	"valued/data"
	"valued/database"

	"crypto/tls"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	KeyFile  = "conf/key"
	CertFile = "conf/cert"
)

type Server struct {
	lSSL                 bool
	OS, cPort, cRedirect string
	KillSwitch           chan bool
	AndroidPipe          chan string
	router               *Router
}

func (this Server) Start(cPort, cRedirect string, lSSL, lPGSQL, lInit bool) {
	this.lSSL = lSSL
	this.cPort = cPort
	this.cRedirect = cRedirect
	this.router = new(Router)

	this.router.Database.OS = this.OS
	this.router.Database.Init = lInit
	this.router.Database.PGSQL = lPGSQL
	if this.OS == "android" {
		this.router.Database.OSfilepath = "/sdcard/com.valued.app/"
	}

	this.router.Database.AndroidPipe = this.AndroidPipe
	_, success := this.router.Connect()
	if success {
		if this.router.Database.Init {
			new(database.Menu).Initialize(this.router.Database)
			new(database.Profile).Initialize(this.router.Database)
			new(database.Media).Initialize(this.router.Database)
			new(database.Industry).Initialize(this.router.Database)
			new(database.Permission).Initialize(this.router.Database)
			new(database.ActivationLink).Initialize(this.router.Database)
			new(database.PendingApproval).Initialize(this.router.Database)

			// new(database.Report).Initialize(this.router.Database)
			//new
			new(database.Profile).Initialize(this.router.Database)
			new(database.Redemption).Initialize(this.router.Database)

			new(database.Favorite).Initialize(this.router.Database)

			new(database.Review).Initialize(this.router.Database)
			new(database.ReviewCategory).Initialize(this.router.Database)
			new(database.ReviewCategoryLink).Initialize(this.router.Database)

			new(database.TelrCode).Initialize(this.router.Database)
			new(database.TelrOrder).Initialize(this.router.Database)
			new(database.TelrCountry).Initialize(this.router.Database)

			new(database.Coupon).Initialize(this.router.Database)
			new(database.Reward).Initialize(this.router.Database)
			new(database.Feedback).Initialize(this.router.Database)
			new(database.Redemption).Initialize(this.router.Database)

			new(database.Groups).Initialize(this.router.Database)
			new(database.EmployerGroup).Initialize(this.router.Database)
			new(database.MemberGroup).Initialize(this.router.Database)
			new(database.RewardGroup).Initialize(this.router.Database)
			new(database.Referral).Initialize(this.router.Database)

			new(database.Category).Initialize(this.router.Database)
			new(database.CategoryLink).Initialize(this.router.Database)
			new(database.RewardStore).Initialize(this.router.Database)
			new(database.RewardScheme).Initialize(this.router.Database)

			new(database.Scheme).Initialize(this.router.Database)
			new(database.Smtp).Initialize(this.router.Database)
			new(database.Store).Initialize(this.router.Database)
			new(database.Subscription).Initialize(this.router.Database)
		}

		switch this.OS {
		case "android", "ios":
			// new(Synchronize).Run(this.router.Database)
			break
		}

		exitChan := make(chan os.Signal, 1)
		signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

		this.startHTTP()
		go this.startWebserver()
		//Blocking startWebserver is in a goroutine

		for range exitChan {
			this.router.Disconnect()
			os.Exit(1)
		}

	} else {
		println("!!!Database Connection Error!!!")
	}
}

func (this Server) redirHTTP(httpRes http.ResponseWriter, httpReq *http.Request) {
	http.Redirect(httpRes, httpReq, "https://"+httpReq.Host+":"+this.cPort+httpReq.RequestURI, http.StatusMovedPermanently)
}

func (this Server) startHTTP() {
	if this.cRedirect != "" {
		go http.ListenAndServe(":"+this.cRedirect, http.HandlerFunc(this.redirHTTP))
	}
}

func (this Server) startWebserver() {
	println("Server Starting...")

	mime.AddExtensionType(".js", "application/x-javascript; charset=utf-8")
	mime.AddExtensionType(".css", "text/css; charset=utf-8")
	mime.AddExtensionType(".svg", "image/svg+xml; charset=utf-8")
	mime.AddExtensionType(".doc", "application/msword; charset=utf-8")
	mime.AddExtensionType(".json", "application/json; charset=utf-8")

	http.HandleFunc("/", this.router.Route)

	server := http.Server{Addr: ":" + this.cPort}

	if this.lSSL {
		certPEMBlock, err := data.Asset("conf/cert")
		if err != nil {
			println("Cert error ... " + err.Error())
			return
		}

		keyPEMBlock, err := data.Asset("conf/key")
		if err != nil {
			println("Key error ... " + err.Error())
			return
		}

		server.TLSConfig = &tls.Config{Certificates: make([]tls.Certificate, 1)}
		server.TLSConfig.Certificates[0], err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			println("X509KeyPair Cert/Key error ... " + err.Error())
			return
		}
		//Load embedded certificate
		if err := server.ListenAndServeTLS("", ""); err != nil {
			println("Server Error-> ", err.Error())
		}

	} else {
		if err := server.ListenAndServe(); err != nil {
			println("Server Error-> ", err.Error())
		}
	}

	println("Server Stoping...")
	this.router.Disconnect()
	os.Exit(1)
}
