package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	logLevl := os.Getenv("VERLOOP_DEBUG")
	if logLevl == "" {
		logLevl = "INFO"

	}
	lv, _ := log.ParseLevel(logLevl)
	log.SetLevel(lv)
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	log.Info("Starting!")
	db := datastore{}
	db.initDatastore()
	setupRoutes(db)
}
