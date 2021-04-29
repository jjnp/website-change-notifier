package pushover

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Config struct {
	Token  string `yaml:"token"`
	User   string `yaml:"user"`
	Device string `yaml:"device"`
}

func New(config *Config) *Pushover {
	logrus.Debugln("creating pushover instance")
	return &Pushover{
		token:  config.Token,
		user:   config.User,
		device: config.Device,
	}
}

type Pushover struct {
	token  string
	user   string
	device string
}

type messageBody struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Device  string `json:"device"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

func (p Pushover) NotifyOfChange(title string, message string, url string) {
	logrus.WithField("url", url).Infoln("notifying user of site change")
	body := messageBody{
		Token:   p.token,
		User:    p.user,
		Device:  p.device,
		Title:   title,
		Message: message,
		Url:     url,
	}
	bodyJson, _ := json.Marshal(body)

	_, err := http.Post("https://api.pushover.net/1/messages.json", "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		logrus.Errorln(err)
	}
}

func (p Pushover) NotifyOfStart(message string, url string) {
	logrus.WithField("url", url).Infoln("notifying user that monitoring started")
	body := messageBody{
		Token:   p.token,
		User:    p.user,
		Device:  p.device,
		Title:   "Started Monitoring...",
		Message: message,
		Url:     url,
	}
	bodyJson, _ := json.Marshal(body)

	_, err := http.Post("https://api.pushover.net/1/messages.json", "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		logrus.Errorln(err)
	}
}

func (p Pushover) SendSummary(title string, message string, url string) {
	logrus.WithField("url", url).Infoln("sending summary to user")
	body := messageBody{
		Token:   p.token,
		User:    p.user,
		Device:  p.device,
		Title:   title,
		Message: message,
		Url:     url,
	}
	bodyJson, _ := json.Marshal(body)

	_, err := http.Post("https://api.pushover.net/1/messages.json", "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		logrus.Errorln(err)
	}
}
