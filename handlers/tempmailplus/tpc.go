package tempmailplus

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

var (
	Domains = []string{
		"mailto.plus",
		"fexpost.com",
		"fexbox.org",
		"mailbox.in.ua",
		"rover.info",
		"chitthi.in",
		"fextemp.com",
		"any.pink",
		"merepost.com",
	}
)

type TmailPlusClient struct {
	Domain   string
	Username string
}
type Inbox struct {
	Count    int `json:"count"`
	FirstId  int `json:"first_id"`
	LastId   int `json:"last_id"`
	Limit    int `json:"limit"`
	MailList []struct {
		AttachmentCount     int    `json:"attachment_count"`
		FirstAttachmentName string `json:"first_attachment_name"`
		FromMail            string `json:"from_mail"`
		FromName            string `json:"from_name"`
		IsNew               bool   `json:"is_new"`
		MailId              int    `json:"mail_id"`
		Subject             string `json:"subject"`
		Time                string `json:"time"`
	} `json:"mail_list"`
	More   bool `json:"more"`
	Result bool `json:"result"`
}

type Email struct {
	Attachments []interface{} `json:"attachments"`
	Date        string        `json:"date"`
	From        string        `json:"from"`
	FromIsLocal bool          `json:"from_is_local"`
	FromMail    string        `json:"from_mail"`
	FromName    string        `json:"from_name"`
	Html        string        `json:"html"`
	IsTls       bool          `json:"is_tls"`
	MailId      int           `json:"mail_id"`
	MessageId   string        `json:"message_id"`
	Result      bool          `json:"result"`
	Subject     string        `json:"subject"`
	Text        string        `json:"text"`
	To          string        `json:"to"`
}

type Account struct {
	Email    string
	Domain   string
	Username string
}

func (tpc TmailPlusClient) CheckInbox() Inbox {
	BaseLink := "https://tempmail.plus/api/mails?email=" + tpc.Username + "@" + tpc.Domain + "&first_id=0&epin="
	fmt.Println("Base > " + BaseLink)
	client := &http.Client{}
	req, err := http.NewRequest("GET", BaseLink, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("cookie", "email=oahycy%40mailto.plus")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://tempmail.plus/en/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var InboxJsonData Inbox
	err = json.Unmarshal(bodyText, &InboxJsonData)
	if err != nil {
		log.Fatal(err)
	}
	// for _, mail := range InboxJsonData.MailList {
	// 	// fmt.Printf("From: %s (%s)\nSubject: %s\nTime: %s\nAttachments: %d (%s)\n\n",
	// 	// 	mail.FromName, mail.FromMail, mail.Subject, mail.Time, mail.AttachmentCount, mail.FirstAttachmentName)
	// 	tpc.GetMessageByID(mail.MailId)
	// 	fmt.Println("=========")
	// }
	return InboxJsonData
}

func (tpc TmailPlusClient) GetMessageByID(id int) Email {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://tempmail.plus/api/mails/"+strconv.Itoa(id)+"?email="+tpc.Username+"@"+tpc.Domain+"&epin=", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("cookie", "email=oahycy%40mailto.plus")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://tempmail.plus/en/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var Email Email
	err = json.Unmarshal(bodyText, &Email)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("From: %s (%s)\nSubject: %s\nDate: %s\nHTML: %s\nText: %s\n", Email.FromName, Email.FromMail, Email.Subject, Email.Date, Email.Html, Email.Text)
	return Email
}

func (tpc TmailPlusClient) NewAccount() Account {
	username := gofakeit.Username()
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(Domains))
	domain := Domains[idx]
	return Account{
		Email:    username + "@" + domain,
		Domain:   domain,
		Username: username,
	}
}
