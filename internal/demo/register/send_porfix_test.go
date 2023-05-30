package register

import (
	"net/mail"
	"testing"
)

func TestSendMail(t *testing.T) {

	err := SendMail("88.218.195.129:25", (&mail.Address{"from name", "from@example.com"}).String(), "Email Subject", "message body", []string{(&mail.Address{"to name", "1942056324@qq.com"}).String()})
	panic(err)
}
