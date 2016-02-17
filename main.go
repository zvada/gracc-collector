package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "gratia.cfg", "config file")
}

func main() {
	flag.Parse()

	log.WithField("file", configFile).Info("reading config")
	config, err := ReadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	log.WithField("level", config.LogLevel).Info("setting log level")
	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	log.WithFields(log.Fields{
		"brokers": config.Kafka.Brokers,
		"topic":   config.Kafka.Topic,
	}).Info("starting kafka collector")
	g, err := NewCollector(config.Kafka.Brokers, config.Kafka.Topic)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/gratia-servlets/rmi", g)

	log.WithFields(log.Fields{
		"address": config.Address,
		"port":    config.Port,
	}).Info("listening")

	log.Fatal(http.ListenAndServe(config.Address+":"+config.Port, nil))
}
