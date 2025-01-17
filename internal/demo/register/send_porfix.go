package register

import (
	"fmt"
	"os/exec"
)

// ex: SendMail("127.0.0.1:25", (&mail.Address{"from name", "from@example.com"}).String(), "Email Subject", "message body", []string{(&mail.Address{"to name", "to@example.com"}).String()})
// func SendMail(addr, from, subject, body string, to []string) error {
// 	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

// 	c, err := smtp.Dial(addr)
// 	if err != nil {
// 		return err
// 	}
// 	defer c.Close()
// 	if err = c.Mail(r.Replace(from)); err != nil {
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

//		_, err = w.Write([]byte(msg))
//		if err != nil {
//			return err
//		}
//		err = w.Close()
//		if err != nil {
//			return err
//		}
//		return c.Quit()
//	}
func SendMail(addr, from, subject, body string, to string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo %s | mail -s "验证码" %s`, body, to))
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%+v", err)
	}
	return err
}
