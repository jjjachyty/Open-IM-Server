package register

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/gomail.v2"
)

func init() {
	os.Setenv("CONFIG_NAME", "/Users/janly/data/go/src/Open-IM-Server/")
	os.Setenv("USUAL_CONFIG_NAME", "/Users/janly/data/go/src/Open-IM-Server/")
}
func TestSendVerificationCode(t *testing.T) {

	m := gomail.NewMessage()
	m.SetHeader(`From`, "push131.com@hotmail.com")
	m.SetHeader(`To`, "jjjachyty@163.com")
	m.SetHeader(`Subject`, "验证码")
	m.SetBody(`text/html`, fmt.Sprintf("%d", 123456))
	if err := gomail.NewDialer("smtp-mail.outlook.com", 587, "push131.com@hotmail.com", "").DialAndSend(m); err != nil {
		panic(err)
	}
}
