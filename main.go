package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	log.Info("Starting!")
	db := datastore{id: 1}
	db.initDatastore()
	setupRoutes(db)
}
