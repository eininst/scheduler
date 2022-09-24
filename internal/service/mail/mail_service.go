package mail

import (
	"encoding/json"
	"github.com/eininst/scheduler/configs"
	"gopkg.in/gomail.v2"
)

type MailConfig struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Cc       []string `json:"cc"`
}

type MailService interface {
	Send(toAddr string, subject string, msg string) error
	IsConfig() bool
}

type mailService struct {
	Config *MailConfig
}

func NewService() MailService {
	var mcfg *MailConfig
	mstr := configs.Get("mail").String()
	if mstr != "" {
		er := json.Unmarshal([]byte(mstr), &mcfg)
		if er != nil {
			mcfg = nil
		} else {
			if mcfg.Port == 0 {
				mcfg.Port = 465
			}

			if mcfg.Host == "" {
				mcfg = nil
			}

			if mcfg.Username == "" {
				mcfg = nil
			}

			if mcfg.Password == "" {
				mcfg = nil
			}
		}
	}

	return &mailService{
		Config: mcfg,
	}
}

func (m *mailService) IsConfig() bool {
	return m.Config != nil
}

func (ms *mailService) Send(toAddr string, subject string, msg string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", ms.Config.Username)
	m.SetHeader("To", toAddr)

	//抄送
	if len(ms.Config.Cc) > 0 {
		ccList := []string{}
		for _, val := range ms.Config.Cc {
			if val == toAddr {
				continue
			}
			ccList = append(ccList, m.FormatAddress(val, ""))
		}
		if len(ccList) > 0 {
			m.SetHeader("Cc", ccList...)
		}
	}

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(ms.Config.Host, ms.Config.Port, ms.Config.Username, ms.Config.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
