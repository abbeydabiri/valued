package functions

import (
	"regexp"
	"valued/data"

	"crypto/tls"
	"fmt"
	"html"
	"strings"

	"gopkg.in/gomail.v2"

	"encoding/json"
	"io/ioutil"
	"strconv"

	"log"
)

func GenerateEmail(emailFrom, emailFromName, emailTo, emailSubject, emailTemplate string, emailCC []string, emailFields map[string]interface{}) bool {

	emailBytes, _ := data.Asset("email/" + emailTemplate)

	if emailBytes != nil {
		emailMessage := string(emailBytes)
		for cFieldname, iFieldvalue := range emailFields {
			tagsToReplace, _ := regexp.Compile(`\[@` + cFieldname + `@\]`)

			if iFieldvalue == nil {
				continue
			}

			sFieldvalue := fmt.Sprintf("%v", iFieldvalue)
			if cFieldname != sFieldvalue {
				emailMessage = tagsToReplace.ReplaceAllString(emailMessage, sFieldvalue)
			} else {
				emailMessage = tagsToReplace.ReplaceAllString(emailMessage, "")
			}
		}

		fileFields, _ := regexp.Compile(`\[@.\S*?@\]`)
		emailMessage = fileFields.ReplaceAllString(emailMessage, "")

		return SendEmail(emailFrom, emailFromName, emailTo, emailSubject, emailMessage, emailCC)
	}

	return false
}

func SendEmail(emailFrom, emailFromName, emailTo, emailSubject, emailMessage string, emailCC []string) bool {

	if emailTo == "" || emailSubject == "" || emailMessage == "" {
		return false
	}

	jsonEmail, _ := ioutil.ReadFile("email.json")
	if len(jsonEmail) == 0 {
		log.Printf(`email.json file is missing`)
		return false
	}

	mapEmail := make(map[string]string)
	json.Unmarshal(jsonEmail, &mapEmail)
	if len(mapEmail) == 0 {
		log.Printf(`email.json file is corrupt`)
		return false
	}

	Port, _ := strconv.Atoi(mapEmail["port"])
	mySMTP := SMTP{
		Port: Port, Server: mapEmail["server"], Username: mapEmail["username"], Password: mapEmail["password"],
	}

	// if emailFrom == "" {
	emailFrom = mySMTP.Username
	// }

	if emailFromName == "" {
		emailFromName = "VALUED MEMBERSHIP"
	}
	emailSender := fmt.Sprintf("%s <%s>", emailFromName, emailFrom)

	emailBCC := ""
	// emailBCC := []string{"general@valued.com"}
	// emailBCC := []string{"info@valued.com", "general@valued.com"}
	// if emailTo == "info@valued.com" {
	// 	emailBCC = []string{"general@valued.com"}
	// }

	var myMsgList []Message
	myMsgList = append(myMsgList,
		Message{
			Attachment: "",
			To:         emailTo,
			From:       emailSender,
			Cc:         emailCC, Bcc: emailBCC, Replyto: "info@valued.com",
			Subject: emailSubject,
			Content: emailMessage,
		})
	mailer := Mailer{mySMTP, myMsgList}

	log.Printf(" - - -- - - - -- - -- - --- - \n Mail:  %v ", myMsgList)

	sMessage := mailer.CheckMail()
	if len(sMessage) > 0 {
		log.Printf(sMessage)
		return false
	}

	sMessage = mailer.SendMail()
	if len(sMessage) > 0 {
		log.Printf(sMessage)
		return false
	}

	return true
}

type SMTP struct {
	Port                       int
	Server, Username, Password string
}

type Message struct {
	To, From, Replyto, Subject,
	Content, Attachment string
	Cc, Bcc []string
}

type Mailer struct {
	SMTP
	MessageList []Message
}

func (this *Mailer) SendMail() (sMessage string) {

	sMessage = ""
	goMsg := gomail.NewMessage()

	goDialer := gomail.NewDialer(this.SMTP.Server, this.SMTP.Port, this.SMTP.Username, this.SMTP.Password)
	goDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	goSender, err := goDialer.Dial()
	if err != nil {
		sMessage = strings.Replace(err.Error(), "\n", "", -1)
		sMessage = html.EscapeString(sMessage)
		return
	}

	for _, Msg := range this.MessageList {

		goMsg.SetHeader("To", Msg.To)
		goMsg.SetHeader("From", Msg.From)

		for _, cc := range Msg.Cc {
			goMsg.SetHeader("Cc", cc)
		}

		for _, bcc := range Msg.Bcc {
			goMsg.SetHeader("Bcc", bcc)
		}

		if Msg.Replyto != "" {
			goMsg.SetHeader("Reply-to", Msg.Replyto)
		}

		goMsg.SetHeader("Subject", Msg.Subject)
		goMsg.SetBody("text/html", Msg.Content)
		if Msg.Attachment != "" {
			goMsg.Attach(Msg.Attachment)
		}

		if err := gomail.Send(goSender, goMsg); err != nil {
			sMessage = strings.Replace(err.Error(), "\n", "", -1)
			sMessage = html.EscapeString(sMessage)
		}
		goMsg.Reset()
	}

	return
}

func (this *Mailer) CheckMail() (sMessage string) {

	sMessage = ""

	if this.SMTP.Port == 0 {
		sMessage += "SMTP.Port is blank <br>"
	}

	if this.SMTP.Server == "" {
		sMessage += "SMTP.Server is blank <br>"
	}

	if this.SMTP.Username == "" {
		sMessage += "SMTP.Username is blank <br>"
	}

	if this.SMTP.Password == "" {
		sMessage += "SMTP.Password is blank <br>"
	}

	for Key, Msg := range this.MessageList {

		if Msg.To == "" {
			sMessage += fmt.Sprintf("Msg %d To is blank <br>", Key)
		}

		if Msg.From == "" {
			sMessage += fmt.Sprintf("Msg %d From is blank <br>", Key)
		}

		if Msg.Subject == "" {
			sMessage += fmt.Sprintf("Msg %d Subject is blank <br>", Key)
		}

		if Msg.Content == "" {
			sMessage += fmt.Sprintf("Msg %d Content is blank <br>", Key)
		}
	}

	return
}
