package main

import (
	"database/sql"
	"encoding/json"
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

func (instance *datastore) initDatastore() {
	log.Debug("initiating datastore")
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
}

func check(e error) {
	if e != nil {
		// use Dave Cheney's errors package to Wrap err with a message and stack trace before propagating it
		panic(errors.Wrap(e, "something failed"))
	}
}

func (instance *datastore) insertTest() {
	if instance.psqlDB == nil {
		instance.connectToDB()
	}
	title := "new title"
	insertStmt := `insert into "storyboard" ("id", "title") values ($1, $2)`
	_, e := instance.psqlDB.Exec(insertStmt, 3, title)
	//_, e := instance.psqlDB.Exec(insertStmt, instance.id, title)
	check(e)
}

func (instance *datastore) CreateStory(data story) {
	if instance.psqlDB == nil {
		instance.connectToDB()
	}
	curTime := (time.Now().UnixNano() / int64(time.Millisecond))
	insertStmt := `insert into "storyboard" ("id", "title", "created_at", "updated_at") values ($1, $2, $3, $4)`
	_, e := instance.psqlDB.Exec(insertStmt, data.Id, data.Title, curTime, curTime)
	check(e)
}

func (instance *datastore) UpdateTitle(data story) {
	if instance.psqlDB == nil {
		instance.connectToDB()
	}
	curTime := (time.Now().UnixNano() / int64(time.Millisecond))
	insertStmt := `update storyboard set title = $1, updated_at = $2 where id = $3`
	_, e := instance.psqlDB.Exec(insertStmt, data.Title, curTime, data.Id)
	check(e)
}

func (instance *datastore) GetNextId() (int64, error) {
	if instance.psqlDB == nil {
		instance.connectToDB()
	}
	rows, e := instance.psqlDB.Query(`SELECT coalesce(MAX(id), 0) AS id FROM storyboard`)
	check(e)
	defer rows.Close()
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		check(err)
		rv := id + 1
		return rv, nil
	}
	return 1, fmt.Errorf("db error")
}

func (instance *datastore) UpdateTime(id int64) {
	curTime := (time.Now().UnixNano() / int64(time.Millisecond))
	updateStmt := `update storyboard set updated_at = $1 where id = $2`
	_, e := instance.psqlDB.Exec(updateStmt, curTime, id)
	check(e)
}

func (instance *datastore) UpdateBody(data story) {
	b, _ := json.Marshal(data.Body)
	curTime := (time.Now().UnixNano() / int64(time.Millisecond))
	updateStmt := `update storyboard set updated_at = $1, body = $2 where id = $3`
	_, e := instance.psqlDB.Exec(updateStmt, curTime, b, data.Id)
	check(e)
	log.Debug("Updated body", string(b))
}

func (instance *datastore) CheckStoryExist(id int64) bool {
	if instance.psqlDB == nil {
		instance.connectToDB()
	}
	queryStmnt := `select COUNT(1) from storyboard where id=$1`
	rows, e := instance.psqlDB.Query(queryStmnt, id)
	check(e)
	defer rows.Close()
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		check(err)
		rv := id == 1
		return rv
	}
	return false
}

func (instance *datastore) GetStoryById(Id int64) (int, RespSingleStory) {
	if instance.psqlDB == nil {
		instance.connectToDB()
	}
	queryStmnt := `SELECT id, title, created_at, updated_at,body FROM public.storyboard where id = $1`
	rows, e := instance.psqlDB.Query(queryStmnt, Id)
	defer rows.Close()
	check(e)
	var output RespSingleStory
	for rows.Next() {
		var id int64
		var title sql.NullString
		var created_at, updated_at int64
		var body sql.NullString
		err := rows.Scan(&id, &title, &created_at, &updated_at, &body)
		output.Id = id
		if title.Valid {
			output.Title = title.String
		}
		output.CreatedTime = time.Unix(0, created_at*int64(time.Millisecond)).String()
		output.UpdateTime = time.Unix(0, updated_at*int64(time.Millisecond)).String()
		if body.Valid {
			output.Paragraphs = body.String
		}

		check(err)

		return http.StatusOK, output
	}

	return http.StatusUnprocessableEntity, output
}

func (instance *datastore) GetStories(limit, offset uint64, sortby, orderby string) (int, MulityStoryResponse) {
	if instance.psqlDB == nil {
		instance.connectToDB()
	}
	log.Debug("limit = ", limit, " offset = ", offset, " sortby = ", sortby, "order by =", orderby)
	var queryStmnt string
	var err error
	var rows *sql.Rows
	var output MulityStoryResponse

	sortKey := ""
	if sortby != "" {
		sortKey = sortby + " " + orderby
		log.Debug("Sort key is ", sortKey)
	}

	if limit != 0 && offset != 0 && sortby != "" {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard ORDER BY $1 LIMIT $2 OFFSET $3`
		rows, err = instance.psqlDB.Query(queryStmnt, sortKey, limit, offset)
	} else if limit != 0 && offset != 0 {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard LIMIT $1 OFFSET $2`
		rows, err = instance.psqlDB.Query(queryStmnt, limit, offset)
	} else if limit != 0 && sortby != "" {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard ORDER BY $1 LIMIT $2`
		rows, err = instance.psqlDB.Query(queryStmnt, sortKey, limit)
	} else if offset != 0 && sortby != "" {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard ORDER BY $1 OFFSET $2`
		rows, err = instance.psqlDB.Query(queryStmnt, sortKey, offset)
	} else if limit != 0 {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard LIMIT $1`
		rows, err = instance.psqlDB.Query(queryStmnt, limit)
	} else if offset != 0 {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard OFFSET $1`
		rows, err = instance.psqlDB.Query(queryStmnt, offset)
	} else if sortby != "" {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard ORDER BY $1`
		rows, err = instance.psqlDB.Query(queryStmnt, sortKey)
	} else {
		queryStmnt = `SELECT id, title, created_at, updated_at FROM storyboard`
		rows, err = instance.psqlDB.Query(queryStmnt)
	}

	check(err)
	for rows.Next() {
		var id int64
		var title sql.NullString
		var created_at, updated_at int64
		err := rows.Scan(&id, &title, &created_at, &updated_at)
		check(err)
		var result ResponseResult
		result.Id = id
		if title.Valid {
			result.Title = title.String
		}
		result.CreatedTime = time.Unix(0, created_at*int64(time.Millisecond)).String()
		result.UpdateTime = time.Unix(0, updated_at*int64(time.Millisecond)).String()
		output.Results = append(output.Results, result)
	}
	output.Count = len(output.Results)
	output.Limit = limit
	output.Offset = offset
	defer rows.Close()
	log.Debug("Completed GetStories successfully")
	return http.StatusOK, output
}
