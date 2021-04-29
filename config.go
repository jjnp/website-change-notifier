package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"time"
	"website-change-notifier/pushover"
)

type SiteConfig struct {
	Name            string        `yaml:"name"`
	Url             string        `yaml:"url"`
	Interval        time.Duration `yaml:"interval"`
	SummaryInterval time.Duration `yaml:"summary-interval"`
}

type LogConfig struct {
	Level logrus.Level `yaml:"level"`
}

type Config struct {
	Site     SiteConfig      `yaml:"site"`
	Pushover pushover.Config `yaml:"pushover"`
	Log      LogConfig       `yaml:"log,omitempty"`
}

func loadConfig() *Config {
	logrus.Debugln("opening config file")
	f, err := openConfigFile()
	if err != nil {
		panic("Couldn't open config file!")
	}
	defer f.Close()
	logrus.Debugln("parsing config file")
	config, err := parseConfigFile(f)
	if err != nil {
		panic("Error parsing config. Make sure it is a valid yaml file!")
	}
	logrus.Debugln("settings defaults where config values are missing")
	config.setDefaultsWhenMissing()
	return config
}

func openConfigFile() (*os.File, error) {
	path, present := os.LookupEnv("CONFIG_FILE")
	if !present {
		path = "/config.yml"
	}
	return os.Open(path)
}

func parseConfigFile(file *os.File) (*Config, error) {
	var config Config
	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(&config)
	return &config, err
}

func (c *Config) setDefaultsWhenMissing() {

}
