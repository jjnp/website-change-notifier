package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Worker struct {
	name            string
	url             string
	interval        time.Duration
	summaryInterval time.Duration
	notifier        Notifier
	lastHash        string
	checkCounter    int
	lastCheck       time.Time
	ticker          *time.Ticker
	summaryTicker   *time.Ticker
	quit            chan struct{}
}

func NewWorker(config *SiteConfig, notifier Notifier) (*Worker, error) {
	logrus.WithField("url", config.Url).WithField("name", config.Name).WithField("interval", config.Interval).WithField("summary interval", config.SummaryInterval).Infoln("creating new worker")
	if config.Interval < time.Second {
		return nil, errors.New("interval is an invalid duration. must be >= 1s")
	}
	if config.Url == "" {
		return nil, errors.New("url was not supplied in site config")
	}
	if config.Name == "" {
		return nil, errors.New("name was not supplied in site config")
	}
	if notifier == nil {
		return nil, errors.New("notifier passed to worker is nil")
	}
	return &Worker{
		name:            config.Name,
		url:             config.Url,
		interval:        config.Interval,
		checkCounter:    0,
		summaryInterval: config.SummaryInterval,
		notifier:        notifier,
	}, nil
}

func (w *Worker) start() {
	logrus.WithField("url", w.url).Debugln("starting worker")
	w.ticker = time.NewTicker(w.interval)
	if w.summaryInterval > time.Second {
		w.summaryTicker = time.NewTicker(w.summaryInterval)
	} else {
		w.summaryTicker = time.NewTicker(time.Duration(7*24) * time.Hour)
	}
	w.quit = make(chan struct{})
	w.initHash()
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.check()
			case <-w.summaryTicker.C:
				w.summary()
			case <-w.quit:
				w.ticker.Stop()
				logrus.WithField("url", w.url).Debugln("stopped worker")
				return
			}
		}
	}()
}

func (w *Worker) initHash() {
	w.lastHash, _ = getSiteHash(w.url)
	w.notifier.NotifyOfStart(fmt.Sprintf("Started monitoring %s", w.name), w.url)
}

func (w *Worker) check() {
	logrus.WithField("url", w.url).Debugln("checking for changes")
	if w.hasChanged() {
		logrus.WithField("url", w.url).Infoln("site changed")
		w.notifyOfChange()
		return
	}
	logrus.WithField("url", w.url).Infoln("no change detected")
}

func (w *Worker) summary() {
	logrus.WithField("url", w.url).Infoln("preparing summary")
	lastCheckDelta := time.Now().Sub(w.lastCheck).String()
	w.notifier.SendSummary(fmt.Sprintf("%s summary", w.name), fmt.Sprintf("%d checks in the last %s. Last check was %s ago", w.checkCounter, w.summaryInterval.String(), lastCheckDelta), w.url)
	w.checkCounter = 0
}

func (w *Worker) notifyOfChange() {
	w.notifier.NotifyOfChange("Site changed!", fmt.Sprintf("%s has changed", w.name), w.url)
}

func (w *Worker) hasChanged() bool {
	nextHash, err := getSiteHash(w.url)
	if err != nil {
		return false
	}
	w.checkCounter++
	w.lastCheck = time.Now()
	if nextHash != w.lastHash {
		w.lastHash = nextHash
		return true
	}
	return false
}

func (w *Worker) stop() {
	logrus.WithField("url", w.url).Debugln("stopping worker")
	w.quit <- struct{}{}
}

func getSiteHash(url string) (string, error) {
	res, err := getSite(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := getBody(res)
	if err != nil {
		return "", err
	}
	return sha256Hash(body), nil
}

func getSite(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	return res, nil
}

func getBody(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	return body, nil
}

func sha256Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
