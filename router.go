package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupRoutes(db datastore) {
	app := iris.New()

	cur_story := story{Id: 0, DBPtr: db}
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	verl := app.Party("/", crs).AllowMethods(iris.MethodOptions)
	{
		verl.Post("/add", func(ctx iris.Context) {
			var msg incomingMsg
			err := ctx.ReadJSON(&msg)
			if err != nil {
				log.WithFields(log.Fields{
					"invalidReq_error": err,
				}).Error("Unable to parse request from client!")

				ctx.Values().Set("message", "Malformed packet or empty terminal id")
				ctx.StatusCode(http.StatusUnprocessableEntity)
				return
			}
			status, resp := cur_story.addWord(msg.Word)
			ctx.StatusCode(status)
			ctx.JSON(resp)
		})

		verl.Get("/stories/{id}", func(ctx iris.Context) {

			id := ctx.Params().Get("id")
			log.Info("got request to get story by id", id)
			i, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				panic(err)
			}

			if !db.CheckStoryExist(i) {
				log.Info("story not found with id = ", i)
				ctx.StatusCode(http.StatusNotFound)
				return
			}
			status, resp := db.GetStoryById(i)
			ctx.StatusCode(status)
			ctx.JSON(resp)
		})

		verl.Get("/stories/", func(ctx iris.Context) {
			fmt.Println("get whole stories")
			var limit, offset uint64 = 0, 0
			var sortby, orderby string = "", "ASC"
			var err error
			if ctx.URLParamExists("limit") {
				tt := ctx.URLParam("limit")
				limit, err = strconv.ParseUint(tt, 10, 64)
				if err != nil {
					panic(err)
				}
				log.Info("limit is ", limit)
			}
			if ctx.URLParamExists("offset") {
				tt := ctx.URLParam("offset")
				offset, err = strconv.ParseUint(tt, 10, 64)
				if err != nil {
					panic(err)
				}
				log.Info("offset is ", offset)
			}
			if ctx.URLParamExists("sort") {
				sortby = ctx.URLParam("sort")
				log.Info("sortby is ", sortby)
			}
			if ctx.URLParamExists("order") {
				orderby = ctx.URLParam("order")
				log.Info("orderby is ", orderby)
			}

			status, resp := db.GetStories(limit, offset, sortby, orderby)
			ctx.StatusCode(status)
			ctx.JSON(resp)

		})

	}
	port := fmt.Sprintf(":%s", viper.GetString("Verloop.Port"))
	app.Run(iris.Addr(port), iris.WithCharset(viper.GetString("CharSet")))
}

func notFoundHandler(ctx iris.Context) {
	ctx.HTML("<b>404 not found</b>")
}
