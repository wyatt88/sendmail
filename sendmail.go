package sendmail

import (
	"encoding/base64"
	"log"
	"net/smtp"
	"strings"
)

func SendPlainMail(addr, from, subject, body string, to []string, cc []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")
	var realTo []string
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	if len(to) > 0 {
		for i := range to {
			realTo = append(realTo,to[i])
		}
	}
	if len(cc) > 0 {
		for i := range cc {
			realTo = append(realTo,cc[i])
		}
	}
	for i := range realTo {
		realTo[i] = r.Replace(realTo[i])
		if err = c.Rcpt(realTo[i]); err != nil {
			return err
		}
	}
	
	w, err := c.Data()
	if err != nil {
		return err
	}
	
	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"Cc: " + strings.Join(cc, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
