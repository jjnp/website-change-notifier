package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"website-change-notifier/pushover"
)

func main() {
	logrus.Debugln("Loading config")
	config := loadConfig()
	logrus.Infof("config loaded")
	configureLogger(&config.Log)
	notifier := pushover.New(&config.Pushover)
	worker, err := NewWorker(&config.Site, notifier)
	if err != nil {
		panic(err)
	}
	worker.start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	worker.stop()
	os.Exit(0)
}

func configureLogger(config *LogConfig) {
	logrus.Debugf("setting log level to %s", config.Level.String())
	logrus.SetLevel(config.Level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
