package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type datastore struct {
	HTTPClient http.Client
	config     readConfig
	psqlDB     *sql.DB
}

type readConfig struct {
	ConfigFile       string
	LogLevel         string
	DeploymentMode   string
	AllowCrossOrigin bool
	CharSet          string
}

func (instance *datastore) initDatastore() {
	fmt.Println("initiating datastore")
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	switch err.(type) {
	case viper.UnsupportedConfigError:
		log.Info("No config file, using defaults")
	default:
		check(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:" + e.Name)
		err = viper.Unmarshal(&instance.config)
		check(err)
	})
	instance.connectToDB()
}

func (instance *datastore) connectToDB() {

	if instance.psqlDB != nil {
		return
	}

	connstring := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
		viper.GetString("Database.IP"),
		viper.GetInt("Database.Port"),
		viper.GetString("Database.UserID"),
		viper.GetString("Database.Passwd"),
		viper.GetString("Database.Name"))

	log.WithFields(log.Fields{
		"connstring": connstring,
	}).Info("Connection String")

	var err error
	instance.psqlDB, err = sql.Open("postgres", connstring)
	if err != nil {
		log.WithFields(log.Fields{
			"db_error": err,
		}).Error("DB Error")
	}
	for instance.psqlDB.Ping() != nil {
		log.Info("Connection to Postgres was lost. Waiting for 5s...")
		instance.psqlDB.Close()
		time.Sleep(5 * time.Second)
		log.Info("Reconnecting...")
		instance.psqlDB, err = sql.Open("postgres", connstring)
		if err != nil {
			log.WithFields(log.Fields{
				"db_error": err,
			}).Error("DB Error")
		}
	}

	log.Info("Successfully connected!")
	queryResp, err := instance.psqlDB.Query(`list tables for all`)
	if queryResp != nil {
		fmt.Println("got resp")
	}
}

func check(e error) {
	if e != nil {
		// use Dave Cheney's errors package to Wrap err with a message and stack trace before propagating it
		panic(errors.Wrap(e, "something failed"))
	}
}

func (instance *datastore) insertTest() {
	for instance.psqlDB == nil {
		instance.connectToDB()
	}
	insertStmt := `insert into "storyboard"("id", "title") values('1', 'renjurenju')`
	_, e := instance.psqlDB.Exec(insertStmt)
	check(e)
}
