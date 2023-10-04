package main

import (
	"fmt"
	"net/smtp"
	"testing"
)

func TestSendMail(t *testing.T) {
	// Sender data.
	from := "mail"
	password := "pushimmail"

	// Receiver email address.
	to := []string{
		"1942056324@qq.com",
	}

	// smtp server configuration.
	smtpHost := "88.218.195.129"
	smtpPort := "25"

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

// func SendMail(addr, from, subject, body string, to []string) error {
// 	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

// 	c, err := smtp.Dial(addr)
// 	if err != nil {
// 		return err
// 	}
// 	defer c.Close()
// 	if err = c.SendMail(r.Replace(from)); err != nil {
// 		return err
// 	}
// 	for i := range to {
// 		to[i] = r.Replace(to[i])
// 		if err = c.Rcpt(to[i]); err != nil {
// 			return err
// 		}
// 	}

// 	w, err := c.Data()
// 	if err != nil {
// 		return err
// 	}

// 	msg := "To: " + strings.Join(to, ",") + "\r\n" +
// 		"From: " + from + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
// 		"Content-Transfer-Encoding: base64\r\n" +
// 		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

// 	_, err = w.Write([]byte(msg))
// 	if err != nil {
// 		return err
// 	}
// 	err = w.Close()
// 	if err != nil {
// 		return err
// 	}
// 	return c.Quit()
// }
